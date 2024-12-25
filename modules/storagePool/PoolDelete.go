package storagePool

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

// Delete the underlying pool resources. This is a non-recoverable operation.
func PoolDelete(connection *libvirt.Connect, poolUuid string, option libvirt.StoragePoolDeleteFlags) (libvirt.Error, bool) {
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
	libvirtError, isError = storagePoolObject.Delete(option).(libvirt.Error)
	if isError {
		libvirtError.Message = fmt.Sprintf("failed to delete pool: %s", libvirtError.Message)
		return libvirtError, isError
	}

	return libvirtError, false
}
