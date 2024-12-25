package storagePool

import (
	"net/http"
	"os"
	"strconv"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules/storagePool"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDetail"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/golang-jwt/jwt"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

func PoolDetail(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		result          libvirtxml.StoragePool
		connection      *libvirt.Connect
		requestBodyData poolDetail.Request
		httpBody        poolDetail.Response
		libvirtError    libvirt.Error
		isError         bool
	)

	connection, libvirtError, isError = storagePool.RequestPrecondition(
		request,
		http.MethodGet,
		&requestBodyData,
		os.Getenv("VIREST_STORAGE_POOL_APPLICATION_NAME"),
		jwt.SigningMethodHS512,
		[]byte(os.Getenv("VIREST_STORAGE_POOL_APPLICATION_JWT_SIGNATURE_KEY")),
	)
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

	option, errorParseOptionToUint := strconv.ParseUint(requestBodyData.Option, 10, 32)
	if errorParseOptionToUint != nil {
		libvirtError = libvirt.Error{
			Code:    libvirt.ERR_INVALID_ARG,
			Domain:  libvirt.FROM_NET,
			Message: "'Option' value not number or not exist",
			Level:   libvirt.ERR_ERROR,
		}
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirt.ERR_INVALID_ARG)
		httpBody.Error = libvirtError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"request unexpected [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			libvirtError.Message,
		)
		return
	}

	result, libvirtError, isError = storagePool.PoolDetail(connection, requestBodyData.Uuid, libvirt.StorageXMLFlags(option))
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirtError.Code)
		httpBody.Error = libvirtError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed list storage pool [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			libvirtError.Message,
		)
		return
	}

	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data = result
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("get pool detail:", requestBodyData.Uuid, "[", request.URL.Path, "]")
}
