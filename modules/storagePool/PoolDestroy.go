package storagePool

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

// Destroy an active storage pool. This will deactivate the pool on the host, but keep any persistent config associated with it.
// If it has a persistent config it can later be restarted with storagePool.PoolCreate().
func PoolDestroy(connection *libvirt.Connect, poolUuid string) (libvirt.Error, bool) {
	var (
		libvirtError libvirt.Error
		isError      bool
	)

	// Get libvirt storage pool object
	storagePoolObject, errorGetStoragePoolObject := connection.LookupStoragePoolByUUIDString(poolUuid)
	libvirtError, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		libvirtError.Message = fmt.Sprintf("failed get storage pool object: %s", libvirtError.Message)
		return libvirtError, isError
	}
	defer storagePoolObject.Free()

	// Build pool
	libvirtError, isError = storagePoolObject.Destroy().(libvirt.Error)
	if isError {
		libvirtError.Message = fmt.Sprintf("failed to destroy pool: %s", libvirtError.Message)
		return libvirtError, isError
	}

	return libvirtError, false
}
