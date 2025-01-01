package structures

import (
	"encoding/xml"
)

type Source struct {
	XMLName xml.Name `xml:"source"`
	Host    Host     `xml:"host"`
}

type Host struct {
	Name string `xml:"name,attr"`
	Port int    `xml:"port,attr"`
}

func (source *Source) Unmarshal(doc string) error {
	return xml.Unmarshal([]byte(doc), source)
}

func (source *Source) Marshal() (string, error) {
	doc, err := xml.MarshalIndent(source, "", "  ")
	if err != nil {
		return "", err
	}
	return string(doc), nil
}
