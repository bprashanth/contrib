/*
Copyright 2015 The Kubernetes Authors All rights reserved.

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

package main

import (
	"fmt"

	"k8s.io/contrib/Ingress/controllers/gce/backends"
	"k8s.io/contrib/Ingress/controllers/gce/healthchecks"
	"k8s.io/contrib/Ingress/controllers/gce/instances"
	"k8s.io/contrib/Ingress/controllers/gce/loadbalancers"
	"k8s.io/contrib/Ingress/controllers/gce/utils"
	"k8s.io/kubernetes/pkg/util/intstr"
	"k8s.io/kubernetes/pkg/util/sets"
)

const (
	// Add used to record additions in a sync pool.
	Add = iota
	// Remove used to record removals from a sync pool.
	Remove
	// Sync used to record syncs of a sync pool.
	Sync
	// Get used to record Get from a sync pool.
	Get
	// Create used to recrod creations in a sync pool.
	Create
	// Update used to record updates in a sync pool.
	Update
	// Delete used to record deltions from a sync pool.
	Delete
	// AddInstances used to record a call to AddInstances.
	AddInstances
	// RemoveInstances used to record a call to RemoveInstances.
	RemoveInstances
)

var (
	testBackendPort = intstr.IntOrString{Type: intstr.Int, IntVal: 80}
	testClusterName = "testcluster"
	testPathMap     = map[string]string{"/foo": defaultBackendName(testClusterName)}
	testIPManager   = testIP{}
)

// ClusterManager fake
type fakeClusterManager struct {
	*ClusterManager
	fakeLbs      *loadbalancers.FakeLoadBalancers
	fakeBackends *backends.FakeBackendServices
	fakeIGs      *instances.FakeInstanceGroups
}

// newFakeClusterManager creates a new fake ClusterManager.
func newFakeClusterManager(clusterName string) *fakeClusterManager {
	fakeLbs := loadbalancers.NewFakeLoadBalancers(clusterName)
	fakeBackends := backends.NewFakeBackendServices()
	fakeIGs := instances.NewFakeInstanceGroups(sets.NewString())
	fakeHCs := healthchecks.NewFakeHealthChecks()

	nodePool := instances.NewNodePool(fakeIGs)
	healthChecker := healthchecks.NewHealthChecker(fakeHCs, "/")
	backendPool := backends.NewBackendPool(
		fakeBackends,
		healthChecker, nodePool)
	l7Pool := loadbalancers.NewLoadBalancerPool(
		fakeLbs,
		// TODO: change this
		backendPool,
		utils.TestDefaultBeNodePort,
	)
	cm := &ClusterManager{
		ClusterName:  clusterName,
		instancePool: nodePool,
		backendPool:  backendPool,
		l7Pool:       l7Pool,
	}
	return &fakeClusterManager{cm, fakeLbs, fakeBackends, fakeIGs}
}

type testIP struct {
	start int
}

func (t *testIP) ip() string {
	t.start++
	return fmt.Sprintf("0.0.0.%v", t.start)
}
