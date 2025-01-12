package storagePool

import (
	"fmt"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

// Request that the pool refresh its list of volumes. This may involve communicating with a remote server,
// and/or initializing new devices at the OS layer.
func PoolRefresh(connection virest.Connection, poolUuid string) (virest.Error, bool) {
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

	// extra flags; not used yet, so callers should always pass 0
	// https://libvirt.org/html/libvirt-libvirt-storage.html#virStoragePoolRefresh
	virestError.Error, isError = storagePoolObject.Refresh(0).(libvirt.Error)
	return virestError, isError
}
