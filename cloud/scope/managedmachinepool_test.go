package scope

import (
	"cloud.google.com/go/container/apiv1/containerpb"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
	clusterv1exp "sigs.k8s.io/cluster-api/exp/api/v1beta1"
)

var TestGCPMMP *v1beta1.GCPManagedMachinePool
var TestMP *clusterv1exp.MachinePool
var TestClusterName string

var _ = Describe("GCPManagedMachinePool Scope", func() {
	BeforeEach(func() {
		TestClusterName = "test-cluster"
		gcpmmpName := "test-gcpmmp"
		nodePoolName := "test-pool"
		namespace := "capg-system"
		replicas := int32(1)

		TestGCPMMP = &v1beta1.GCPManagedMachinePool{
			ObjectMeta: metav1.ObjectMeta{
				Name:      gcpmmpName,
				Namespace: namespace,
			},
			Spec: v1beta1.GCPManagedMachinePoolSpec{
				NodePoolName: nodePoolName,
			},
		}
		TestMP = &clusterv1exp.MachinePool{
			Spec: clusterv1exp.MachinePoolSpec{
				Replicas: &replicas,
			},
		}
	})

	Context("Test NodePoolResourceLabels", func() {
		It("should append cluster owned label", func() {
			labels := infrav1.Labels{"test-key": "test-value"}

			Expect(NodePoolResourceLabels(labels, TestClusterName)).To(Equal(infrav1.Labels{
				"test-key":                             "test-value",
				infrav1.ClusterTagKey(TestClusterName): string(infrav1.ResourceLifecycleOwned),
			}))
		})
	})

	Context("Test ConvertToSdkNodePool", func() {
		It("should convert to SDK node pool with default values", func() {
			sdkNodePool := ConvertToSdkNodePool(*TestGCPMMP, *TestMP, 3, TestClusterName)

			Expect(sdkNodePool).To(Equal(&containerpb.NodePool{
				Name:             TestGCPMMP.Spec.NodePoolName,
				InitialNodeCount: *TestMP.Spec.Replicas,
				Config: &containerpb.NodeConfig{
					ResourceLabels: NodePoolResourceLabels(nil, TestClusterName),
				},
			}))
		})

		It("should convert to SDK node pool node count in a regional cluster", func() {
			replicas := int32(6)
			TestMP.Spec.Replicas = &replicas

			sdkNodePool := ConvertToSdkNodePool(*TestGCPMMP, *TestMP, 3, TestClusterName)

			Expect(sdkNodePool).To(Equal(&containerpb.NodePool{
				Name:             TestGCPMMP.Spec.NodePoolName,
				InitialNodeCount: 2,
				Config: &containerpb.NodeConfig{
					ResourceLabels: NodePoolResourceLabels(nil, TestClusterName),
				},
			}))
		})

		It("should convert to SDK node pool using GCPManagedMachinePool", func() {
			machineType := "n1-standard-1"
			diskSizeGb := int32(128)
			serviceAccount := "test@project.iam.gserviceaccount.com"
			imageType := "ubuntu_containerd"
			localSsdCount := int32(2)
			diskType := "pd-ssd"
			maxPodsConstraint := int64(20)
			enableAutoscaling := false
			scaling := v1beta1.NodePoolAutoScaling{
				EnableAutoscaling: &enableAutoscaling,
			}
			labels := infrav1.Labels{"test-key": "test-value"}
			taints := v1beta1.Taints{
				{
					Key:    "test-key",
					Value:  "test-value",
					Effect: "NoSchedule",
				},
			}
			tags := []string{"test-tag"}
			resourceLabels := infrav1.Labels{"test-key": "test-value"}
			locations := []string{"us-central1-a"}
			reservation := "gpu-reservation"

			TestGCPMMP.Spec.MachineType = &machineType
			TestGCPMMP.Spec.DiskSizeGb = &diskSizeGb
			TestGCPMMP.Spec.ServiceAccount = &serviceAccount
			TestGCPMMP.Spec.ImageType = &imageType
			TestGCPMMP.Spec.LocalSsdCount = &localSsdCount
			TestGCPMMP.Spec.DiskType = &diskType
			TestGCPMMP.Spec.Scaling = &scaling
			TestGCPMMP.Spec.MaxPodsConstraint = &maxPodsConstraint
			TestGCPMMP.Spec.KubernetesLabels = labels
			TestGCPMMP.Spec.KubernetesTaints = taints
			TestGCPMMP.Spec.NetworkTags = tags
			TestGCPMMP.Spec.AdditionalLabels = resourceLabels
			TestGCPMMP.Spec.Locations = locations
			TestGCPMMP.Spec.Reservation = &reservation

			sdkNodePool := ConvertToSdkNodePool(*TestGCPMMP, *TestMP, 1, TestClusterName)

			Expect(sdkNodePool).To(Equal(&containerpb.NodePool{
				Name:             TestGCPMMP.Spec.NodePoolName,
				InitialNodeCount: *TestMP.Spec.Replicas,
				Config: &containerpb.NodeConfig{
					Labels:         labels,
					Taints:         v1beta1.ConvertToSdkTaint(taints),
					ResourceLabels: NodePoolResourceLabels(resourceLabels, TestClusterName),
					Tags:           tags,
					MachineType:    machineType,
					DiskSizeGb:     diskSizeGb,
					ServiceAccount: serviceAccount,
					ImageType:      imageType,
					LocalSsdCount:  localSsdCount,
					DiskType:       diskType,
					ReservationAffinity: &containerpb.ReservationAffinity{
						ConsumeReservationType: containerpb.ReservationAffinity_SPECIFIC_RESERVATION,
						Key:                    "compute.googleapis.com/reservation-name",
						Values:                 []string{reservation},
					},
				},
				Autoscaling: v1beta1.ConvertToSdkAutoscaling(&scaling),
				MaxPodsConstraint: &containerpb.MaxPodsConstraint{
					MaxPodsPerNode: maxPodsConstraint,
				},
				Locations: locations,
			}))
		})
	})

	Context("Test NumZones", func() {
		It("should calculate the correct number for a regular node pool in a regional cluster", func() {
			scope := &ManagedMachinePoolScope{
				GCPManagedControlPlane: &v1beta1.GCPManagedControlPlane{
					Spec: v1beta1.GCPManagedControlPlaneSpec{
						Location: "us-central1",
					},
				},
				GCPManagedMachinePool: &v1beta1.GCPManagedMachinePool{
					Spec: v1beta1.GCPManagedMachinePoolSpec{},
				},
			}

			Expect(scope.NumZones()).To(Equal(3))
		})

		It("should calculate the correct number for a zonal node pool in a regional cluster", func() {
			scope := &ManagedMachinePoolScope{
				GCPManagedControlPlane: &v1beta1.GCPManagedControlPlane{
					Spec: v1beta1.GCPManagedControlPlaneSpec{
						Location: "us-central1",
					},
				},
				GCPManagedMachinePool: &v1beta1.GCPManagedMachinePool{
					Spec: v1beta1.GCPManagedMachinePoolSpec{
						Locations: []string{"us-central1-a"},
					},
				},
			}

			Expect(scope.NumZones()).To(Equal(1))
		})

		It("should calculate the correct number for a regular node pool in a zonal cluster", func() {
			scope := &ManagedMachinePoolScope{
				GCPManagedControlPlane: &v1beta1.GCPManagedControlPlane{
					Spec: v1beta1.GCPManagedControlPlaneSpec{
						Location: "us-central1-a",
					},
				},
				GCPManagedMachinePool: &v1beta1.GCPManagedMachinePool{
					Spec: v1beta1.GCPManagedMachinePoolSpec{},
				},
			}

			Expect(scope.NumZones()).To(Equal(1))
		})

		It("should calculate the correct number for a zonal node pool in a zonal cluster", func() {
			scope := &ManagedMachinePoolScope{
				GCPManagedControlPlane: &v1beta1.GCPManagedControlPlane{
					Spec: v1beta1.GCPManagedControlPlaneSpec{
						Location: "us-central1-a",
					},
				},
				GCPManagedMachinePool: &v1beta1.GCPManagedMachinePool{
					Spec: v1beta1.GCPManagedMachinePoolSpec{
						Locations: []string{"us-central1-a", "us-central1-b", "us-central1-c"},
					},
				},
			}

			Expect(scope.NumZones()).To(Equal(3))
		})
	})
})
