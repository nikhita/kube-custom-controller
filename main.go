package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"golang.org/x/oauth2"

	"github.com/google/go-github/github"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/workqueue"

	"github.com/nikhita/kube-custom-controller/pkg/apis/github/v1"
	"github.com/nikhita/kube-custom-controller/pkg/client"
	factory "github.com/nikhita/kube-custom-controller/pkg/informers/externalversions"
)

var (
	ctx context.Context

	githubClient *github.Client

	queue = workqueue.NewRateLimitingQueue(workqueue.NewItemExponentialFailureRateLimiter(time.Second*5, time.Minute))

	stopCh = make(chan struct{})

	sharedFactory factory.SharedInformerFactory

	cl client.Interface
)

func main() {
	kubeconfig := ""
	flag.StringVar(&kubeconfig, "kubeconfig", kubeconfig, "kubeconfig file")

	githubToken := ""
	flag.StringVar(&githubToken, "token", githubToken, "Github API token")

	flag.Parse()

	// set kubeconfig
	if kubeconfig == "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	}

	var (
		config *rest.Config
		err    error
	)

	if kubeconfig != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	} else {
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating client: %v", err)
		os.Exit(1)
	}

	// create the Kubernetes client
	cl = client.NewForConfigOrDie(config)

	// set github API token
	if githubToken == "" {
		githubToken = os.Getenv("TOKEN")
	}

	// Create an authenticated Github client.
	ctx = context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: githubToken},
	)
	tc := oauth2.NewClient(ctx, ts)
	githubClient = github.NewClient(tc)

	// we use a shared informer from the informer factory, to save calls to the
	// API as we grow our application and so state is consistent between our
	// control loops. We set a resync period of 30 seconds, in case any
	// create/replace/update/delete operations are missed when watching
	sharedFactory = factory.NewSharedInformerFactory(cl, time.Second*30)

	informer := sharedFactory.Github().V1().Comments().Informer()
	informer.AddEventHandler(
		cache.ResourceEventHandlerFuncs{
			AddFunc: enqueue,
			UpdateFunc: func(old, cur interface{}) {
				if !reflect.DeepEqual(old, cur) {
					enqueue(cur)
				}
			},
			DeleteFunc: enqueue,
		},
	)

	// start the informer.
	sharedFactory.Start(stopCh)
	log.Printf("Started informer factory.")

	// wait for the informer cache to finish performing it's initial sync of
	// resources
	if !cache.WaitForCacheSync(stopCh, informer.HasSynced) {
		log.Fatalf("error waiting for informer cache to sync: %s", err.Error())
	}

	log.Printf("Finished populating shared informer cache.")
	// here we start just one worker reading objects off the queue. If you
	// wanted to parallelize this, you could start many instances of the worker
	// function, then ensure your application handles concurrency correctly.
	work()
}

// sync will attempt to 'Sync' a resource. It checks to see if the comment
// has already been created, and if not will send it and update the resource
// accordingly. This method is called whenever this controller starts, and
// whenever the resource changes, and also periodically every resyncPeriod.
func sync(notif *v1.Comment) error {
	// If the comment has already been created, we exit with no error
	if notif.Status.Created {
		log.Printf("Skipping already Sent alert '%s/%s'", notif.Namespace, notif.Name)
		return nil
	}

	// send the comment now
	if err := sendComment(ctx, githubClient, notif.Spec.Message); err != nil {
		return err
	}

	log.Printf("Sent github comment!")
	log.Printf(notif.Spec.Message)

	// mark it as created
	notif.Status.Created = true
	if _, err := cl.GithubV1().Comments(notif.Namespace).Update(notif); err != nil {
		return fmt.Errorf("error saving update to Comment resource: %s", err.Error())
	}
	log.Printf("Finished saving update to Comment resource '%s/%s'", notif.Namespace, notif.Name)

	// we didn't encounter any errors, so we return nil to allow the callee
	// to 'forget' this item from the queue altogether.
	return nil
}

func work() {
	for {
		// we read a message off the queue
		key, shutdown := queue.Get()

		// if the queue has been shut down, we should exit the work queue here
		if shutdown {
			stopCh <- struct{}{}
			return
		}

		// convert the queue item into a string. If it's not a string, we'll
		// simply discard it as invalid data and log a message.
		var strKey string
		var ok bool
		if strKey, ok = key.(string); !ok {
			runtime.HandleError(fmt.Errorf("key in queue should be of type string but got %T. discarding", key))
			return
		}

		// we define a function here to process a queue item, so that we can
		// use 'defer' to make sure the message is marked as Done on the queue
		func(key string) {
			defer queue.Done(key)

			// attempt to split the 'key' into namespace and object name
			namespace, name, err := cache.SplitMetaNamespaceKey(strKey)

			if err != nil {
				runtime.HandleError(fmt.Errorf("error splitting meta namespace key into parts: %s", err.Error()))
				return
			}

			log.Printf("Read item '%s/%s' off workqueue. Processing...", namespace, name)

			// retrieve the latest version in the cache of this comment
			obj, err := sharedFactory.Github().V1().Comments().Lister().Comments(namespace).Get(name)
			if err != nil {
				runtime.HandleError(fmt.Errorf("error getting object '%s/%s' from api: %s", namespace, name, err.Error()))
				return
			}

			log.Printf("Got most up to date version of '%s/%s'. Syncing...", namespace, name)

			// attempt to sync the current state of the world with the desired!
			// If sync returns an error, we skip calling `queue.Forget`,
			// thus causing the resource to be requeued at a later time.
			if err := sync(obj); err != nil {
				runtime.HandleError(fmt.Errorf("error processing item '%s/%s': %s", namespace, name, err.Error()))
				return
			}

			log.Printf("Finished processing '%s/%s' successfully! Removing from queue.", namespace, name)

			// as we managed to process this successfully, we can forget it
			// from the work queue altogether.
			queue.Forget(key)
		}(strKey)
	}
}

// enqueue will add an object 'obj' into the workqueue. The object being added
// must be of type metav1.Object, metav1.ObjectAccessor or cache.ExplicitKey.
func enqueue(obj interface{}) {
	// DeletionHandlingMetaNamespaceKeyFunc will convert an object into a
	// 'namespace/name' string. We do this because our item may be processed
	// much later than now, and so we want to ensure it gets a fresh copy of
	// the resource when it starts. Also, this allows us to keep adding the
	// same item into the work queue without duplicates building up.
	key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
	if err != nil {
		runtime.HandleError(fmt.Errorf("error obtaining key for object being enqueue: %s", err.Error()))
		return
	}
	// add the item to the queue
	queue.Add(key)
}

func sendComment(ctx context.Context, client *github.Client, message string) error {
	comment := &github.IssueComment{
		Body: &message,
	}
	_, _, err := client.Issues.CreateComment(ctx, "nikhita", "kube-custom-controller", 1, comment)
	if err != nil {
		return err
	}
	return nil
}
