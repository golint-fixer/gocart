package hostpool

import (
	"encoding/xml"
	"io"

	"github.com/marthjod/gocart/host"
	"github.com/marthjod/gocart/vmpool"
)

// TODO use real enum like for VM states
const (
	Init                int = iota
	MonitoringMonitored int = iota // Currently monitoring, previously MONITORED
	Monitored           int = iota
	Error               int = iota
	Disabled            int = iota
	MonitoringError     int = iota // Currently monitoring, previously ERROR
	MonitoringInit      int = iota // Currently monitoring, previously initialized
	MonitoringDisabled  int = iota // Currently monitoring, previously DISABLED
)

// HostPool represents a host pool.
type HostPool struct {
	XMLName xml.Name `xml:"HOST_POOL"`
	Hosts   []*Host  `xml:"HOST"`
}

// APIMethod implements the api.Endpointer interface
func (p *HostPool) APIMethod() string {
	return "one.hostpool.info"
}

// APIArgs implements the api.Endpointer interface
// API parameter documentation: http://docs.opennebula.org/4.10/integration/system_interfaces/api.html#one-hostpool-info
func (p *HostPool) APIArgs(authstring string) []interface{} {
	return []interface{}{authstring}
}

// Unmarshal unmarshals into a host pool.
func (p *HostPool) Unmarshal(data []byte) error {
	return xml.Unmarshal(data, p)
}

// MapVMs ...
func (p *HostPool) MapVMs(vmpool *vmpool.VMPool) {
	for _, host := range p.Hosts {
		host.MapVMs(vmpool)
	}
}

// Host represents an OpenNebula node/host.
type Host struct {
	*host.Host
	VMPool *vmpool.VMPool
}

// String returns a host's short strings representation.
func (h *Host) String() string {
	return h.Name
}

// NewHostPool returns a new host pool.
func NewHostPool() *HostPool {
	return &HostPool{}
}

// FromReader reads into a host pool.
func FromReader(r io.Reader) (*HostPool, error) {
	pool := HostPool{}
	dec := xml.NewDecoder(r)
	if err := dec.Decode(&pool); err != nil {
		return nil, err
	}
	return &pool, nil
}

// GetHostsInCluster returns a pool of hosts in the provided cluster.
func (p *HostPool) GetHostsInCluster(cluster string) *HostPool {
	var (
		hostpool HostPool
	)
	for _, host := range p.Hosts {
		if host.Cluster == cluster {
			hostpool.Hosts = append(hostpool.Hosts, host)
		}
	}
	return &hostpool
}

// FilterHostsByStates returns host pool containing only hosts in one of the provided states.
func (p *HostPool) FilterHostsByStates(states ...int) *HostPool {
	var (
		hp HostPool
	)
	for _, host := range p.Hosts {
		for _, state := range states {
			if host.State == state {
				hp.Hosts = append(hp.Hosts, host)
				continue
			}
		}
	}
	return &hp
}

// FilterOutEmptyHosts filters out hosts without VMs.
func (p *HostPool) FilterOutEmptyHosts() *HostPool {
	var (
		hp HostPool
	)
	for _, host := range p.Hosts {
		if !host.IsEmpty() {
			hp.Hosts = append(hp.Hosts, host)
		}
	}
	return &hp
}

// MapVMs ...
func (h *Host) MapVMs(vmpool *vmpool.VMPool) {
	h.VMPool = vmpool.GetVMsByID(h.VMIDs...)
}
