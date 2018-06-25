package vnetpool

import (
	"encoding/xml"
	"fmt"

	"github.com/marthjod/gocart/api"
	"github.com/marthjod/gocart/vnet"
)

// VNetPool is a list of Virtual Networks
type VNetPool struct {
	XMLName  xml.Name     `xml:"VNET_POOL"`
	Networks []*vnet.VNet `xml:"VNET"`
}

// ExistsName does this vnet exists?
func (vt *VNetPool) ExistsName(n string) bool {
	for _, net := range vt.Networks {
		if net.Name == n {
			return true
		}
	}
	return false
}

// GetIDByName returns vnet ID
func (vt *VNetPool) GetIDByName(n string) (c int, err error) {
	for _, net := range vt.Networks {
		if net.Name == n {
			return net.ID, nil
		}
	}
	return -1, fmt.Errorf("could not find network with name %s", n)
}

// GetVNetByName returns the VNet
func (vt *VNetPool) GetVNetByName(n string) (vnet *vnet.VNet, err error) {
	for _, net := range vt.Networks {
		if net.Name == n {
			return net, nil
		}
	}
	return vnet, fmt.Errorf("could not find network with name %s", n)
}

// GetVNetByID returns the VNet
func (vt *VNetPool) GetVNetByID(i int) (vnet *vnet.VNet, err error) {
	for _, net := range vt.Networks {
		if net.ID == i {
			return net, nil
		}
	}
	return vnet, fmt.Errorf("could not find network with ID %d", i)
}

// ExistsID does this vnet exists?
func (vt *VNetPool) ExistsID(n int) bool {
	for _, net := range vt.Networks {
		if net.ID == n {
			return true
		}
	}
	return false
}

// Info http://docs.opennebula.org/4.12/integration/system_interfaces/api.html#one-vnpool-info
func (vt *VNetPool) Info(c *api.RPC) error {
	return c.Call(vt, "one.vnpool.info", []interface{}{c.AuthString, -2, -1, -1})
}
