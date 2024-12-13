/*
Copyright 2024 The Kubernetes Authors.

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

package nodepool

import (
	"math"
	"sigs.k8s.io/cluster-api-provider-gcp/cloud"
)

func PerLocationReplicaCount(totalReplicas int32, isRegional bool) int32 {
	replicas := totalReplicas
	if isRegional {
		// If it's regional, we need to ensure that the total number of replicas is greater or equal to the specified replicas.
		// For example, when we are using regional node pools and want 4 nodes, we should specify 2 replicas (per zone) and let GKE to scale it down back to 4
		replicas = int32(math.Ceil(float64(replicas) / cloud.DefaultNumRegionsPerZone))
	}
	return replicas
}
