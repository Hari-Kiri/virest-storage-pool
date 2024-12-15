package storagePool

import (
	"net/http"
	"os"
	"strconv"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules/storagePool"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolList"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"libvirt.org/go/libvirt"
)

func PoolList(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		result          []poolList.Data
		connection      *libvirt.Connect
		requestBodyData poolList.Request
		httpBody        poolList.Response
		libvirtError    libvirt.Error
		isError         bool
	)

	connection, libvirtError, isError = storagePool.RequestPrecondition(request, http.MethodGet,
		os.Getenv("VIREST_STORAGE_POOL_CONNECTION_URI"), &requestBodyData)
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

	inactive, errorParseInactiveToUint := strconv.ParseUint(requestBodyData.Inactive, 10, 32)
	if errorParseInactiveToUint != nil {
		libvirtError = libvirt.Error{
			Code:    libvirt.ERR_INVALID_ARG,
			Domain:  libvirt.FROM_NET,
			Message: "'Inactive' value not number or not exist",
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

	result, libvirtError, isError = storagePool.PoolList(connection, libvirt.ConnectListAllStoragePoolsFlags(option),
		libvirt.StorageXMLFlags(inactive))
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
	temboLog.InfoLogging("listing pool on hypervisor", os.Getenv("VIREST_STORAGE_POOL_CONNECTION_URI"), "[", request.URL.Path, "]")
}
