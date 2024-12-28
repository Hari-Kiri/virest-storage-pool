package storagePool

import (
	"fmt"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

// Configure the storage pool to be automatically started when the host machine boots.
func PoolAutostart(connection virest.Connection, poolUuid string, autostart bool) (virest.Error, bool) {
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

	// Build pool
	virestError.Error, isError = storagePoolObject.SetAutostart(autostart).(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed to set pool '%s' autostart status: %s", poolUuid, virestError.Message)
		return virestError, isError
	}

	return virestError, false
}
