/*
Copyright 2019 Replicated, Inc..

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
// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1beta1 "github.com/replicatedhq/kots/kotskinds/apis/kots/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeLintConfigs implements LintConfigInterface
type FakeLintConfigs struct {
	Fake *FakeKotsV1beta1
	ns   string
}

var lintconfigsResource = schema.GroupVersionResource{Group: "kots.io", Version: "v1beta1", Resource: "lintconfigs"}

var lintconfigsKind = schema.GroupVersionKind{Group: "kots.io", Version: "v1beta1", Kind: "LintConfig"}

// Get takes name of the lintConfig, and returns the corresponding lintConfig object, and an error if there is any.
func (c *FakeLintConfigs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.LintConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(lintconfigsResource, c.ns, name), &v1beta1.LintConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.LintConfig), err
}

// List takes label and field selectors, and returns the list of LintConfigs that match those selectors.
func (c *FakeLintConfigs) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.LintConfigList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(lintconfigsResource, lintconfigsKind, c.ns, opts), &v1beta1.LintConfigList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.LintConfigList{ListMeta: obj.(*v1beta1.LintConfigList).ListMeta}
	for _, item := range obj.(*v1beta1.LintConfigList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested lintConfigs.
func (c *FakeLintConfigs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(lintconfigsResource, c.ns, opts))

}

// Create takes the representation of a lintConfig and creates it.  Returns the server's representation of the lintConfig, and an error, if there is any.
func (c *FakeLintConfigs) Create(ctx context.Context, lintConfig *v1beta1.LintConfig, opts v1.CreateOptions) (result *v1beta1.LintConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(lintconfigsResource, c.ns, lintConfig), &v1beta1.LintConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.LintConfig), err
}

// Update takes the representation of a lintConfig and updates it. Returns the server's representation of the lintConfig, and an error, if there is any.
func (c *FakeLintConfigs) Update(ctx context.Context, lintConfig *v1beta1.LintConfig, opts v1.UpdateOptions) (result *v1beta1.LintConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(lintconfigsResource, c.ns, lintConfig), &v1beta1.LintConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.LintConfig), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeLintConfigs) UpdateStatus(ctx context.Context, lintConfig *v1beta1.LintConfig, opts v1.UpdateOptions) (*v1beta1.LintConfig, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(lintconfigsResource, "status", c.ns, lintConfig), &v1beta1.LintConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.LintConfig), err
}

// Delete takes name of the lintConfig and deletes it. Returns an error if one occurs.
func (c *FakeLintConfigs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(lintconfigsResource, c.ns, name), &v1beta1.LintConfig{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeLintConfigs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(lintconfigsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.LintConfigList{})
	return err
}

// Patch applies the patch and returns the patched lintConfig.
func (c *FakeLintConfigs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.LintConfig, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(lintconfigsResource, c.ns, name, pt, data, subresources...), &v1beta1.LintConfig{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.LintConfig), err
}
