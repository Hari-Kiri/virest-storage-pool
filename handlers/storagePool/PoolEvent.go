package storagePool

import (
	"net/http"
	"os"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules/storagePool"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolEvent"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/golang-jwt/jwt"
)

func PoolEvent(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		result          poolEvent.Event
		requestBodyData poolEvent.Request
		httpBody        poolEvent.Response
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

	types, errorParseTypesToUint, isErrorParseTypesToUint := utils.StringToUint(requestBodyData.Types)
	if isErrorParseTypesToUint {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(errorParseTypesToUint.Code)
		httpBody.Error = errorParseTypesToUint
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"request unexpected [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			errorParseTypesToUint.Message,
		)
		return
	}

	errorGetStoragePoolEvent, isErrorGetStoragePoolEvent := storagePool.PoolEvent(
		connection,
		requestBodyData.Uuid,
		types,
	)
	if isErrorGetStoragePoolEvent {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(errorGetStoragePoolEvent.Code)
		httpBody.Error = errorGetStoragePoolEvent
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed get pool event [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			errorGetStoragePoolEvent.Message,
		)
		return
	}

	for {
		if storagePool.PoolEventProbingResult.EventRefresh > 0 {
			result.EventRefresh = storagePool.PoolEventProbingResult.EventRefresh
			break
		}
		if storagePool.PoolEventProbingResult.EventLifecycle.Event < 6 {
			result.EventLifecycle = storagePool.PoolEventProbingResult.EventLifecycle
			break
		}
	}

	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data = result
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("get pool event with uuid:", requestBodyData.Uuid, "on hypervisor", request.Header["Hypervisor-Uri"][0], "[", request.URL.Path, "]")
}
