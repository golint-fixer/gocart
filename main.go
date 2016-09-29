// main.go
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"ocatypes"
	"os"
	"time"
)

func main() {
	var (
		err        error
		verbose    bool
		elapsed    time.Duration
		vmPoolFile string
		hostPoolFile string
		xmlFile    *os.File
		vmPool     *ocatypes.VmPool
		hostPool *ocatypes.HostPool
	)

	flag.StringVar(&vmPoolFile, "vm-pool", "", `VM pool XML dump file path`)
	flag.StringVar(&hostPoolFile, "host-pool", "", `Host pool XML dump file path`)
	flag.BoolVar(&verbose, "v", false, "Verbose mode")
	flag.Parse()

	if vmPoolFile == "" && hostPoolFile == "" {
		flag.PrintDefaults()
		return
	}

	if vmPoolFile != "" {

		vmPool = ocatypes.NewVmPool()

		xmlFile, err = os.Open(vmPoolFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer xmlFile.Close()

		data, err := ioutil.ReadAll(xmlFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		if elapsed, err = vmPool.Read(data); err != nil {
			fmt.Println("Error during unmarshaling:", err)
			return
		}

		fmt.Printf("Read in VM pool of length %v in %v\n", len(vmPool.Vms), elapsed)
		if verbose {
			for i := 0; i < len(vmPool.Vms); i++ {
				vm := vmPool.Vms[i]
				fmt.Printf("%v %v (CPU: %v, template/mem: %v)\n",
					vm.Id, vm.Name, vm.Cpu, vm.Template.Memory)
			}
		}


	} else {


		hostPool = ocatypes.NewHostPool()

		xmlFile, err = os.Open(hostPoolFile)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer xmlFile.Close()

		data, err := ioutil.ReadAll(xmlFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		if elapsed, err = hostPool.Read(data); err != nil {
			fmt.Println("Error during unmarshaling:", err)
			return
		}

		fmt.Printf("Read in host pool of length %v in %v\n", len(hostPool.Hosts), elapsed)
		if verbose {
			for i := 0; i < len(hostPool.Hosts); i++ {
				host := hostPool.Hosts[i]
				fmt.Printf("%v %v\n", host.Id, host.Template.Datacenter)

				// fmt.Printf("%v %v", )
			}
		}

	}

}