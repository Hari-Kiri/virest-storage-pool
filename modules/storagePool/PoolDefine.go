package storagePool

import (
	"fmt"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDefine"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

// Define new storage pool using json formatted data with option as define flags.
// The option with UInteger 1 will validate the JSON document against libvirt schema, while the option with UInteger 0 does nothing.
// Upon success, the UUID of the newly defined pool will be returned.
func PoolDefine(connection virest.Connection, storagePool libvirtxml.StoragePool, option libvirt.StoragePoolDefineFlags) (poolDefine.Uuid, virest.Error, bool) {
	var (
		virestError virest.Error
		isError     bool
	)

	// Convert request body to libvirt xml
	libvirtXml, errorGetLibvirtXml := storagePool.Marshal()
	virestError.Error, isError = errorGetLibvirtXml.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed to create pool config xml: %s", virestError.Message)
		return poolDefine.Uuid{}, virestError, isError
	}

	// Define pool
	definePool, errorDefinePool := connection.StoragePoolDefineXML(libvirtXml, option)
	virestError.Error, isError = errorDefinePool.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed to define new pool: %s", virestError.Message)
		return poolDefine.Uuid{}, virestError, isError
	}
	defer definePool.Free()

	// Get defined pool UUID
	definedPoolUuid, errorGetDefinedPoolUuid := definePool.GetUUIDString()
	virestError.Error, isError = errorGetDefinedPoolUuid.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed to get defined pool UUID: %s", virestError.Message)
		return poolDefine.Uuid{}, virestError, isError
	}

	return poolDefine.Uuid{
		Uuid: definedPoolUuid,
	}, virestError, false
}
