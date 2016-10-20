package vmpool

import (
	"encoding/xml"
	"io"
	"regexp"
	"strings"

	"github.com/marthjod/gocart/ocatypes"
)

type VmPool struct {
	XMLName xml.Name       `xml:"VM_POOL"`
	Vms     []*ocatypes.Vm `xml:"VM"`
}

// ApiMethod implements the api.Endpointer interface
func (vmpool *VmPool) ApiMethod() string {
	return "one.vmpool.info"
}

// ApiArgs implements the api.Endpointer interface
// API parameter documentation: http://docs.opennebula.org/4.10/integration/system_interfaces/api.html#one-vmpool-info
func (vmpool *VmPool) ApiArgs(authstring string) []interface{} {
	return []interface{}{authstring, -2, -1, -1, -1}
}

func (vmpool *VmPool) Unmarshal(data []byte) error {
	err := xml.Unmarshal(data, vmpool)
	return err
}

func NewVmPool() *VmPool {
	p := new(VmPool)
	return p
}

func (vmPool *VmPool) String() string {
	var (
		vmNames []string
	)

	for _, vm := range vmPool.Vms {
		vmNames = append(vmNames, vm.Name)
	}

	return strings.Join(vmNames, ", ")
}

func FromReader(r io.Reader) (*VmPool, error) {
	pool := VmPool{}
	dec := xml.NewDecoder(r)
	if err := dec.Decode(&pool); err != nil {
		return nil, err
	}
	return &pool, nil
}

func (vmPool *VmPool) GetVmsById(ids ...int) *VmPool {
	var (
		pool VmPool
	)
	for _, vm := range vmPool.Vms {
		for _, id := range ids {
			if vm.Id == id {
				pool.Vms = append(pool.Vms, vm)
			}
		}
	}
	return &pool
}

func (vmPool *VmPool) GetVmsByName(matchPattern string) (*VmPool, error) {
	var pool VmPool
	for _, vm := range vmPool.Vms {
		match, err := regexp.MatchString(matchPattern, vm.Name)
		if err != nil {
			return &pool, err
		}
		if match {
			pool.Vms = append(pool.Vms, vm)
		}
	}
	return &pool, nil
}
