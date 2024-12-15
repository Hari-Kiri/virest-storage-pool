package storagePool

import (
	"fmt"

	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

// Fetch an XML document describing all aspects of the storage pool.
//
// Get detail when pool in active state:
//
//	libvirtStorageXMLFlags = 0
//
// Get detail when pool in inactive state:
//
//	libvirtStorageXMLFlags = 1
func poolDetail(libvirtStoragePoolObject libvirt.StoragePool, libvirtStorageXMLFlags libvirt.StorageXMLFlags) (libvirtxml.StoragePool, libvirt.Error, bool) {
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
//
// Get detail when pool in active state:
//
//	libvirtStorageXMLFlags = 0
//
// Get detail when pool in inactive state:
//
//	libvirtStorageXMLFlags = 1
func PoolDetail(connection *libvirt.Connect, poolUuid string, option libvirt.StorageXMLFlags) (libvirtxml.StoragePool, libvirt.Error, bool) {
	var (
		libvirtError libvirt.Error
		isError      bool
	)

	storagePoolObject, errorGetStoragePoolObject := connection.LookupStoragePoolByUUIDString(poolUuid)
	libvirtError, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		libvirtError.Message = fmt.Sprintf("failed get storage pool object: %s", libvirtError.Message)
		return libvirtxml.StoragePool{}, libvirtError, isError
	}
	defer storagePoolObject.Free()

	var result libvirtxml.StoragePool
	result, libvirtError, isError = poolDetail(*storagePoolObject, option)
	if isError {
		libvirtError.Message = fmt.Sprintf("failed get storage pool using object: %s", libvirtError.Message)
		return libvirtxml.StoragePool{}, libvirtError, isError
	}

	return result, libvirtError, false
}
