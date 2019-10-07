/*
Copyright The Kubernetes Authors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1 "kinnylee.com/crds-controller-demo/pkg/apis/control/v1"
)

// ScalingLister helps list Scalings.
type ScalingLister interface {
	// List lists all Scalings in the indexer.
	List(selector labels.Selector) (ret []*v1.Scaling, err error)
	// Scalings returns an object that can list and get Scalings.
	Scalings(namespace string) ScalingNamespaceLister
	ScalingListerExpansion
}

// scalingLister implements the ScalingLister interface.
type scalingLister struct {
	indexer cache.Indexer
}

// NewScalingLister returns a new ScalingLister.
func NewScalingLister(indexer cache.Indexer) ScalingLister {
	return &scalingLister{indexer: indexer}
}

// List lists all Scalings in the indexer.
func (s *scalingLister) List(selector labels.Selector) (ret []*v1.Scaling, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Scaling))
	})
	return ret, err
}

// Scalings returns an object that can list and get Scalings.
func (s *scalingLister) Scalings(namespace string) ScalingNamespaceLister {
	return scalingNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ScalingNamespaceLister helps list and get Scalings.
type ScalingNamespaceLister interface {
	// List lists all Scalings in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1.Scaling, err error)
	// Get retrieves the Scaling from the indexer for a given namespace and name.
	Get(name string) (*v1.Scaling, error)
	ScalingNamespaceListerExpansion
}

// scalingNamespaceLister implements the ScalingNamespaceLister
// interface.
type scalingNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all Scalings in the indexer for a given namespace.
func (s scalingNamespaceLister) List(selector labels.Selector) (ret []*v1.Scaling, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.Scaling))
	})
	return ret, err
}

// Get retrieves the Scaling from the indexer for a given namespace and name.
func (s scalingNamespaceLister) Get(name string) (*v1.Scaling, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("scaling"), name)
	}
	return obj.(*v1.Scaling), nil
}
