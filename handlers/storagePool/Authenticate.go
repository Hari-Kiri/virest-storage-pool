package storagePool

import (
	"net/http"
	"os"
	"time"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules/storagePool"
	"github.com/Hari-Kiri/virest-storage-pool/structures/authenticate"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/golang-jwt/jwt"
	"libvirt.org/go/libvirt"
)

func Authenticate(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		httpBody     authenticate.Response
		result       string
		libvirtError libvirt.Error
		isError      bool
	)

	result, libvirtError, isError = storagePool.Authenticate(
		request,
		os.Getenv("VIREST_STORAGE_POOL_APPLICATION_NAME"),
		time.Minute*1,
		jwt.SigningMethodHS256,
		[]byte("lorenzo lamas"),
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

	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data.Token = result
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("listing pool on hypervisor", request.Header["Hypervisor-Uri"][0], "[", request.URL.Path, "]")
}
