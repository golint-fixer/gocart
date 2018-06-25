package dspool

import (
	"encoding/xml"
	"fmt"

	"github.com/marthjod/gocart/datastore"
)

// DSPool is a list of Datastores.
type DSPool struct {
	XMLName    xml.Name               `xml:"DATASTORE_POOL"`
	Datastores []*datastore.Datastore `xml:"DATASTORE"`
}

// ExistsName determines if a datastore with a given name exists.
func (vt *DSPool) ExistsName(n string) bool {
	for _, ds := range vt.Datastores {
		if ds.Name == n {
			return true
		}
	}
	return false
}

// ExistsID determines if a datastore with a given ID exists.
func (vt *DSPool) ExistsID(n int) bool {
	for _, ds := range vt.Datastores {
		if ds.ID == n {
			return true
		}
	}
	return false
}

// GetClusterIDByName returns the cluster ID.
func (vt *DSPool) GetClusterIDByName(n string) (c int, err error) {
	for _, ds := range vt.Datastores {
		if ds.Name == n {
			return ds.ClusterID, nil
		}
	}
	return -1, fmt.Errorf("could not find datastore with name %s", n)
}

// APIMethod implements the api.Endpointer interface
func (vt *DSPool) APIMethod() string {
	return "one.datastorepool.info"
}

// Unmarshal implements the api.Endpointer interface
func (vt *DSPool) Unmarshal(data []byte) error {
	return xml.Unmarshal(data, vt)
}

// APIArgs implements the api.Endpointer interface.
// API parameter documentation: http://docs.opennebula.org/4.10/integration/system_interfaces/api.html#one-template-info
func (vt *DSPool) APIArgs(authstring string) []interface{} {
	return []interface{}{authstring, -2, -1, -1}
}

// NewDSPool creates a new DSPool.
func NewDSPool() *DSPool {
	return &DSPool{}
}
