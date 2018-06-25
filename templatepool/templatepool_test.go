package templatepool

import (
	"encoding/xml"
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
	c := getFileContents(path)
	_ = xml.Unmarshal(c, &pool)
	return pool
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
