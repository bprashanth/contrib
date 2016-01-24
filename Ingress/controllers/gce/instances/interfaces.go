package instances

import (
	compute "google.golang.org/api/compute/v1"
)

// NodePool is an interface to manage a pool of kubernetes nodes synced with vm instances in the cloud
// through the InstanceGroups interface.
type NodePool interface {
	AddInstanceGroup(name string, port int64) (*compute.InstanceGroup, *compute.NamedPort, error)
	DeleteInstanceGroup(name string) error

	// TODO: Refactor for modularity
	Add(groupName string, nodeNames []string) error
	Remove(groupName string, nodeNames []string) error
	Sync(nodeNames []string) error
	Get(name string) (*compute.InstanceGroup, error)
}

// InstanceGroups is an interface for managing gce instances groups, and the instances therein.
type InstanceGroups interface {
	GetInstanceGroup(name string) (*compute.InstanceGroup, error)
	CreateInstanceGroup(name string) (*compute.InstanceGroup, error)
	DeleteInstanceGroup(name string) error

	// TODO: Refactor for modulatiry.
	ListInstancesInInstanceGroup(name string, state string) (*compute.InstanceGroupsListInstances, error)
	AddInstancesToInstanceGroup(name string, instanceNames []string) error
	RemoveInstancesFromInstanceGroup(name string, instanceName []string) error
	AddPortToInstanceGroup(ig *compute.InstanceGroup, port int64) (*compute.NamedPort, error)
}
