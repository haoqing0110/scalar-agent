# Cluster Score API

## Release Signoff Checklist

- [ ] Enhancement is `implementable`
- [ ] Design details are appropriately documented from clear requirements
- [ ] Test plan is defined
- [ ] Graduation criteria for dev preview, tech preview, GA
- [ ] User-facing documentation is created in [website](https://github.com/open-cluster-management-io/website/)

## Summary

The proposed work provide a API to represents metrics value of the managed cluster expansively.

## Motivation

When implementing placement resource based scheduling, we find some prioritizer need more metrics to calculate the score of the managed cluster, more than the default value provided by API `ManagedCluster` and `ManagedClusterInfo`.
For example, prioritizer ResourceRatioCPU depends on knowing the real-time avaliable CPU, which is not provided by `ManagedCluster`.
Considering below user stories, we want a more extensible mechanism to support this.

### User Stories

#### Story 1: In placement resource based scheduling, placement controller could read each cluster score directly from the new API.
  - An agent could calculate each cluster score based on resource collection, it could be running on hub or spoke.
  - Placement prioritizer plugins could read cluster score directly from API.

#### Story 2: In disaster recovery, the controller could read each cluster score directly from the new API.
  - An agent could calculate each cluster score based on the cluster health.
  - The controller could read cluster score directly from API.

### Goals

- TODO

### Non-Goals

- TODO

## Proposal

### ManagedClusterScore API
ManagedClusterScore is the new API we want to add, to represents a scalable value (aka score) of one managed cluster.
```go
// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +kubebuilder:resource:scope="Namespaced"
// +kubebuilder:subresource:status

// ManagedClusterScore represents a scalable value (aka score) of one managed cluster.
// Each ManagedClusterScore only represents the score for one specific calculator type.
// ManagedClusterScore is a namesapce scoped resource.
//
// The ManagedClusterScore name should follow the format {cluster name}-{calculator name}.
// For example, a calculator named ResourceAllocatableMemory can calculate the totale allocatable memory
// of one cluster.
// So for cluster1, the corresponding ManagedClusterScore name is cluster1-resourceallocatablememory.
type ManagedClusterScore struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Status represents the status of the ManagedClusterScore.
	// +optional
	Status ManagedClusterScoreStatus `json:"status,omitempty"`
}

// ManagedClusterScoreStatus represents the current status of ManagedClusterScore.
type ManagedClusterScoreStatus struct {
	// Conditions contains the different condition statuses for this managed cluster score.
	Conditions []ManagedClusterScoreCondition `json:"conditions"`

	// Score contains a scalable value of this managed cluster.
	Score int64 `json:"score,omitempty"`
}

// ManagedClusterScoreCondition represents the condition of ManagedClusterScore.
type ManagedClusterScoreCondition struct {
	metav1.Condition `json:",inline"`

	// lastUpdateTime is the last time the statue score updated.
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Type=string
	// +kubebuilder:validation:Format=date-time
	LastUpdateTime metav1.Time `json:"lastUpdateTime" protobuf:"bytes,4,opt,name=lastTransitionTime"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ManagedClusterScoreList is a collection of managed cluster score.
type ManagedClusterScoreList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`

	// Items is a list of managed clusters
	Items []ManagedClusterScore `json:"items"`
}
```

## Examples

### 1. An agent update ResourceRatioCPU score to API, and placement prioritizer plugin read the score for each cluster. (Use Story 1)
```yaml
apiVersion: cluster.open-cluster-management.io/v1alpha1
kind: ManagedClusterScore
metadata:
  name: cluster1-resourceratiocpu
  namespace: cluster1
status:
  conditions:
  - lastTransitionTime: "2021-10-28T08:31:39Z"
    lastUpdateTime: "2021-10-29T18:31:39Z"
    message: ManagedClusterScore updated successfully
    reason: ManagedClusterScoreUpdated
    status: "True"
    type: ManagedClusterScoreUpdated
  score: 74
```

## Test Plan

- Unit tests cover the new plugin;
- Unit tests cover the 4 new prioritizers;
- Integration tests cover user story 1-4;

## Graduation Criteria
#### Alpha
1. The new APIs is reviewed and accepted;
2. Implementation is completed to support the functionalities;
3. Develop test cases to demonstrate that the above user stories work correctly;

#### Beta
1. Need to revisit the API shape before upgrade to beta based on user’s feedback.

## Upgrade / Downgrade Strategy
N/A

## Version Skew Strategy
N/A

## Appendix
N/A

#### Scale up
N/A

#### Scale down
N/A