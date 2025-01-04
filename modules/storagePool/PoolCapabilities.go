package storagePool

import (
	"fmt"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolCapabilities"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/libvirtxml"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

// Prior creating a storage pool it may be suitable to know what pool types are supported along with the file/disk
// formats for each pool.
func PoolCapabilities(connection virest.Connection) (poolCapabilities.StoragepoolCapabilities, virest.Error, bool) {
	var (
		virestError virest.Error
		isError     bool
	)

	// extra flags; not used yet, so callers should always pass 0
	// https://libvirt.org/html/libvirt-libvirt-storage.html#virConnectGetStoragePoolCapabilities
	storagepoolCapabilities, errorGetStoragePoolCapabilities := connection.GetStoragePoolCapabilities(0)
	virestError.Error, isError = errorGetStoragePoolCapabilities.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed get storage pool capabilities: %s", virestError.Message)
		return poolCapabilities.StoragepoolCapabilities{}, virestError, isError
	}

	var result libvirtxml.StoragepoolCapabilities
	virestError, isError = result.Unmarshal(storagepoolCapabilities)
	if isError {
		virestError.Message = fmt.Sprintf("failed parse storage pool capabilities: %s", virestError.Message)
		return poolCapabilities.StoragepoolCapabilities{}, virestError, isError
	}
	return poolCapabilities.StoragepoolCapabilities{
		StoragepoolCapabilities: result,
	}, virestError, false
}
