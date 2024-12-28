package storagePool

import (
	"fmt"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

// Starts an inactive storage pool. Return libvirt.error nil on success, or libvirt.error not nil upon failure
func PoolCreate(connection virest.Connection, poolUuid string, option libvirt.StoragePoolCreateFlags) (virest.Error, bool) {
	var (
		virestError virest.Error
		isError     bool
	)

	// Get libvirt storage pool object
	storagePoolObject, errorGetStoragePoolObject := connection.LookupStoragePoolByUUIDString(poolUuid)
	virestError.Error, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed get storage pool object: %s", virestError.Message)
		return virestError, isError
	}
	defer storagePoolObject.Free()

	// Create pool
	virestError.Error, isError = storagePoolObject.Create(option).(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed to create pool: %s", virestError.Message)
		return virestError, isError
	}

	return virestError, false
}
