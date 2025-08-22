package storagePool

import (
	"fmt"
	"net/http"

	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/Hari-Kiri/virest-utilities/utils/auth"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"github.com/golang-jwt/jwt"
	"libvirt.org/go/libvirt"
)

// Connection pointer to the targeted hypervisor for storage pool management purposes.
type poolConnection virest.Connection

// This function simplifies the preconditioning process after a HTTP request from a client occurs:
//
//   - Validate user request using given bearer token which is generated JWT by ViRest Utilities 'auth.BasicAuth()' module.
//
//   - Look up hypervisor URI on 'Hypervisor-Uri' request header field.
//
//   - Connect to hypervisor via SSH tunnel and check the expected HTTP request method then convert the JSON request body to structure if any.
//
//   - SSH tunnel work with Key-Based authentication (Please, create SSH Key on the host and copy it on the remote libvirt-daemon host '~/.ssh/authorized_keys')
//
// Notes for HTTP GET method:
//
//   - Query parameter and structure field will be compared in case sensitive.
//
//   - Every structure field data type must be string, so You must convert it to the right data type before You use it.
//
//   - Untested for array query argument.
//
// Notes for HTTP POST, PUT, PATCH and DELETE method:
//
//   - This function always looking for request body for data and parse them to 'structure' parameter.
func HttpRequestPrecondition[RequestStructure utils.RequestStructure](
	httpRequest *http.Request,
	expectedRequestMethod string,
	structure *RequestStructure,
	applicationName string,
	jwtSigningMethod *jwt.SigningMethodHMAC,
	jwtSignatureKey []byte,
) (poolConnection, virest.Error, bool) {
	libvirtErrorAuth, isErrorAuth := auth.BearerTokenAuth(
		httpRequest,
		applicationName,
		jwtSigningMethod,
		jwtSignatureKey,
	)
	if isErrorAuth {
		return poolConnection{}, virest.Error{Error: libvirt.Error{
			Code:    libvirt.ERR_AUTH_FAILED,
			Domain:  libvirt.FROM_NET,
			Message: fmt.Sprintf("authentication failed: %s", libvirtErrorAuth.Message),
			Level:   libvirt.ERR_ERROR,
		}}, true
	}

	virestConnectionChannel := make(chan virest.Connection)
	errorConnectChannel := make(chan virest.Error)
	isErrorConnectChannel := make(chan bool)
	go func() {
		var (
			virestConnection virest.Connection
			errorConnect     virest.Error
			isErrorConnect   bool
		)

		if len(httpRequest.Header["Hypervisor-Uri"]) == 0 {
			virestConnection = virest.Connection{}
			errorConnect = virest.Error{Error: libvirt.Error{
				Code:    libvirt.ERR_INVALID_CONN,
				Domain:  libvirt.FROM_NET,
				Message: "hypervisor uri not exist on request header",
				Level:   libvirt.ERR_ERROR,
			}}
			isErrorConnect = true

			virestConnectionChannel <- virestConnection
			errorConnectChannel <- errorConnect
			isErrorConnectChannel <- isErrorConnect

			return
		}

		virestConnection, errorConnect, isErrorConnect = utils.NewConnectWithAuth(httpRequest.Header["Hypervisor-Uri"][0], nil, 0)

		virestConnectionChannel <- virestConnection
		errorConnectChannel <- errorConnect
		isErrorConnectChannel <- isErrorConnect
	}()

	errorPrepareRequestChannel := make(chan virest.Error)
	isErrorPrepareRequestChannel := make(chan bool)
	go func() {
		errorPrepareRequest, isErrorPrepareRequest := utils.CheckRequest(httpRequest, expectedRequestMethod, structure)

		errorPrepareRequestChannel <- errorPrepareRequest
		isErrorPrepareRequestChannel <- isErrorPrepareRequest
	}()

	virestConnection := <-virestConnectionChannel
	errorConnect := <-errorConnectChannel
	isErrorConnect := <-isErrorConnectChannel
	if isErrorConnect {
		return poolConnection{}, errorConnect, isErrorConnect
	}

	errorPrepareRequest := <-errorPrepareRequestChannel
	isErrorPrepareRequest := <-isErrorPrepareRequestChannel
	if isErrorPrepareRequest {
		virestConnection.Close()
		return poolConnection{}, errorPrepareRequest, isErrorPrepareRequest
	}

	return poolConnection{virestConnection.Connect}, virest.Error{}, false
}
