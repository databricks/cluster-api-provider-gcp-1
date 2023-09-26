/*
Copyright 2022 The Kubernetes Authors.

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

package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// ManagedMachinePoolFinalizer allows Reconcile to clean up GCP resources associated with the GCPManagedMachinePool before
	// removing it from the apiserver.
	ManagedMachinePoolFinalizer = "gcpmanagedmachinepool.infrastructure.cluster.x-k8s.io"
)

// GCPManagedMachinePoolSpec defines the desired state of GCPManagedMachinePool.
type GCPManagedMachinePoolSpec struct {
	// NodePoolName specifies the name of the GKE node pool corresponding to this MachinePool. If you don't specify a name
	// then a default name will be created based on the namespace and name of the managed machine pool.
	// +optional
	NodePoolName string `json:"nodePoolName,omitempty"`
	// MachineType is the name of a Google Compute Engine [machine
	// type](https://cloud.google.com/compute/docs/machine-types).
	// If unspecified, the default machine type is `e2-medium`.
	// +optional
	MachineType *string `json:"machineType,omitempty"`
	// DiskSizeGb is the size of the disk attached to each node, specified in GB.
	// The smallest allowed disk size is 10GB. If unspecified, the default disk size is 100GB.
	// +optional
	DiskSizeGb *int32 `json:"diskSizeGb,omitempty"`
	// ServiceAccount is the Google Cloud Platform Service Account to be used by the node VMs.
	// Specify the email address of the Service Account; otherwise, if no Service
	// Account is specified, the "default" service account is used.
	// +optional
	ServiceAccount *string `json:"serviceAccount,omitempty"`
	// ImageType is the type of image to use for this node. Note that for a given image type,
	// the latest version of it will be used. Please see
	// https://cloud.google.com/kubernetes-engine/docs/concepts/node-images for
	// available image types.
	// +optional
	ImageType *string `json:"imageType,omitempty"`
	// LocalSsdCount is the number of local SSD disks to be attached to the node.
	// +optional
	LocalSsdCount *int32 `json:"localSsdCount,omitempty"`
	// The list of instance tags applied to all nodes. Tags are used to identify
	// valid sources or targets for network firewalls and are specified by
	// the client during cluster or node pool creation. Each tag within the list
	// must comply with RFC1035.
	NetworkTags []string `json:"networkTags,omitempty"`
	// DiskType is the of the disk attached to each node
	// +optional
	DiskType *string `json:"diskType,omitempty"`
	// Scaling specifies scaling for the node pool
	// +optional
	Scaling *NodePoolAutoScaling `json:"scaling,omitempty"`
	// KubernetesLabels specifies the labels to apply to the nodes of the node pool.
	// +optional
	KubernetesLabels infrav1.Labels `json:"kubernetesLabels,omitempty"`
	// KubernetesTaints specifies the taints to apply to the nodes of the node pool.
	// +optional
	KubernetesTaints Taints `json:"kubernetesTaints,omitempty"`
	// AdditionalLabels is an optional set of tags to add to GCP resources managed by the GCP provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalLabels infrav1.Labels      `json:"additionalLabels,omitempty"`
	Management       *NodePoolManagement `json:"management,omitempty"`
	// MaxPodsConstraint is the maximum number of pods that can be run
	// simultaneously on a node in the node pool.
	// +optional
	MaxPodsConstraint *int64 `json:"maxPodsConstraint,omitempty"`
	// ProviderIDList are the provider IDs of instances in the
	// managed instance group corresponding to the nodegroup represented by this
	// machine pool
	// +optional
	ProviderIDList []string `json:"providerIDList,omitempty"`
}

// GCPManagedMachinePoolStatus defines the observed state of GCPManagedMachinePool.
type GCPManagedMachinePoolStatus struct {
	// Ready denotes that the GCPManagedMachinePool has joined the cluster
	// +kubebuilder:default=false
	Ready bool `json:"ready"`
	// Replicas is the most recently observed number of replicas.
	// +optional
	Replicas int32 `json:"replicas"`
	// Conditions specifies the cpnditions for the managed machine pool
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready"
// +kubebuilder:printcolumn:name="Replicas",type="string",JSONPath=".status.replicas"
// +kubebuilder:resource:path=gcpmanagedmachinepools,scope=Namespaced,categories=cluster-api,shortName=gcpmmp
// +kubebuilder:storageversion

// GCPManagedMachinePool is the Schema for the gcpmanagedmachinepools API.
type GCPManagedMachinePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   GCPManagedMachinePoolSpec   `json:"spec,omitempty"`
	Status GCPManagedMachinePoolStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// GCPManagedMachinePoolList contains a list of GCPManagedMachinePool.
type GCPManagedMachinePoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []GCPManagedMachinePool `json:"items"`
}

// NodePoolAutoScaling specifies scaling options.
type NodePoolAutoScaling struct {
	// MinCount specifies the minimum number of nodes in the node pool
	// +optional
	MinCount *int32 `json:"minCount,omitempty"`
	// MaxCount specifies the maximum number of nodes in the node pool
	// +optional
	MaxCount *int32 `json:"maxCount,omitempty"`
	// Is autoscaling enabled for this node pool. If unspecified, the default value is true.
	// +optional
	EnableAutoscaling *bool `json:"enableAutoscaling,omitempty"`
	// Location policy used when scaling up a nodepool.
	// +kubebuilder:validation:Enum=balanced;any
	// +optional
	LocationPolicy *ManagedNodePoolLocationPolicy `json:"locationPolicy,omitempty"`
}

// NodePoolManagement specifies auto-upgrade and auto-repair options.
type NodePoolManagement struct {
	// AutoUpgrade specifies whether node auto-upgrade is enabled for the node
	// pool. If enabled, node auto-upgrade helps keep the nodes in your node pool
	// up to date with the latest release version of Kubernetes.
	AutoUpgrade bool `json:"autoUpgrade,omitempty"`
	// AutoRepair specifies whether the node auto-repair is enabled for the node
	// pool. If enabled, the nodes in this node pool will be monitored and, if
	// they fail health checks too many times, an automatic repair action will be
	// triggered.
	AutoRepair bool `json:"autoRepair,omitempty"`
}

// ManagedNodePoolLocationPolicy specifies the location policy of the node pool when autoscaling is enabled.
type ManagedNodePoolLocationPolicy string

const (
	// ManagedNodePoolLocationPolicyBalanced aims to balance the sizes of different zones.
	ManagedNodePoolLocationPolicyBalanced ManagedNodePoolLocationPolicy = "balanced"
	// ManagedNodePoolLocationPolicyAny picks zones that have the highest capacity available.
	ManagedNodePoolLocationPolicyAny ManagedNodePoolLocationPolicy = "any"
)

// GetConditions returns the machine pool conditions.
func (r *GCPManagedMachinePool) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the status conditions for the GCPManagedMachinePool.
func (r *GCPManagedMachinePool) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&GCPManagedMachinePool{}, &GCPManagedMachinePoolList{})
}
