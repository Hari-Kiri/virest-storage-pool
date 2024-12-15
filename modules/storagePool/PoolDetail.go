package storagePool

import (
	"github.com/Hari-Kiri/temboLog"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

// Fetch an XML document describing all aspects of the storage pool.
// This is suitable for later feeding back into the virStoragePoolCreateXML method.
//
// Get detail when pool in active state:
//
//	libvirtStorageXMLFlags = 0
//
// Get detail when pool in inactive state:
//
//	libvirtStorageXMLFlags = 1
func PoolDetail(libvirtStoragePool libvirt.StoragePool, libvirtStorageXMLFlags libvirt.StorageXMLFlags) (libvirtxml.StoragePool, libvirt.Error, bool) {
	var (
		libvirtError libvirt.Error
		isError      bool
	)

	storagePoolXml, errorGetStoragePoolXml := libvirtStoragePool.GetXMLDesc(libvirtStorageXMLFlags)
	libvirtError, isError = errorGetStoragePoolXml.(libvirt.Error)
	if errorGetStoragePoolXml != nil {
		temboLog.ErrorLogging("failed get XML of pool", errorGetStoragePoolXml)
		return libvirtxml.StoragePool{}, libvirtError, isError
	}

	var result libvirtxml.StoragePool
	errorUnmarshallStoragePool := result.Unmarshal(storagePoolXml)
	libvirtError, isError = errorUnmarshallStoragePool.(libvirt.Error)
	if errorUnmarshallStoragePool != nil {
		temboLog.ErrorLogging("failed Unmarshal XML of pool", errorUnmarshallStoragePool)
		return libvirtxml.StoragePool{}, libvirtError, isError
	}

	return result, libvirtError, false
}
