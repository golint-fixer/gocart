package datastore

import "encoding/xml"

// Datastore represents a datastore.
type Datastore struct {
	XMLName   xml.Name `xml:"DATASTORE"`
	Name      string   `xml:"NAME"`
	ID        int      `xml:"ID"`
	ClusterID int      `xml:"CLUSTER_ID"`
	Cluster   string   `xml:"CLUSTER"`
}
