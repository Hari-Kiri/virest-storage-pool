package storagePool

import (
	"net/http"
	"os"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules/storagePool"
	"github.com/Hari-Kiri/virest-storage-pool/structures/findStoragePoolSources"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/golang-jwt/jwt"
)

func FindStoragePoolSource(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		requestBodyData findStoragePoolSources.Request
		httpBody        findStoragePoolSources.Response
	)

	connection, errorRequestPrecondition, isError := storagePool.RequestPrecondition(
		request,
		http.MethodPost,
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

	findStoragePoolSource, errorFindStoragePoolSource, isErrorFindStoragePoolSource := storagePool.FindStoragePoolSource(
		connection,
		requestBodyData.Type,
		requestBodyData.SrcSpec.Source,
	)
	if isErrorFindStoragePoolSource {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(errorFindStoragePoolSource.Code)
		httpBody.Error = errorFindStoragePoolSource
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed to find pool sources [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			errorFindStoragePoolSource.Message,
		)
		return
	}

	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data = findStoragePoolSource
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("find potential pool", requestBodyData.Type, "sources [", request.URL.Path, "]")
}
