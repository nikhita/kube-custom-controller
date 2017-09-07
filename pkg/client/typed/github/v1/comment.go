/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	v1 "github.com/nikhita/kube-custom-controller/pkg/apis/github/v1"
	scheme "github.com/nikhita/kube-custom-controller/pkg/client/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// CommentsGetter has a method to return a CommentInterface.
// A group's client should implement this interface.
type CommentsGetter interface {
	Comments(namespace string) CommentInterface
}

// CommentInterface has methods to work with Comment resources.
type CommentInterface interface {
	Create(*v1.Comment) (*v1.Comment, error)
	Update(*v1.Comment) (*v1.Comment, error)
	UpdateStatus(*v1.Comment) (*v1.Comment, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.Comment, error)
	List(opts meta_v1.ListOptions) (*v1.CommentList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Comment, err error)
	CommentExpansion
}

// comments implements CommentInterface
type comments struct {
	client rest.Interface
	ns     string
}

// newComments returns a Comments
func newComments(c *GithubV1Client, namespace string) *comments {
	return &comments{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the comment, and returns the corresponding comment object, and an error if there is any.
func (c *comments) Get(name string, options meta_v1.GetOptions) (result *v1.Comment, err error) {
	result = &v1.Comment{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("comments").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Comments that match those selectors.
func (c *comments) List(opts meta_v1.ListOptions) (result *v1.CommentList, err error) {
	result = &v1.CommentList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("comments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested comments.
func (c *comments) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("comments").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a comment and creates it.  Returns the server's representation of the comment, and an error, if there is any.
func (c *comments) Create(comment *v1.Comment) (result *v1.Comment, err error) {
	result = &v1.Comment{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("comments").
		Body(comment).
		Do().
		Into(result)
	return
}

// Update takes the representation of a comment and updates it. Returns the server's representation of the comment, and an error, if there is any.
func (c *comments) Update(comment *v1.Comment) (result *v1.Comment, err error) {
	result = &v1.Comment{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("comments").
		Name(comment.Name).
		Body(comment).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *comments) UpdateStatus(comment *v1.Comment) (result *v1.Comment, err error) {
	result = &v1.Comment{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("comments").
		Name(comment.Name).
		SubResource("status").
		Body(comment).
		Do().
		Into(result)
	return
}

// Delete takes name of the comment and deletes it. Returns an error if one occurs.
func (c *comments) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("comments").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *comments) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("comments").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched comment.
func (c *comments) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Comment, err error) {
	result = &v1.Comment{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("comments").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
