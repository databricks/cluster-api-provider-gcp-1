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

// AcceleratorConfig defines the AcceleratorConfig represents a GPU Accelerator request.
// https://pkg.go.dev/cloud.google.com/go/container/apiv1/containerpb#AcceleratorConfig
type AcceleratorConfig struct {
	// The number of the accelerator cards exposed to an instance.
	// +optional
	AcceleratorCount int64 `json:"acceleratorCount,omitempty"`
	// The accelerator type resource name. List of supported accelerators
	// [here](https://cloud.google.com/compute/docs/gpus)
	// +optional
	AcceleratorType string `json:"acceleratorType,omitempty"`
}

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
	AdditionalLabels infrav1.Labels `json:"additionalLabels,omitempty"`
	// NodeManagement configuration for this NodePool.
	// +optional
	Management *NodePoolManagement `json:"management,omitempty"`
	// Networking configuration for this NodePool. If specified, it overrides the
	// cluster-level defaults.
	// +optional
	NetworkConfig *NodeNetworkConfig `json:"networkConfig,omitempty"`
	// MaxPodsConstraint is the maximum number of pods that can be run
	// simultaneously on a node in the node pool.
	// +optional
	MaxPodsConstraint *int64 `json:"maxPodsConstraint,omitempty"`
	// ProviderIDList are the provider IDs of instances in the
	// managed instance group corresponding to the nodegroup represented by this
	// machine pool
	// +optional
	ProviderIDList []string `json:"providerIDList,omitempty"`
	// Locations are the zones in which the node pool's nodes should be located.
	// The list of Google Compute Engine
	// [zones](https://cloud.google.com/compute/docs/zones#available) in which the
	// NodePool's nodes should be located.
	// If this value is unspecified during node pool creation, the
	// [Cluster.Locations](https://cloud.google.com/kubernetes-engine/docs/reference/rest/v1/projects.locations.clusters#Cluster.FIELDS.locations)
	// value will be used, instead.
	// +optional
	Locations []string `json:"locations,omitempty"`
	// Reservation is the name of GCP compute specific reservation.
	// It usually specifies zone information as well.
	// +optional
	Reservation *string `json:"reservation,omitempty"`
	// Accelerators is a list of hardware accelerators to be attached to each node.
	// See https://cloud.google.com/compute/docs/gpus for more information about
	// support for GPUs.
	// +optional
	Accelerators []*AcceleratorConfig `json:"accelerators,omitempty"`
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

// NodeNetworkConfig specifies networking options for node pool.
type NodeNetworkConfig struct {
	// Whether to create a new range for pod IPs in this node pool.
	// Defaults are provided for `pod_range` and `pod_ipv4_cidr_block` if they
	// are not specified.
	//
	// If neither `create_pod_range` or `pod_range` are specified, the
	// cluster-level default (`ip_allocation_policy.cluster_ipv4_cidr_block`) is
	// used.
	//
	// Only applicable if `ip_allocation_policy.use_ip_aliases` is true.
	//
	// This field cannot be changed after the node pool has been created.
	CreatePodRange bool `json:"createPodRange,omitempty"`

	// The ID of the secondary range for pod IPs.
	// If `create_pod_range` is true, this ID is used for the new range.
	// If `create_pod_range` is false, uses an existing secondary range with this
	// ID.
	//
	// Only applicable if `ip_allocation_policy.use_ip_aliases` is true.
	//
	// This field cannot be changed after the node pool has been created.
	PodRange string `json:"podRange,omitempty"`

	// The IP address range for pod IPs in this node pool.
	//
	// Only applicable if `create_pod_range` is true.
	//
	// Set to blank to have a range chosen with the default size.
	//
	// Set to /netmask (e.g. `/14`) to have a range chosen with a specific
	// netmask.
	//
	// Set to a
	// [CIDR](https://en.wikipedia.org/wiki/Classless_Inter-Domain_Routing)
	// notation (e.g. `10.96.0.0/14`) to pick a specific range to use.
	//
	// Only applicable if `ip_allocation_policy.use_ip_aliases` is true.
	//
	// This field cannot be changed after the node pool has been created.
	PodIpv4CidrBlock string `json:"podIpv4CidrBlock,omitempty"`
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
