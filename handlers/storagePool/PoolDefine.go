package storagePool

import (
	"net/http"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules/utils"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDefine"
	"libvirt.org/go/libvirt"
)

func PoolDefine(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		qemuConnection  *libvirt.Connect
		requestBodyData poolDefine.Request
		httpBody        poolDefine.Response
		libvirtError    libvirt.Error
		isError         bool
	)

	qemuConnection, libvirtError, isError = utils.RequestPrecondition(request, http.MethodPost, &requestBodyData)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirtError.Code)
		httpBody.Error = libvirtError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		return
	}
	defer qemuConnection.Close()

	// Convert request body to libvirt xml
	libvirtXml, errorGetLibvirtXml := requestBodyData.StoragePool.Marshal()
	libvirtError, isError = errorGetLibvirtXml.(libvirt.Error)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirtError.Code)
		httpBody.Error = libvirtError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed to create pool config xml [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			libvirtError.Message,
		)
		return
	}

	// Define pool
	definePool, errorDefinePool := qemuConnection.StoragePoolDefineXML(libvirtXml, requestBodyData.Option)
	libvirtError, isError = errorDefinePool.(libvirt.Error)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirtError.Code)
		httpBody.Error = libvirtError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed to define new pool [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			libvirtError.Message,
		)
		return
	}
	defer definePool.Free()

	// Get defined pool UUID
	definedPoolUuid, errorGetDefinedPoolUuid := definePool.GetUUIDString()
	libvirtError, isError = errorGetDefinedPoolUuid.(libvirt.Error)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirtError.Code)
		httpBody.Error = libvirtError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed to get defined pool UUID [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			libvirtError.Message,
		)
		return
	}

	// Http ok response
	httpBody.Response = true
	httpBody.Code = http.StatusCreated
	httpBody.Data.Uuid = definedPoolUuid
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("new pool defined with uuid:", definedPoolUuid, "[", request.URL.Path, "]")
}
