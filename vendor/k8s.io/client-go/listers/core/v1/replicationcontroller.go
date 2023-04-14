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
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// ReplicationControllerLister helps list ReplicationControllers.
// All objects returned here must be treated as read-only.
type ReplicationControllerLister interface {
	// List lists all ReplicationControllers in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.ReplicationController, err error)
	// ReplicationControllers returns an object that can list and get ReplicationControllers.
	ReplicationControllers(namespace string) ReplicationControllerNamespaceLister
	ReplicationControllerListerExpansion
}

// replicationControllerLister implements the ReplicationControllerLister interface.
type replicationControllerLister struct {
	indexer cache.Indexer
}

// NewReplicationControllerLister returns a new ReplicationControllerLister.
func NewReplicationControllerLister(indexer cache.Indexer) ReplicationControllerLister {
	return &replicationControllerLister{indexer: indexer}
}

// List lists all ReplicationControllers in the indexer.
func (s *replicationControllerLister) List(selector labels.Selector) (ret []*v1.ReplicationController, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ReplicationController))
	})
	return ret, err
}

// ReplicationControllers returns an object that can list and get ReplicationControllers.
func (s *replicationControllerLister) ReplicationControllers(namespace string) ReplicationControllerNamespaceLister {
	return replicationControllerNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ReplicationControllerNamespaceLister helps list and get ReplicationControllers.
// All objects returned here must be treated as read-only.
type ReplicationControllerNamespaceLister interface {
	// List lists all ReplicationControllers in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1.ReplicationController, err error)
	// Get retrieves the ReplicationController from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1.ReplicationController, error)
	ReplicationControllerNamespaceListerExpansion
}

// replicationControllerNamespaceLister implements the ReplicationControllerNamespaceLister
// interface.
type replicationControllerNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ReplicationControllers in the indexer for a given namespace.
func (s replicationControllerNamespaceLister) List(selector labels.Selector) (ret []*v1.ReplicationController, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1.ReplicationController))
	})
	return ret, err
}

// Get retrieves the ReplicationController from the indexer for a given namespace and name.
func (s replicationControllerNamespaceLister) Get(name string) (*v1.ReplicationController, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1.Resource("replicationcontroller"), name)
	}
	return obj.(*v1.ReplicationController), nil
}
