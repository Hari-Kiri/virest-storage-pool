package modules

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

func PoolUndefine(qemuConnection *libvirt.Connect, poolUuid string) (libvirt.Error, bool) {
	var (
		libvirtError libvirt.Error
		isError      bool
	)

	// Get libvirt storage pool object
	storagePoolObject, errorGetStoragePoolObject := qemuConnection.LookupStoragePoolByUUIDString(poolUuid)
	libvirtError, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		libvirtError.Message = fmt.Sprintf("failed get storage pool object: %s", libvirtError.Message)
		return libvirtError, isError
	}
	defer storagePoolObject.Free()

	// Undefine pool
	libvirtError, isError = storagePoolObject.Undefine().(libvirt.Error)
	if isError {
		libvirtError.Message = fmt.Sprintf("failed to undefine pool: %s", libvirtError.Message)
		return libvirtError, isError
	}

	return libvirtError, false
}
