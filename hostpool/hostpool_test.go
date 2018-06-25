package hostpool_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/marthjod/gocart/hostpool"
)

var expectedPools = []struct {
	cluster string
	poolLen int
}{
	{"default", 2},
	{"non-existent", 0},
}

func getFileContents(path string) []byte {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}
	}
	return f
}

func getHostPoolFromFile(path string) *hostpool.HostPool {
	pool := hostpool.NewHostPool()
	xml := getFileContents(path)
	_ = pool.Unmarshal(xml)
	return pool
}

func TestApiMethod(t *testing.T) {
	const apiMethod = "one.hostpool.info"

	if hostpool.NewHostPool().APIMethod() != apiMethod {
		t.Fatalf("API method differs from %s", apiMethod)
	}
}

func TestApiArgs(t *testing.T) {
	const (
		authstring = "user:pass"
		expected   = "[user:pass]"
	)

	args := hostpool.NewHostPool().APIArgs(authstring)
	argsStr := fmt.Sprintf("%s", args)
	if argsStr != expected {
		t.Fatalf("Mismatch: %s != %s", argsStr, expected)
	}
}

func TestUnmarshalValidHostpool(t *testing.T) {
	pool := hostpool.NewHostPool()
	xml := getFileContents("testdata/hostpool.xml")

	err := pool.Unmarshal(xml)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestUnmarshalInvalidXML(t *testing.T) {
	var expected = "expected element type <HOST_POOL> but have <INVALID_HOST_POOL>"
	pool := hostpool.NewHostPool()
	xml := getFileContents("testdata/invalid-hostpool.xml")

	err := pool.Unmarshal(xml)
	if err == nil {
		t.Fatal("Call did not throw an error")
	} else {
		if err.Error() != expected {
			t.Fatalf("Errors do not match: %s != %s", err.Error(), expected)
		}
	}
}

func TestFromReader(t *testing.T) {
	var expected = 2

	f, err := os.Open("testdata/hostpool.xml")
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

	if len(pool.Hosts) != expected {
		t.Errorf("Host pool has length %d, expected %d", len(pool.Hosts), expected)
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

func TestGetHostsInCluster(t *testing.T) {
	pool := getHostPoolFromFile("testdata/hostpool.xml")

	for _, expected := range expectedPools {
		clusterPool := pool.GetHostsInCluster(expected.cluster)
		if len(clusterPool.Hosts) != expected.poolLen {
			t.Fatalf("Unexpected pool size (%d != %d) for cluster %s",
				len(clusterPool.Hosts), expected.poolLen, expected.cluster)
		}
	}
}

func TestFilterHostsByStates(t *testing.T) {
	pool := getHostPoolFromFile("testdata/hostpool.xml")

	disabledHosts := pool.FilterHostsByStates(hostpool.Disabled)
	if len(disabledHosts.Hosts) != 1 {
		t.Fatalf("Expected 1 disabled host, found %d", len(disabledHosts.Hosts))
	}
	if disabledHosts.Hosts[0].Name != "thost" {
		t.Fatalf("Found wrong disabled host %s", disabledHosts.Hosts[0].Name)
	}

	emptyPool := pool.FilterHostsByStates(hostpool.Error)
	if len(emptyPool.Hosts) > 0 {
		t.Fatalf("Found more than 0 hosts in state ERROR")
	}

	twoStates := pool.FilterHostsByStates(hostpool.Disabled, hostpool.Monitored)
	if len(twoStates.Hosts) != 2 {
		t.Fatalf("Expected 2 hosts for states Disabled + ERROR, found %d", len(twoStates.Hosts))
	}
}

func TestFilterOutEmptyHosts(t *testing.T) {
	pool := getHostPoolFromFile("testdata/hostpool.xml")
	hostsRunningVMs := pool.FilterOutEmptyHosts()
	if len(hostsRunningVMs.Hosts) > 1 {
		t.Error("Found too many non-empty hosts")
	}
}

func TestFilterChain(t *testing.T) {
	pool := getHostPoolFromFile("testdata/hostpool.xml")
	disabledHosts := pool.FilterHostsByStates(hostpool.Disabled)
	if len(disabledHosts.Hosts) != 1 {
		t.Errorf("Expected 1 disabled host, found %d", len(disabledHosts.Hosts))
	}

	filtered := pool.FilterHostsByStates(hostpool.Disabled).FilterOutEmptyHosts()
	if len(filtered.Hosts) > 1 {
		t.Error("Found too many hosts while chaining filters")
	}
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
