package dspool

import (
	"encoding/xml"
	"fmt"

	"github.com/marthjod/gocart/api"
	"github.com/marthjod/gocart/datastore"
)

// DSPool is a list of Datastores.
type DSPool struct {
	XMLName    xml.Name               `xml:"DATASTORE_POOL"`
	Datastores []*datastore.Datastore `xml:"DATASTORE"`
}

// ExistsName determines if a datastore with a given name exists.
func (p *DSPool) ExistsName(n string) bool {
	for _, ds := range p.Datastores {
		if ds.Name == n {
			return true
		}
	}
	return false
}

// ExistsID determines if a datastore with a given ID exists.
func (p *DSPool) ExistsID(n int) bool {
	for _, ds := range p.Datastores {
		if ds.ID == n {
			return true
		}
	}
	return false
}

// GetClusterIDByName returns the cluster ID.
func (p *DSPool) GetClusterIDByName(n string) (c int, err error) {
	for _, ds := range p.Datastores {
		if ds.Name == n {
			return ds.ClusterID, nil
		}
	}
	return -1, fmt.Errorf("could not find datastore with name %s", n)
}

// Info http://docs.opennebula.org/4.12/integration/system_interfaces/api.html#one-datastorepool-info
func (p *DSPool) Info(c *api.RPC) error {
	return c.Call(p, "one.datastorepool.info", []interface{}{c.AuthString, -2, -1, -1})
}
