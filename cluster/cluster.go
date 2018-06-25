package cluster

import "encoding/xml"

// Cluster represents a cluster.
type Cluster struct {
	XMLName      xml.Name `xml:"CLUSTER"`
	Name         string   `xml:"NAME"`
	ID           int      `xml:"ID"`
	DatastoreIDs []int    `xml:"DATASTORES"`
	VnetIDs      []int    `xml:"VNETS"`
}
