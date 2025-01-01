package structures

import (
	"encoding/xml"

	"libvirt.org/go/libvirtxml"
)

type Sources struct {
	XMLName xml.Name                       `xml:"sources"`
	Source  []libvirtxml.StoragePoolSource `xml:"source"`
}

func (sources *Sources) Unmarshal(doc string) error {
	return xml.Unmarshal([]byte(doc), sources)
}

func (sources *Sources) Marshal() (string, error) {
	doc, err := xml.MarshalIndent(sources, "", "  ")
	if err != nil {
		return "", err
	}
	return string(doc), nil
}
