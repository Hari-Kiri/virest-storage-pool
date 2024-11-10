package modules

import (
	"fmt"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDefine"
	"libvirt.org/go/libvirt"
)

func PoolDefine(qemuConnection *libvirt.Connect, requestBodyData poolDefine.Request) (string, libvirt.Error, bool) {
	var (
		libvirtError libvirt.Error
		isError      bool
	)

	// Convert request body to libvirt xml
	libvirtXml, errorGetLibvirtXml := requestBodyData.StoragePool.Marshal()
	libvirtError, isError = errorGetLibvirtXml.(libvirt.Error)
	if isError {
		libvirtError.Message = fmt.Sprintf("failed to create pool config xml: %s", libvirtError.Message)
		return "", libvirtError, isError
	}

	// Define pool
	definePool, errorDefinePool := qemuConnection.StoragePoolDefineXML(libvirtXml, requestBodyData.Option)
	libvirtError, isError = errorDefinePool.(libvirt.Error)
	if isError {
		libvirtError.Message = fmt.Sprintf("failed to define new pool: %s", libvirtError.Message)
		return "", libvirtError, isError
	}
	defer definePool.Free()

	// Get defined pool UUID
	definedPoolUuid, errorGetDefinedPoolUuid := definePool.GetUUIDString()
	libvirtError, isError = errorGetDefinedPoolUuid.(libvirt.Error)
	if isError {
		libvirtError.Message = fmt.Sprintf("failed to get defined pool UUID: %s", libvirtError.Message)
		return "", libvirtError, isError
	}

	return definedPoolUuid, libvirtError, false
}
