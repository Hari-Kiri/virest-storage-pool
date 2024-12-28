package storagePool

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/structures/authenticate"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/Hari-Kiri/virest-utilities/utils/auth"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"github.com/golang-jwt/jwt"
	"libvirt.org/go/libvirt"
)

func Authenticate(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		httpBody    authenticate.Response
		result      string
		virestError virest.Error
		isError     bool
	)

	jwtLifetimeSeconds, errorParseJwtLifetimeSeconds := strconv.Atoi(os.Getenv("VIREST_STORAGE_POOL_APPLICATION_JWT_LIFETIME_SECONDS"))
	if errorParseJwtLifetimeSeconds != nil {
		virestError = virest.Error{Error: libvirt.Error{
			Code:    libvirt.ERR_INTERNAL_ERROR,
			Domain:  libvirt.FROM_AUTH,
			Message: "environment variable for JWT lifetime not number or not exist",
			Level:   libvirt.ERR_ERROR,
		}}
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirt.ERR_INVALID_ARG)
		httpBody.Error = virestError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"request unexpected [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			virestError.Message,
		)
		return
	}

	result, virestError, isError = auth.BasicAuth(
		request,
		os.Getenv("VIREST_STORAGE_POOL_APPLICATION_NAME"),
		time.Second*time.Duration(jwtLifetimeSeconds),
		jwt.SigningMethodHS512,
		[]byte(os.Getenv("VIREST_STORAGE_POOL_APPLICATION_JWT_SIGNATURE_KEY")),
	)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(virestError.Code)
		httpBody.Error = virestError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"request unexpected [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			virestError.Message,
		)
		return
	}

	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data.Token = result
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("authenticate user", "[", request.URL.Path, "]")
}
