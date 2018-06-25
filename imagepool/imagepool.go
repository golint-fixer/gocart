package imagepool

import (
	"encoding/xml"

	"github.com/marthjod/gocart/image"
)

// ImagePool is a list of Images.
type ImagePool struct {
	XMLName xml.Name       `xml:"IMAGE_POOL"`
	Images  []*image.Image `xml:"IMAGE"`
}

// APIMethod implements the api.Endpointer interface
func (vt *ImagePool) APIMethod() string {
	return "one.imagepool.info"
}

// Unmarshal implements the api.Endpointer interface
func (vt *ImagePool) Unmarshal(data []byte) error {
	err := xml.Unmarshal(data, vt)
	return err
}

// APIArgs implements the api.Endpointer interface
// API parameter documentation: http://docs.opennebula.org/4.10/integration/system_interfaces/api.html#one-template-info
func (vt *ImagePool) APIArgs(authstring string) []interface{} {
	return []interface{}{authstring, -2, -1, -1}
}

// NewImagePool creates a new ImagePool
func NewImagePool() *ImagePool {
	p := new(ImagePool)
	return p
}

// ExistsName does this image exists?
func (vt *ImagePool) ExistsName(n string) bool {
	for _, image := range vt.Images {
		if image.Name == n {
			return true
		}
	}
	return false
}

// ExistsID does this image exists?
func (vt *ImagePool) ExistsID(i int) bool {
	for _, image := range vt.Images {
		if image.ID == i {
			return true
		}
	}
	return false
}
