package storagePool

import (
	"net/http"
	"os"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules/storagePool"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDetail"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/golang-jwt/jwt"
)

func PoolDetail(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		requestBodyData poolDetail.Request
		httpBody        poolDetail.Response
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

	option, errorParseOptionToUint, isErrorParseOptionToUint := utils.StringToUint(requestBodyData.Option)
	if isErrorParseOptionToUint {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(errorParseOptionToUint.Code)
		httpBody.Error = errorParseOptionToUint
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"request unexpected [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			errorParseOptionToUint.Message,
		)
		return
	}

	result, errorGetPoolDetail, isErrorGetPoolDetail := storagePool.PoolDetail(connection, requestBodyData.Uuid, option)
	if isErrorGetPoolDetail {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(errorGetPoolDetail.Code)
		httpBody.Error = errorGetPoolDetail
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed get pool detail [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			errorGetPoolDetail.Message,
		)
		return
	}

	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data = result
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("get pool detail:", requestBodyData.Uuid, "on hypervisor", request.Header["Hypervisor-Uri"][0], "[", request.URL.Path, "]")
}
