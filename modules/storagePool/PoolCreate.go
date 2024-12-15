package storagePool

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

// Starts an inactive storage pool. Return libvirt.error nil on success, or libvirt.error not nil upon failure
func PoolCreate(connection *libvirt.Connect, poolUuid string, option libvirt.StoragePoolCreateFlags) (libvirt.Error, bool) {
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

	// Create pool
	libvirtError, isError = storagePoolObject.Create(option).(libvirt.Error)
	if isError {
		libvirtError.Message = fmt.Sprintf("failed to create pool: %s", libvirtError.Message)
		return libvirtError, isError
	}

	return libvirtError, false
}
