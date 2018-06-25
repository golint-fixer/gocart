package vmpool_test

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"testing"

	"github.com/marthjod/gocart/vm"
	"github.com/marthjod/gocart/vmpool"
)

func getFileContents(path string) []byte {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}
	}
	return f
}

func getVMPoolFromFile(path string) *vmpool.VMPool {
	pool := &vmpool.VMPool{}
	c := getFileContents(path)
	_ = xml.Unmarshal(c, &pool)
	return pool
}

func TestFromReader(t *testing.T) {
	var expected = 3

	f, err := os.Open("testdata/vmpool.xml")
	defer func() {
		_ = f.Close()
	}()
	if err != nil {
		t.Fatal(err.Error())
	}

	pool, err := vmpool.FromReader(f)
	if err != nil {
		t.Fatal(err.Error())
	}

	if len(pool.VMs) != expected {
		t.Errorf("VM pool has length %d, expected %d", len(pool.VMs), expected)
	}

}

func TestGetVMsById(t *testing.T) {
	var expectedVMsById = map[int]int{
		6:    1,
		7:    1,
		8:    1,
		1337: 0,
	}

	pool := getVMPoolFromFile("testdata/vmpool.xml")

	for id, vmsExpected := range expectedVMsById {
		vmsFound := pool.GetVMsByID(id)
		if len(vmsFound.VMs) != vmsExpected {
			t.Errorf("Expected %d VM(s), found %d", vmsExpected, len(vmsFound.VMs))
		}
	}

}

func TestGetVMsByName(t *testing.T) {
	var expectedVMsByPattern = map[string]int{
		"^vm":      3,
		"vm-in$":   2,
		"exa":      1,
		"dummy-vm": 0,
	}

	pool := getVMPoolFromFile("testdata/vmpool.xml")

	for pattern, vmsExpected := range expectedVMsByPattern {
		vmsFound, err := pool.GetVMsByName(pattern)
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(vmsFound.VMs) != vmsExpected {
			t.Errorf("%q: expected %d VM(s), found %d", pattern, vmsExpected, len(vmsFound.VMs))
		}

	}
}

func TestGetDistinctVmNamePatterns(t *testing.T) {
	var prefix = "^"
	var infix = ".+"
	var suffix = "$"
	var expectedDistinctPatterns = []struct {
		filter  string
		pattern string
	}{
		{"^([vm]{2})-([ni]{2})$", "^vm.+in$"},
		{"^vm-(ex)....(e)$", "^ex.+e$"},
		{"^vm-([a-z]{4}).(.).*", "^exam.+l$"},
	}

	pool := getVMPoolFromFile("testdata/vmpool.xml")

	for _, expected := range expectedDistinctPatterns {
		distinctPatterns := pool.GetDistinctVMNamePatterns(expected.filter, prefix, infix, suffix)

		if !distinctPatterns[expected.pattern] {
			t.Errorf("Expected distinct pattern %q not extracted by filter %q", expected.pattern, expected.filter)
		}
	}
}

func TestGetDistinctVmNamePatternsExtractHostname(t *testing.T) {
	var prefix = "^"
	var infix = ".+"
	var suffix = "$"
	var expectedDistinctPatterns = []struct {
		filter  string
		pattern string
	}{
		{".{2}(amp).*(amp)", "^amp.+amp$"},
		{"^..-(..).*(m)$", "^vm.+m$"},
		{"^(.{2}).*(example.com)", "^in.+example.com$"},
	}

	pool := getVMPoolFromFile("testdata/vmpool.xml")
	fqdnExtractor := func(vm *vm.VM) string {
		h, _ := vm.UserTemplate.Items.GetCustom("CUSTOM_FQDN")
		return h
	}

	for _, expected := range expectedDistinctPatterns {
		distinctPatterns := pool.GetDistinctVMNamePatternsExtractHostname(expected.filter, prefix, infix, suffix, fqdnExtractor)

		if !distinctPatterns[expected.pattern] {
			t.Errorf("Expected distinct pattern %q not extracted by filter %q", expected.pattern, expected.filter)
		}
	}
}
