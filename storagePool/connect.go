package storagePool

import (
	"errors"

	"github.com/Hari-Kiri/virest/virestUtilities"
)

// errNoNativePool is returned when event registration needs a real libvirt pool.
var errNoNativePool = errors.New("storage pool has no native libvirt handle")

// errConnectionClosed is returned when methods are called after Close or on a nil Connection.
var errConnectionClosed = errors.New("connection is closed")

// errNilLibvirtConnect is returned when ConnectWithAuth succeeds but yields a nil *libvirt.Connect.
var errNilLibvirtConnect = errors.New("connect: nil libvirt connection")

// connectWithAuth dials the hypervisor; overridden in unit tests.
var connectWithAuth = virestUtilities.ConnectWithAuth

// Connection is a handle to a libvirt hypervisor used for storage-pool operations.
type Connection struct {
	hv hypervisor
}

// Connect opens a connection to the hypervisor at uri.
// See https://libvirt.org/uri.html for URI formats.
func Connect(uri string) (*Connection, error) {
	virestConn, err := connectWithAuth(uri, nil, 0)
	if err != nil {
		return nil, wrap("connect", err)
	}
	if virestConn.Connect == nil {
		return nil, errNilLibvirtConnect
	}
	return &Connection{hv: &libvirtHypervisor{conn: virestConn.Connect}}, nil
}

// Close releases the underlying libvirt connection.
func (c *Connection) Close() error {
	if c == nil || c.hv == nil {
		return nil
	}
	_, err := c.hv.Close()
	c.hv = nil
	return wrap("close", err)
}

// lookupPool looks up a pool by UUID and returns it. Caller must Free.
func (c *Connection) lookupPool(uuid string) (poolHandle, error) {
	if c == nil || c.hv == nil {
		return nil, errConnectionClosed
	}
	pool, err := c.hv.LookupStoragePoolByUUIDString(uuid)
	if err != nil {
		return nil, wrap("lookup storage pool", err)
	}
	return pool, nil
}
