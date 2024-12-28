package storagePool

import (
	"fmt"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

// Undefine storage pool using their UUID string. Return libvirt.error nil on success, or libvirt.error not nil upon failure.
func PoolUndefine(connection virest.Connection, poolUuid string) (virest.Error, bool) {
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

	// Undefine pool
	virestError.Error, isError = storagePoolObject.Undefine().(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed to undefine pool: %s", virestError.Message)
		return virestError, isError
	}

	return virestError, false
}
