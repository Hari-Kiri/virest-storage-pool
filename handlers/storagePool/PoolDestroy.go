package storagePool

import (
	"net/http"
	"os"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules/storagePool"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDestroy"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/golang-jwt/jwt"
)

func PoolDestroy(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		requestBodyData poolDestroy.Request
		httpBody        poolDestroy.Response
	)

	connection, errorRequestPrecondition, isError := storagePool.RequestPrecondition(
		request,
		http.MethodPatch,
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

	errorPoolDestroy, isErrorPoolDestroy := storagePool.PoolDestroy(connection, requestBodyData.Uuid)
	if isErrorPoolDestroy {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(errorPoolDestroy.Code)
		httpBody.Error = errorPoolDestroy
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed to destroy pool [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			errorPoolDestroy.Message,
		)
		return
	}

	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data.Uuid = requestBodyData.Uuid
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("pool", requestBodyData.Uuid, "destroyed on hypervisor", request.Header["Hypervisor-Uri"][0], "[", request.URL.Path, "]")
}
