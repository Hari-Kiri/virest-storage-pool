package storagePool

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

// Configure the storage pool to be automatically started when the host machine boots.
func PoolAutostart(connection *libvirt.Connect, poolUuid string, autostart bool) (libvirt.Error, bool) {
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
	libvirtError, isError = storagePoolObject.SetAutostart(autostart).(libvirt.Error)
	if isError {
		libvirtError.Message = fmt.Sprintf("failed to set pool '%s' autostart status: %s", poolUuid, libvirtError.Message)
		return libvirtError, isError
	}

	return libvirtError, false
}
