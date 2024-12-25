package storagePool

import (
	"net/http"
	"os"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules/storagePool"
	"github.com/Hari-Kiri/virest-storage-pool/structures/getUid"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/golang-jwt/jwt"
	"libvirt.org/go/libvirt"
)

func GetUid(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		connection      *libvirt.Connect
		requestBodyData getUid.Request
		httpBody        getUid.Response
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

	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data.Uid = os.Getuid()
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("show uid for user", httpBody.Data.Uid, "[", request.URL.Path, "]")
}
