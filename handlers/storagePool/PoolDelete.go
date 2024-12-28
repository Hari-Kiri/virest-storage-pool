package storagePool

import (
	"net/http"
	"os"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules/storagePool"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDelete"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/golang-jwt/jwt"
)

func PoolDelete(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		requestBodyData poolDelete.Request
		httpBody        poolDelete.Response
	)

	connection, errorRequestPrecondition, isError := storagePool.RequestPrecondition(
		request,
		http.MethodDelete,
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

	errorPoolDelete, isError := storagePool.PoolDelete(connection, requestBodyData.Uuid, requestBodyData.Option)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(errorPoolDelete.Code)
		httpBody.Error = errorPoolDelete
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed to delete pool [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			errorPoolDelete.Message,
		)
		return
	}

	utils.NoContentResponseBuilder(responseWriter)
	temboLog.InfoLogging("pool", requestBodyData.Uuid, "deleted [", request.URL.Path, "]")
}
