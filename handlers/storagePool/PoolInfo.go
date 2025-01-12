package storagePool

import (
	"net/http"
	"os"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules/storagePool"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolInfo"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/golang-jwt/jwt"
)

func PoolInfo(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		requestBodyData poolInfo.Request
		httpBody        poolInfo.Response
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

	result, errorGetPoolInfo, isErrorGetPoolInfo := storagePool.PoolInfo(connection, requestBodyData.Uuid)
	if isErrorGetPoolInfo {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(errorGetPoolInfo.Code)
		httpBody.Error = errorGetPoolInfo
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed to get pool info '"+requestBodyData.Uuid+"' [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			errorGetPoolInfo.Message,
		)
		return
	}

	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data = result
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("get pool info with uuid:", result.Uuid, "on hypervisor", request.Header["Hypervisor-Uri"][0], "[", request.URL.Path, "]")
}
