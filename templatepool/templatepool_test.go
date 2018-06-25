package templatepool

import (
	"fmt"
	"io/ioutil"
	"testing"
)

func getFileContents(path string) []byte {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}
	}
	return f
}

func getTemplatePoolFromFile(path string) *TemplatePool {
	pool := &TemplatePool{}
	xml := getFileContents(path)
	_ = pool.Unmarshal(xml)
	return pool
}

func TestAPIMethod(t *testing.T) {
	const apiMethod = "one.templatepool.info"

	p := &TemplatePool{}

	if p.APIMethod() != apiMethod {
		t.Fatalf("API method differs from %s", apiMethod)
	}
}

func TestAPIArgs(t *testing.T) {
	const (
		authstring = "user:pass"
		expected   = "[user:pass %!s(int=-2) %!s(int=-1) %!s(int=-1)]"
	)

	p := &TemplatePool{}

	args := p.APIArgs(authstring)
	argsStr := fmt.Sprintf("%s", args)
	if argsStr != expected {
		t.Fatalf("Mismatch: %s != %s", argsStr, expected)
	}
}

func TestUnmarshalValidVmpool(t *testing.T) {
	pool := &TemplatePool{}
	xml := getFileContents("testdata/templatepool.xml")

	err := pool.Unmarshal(xml)
	if err != nil {
		t.Fatal(err.Error())
	}
}

func TestUnmarshalInvalidXML(t *testing.T) {
	var expected = "expected element type <VMTEMPLATE_POOL> but have <INVALID_TEMPLATE_POOL>"
	pool := &TemplatePool{}
	xml := getFileContents("testdata/invalid-templatepool.xml")

	err := pool.Unmarshal(xml)
	if err == nil {
		t.Fatal("Call did not throw an error")
	} else {
		if err.Error() != expected {
			t.Fatalf("Errors do not match: %s != %s", err.Error(), expected)
		}
	}
}

func TestGetTemplatesByName(t *testing.T) {
	var expectedTemplatesByPattern = map[string]int{
		"^test": 2,
		"test":  2,
		"st1$":  1,
		"foo":   0,
	}

	pool := getTemplatePoolFromFile("testdata/templatepool.xml")

	for pattern, expected := range expectedTemplatesByPattern {
		found, err := pool.GetTemplatesByName(pattern)
		if err != nil {
			t.Fatal(err.Error())
		}
		if len(found.Templates) != expected {
			t.Errorf("%q: expected %d VM(s), found %d", pattern, expected, len(found.Templates))
		}

	}
}
