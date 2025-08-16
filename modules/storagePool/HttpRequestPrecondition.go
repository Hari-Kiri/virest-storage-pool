package storagePool

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/Hari-Kiri/temboLog"
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

	var (
		virestConnection                      virest.Connection
		waitGroup                             sync.WaitGroup
		errorConnect, errorPrepareRequest     virest.Error
		isErrorConnect, isErrorPrepareRequest bool
	)
	waitGroup.Add(2)
	go func() {
		defer waitGroup.Done()

		if len(httpRequest.Header["Hypervisor-Uri"]) == 0 {
			isErrorConnect = true
			errorConnect.Code = libvirt.ERR_INVALID_CONN
			errorConnect.Domain = libvirt.FROM_NET
			errorConnect.Message = "hypervisor uri not exist on request header"
			errorConnect.Level = libvirt.ERR_ERROR
			temboLog.ErrorLogging(
				"failed connect to hypervisor [ "+httpRequest.URL.Path+" ], requested from "+httpRequest.RemoteAddr+":",
				errorConnect.Message,
			)
			return
		}

		virestConnection, errorConnect, isErrorConnect = utils.NewConnectWithAuth(httpRequest.Header["Hypervisor-Uri"][0], nil, 0)
		if isErrorConnect {
			temboLog.ErrorLogging(
				"failed connect to hypervisor [ "+httpRequest.URL.Path+" ], requested from "+httpRequest.RemoteAddr+":",
				errorConnect.Message,
			)
		}
	}()
	go func() {
		defer waitGroup.Done()

		errorPrepareRequest, isErrorPrepareRequest = utils.CheckRequest(httpRequest, expectedRequestMethod, structure)
		if isErrorPrepareRequest {
			temboLog.ErrorLogging(
				"failed preparing request [ "+httpRequest.URL.Path+" ], requested from "+httpRequest.RemoteAddr+":",
				errorPrepareRequest.Message,
			)
		}
	}()
	waitGroup.Wait()

	if isErrorConnect {
		return poolConnection{}, errorConnect, isErrorConnect
	}
	if isErrorPrepareRequest {
		virestConnection.Close()
		return poolConnection{}, errorPrepareRequest, isErrorPrepareRequest
	}

	return poolConnection{virestConnection.Connect}, virest.Error{}, false
}
