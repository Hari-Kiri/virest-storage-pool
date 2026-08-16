package storagePool

import (
	"encoding/xml"

	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

// Detail returns the full XML description of a storage pool.
func (c *Connection) Detail(uuid string, xmlFlags uint) (detail Detail, err error) {
	pool, err := c.lookupPool(uuid)
	if err != nil {
		return Detail{}, err
	}
	defer finishFree(pool, &err)

	model, err := poolXML(pool, libvirt.StorageXMLFlags(xmlFlags))
	if err != nil {
		return Detail{}, err
	}
	return Detail{StoragePool: model}, nil
}

func poolXML(pool poolHandle, flags libvirt.StorageXMLFlags) (libvirtxml.StoragePool, error) {
	doc, err := pool.GetXMLDesc(flags)
	if err != nil {
		return libvirtxml.StoragePool{}, wrap("get pool xml", err)
	}
	var model libvirtxml.StoragePool
	err = model.Unmarshal(doc)
	if err != nil {
		return libvirtxml.StoragePool{}, wrap("unmarshal pool xml", err)
	}
	return model, nil
}

// Capabilities returns supported storage pool types and formats.
func (c *Connection) Capabilities() (Capabilities, error) {
	doc, err := c.hv.GetStoragePoolCapabilities(0)
	if err != nil {
		return Capabilities{}, wrap("get storage pool capabilities", err)
	}
	var caps Capabilities
	err = xml.Unmarshal([]byte(doc), &caps)
	if err != nil {
		return Capabilities{}, wrap("unmarshal storage pool capabilities", err)
	}
	return caps, nil
}

// FindSources discovers available storage pool sources for the given pool type.
func (c *Connection) FindSources(poolType string, src SourceSpec) (Sources, error) {
	srcXML, err := xml.MarshalIndent(src, "", "  ")
	if err != nil {
		return Sources{}, wrap("marshal source spec", err)
	}
	doc, err := c.hv.FindStoragePoolSources(poolType, string(srcXML), 0)
	if err != nil {
		return Sources{}, wrap("find storage pool sources", err)
	}
	var sources Sources
	err = xml.Unmarshal([]byte(doc), &sources)
	if err != nil {
		return Sources{}, wrap("unmarshal storage pool sources", err)
	}
	return sources, nil
}
