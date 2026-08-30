package storagePool

import (
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

// Define creates a new inactive storage pool from the given XML model.
// On success it returns the UUID of the defined pool.
func (c *Connection) Define(storagePool libvirtxml.StoragePool, flags libvirt.StoragePoolDefineFlags) (uuid string, err error) {
	xmlConfig, err := storagePool.Marshal()
	if err != nil {
		return "", wrap("marshal pool config", err)
	}

	defined, err := c.hv.StoragePoolDefineXML(xmlConfig, flags)
	if err != nil {
		return "", wrap("define storage pool", err)
	}
	defer finishFree(defined, &err)

	uuid, err = defined.GetUUIDString()
	if err != nil {
		return "", wrap("get defined pool uuid", err)
	}
	return uuid, nil
}

// CreateTransient creates and starts a transient (non-persistent) pool from the model.
// Equivalent to virsh pool-create / libvirt StoragePoolCreateXML.
func (c *Connection) CreateTransient(storagePool libvirtxml.StoragePool, flags libvirt.StoragePoolCreateFlags) (uuid string, err error) {
	config, err := storagePool.Marshal()
	if err != nil {
		return "", wrap("marshal pool config", err)
	}

	created, err := c.hv.CreateTransient(config, flags)
	if err != nil {
		return "", wrap("create transient storage pool", err)
	}
	defer finishFree(created, &err)

	uuid, err = created.GetUUIDString()
	if err != nil {
		return "", wrap("get transient pool uuid", err)
	}
	return uuid, nil
}

// Build builds the underlying storage for a defined pool.
func (c *Connection) Build(uuid string, flags libvirt.StoragePoolBuildFlags) (err error) {
	pool, err := c.lookupPool(uuid)
	if err != nil {
		return err
	}
	defer finishFree(pool, &err)
	return wrap("build storage pool", pool.Build(flags))
}

// Create starts an inactive storage pool.
func (c *Connection) Create(uuid string, flags libvirt.StoragePoolCreateFlags) (err error) {
	pool, err := c.lookupPool(uuid)
	if err != nil {
		return err
	}
	defer finishFree(pool, &err)
	return wrap("create storage pool", pool.Create(flags))
}

// Destroy stops an active storage pool.
func (c *Connection) Destroy(uuid string) (err error) {
	pool, err := c.lookupPool(uuid)
	if err != nil {
		return err
	}
	defer finishFree(pool, &err)
	return wrap("destroy storage pool", pool.Destroy())
}

// Delete deletes the underlying storage resources of a pool.
func (c *Connection) Delete(uuid string, flags libvirt.StoragePoolDeleteFlags) (err error) {
	pool, err := c.lookupPool(uuid)
	if err != nil {
		return err
	}
	defer finishFree(pool, &err)
	return wrap("delete storage pool", pool.Delete(flags))
}

// Undefine removes a storage pool definition.
func (c *Connection) Undefine(uuid string) (err error) {
	pool, err := c.lookupPool(uuid)
	if err != nil {
		return err
	}
	defer finishFree(pool, &err)
	return wrap("undefine storage pool", pool.Undefine())
}

// Refresh refreshes the list of volumes in the storage pool.
func (c *Connection) Refresh(uuid string) (err error) {
	pool, err := c.lookupPool(uuid)
	if err != nil {
		return err
	}
	defer finishFree(pool, &err)
	return wrap("refresh storage pool", pool.Refresh(0))
}

// Autostart configures whether the pool starts with the host.
func (c *Connection) Autostart(uuid string, autostart bool) (err error) {
	pool, err := c.lookupPool(uuid)
	if err != nil {
		return err
	}
	defer finishFree(pool, &err)
	return wrap("set storage pool autostart", pool.SetAutostart(autostart))
}
