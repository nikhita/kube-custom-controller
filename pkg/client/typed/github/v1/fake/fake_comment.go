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

package fake

import (
	github_v1 "github.com/nikhita/kube-custom-controller/pkg/apis/github/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeComments implements CommentInterface
type FakeComments struct {
	Fake *FakeGithubV1
	ns   string
}

var commentsResource = schema.GroupVersionResource{Group: "github.k8s.io", Version: "v1", Resource: "comments"}

var commentsKind = schema.GroupVersionKind{Group: "github.k8s.io", Version: "v1", Kind: "Comment"}

// Get takes name of the comment, and returns the corresponding comment object, and an error if there is any.
func (c *FakeComments) Get(name string, options v1.GetOptions) (result *github_v1.Comment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(commentsResource, c.ns, name), &github_v1.Comment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*github_v1.Comment), err
}

// List takes label and field selectors, and returns the list of Comments that match those selectors.
func (c *FakeComments) List(opts v1.ListOptions) (result *github_v1.CommentList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(commentsResource, commentsKind, c.ns, opts), &github_v1.CommentList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &github_v1.CommentList{}
	for _, item := range obj.(*github_v1.CommentList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested comments.
func (c *FakeComments) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(commentsResource, c.ns, opts))

}

// Create takes the representation of a comment and creates it.  Returns the server's representation of the comment, and an error, if there is any.
func (c *FakeComments) Create(comment *github_v1.Comment) (result *github_v1.Comment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(commentsResource, c.ns, comment), &github_v1.Comment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*github_v1.Comment), err
}

// Update takes the representation of a comment and updates it. Returns the server's representation of the comment, and an error, if there is any.
func (c *FakeComments) Update(comment *github_v1.Comment) (result *github_v1.Comment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(commentsResource, c.ns, comment), &github_v1.Comment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*github_v1.Comment), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeComments) UpdateStatus(comment *github_v1.Comment) (*github_v1.Comment, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(commentsResource, "status", c.ns, comment), &github_v1.Comment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*github_v1.Comment), err
}

// Delete takes name of the comment and deletes it. Returns an error if one occurs.
func (c *FakeComments) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(commentsResource, c.ns, name), &github_v1.Comment{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeComments) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(commentsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &github_v1.CommentList{})
	return err
}

// Patch applies the patch and returns the patched comment.
func (c *FakeComments) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *github_v1.Comment, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(commentsResource, c.ns, name, data, subresources...), &github_v1.Comment{})

	if obj == nil {
		return nil, err
	}
	return obj.(*github_v1.Comment), err
}
