package imagepool

import (
	"encoding/xml"

	"github.com/marthjod/gocart/api"
	"github.com/marthjod/gocart/image"
)

// ImagePool is a list of Images.
type ImagePool struct {
	XMLName xml.Name       `xml:"IMAGE_POOL"`
	Images  []*image.Image `xml:"IMAGE"`
}

// Info http://docs.opennebula.org/4.12/integration/system_interfaces/api.html#one-imagepool-info
func (p *ImagePool) Info(c *api.RPC) error {
	return c.Call(p, "one.imagepool.info", []interface{}{c.AuthString, -2, -1, -1})
}

// ExistsName does this image exist?
func (p *ImagePool) ExistsName(n string) bool {
	for _, image := range p.Images {
		if image.Name == n {
			return true
		}
	}
	return false
}

// ExistsID does this image exist?
func (p *ImagePool) ExistsID(i int) bool {
	for _, image := range p.Images {
		if image.ID == i {
			return true
		}
	}
	return false
}
