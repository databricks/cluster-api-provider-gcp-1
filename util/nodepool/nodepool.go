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
	"sigs.k8s.io/cluster-api-provider-gcp/cloud/services/shared"
)

// NumZones returns number of zones used by the GKE node pool.
func NumZones(nodePoolLocations []string, clusterLocation string) int {
	if len(nodePoolLocations) == 0 {
		// When locations are not specified for the node pool, a node pool of a regional cluster has 3 zones
		if shared.IsRegional(clusterLocation) {
			return cloud.DefaultNumRegionsPerZone
		}
		// When locations are not specified for the node pool, a node pool of a zonal cluster has 1 zone
		return 1
	}
	// When locations are specified for the node pool, the number of zones of the node pool is the number of locations specified
	return len(nodePoolLocations)
}

func PerLocationReplicaCount(totalReplicas int32, numZones int) int32 {
	return int32(math.Ceil(float64(totalReplicas) / float64(numZones)))
}
