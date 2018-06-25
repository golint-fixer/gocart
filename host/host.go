package host

import (
	"encoding/xml"

	"github.com/marthjod/gocart/template"
)

// Host represents an OpenNebula host.
type Host struct {
	XMLName   xml.Name              `xml:"HOST"`
	ID        int                   `xml:"ID"`
	Name      string                `xml:"NAME"`
	State     int                   `xml:"STATE"`
	Cluster   string                `xml:"CLUSTER"`
	ClusterID int                   `xml:"CLUSTER_ID"`
	Template  template.HostTemplate `xml:"TEMPLATE"`
	VMIDs     []int                 `xml:"VMS>ID"`
}

// IsEmpty checks if a host has no VMs.
func (h *Host) IsEmpty() bool {
	return len(h.VMIDs) == 0
}
