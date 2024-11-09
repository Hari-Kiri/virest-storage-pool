package storagePool

import (
	"net/http"
	"os"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules"
	"github.com/Hari-Kiri/virest-storage-pool/modules/utils"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolUndefine"
	"libvirt.org/go/libvirt"
)

func PoolUndefine(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		qemuConnection  *libvirt.Connect
		requestBodyData poolUndefine.Request
		httpBody        poolUndefine.Response
		libvirtError    libvirt.Error
		isError         bool
	)

	qemuConnection, libvirtError, isError = modules.RequestPrecondition(request, http.MethodPatch,
		os.Getenv("VIREST_STORAGE_POOL_CONNECTION_URI"), &requestBodyData)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirtError.Code)
		httpBody.Error = libvirtError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		return
	}
	defer qemuConnection.Close()

	// Get libvirt storage pool object
	storagePoolObject, errorGetStoragePoolObject := qemuConnection.LookupStoragePoolByUUIDString(requestBodyData.Uuid)
	libvirtError, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirtError.Code)
		httpBody.Error = libvirtError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed get storage pool object [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			libvirtError.Message,
		)
		return
	}
	defer storagePoolObject.Free()

	// Undefine pool
	libvirtError, isError = storagePoolObject.Undefine().(libvirt.Error)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirtError.Code)
		httpBody.Error = libvirtError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed to undefine pool [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			libvirtError.Message,
		)
		return
	}

	// Http ok response
	utils.NoContentResponseBuilder(responseWriter)
	temboLog.InfoLogging("pool", requestBodyData.Uuid, "undefined [", request.URL.Path, "]")
}
