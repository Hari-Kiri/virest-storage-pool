package storagePool

import (
	"fmt"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

// Destroy an active storage pool. This will deactivate the pool on the host, but keep any persistent config associated with it.
// If it has a persistent config it can later be restarted with storagePool.PoolCreate().
func PoolDestroy(connection virest.Connection, poolUuid string) (virest.Error, bool) {
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

	virestError.Error, isError = storagePoolObject.Destroy().(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed to destroy pool: %s", virestError.Message)
		return virestError, isError
	}

	return virestError, false
}
