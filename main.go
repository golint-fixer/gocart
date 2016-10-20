package main

import (
	"crypto/tls"
	"flag"
	"log"
	"net/http"
	"os"
	"runtime/pprof"

	"github.com/marthjod/gocart/api"
	"github.com/marthjod/gocart/hostpool"
	"github.com/marthjod/gocart/vmpool"
)

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func CheckPlacementAcrossDatacenters(hostPool *hostpool.HostPool, datacenters []string, vmNamePattern string) bool {
	var (
		totalCount int
		err        error
	)

	var placementOk = true
	var perDatacenterCount = make(map[string]int, len(datacenters))

	for _, host := range hostPool.Hosts {
		host.VmPool, err = host.VmPool.GetVmsByName(vmNamePattern)
		if err != nil {
			panic(err)
		}
		if len(host.VmPool.Vms) > 0 {
			log.Printf("%v runs %d '%v' VMs in %v\n", host.Name, len(host.VmPool.Vms), vmNamePattern, host.Template.Datacenter)

			perDatacenterCount[host.Template.Datacenter] += len(host.VmPool.Vms)
			totalCount += len(host.VmPool.Vms)
		}
	}

	log.Printf("Total: %v\n", totalCount)

	for datacenter, vmCount := range perDatacenterCount {
		placementOk = vmCount >= totalCount/len(datacenters)
		log.Printf("%v: %v (OK?: %v)\n", datacenter, vmCount, placementOk)

	}

	return placementOk
}

func main() {
	var (
		verbose       bool
		datacenters   arrayFlags
		cpuprofile    string
		user          string
		password      string
		url           string
		skipVerifySSL bool
		vmNamePattern string
	)

	flag.Var(&datacenters, "datacenter", "Datacenters")
	flag.BoolVar(&verbose, "v", false, "Verbose mode")
	flag.StringVar(&cpuprofile, "cpuprofile", "", "Write CPU profile to file")
	flag.StringVar(&user, "user", "", "OpenNebula API user (mandatory)")
	flag.StringVar(&password, "password", "", "OpenNebula API password (mandatory)")
	flag.StringVar(&url, "url", "https://localhost:61443/RPC2", "OpenNebula XML-RPC API URL")
	flag.BoolVar(&skipVerifySSL, "skip-verify-ssl", true, "Skip verification of OpenNebula API SSL cert")
	flag.StringVar(&vmNamePattern, "vm-name-pattern", "", "VM name pattern (regexp) (mandatory)")

	flag.Parse()

	if cpuprofile != "" {
		f, err := os.Create(cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: skipVerifySSL},
	}
	apiClient, err := api.NewClient(url, user, password, tr)
	if err != nil {
		panic(err)
	}

	vmPool := vmpool.NewVmPool()
	if err := apiClient.Call(vmPool); err != nil {
		panic(err)
	}

	hostPool := hostpool.NewHostPool()
	if err := apiClient.Call(hostPool); err != nil {
		panic(err)
	}

	hostPool.MapVms(vmPool)

	if !CheckPlacementAcrossDatacenters(hostPool, datacenters, vmNamePattern) {
		os.Exit(3)
	}

}
