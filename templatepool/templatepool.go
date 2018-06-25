package templatepool

import (
	"encoding/xml"
	"regexp"

	"github.com/marthjod/gocart/template"
)

// TemplatePool is a list of Templates
type TemplatePool struct {
	XMLName   xml.Name             `xml:"VMTEMPLATE_POOL"`
	Templates []*template.Template `xml:"VMTEMPLATE"`
}

// APIMethod implements the api.Endpointer interface
func (t *TemplatePool) APIMethod() string {
	return "one.templatepool.info"
}

// Unmarshal implements the api.Endpointer interface
func (t *TemplatePool) Unmarshal(data []byte) error {
	err := xml.Unmarshal(data, t)
	return err
}

// APIArgs implements the api.Endpointer interface
// API parameter documentation: http://docs.opennebula.org/4.10/integration/system_interfaces/api.html#one-template-info
func (t *TemplatePool) APIArgs(authstring string) []interface{} {
	return []interface{}{authstring, -2, -1, -1}
}

// GetTemplatesByName returns a VM template pool based on matching template names.
func (t *TemplatePool) GetTemplatesByName(matchPattern string) (*TemplatePool, error) {
	var p TemplatePool
	for _, tpl := range t.Templates {
		match, err := regexp.MatchString(matchPattern, tpl.Name)
		if err != nil {
			return &p, err
		}
		if match {
			p.Templates = append(p.Templates, tpl)
		}
	}
	return &p, nil
}
