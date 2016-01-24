package utils

import (
	"fmt"
	"strings"

	compute "google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
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

	// This allows sharing of backends across loadbalancers.
	backendPrefix = "k8s-be"

	// A delimiter used for clarity in naming GCE resources.
	clusterNameDelimiter = "--"

	// Arbitrarily chosen alphanumeric character to use in constructing resource
	// names, eg: to avoid cases where we end up with a name ending in '-'.
	alphaNumericChar = "0"

	// Names longer than this are truncated, because of GCE restrictions.
	nameLenLimit = 62

	// DefaultBackendKey is the key used to transmit the defaultBackend through
	// a urlmap. It's not a valid subdomain, and it is a catch all path.
	// TODO: Find a better way to transmit this, once we've decided on default
	// backend semantics (i.e do we want a default per host, per lb etc).
	DefaultBackendKey     = "DefaultBackend"
	K8sAnnotationPrefix   = "ingress.kubernetes.io"
	TestDefaultBeNodePort = int64(3000)
)

type FakeIngressRuleValueMap map[string]string

type Namer struct {
	ClusterName string
}

// BeName constructs the name for a backend.
func (n *Namer) BeName(port int64) string {
	// TODO: Pipe the clusterName through, for now it saves code churn to just
	// grab it globally, especially since we haven't decided how to handle
	// namespace conflicts in the Ubernetes context.
	if n.ClusterName == "" {
		return fmt.Sprintf("%v-%d", backendPrefix, port)
	}
	return n.Truncate(fmt.Sprintf("%v-%d%v%v", backendPrefix, port, clusterNameDelimiter, n.ClusterName))
}

func (n *Namer) LBName(key string) string {
	// TODO: Pipe the clusterName through, for now it saves code churn to just
	// grab it globally, especially since we haven't decided how to handle
	// namespace conflicts in the Ubernetes context.
	parts := strings.Split(key, clusterNameDelimiter)
	scrubbedName := strings.Replace(key, "/", "-", -1)
	if n.ClusterName == "" || parts[len(parts)-1] == n.ClusterName {
		return scrubbedName
	}
	return n.Truncate(fmt.Sprintf("%v%v%v", scrubbedName, clusterNameDelimiter, n.ClusterName))
}

// Truncate truncates the given key to a GCE length limit.
func (n *Namer) Truncate(key string) string {
	if len(key) > nameLenLimit {
		// GCE requires names to end with an albhanumeric, but allows characters
		// like '-', so make sure the trucated name ends legally.
		return fmt.Sprintf("%v%v", key[:nameLenLimit], alphaNumericChar)
	}
	return key
}

// GCEURLMap is a nested map of hostname->path regex->backend
type GCEURLMap map[string]map[string]*compute.BackendService

// getDefaultBackend performs a destructive read and returns the default
// backend of the urlmap.
func (g GCEURLMap) GetDefaultBackend() *compute.BackendService {
	var d *compute.BackendService
	var exists bool
	if h, ok := g[DefaultBackendKey]; ok {
		if d, exists = h[DefaultBackendKey]; exists {
			delete(h, DefaultBackendKey)
		}
		delete(g, DefaultBackendKey)
	}
	return d
}

// String implements the string interface for the GCEURLMap.
func (g GCEURLMap) String() string {
	msg := ""
	for host, um := range g {
		msg += fmt.Sprintf("%v\n", host)
		for url, be := range um {
			msg += fmt.Sprintf("\t%v: ", url)
			if be == nil {
				msg += fmt.Sprintf("No backend\n")
			} else {
				msg += fmt.Sprintf("%v\n", be.Name)
			}
		}
	}
	return msg
}

// PutDefaultBackend performs a destructive write replacing the
// default backend of the url map with the given backend.
func (g GCEURLMap) PutDefaultBackend(d *compute.BackendService) {
	g[DefaultBackendKey] = map[string]*compute.BackendService{
		DefaultBackendKey: d,
	}
}

// IsHTTPErrorCode checks if the given error matches the given HTTP Error code.
// For this to work the error must be a googleapi Error.
func IsHTTPErrorCode(err error, code int) bool {
	apiErr, ok := err.(*googleapi.Error)
	return ok && apiErr.Code == code
}

// CompareLinks returns true if the 2 self links are equal.
func CompareLinks(l1, l2 string) bool {
	// TODO: These can be partial links
	return l1 == l2 && l1 != ""
}
