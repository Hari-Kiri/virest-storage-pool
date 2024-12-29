package storagePool

import (
	"fmt"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

// Starts an inactive storage pool.
func PoolCreate(connection virest.Connection, poolUuid string, option libvirt.StoragePoolCreateFlags) (virest.Error, bool) {
	var (
		virestError virest.Error
		isError     bool
	)

	storagePoolObject, errorGetStoragePoolObject := connection.LookupStoragePoolByUUIDString(poolUuid)
	virestError.Error, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed get storage pool object: %s", virestError.Message)
		return virestError, isError
	}
	defer storagePoolObject.Free()

	virestError.Error, isError = storagePoolObject.Create(option).(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed to create pool: %s", virestError.Message)
		return virestError, isError
	}

	return virestError, false
}
