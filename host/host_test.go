package host_test

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/marthjod/gocart/hostpool"
)

func getFileContents(path string) []byte {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}
	}
	return f
}

func getHostPoolFromFile(path string) *hostpool.HostPool {
	pool := &hostpool.HostPool{}
	c := getFileContents(path)
	_ = xml.Unmarshal(c, &pool)
	return pool
}

func TestHostIsEmpty(t *testing.T) {
	fixture := "testdata/hostpool.xml"
	f, err := os.Open(fixture)
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		t.Fatal(err.Error())
	}

	pool, err := hostpool.FromReader(f)
	if err != nil {
		t.Fatal(err.Error())
	}

	emptyHostFound := false
	for _, host := range pool.Hosts {
		if len(host.VMIDs) == 0 {
			emptyHostFound = true
			if !host.IsEmpty() {
				t.Error("IsEmpty() should return true for host without VMs")
			}
		} else {
			if host.IsEmpty() {
				t.Error("IsEmpty() should return false for host with VMs")
			}
		}
	}

	if !emptyHostFound {
		t.Errorf("No empty host found in %q for testing IsEmpty()", fixture)
	}
}

func TestHost_String(t *testing.T) {
	pool := getHostPoolFromFile("testdata/hostpool.xml")
	for _, host := range pool.Hosts {
		if host.String() != fmt.Sprintf(host.Name) {
			t.Errorf("String() does not match host name")
			break
		}
	}
}
