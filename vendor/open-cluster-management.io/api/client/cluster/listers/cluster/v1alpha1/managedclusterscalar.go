// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1alpha1 "open-cluster-management.io/api/cluster/v1alpha1"
)

// ManagedClusterScalarLister helps list ManagedClusterScalars.
// All objects returned here must be treated as read-only.
type ManagedClusterScalarLister interface {
	// List lists all ManagedClusterScalars in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ManagedClusterScalar, err error)
	// ManagedClusterScalars returns an object that can list and get ManagedClusterScalars.
	ManagedClusterScalars(namespace string) ManagedClusterScalarNamespaceLister
	ManagedClusterScalarListerExpansion
}

// managedClusterScalarLister implements the ManagedClusterScalarLister interface.
type managedClusterScalarLister struct {
	indexer cache.Indexer
}

// NewManagedClusterScalarLister returns a new ManagedClusterScalarLister.
func NewManagedClusterScalarLister(indexer cache.Indexer) ManagedClusterScalarLister {
	return &managedClusterScalarLister{indexer: indexer}
}

// List lists all ManagedClusterScalars in the indexer.
func (s *managedClusterScalarLister) List(selector labels.Selector) (ret []*v1alpha1.ManagedClusterScalar, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ManagedClusterScalar))
	})
	return ret, err
}

// ManagedClusterScalars returns an object that can list and get ManagedClusterScalars.
func (s *managedClusterScalarLister) ManagedClusterScalars(namespace string) ManagedClusterScalarNamespaceLister {
	return managedClusterScalarNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// ManagedClusterScalarNamespaceLister helps list and get ManagedClusterScalars.
// All objects returned here must be treated as read-only.
type ManagedClusterScalarNamespaceLister interface {
	// List lists all ManagedClusterScalars in the indexer for a given namespace.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.ManagedClusterScalar, err error)
	// Get retrieves the ManagedClusterScalar from the indexer for a given namespace and name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.ManagedClusterScalar, error)
	ManagedClusterScalarNamespaceListerExpansion
}

// managedClusterScalarNamespaceLister implements the ManagedClusterScalarNamespaceLister
// interface.
type managedClusterScalarNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all ManagedClusterScalars in the indexer for a given namespace.
func (s managedClusterScalarNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.ManagedClusterScalar, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.ManagedClusterScalar))
	})
	return ret, err
}

// Get retrieves the ManagedClusterScalar from the indexer for a given namespace and name.
func (s managedClusterScalarNamespaceLister) Get(name string) (*v1alpha1.ManagedClusterScalar, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("managedclusterscalar"), name)
	}
	return obj.(*v1alpha1.ManagedClusterScalar), nil
}
