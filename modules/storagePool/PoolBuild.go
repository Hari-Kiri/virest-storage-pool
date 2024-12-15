package storagePool

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

// Build the underlying storage pool.
//
// connection:
//
//	pointer to hypervisor connection
//
// option:
//
//	bitwise-OR of libvirt.StoragePoolBuildFlags
//
// Returns:
//
//	libvirt.error nil on success, or libvirt.error not nil upon failure
func PoolBuild(connection *libvirt.Connect, poolUuid string, option libvirt.StoragePoolBuildFlags) (libvirt.Error, bool) {
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
	libvirtError, isError = storagePoolObject.Build(option).(libvirt.Error)
	if isError {
		libvirtError.Message = fmt.Sprintf("failed to build pool: %s", libvirtError.Message)
		return libvirtError, isError
	}

	return libvirtError, false
}
