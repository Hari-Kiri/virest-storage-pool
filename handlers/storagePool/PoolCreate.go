package storagePool

import (
	"net/http"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules/storagePool"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolCreate"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"libvirt.org/go/libvirt"
)

func PoolCreate(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		connection      *libvirt.Connect
		requestBodyData poolCreate.Request
		httpBody        poolCreate.Response
		libvirtError    libvirt.Error
		isError         bool
	)

	connection, libvirtError, isError = storagePool.RequestPrecondition(request, http.MethodPatch, &requestBodyData)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirtError.Code)
		httpBody.Error = libvirtError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"request unexpected [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			libvirtError.Message,
		)
		return
	}
	defer connection.Close()

	libvirtError, isError = storagePool.PoolCreate(connection, requestBodyData.Uuid, requestBodyData.Option)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirtError.Code)
		httpBody.Error = libvirtError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed to build pool [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			libvirtError.Message,
		)
		return
	}

	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data.Uuid = requestBodyData.Uuid
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("pool", requestBodyData.Uuid, "created [", request.URL.Path, "]")
}
