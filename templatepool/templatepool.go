package templatepool

import (
	"encoding/xml"
	"regexp"

	"github.com/marthjod/gocart/api"
	"github.com/marthjod/gocart/template"
)

// TemplatePool is a list of Templates
type TemplatePool struct {
	XMLName   xml.Name             `xml:"VMTEMPLATE_POOL"`
	Templates []*template.Template `xml:"VMTEMPLATE"`
}

// Info http://docs.opennebula.org/4.12/integration/system_interfaces/api.html#one-templatepool-info
func (t *TemplatePool) Info(c *api.RPC) error {
	return c.Call(t, "one.templatepool.info", []interface{}{c.AuthString, -2, -1, -1})
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
