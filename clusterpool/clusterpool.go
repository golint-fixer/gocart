package clusterpool

import (
	"encoding/xml"
	"fmt"

	"github.com/marthjod/gocart/cluster"
)

// ClusterPool is a list of Clusters.
type ClusterPool struct {
	XMLName  xml.Name           `xml:"CLUSTER_POOL"`
	Clusters []*cluster.Cluster `xml:"CLUSTER"`
}

// APIMethod implements the api.Endpointer interface.
func (vt *ClusterPool) APIMethod() string {
	return "one.clusterpool.info"
}

// Unmarshal implements the api.Endpointer interface
func (vt *ClusterPool) Unmarshal(data []byte) error {
	return xml.Unmarshal(data, vt)
}

// APIArgs implements the api.Endpointer interface.
// API parameter documentation: http://docs.opennebula.org/4.10/integration/system_interfaces/api.html
func (vt *ClusterPool) APIArgs(authstring string) []interface{} {
	return []interface{}{authstring, -2, -1, -1}
}

// NewClusterPool creates a new ClusterPool
func NewClusterPool() *ClusterPool {
	return &ClusterPool{}
}

// ExistsName determines if a cluster pool with a given name exists.
func (vt *ClusterPool) ExistsName(n string) bool {
	for _, cl := range vt.Clusters {
		if cl.Name == n {
			return true
		}
	}
	return false
}

// ExistsID determines if a cluster pool with a given ID exists.
func (vt *ClusterPool) ExistsID(n int) bool {
	for _, cl := range vt.Clusters {
		if cl.ID == n {
			return true
		}
	}
	return false
}

// GetIDByName returns ID of Cluster or -1 and error
func (vt *ClusterPool) GetIDByName(name string) (i int, err error) {
	for _, cl := range vt.Clusters {
		if cl.Name == name {
			return cl.ID, nil
		}
	}
	return -1, fmt.Errorf("could not find cluster with name %s", name)
}

// GetNameByID returns Name of Cluster if exits
func (vt *ClusterPool) GetNameByID(id int) (name string, err error) {
	for _, cl := range vt.Clusters {
		if cl.ID == id {
			return cl.Name, nil
		}
	}
	return "", fmt.Errorf("could not find cluster with id %d", id)
}
