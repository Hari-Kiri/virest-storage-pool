package storagePool

import (
	"net/http"
	"os"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules"
	"github.com/Hari-Kiri/virest-storage-pool/modules/utils"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDefine"
	"libvirt.org/go/libvirt"
)

func PoolDefine(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		result          string
		qemuConnection  *libvirt.Connect
		requestBodyData poolDefine.Request
		httpBody        poolDefine.Response
		libvirtError    libvirt.Error
		isError         bool
	)

	qemuConnection, libvirtError, isError = modules.RequestPrecondition(request, http.MethodPost,
		os.Getenv("VIREST_STORAGE_POOL_CONNECTION_URI"), &requestBodyData)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirtError.Code)
		httpBody.Error = libvirtError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed connecting to hypervisor [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			libvirtError.Message,
		)
		return
	}
	defer qemuConnection.Close()

	result, libvirtError, isError = modules.PoolDefine(qemuConnection, requestBodyData.StoragePool, requestBodyData.Option)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirtError.Code)
		httpBody.Error = libvirtError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed to define pool [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			libvirtError.Message,
		)
		return
	}

	httpBody.Response = true
	httpBody.Code = http.StatusCreated
	httpBody.Data.Uuid = result
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("new pool defined with uuid:", result, "[", request.URL.Path, "]")
}
