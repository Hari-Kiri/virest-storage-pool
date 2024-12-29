package storagePool

import (
	"net/http"
	"os"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules/storagePool"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolCapabilities"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/golang-jwt/jwt"
)

func PoolCapabilities(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		requestBodyData poolCapabilities.Request
		httpBody        poolCapabilities.Response
	)

	connection, errorRequestPrecondition, isError := storagePool.RequestPrecondition(
		request,
		http.MethodGet,
		&requestBodyData,
		os.Getenv("VIREST_STORAGE_POOL_APPLICATION_NAME"),
		jwt.SigningMethodHS512,
		[]byte(os.Getenv("VIREST_STORAGE_POOL_APPLICATION_JWT_SIGNATURE_KEY")),
	)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(errorRequestPrecondition.Code)
		httpBody.Error = errorRequestPrecondition
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"request unexpected [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			errorRequestPrecondition.Message,
		)
		return
	}
	defer connection.Close()

	poolCapabilities, errorGetPoolCapabilities, isErrorGetPoolCapabilities := storagePool.PoolCapabilities(connection)
	if isErrorGetPoolCapabilities {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(errorGetPoolCapabilities.Code)
		httpBody.Error = errorGetPoolCapabilities
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed to get pool capabilities [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			errorGetPoolCapabilities.Message,
		)
		return
	}

	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data = poolCapabilities
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("get hypervisor", request.Header["Hypervisor-Uri"], "storage pool capabilities [", request.URL.Path, "]")
}
