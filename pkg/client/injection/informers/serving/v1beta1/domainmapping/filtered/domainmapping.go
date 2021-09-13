/*
Copyright 2021 The Knative Authors

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

// Code generated by injection-gen. DO NOT EDIT.

package filtered

import (
	context "context"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	cache "k8s.io/client-go/tools/cache"
	controller "knative.dev/pkg/controller"
	injection "knative.dev/pkg/injection"
	logging "knative.dev/pkg/logging"
	apisservingv1beta1 "knative.dev/serving/pkg/apis/serving/v1beta1"
	versioned "knative.dev/serving/pkg/client/clientset/versioned"
	v1beta1 "knative.dev/serving/pkg/client/informers/externalversions/serving/v1beta1"
	client "knative.dev/serving/pkg/client/injection/client"
	filtered "knative.dev/serving/pkg/client/injection/informers/factory/filtered"
	servingv1beta1 "knative.dev/serving/pkg/client/listers/serving/v1beta1"
)

func init() {
	injection.Default.RegisterFilteredInformers(withInformer)
	injection.Dynamic.RegisterDynamicInformer(withDynamicInformer)
}

// Key is used for associating the Informer inside the context.Context.
type Key struct {
	Selector string
}

func withInformer(ctx context.Context) (context.Context, []controller.Informer) {
	untyped := ctx.Value(filtered.LabelKey{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch labelkey from context.")
	}
	labelSelectors := untyped.([]string)
	infs := []controller.Informer{}
	for _, selector := range labelSelectors {
		f := filtered.Get(ctx, selector)
		inf := f.Serving().V1beta1().DomainMappings()
		ctx = context.WithValue(ctx, Key{Selector: selector}, inf)
		infs = append(infs, inf.Informer())
	}
	return ctx, infs
}

func withDynamicInformer(ctx context.Context) context.Context {
	untyped := ctx.Value(filtered.LabelKey{})
	if untyped == nil {
		logging.FromContext(ctx).Panic(
			"Unable to fetch labelkey from context.")
	}
	labelSelectors := untyped.([]string)
	for _, selector := range labelSelectors {
		inf := &wrapper{client: client.Get(ctx), selector: selector}
		ctx = context.WithValue(ctx, Key{Selector: selector}, inf)
	}
	return ctx
}

// Get extracts the typed informer from the context.
func Get(ctx context.Context, selector string) v1beta1.DomainMappingInformer {
	untyped := ctx.Value(Key{Selector: selector})
	if untyped == nil {
		logging.FromContext(ctx).Panicf(
			"Unable to fetch knative.dev/serving/pkg/client/informers/externalversions/serving/v1beta1.DomainMappingInformer with selector %s from context.", selector)
	}
	return untyped.(v1beta1.DomainMappingInformer)
}

type wrapper struct {
	client versioned.Interface

	namespace string

	selector string
}

var _ v1beta1.DomainMappingInformer = (*wrapper)(nil)
var _ servingv1beta1.DomainMappingLister = (*wrapper)(nil)

func (w *wrapper) Informer() cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(nil, &apisservingv1beta1.DomainMapping{}, 0, nil)
}

func (w *wrapper) Lister() servingv1beta1.DomainMappingLister {
	return w
}

func (w *wrapper) DomainMappings(namespace string) servingv1beta1.DomainMappingNamespaceLister {
	return &wrapper{client: w.client, namespace: namespace, selector: w.selector}
}

func (w *wrapper) List(selector labels.Selector) (ret []*apisservingv1beta1.DomainMapping, err error) {
	reqs, err := labels.ParseToRequirements(w.selector)
	if err != nil {
		return nil, err
	}
	selector = selector.Add(reqs...)
	lo, err := w.client.ServingV1beta1().DomainMappings(w.namespace).List(context.TODO(), v1.ListOptions{
		LabelSelector: selector.String(),
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
	if err != nil {
		return nil, err
	}
	for idx := range lo.Items {
		ret = append(ret, &lo.Items[idx])
	}
	return ret, nil
}

func (w *wrapper) Get(name string) (*apisservingv1beta1.DomainMapping, error) {
	// TODO(mattmoor): Check that the fetched object matches the selector.
	return w.client.ServingV1beta1().DomainMappings(w.namespace).Get(context.TODO(), name, v1.GetOptions{
		// TODO(mattmoor): Incorporate resourceVersion bounds based on staleness criteria.
	})
}