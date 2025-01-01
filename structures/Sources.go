package structures

import (
	"encoding/xml"

	"libvirt.org/go/libvirtxml"
)

type Sources struct {
	XMLName xml.Name                       `xml:"sources"`
	Source  []libvirtxml.StoragePoolSource `xml:"source"`
}
