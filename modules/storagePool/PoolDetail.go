package storagePool

import (
	"fmt"

	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDetail"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
)

func getPoolDetail(libvirtStoragePoolObject libvirt.StoragePool, libvirtStorageXMLFlags libvirt.StorageXMLFlags) (libvirtxml.StoragePool, libvirt.Error, bool) {
	var (
		libvirtError libvirt.Error
		isError      bool
	)

	storagePoolXml, errorGetStoragePoolXml := libvirtStoragePoolObject.GetXMLDesc(libvirtStorageXMLFlags)
	libvirtError, isError = errorGetStoragePoolXml.(libvirt.Error)
	if errorGetStoragePoolXml != nil {
		libvirtError.Message = fmt.Sprintf("failed get XML of pool: %s", libvirtError.Message)
		return libvirtxml.StoragePool{}, libvirtError, isError
	}

	var result libvirtxml.StoragePool
	errorUnmarshallStoragePool := result.Unmarshal(storagePoolXml)
	libvirtError, isError = errorUnmarshallStoragePool.(libvirt.Error)
	if errorUnmarshallStoragePool != nil {
		libvirtError.Message = fmt.Sprintf("failed Unmarshal XML of pool: %s", libvirtError.Message)
		return libvirtxml.StoragePool{}, libvirtError, isError
	}

	return result, libvirtError, false
}

// Fetch an XML document describing all aspects of the storage pool.
func PoolDetail(connection virest.Connection, poolUuid string, option uint) (poolDetail.Detail, virest.Error, bool) {
	var (
		virestError virest.Error
		isError     bool
	)

	storagePoolObject, errorGetStoragePoolObject := connection.LookupStoragePoolByUUIDString(poolUuid)
	virestError.Error, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed get storage pool object: %s", virestError.Message)
		return poolDetail.Detail{}, virestError, isError
	}
	defer storagePoolObject.Free()

	var result libvirtxml.StoragePool
	result, virestError.Error, isError = getPoolDetail(*storagePoolObject, libvirt.StorageXMLFlags(option))
	if isError {
		virestError.Message = fmt.Sprintf("failed get storage pool using object: %s", virestError.Message)
		return poolDetail.Detail{}, virestError, isError
	}

	return poolDetail.Detail{
		StoragePool: result,
	}, virestError, false
}
