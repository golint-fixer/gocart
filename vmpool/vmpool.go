package vmpool

import (
	"encoding/xml"
	"fmt"
	"io"
	"regexp"

	"github.com/marthjod/gocart/vm"
)

// VMPool represents a VM pool.
type VMPool struct {
	XMLName xml.Name `xml:"VM_POOL"`
	VMs     []*vm.VM `xml:"VM"`
}

// String returns a string representation of a VM pool.
func (p *VMPool) String() string {
	var list = []string{}

	for _, vm := range p.VMs {
		list = append(list, vm.Name)
	}

	return fmt.Sprintf("%s", list)
}

// APIMethod implements the api.Endpointer interface
func (p *VMPool) APIMethod() string {
	return "one.vmpool.info"
}

// APIArgs implements the api.Endpointer interface
// API parameter documentation: http://docs.opennebula.org/4.10/integration/system_interfaces/api.html#one-vmpool-info
func (p *VMPool) APIArgs(authstring string) []interface{} {
	return []interface{}{authstring, -2, -1, -1, -1}
}

// Unmarshal unmarshals into a VM pool.
func (p *VMPool) Unmarshal(data []byte) error {
	err := xml.Unmarshal(data, p)
	return err
}

// NewVMPool returns a new VM pool.
func NewVMPool() *VMPool {
	return &VMPool{}
}

// FromReader reads into a VM pool.
func FromReader(r io.Reader) (*VMPool, error) {
	var pool = VMPool{}
	dec := xml.NewDecoder(r)
	if err := dec.Decode(&pool); err != nil {
		return nil, err
	}
	return &pool, nil
}

// GetVMsByID returns a VM pool based on matching VM IDs.
func (p *VMPool) GetVMsByID(ids ...int) *VMPool {
	var pool VMPool
	for _, vm := range p.VMs {
		for _, id := range ids {
			if vm.ID == id {
				pool.VMs = append(pool.VMs, vm)
			}
		}
	}
	return &pool
}

// GetVMsByName returns a VM pool based on matching VM names.
func (p *VMPool) GetVMsByName(matchPattern string) (*VMPool, error) {
	var pool VMPool
	for _, vm := range p.VMs {
		match, err := regexp.MatchString(matchPattern, vm.Name)
		if err != nil {
			return &pool, err
		}
		if match {
			pool.VMs = append(pool.VMs, vm)
		}
	}
	return &pool, nil
}

// GetDistinctVMNamePatterns returns a set of distinct VM name patterns.
func (p *VMPool) GetDistinctVMNamePatterns(filter, prefix, infix, suffix string) map[string]bool {
	vmNameExtractor := func(vm *vm.VM) string {
		return vm.Name
	}

	return p.GetDistinctVMNamePatternsExtractHostname(filter, prefix, infix, suffix, vmNameExtractor)
}

// GetDistinctVMNamePatternsExtractHostname returns a set of distinct VM name patterns where hostname != VM name.
func (p *VMPool) GetDistinctVMNamePatternsExtractHostname(filter, prefix, infix, suffix string,
	hostNameExtractor func(vm *vm.VM) string) map[string]bool {

	var (
		distinctPatterns = make(map[string]bool)
		pattern          string
	)

	re := regexp.MustCompile(filter)

	for _, vm := range p.VMs {

		groups := re.FindStringSubmatch(hostNameExtractor(vm))
		if groups == nil {
			continue
		}

		if len(groups) >= 3 {
			pattern = fmt.Sprintf("%s%s%s%s%s", prefix, groups[1], infix, groups[2], suffix)
			distinctPatterns[pattern] = true
		} else {
			distinctPatterns[vm.Name] = true
		}
	}

	return distinctPatterns
}

// GetVMsByLCMStates returns a VM pool based on matching LCM states.
func (p *VMPool) GetVMsByLCMStates(states ...vm.LCMState) (*VMPool, error) {
	var pool VMPool
	for _, vm := range p.VMs {
		for _, state := range states {
			if vm.LCMState == state {
				pool.VMs = append(pool.VMs, vm)
				break
			}
		}
	}
	return &pool, nil
}
