package storagePool

import (
	"fmt"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

// Undefine storage pool using their UUID string.
func (poolConnection *poolConnection) PoolUndefine(poolUuid string) (virest.Error, bool) {
	var (
		virestError virest.Error
		isError     bool
	)

	storagePoolObject, errorGetStoragePoolObject := poolConnection.LookupStoragePoolByUUIDString(poolUuid)
	virestError.Error, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed get storage pool object: %s", virestError.Message)
		return virestError, isError
	}
	defer storagePoolObject.Free()

	virestError.Error, isError = storagePoolObject.Undefine().(libvirt.Error)
	return virestError, isError
}
