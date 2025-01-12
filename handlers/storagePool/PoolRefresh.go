package storagePool

import (
	"net/http"
	"os"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules/storagePool"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolRefresh"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/golang-jwt/jwt"
)

func PoolRefresh(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		requestBodyData poolRefresh.Request
		httpBody        poolRefresh.Response
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

	errorPoolRefresh, isErrorPoolRefresh := storagePool.PoolRefresh(connection, requestBodyData.Uuid)
	if isErrorPoolRefresh {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(errorPoolRefresh.Code)
		httpBody.Error = errorPoolRefresh
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed to refresh pool [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			errorPoolRefresh.Message,
		)
		return
	}

	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data.Uuid = requestBodyData.Uuid
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("pool", requestBodyData.Uuid, "have been refreshed on hypervisor", request.Header["Hypervisor-Uri"][0], "[", request.URL.Path, "]")
}
