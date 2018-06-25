package clusterpool

import (
	"encoding/xml"
	"fmt"

	"github.com/marthjod/gocart/api"
	"github.com/marthjod/gocart/cluster"
)

// ClusterPool is a list of Clusters.
type ClusterPool struct {
	XMLName  xml.Name           `xml:"CLUSTER_POOL"`
	Clusters []*cluster.Cluster `xml:"CLUSTER"`
}

// Info http://docs.opennebula.org/4.12/integration/system_interfaces/api.html#one-clusterpool-info
func (p *ClusterPool) Info(c *api.RPC) error {
	return c.Call(p, "one.clusterpool.info", []interface{}{c.AuthString, -2, -1, -1})
}

// ExistsName determines if a cluster pool with a given name exists.
func (p *ClusterPool) ExistsName(n string) bool {
	for _, cl := range p.Clusters {
		if cl.Name == n {
			return true
		}
	}
	return false
}

// ExistsID determines if a cluster pool with a given ID exists.
func (p *ClusterPool) ExistsID(n int) bool {
	for _, cl := range p.Clusters {
		if cl.ID == n {
			return true
		}
	}
	return false
}

// GetIDByName returns ID of Cluster or -1 and error
func (p *ClusterPool) GetIDByName(name string) (i int, err error) {
	for _, cl := range p.Clusters {
		if cl.Name == name {
			return cl.ID, nil
		}
	}
	return -1, fmt.Errorf("could not find cluster with name %s", name)
}

// GetNameByID returns Name of Cluster if exits
func (p *ClusterPool) GetNameByID(id int) (name string, err error) {
	for _, cl := range p.Clusters {
		if cl.ID == id {
			return cl.Name, nil
		}
	}
	return "", fmt.Errorf("could not find cluster with id %d", id)
}
