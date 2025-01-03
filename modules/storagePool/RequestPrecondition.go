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

// Validate user request using given bearer token which is generated JWT by ViRest Utilities 'auth.BasicAuth()' module. Look up hypervisor
// URI on 'Hypervisor-Uri' request header field. Connect to hypervisor via SSH tunnel and check the expected HTTP request method then convert
// the JSON request body to structure if any. SSH tunnel work with Key-Based authentication. Please, create SSH Key on the host and copy it
// on the remote libvirt-daemon host
//
//	~/.ssh/authorized_keys
//
// Notes for HTTP GET method:
//
// - Query parameter and structure field will be compared in case sensitive.
//
// - Every structure field data type must be string, so You must convert it to the right data type before You use it.
//
// - Untested for array query argument.
//
// Notes for HTTP POST, PUT, PATCH and DELETE method:
//
// - This function always looking for request body for data and parse them to 'structure' parameter.
func RequestPrecondition[RequestStructure utils.RequestStructure](
	httpRequest *http.Request,
	expectedRequestMethod string,
	structure *RequestStructure,
	applicationName string,
	jwtSigningMethod *jwt.SigningMethodHMAC,
	jwtSignatureKey []byte,
) (virest.Connection, virest.Error, bool) {
	libvirtErrorAuth, isErrorAuth := auth.BearerTokenAuth(
		httpRequest,
		applicationName,
		jwtSigningMethod,
		jwtSignatureKey,
	)
	if isErrorAuth {
		return virest.Connection{}, virest.Error{Error: libvirt.Error{
			Code:    libvirt.ERR_AUTH_FAILED,
			Domain:  libvirt.FROM_NET,
			Message: fmt.Sprintf("authentication failed: %s", libvirtErrorAuth.Message),
			Level:   libvirt.ERR_ERROR,
		}}, true
	}

	if len(httpRequest.Header["Hypervisor-Uri"]) == 0 {
		return virest.Connection{}, virest.Error{Error: libvirt.Error{
			Code:    libvirt.ERR_INVALID_CONN,
			Domain:  libvirt.FROM_NET,
			Message: "hypervisor uri not exist on request header",
			Level:   libvirt.ERR_ERROR,
		}}, true
	}

	var (
		result                                virest.Connection
		waitGroup                             sync.WaitGroup
		errorConnect, errorPrepareRequest     virest.Error
		isErrorConnect, isErrorPrepareRequest bool
	)
	waitGroup.Add(2)
	go func() {
		result, errorConnect, isErrorConnect = utils.NewConnectWithAuth(httpRequest.Header["Hypervisor-Uri"][0], nil, 0)
		if isErrorConnect {
			temboLog.ErrorLogging(
				"failed connect to hypervisor [ "+httpRequest.URL.Path+" ], requested from "+httpRequest.RemoteAddr+":",
				errorConnect.Message,
			)
		}
		defer waitGroup.Done()
	}()
	go func() {
		// Prepare request
		errorPrepareRequest, isErrorPrepareRequest = utils.CheckRequest(httpRequest, expectedRequestMethod, structure)
		if isErrorPrepareRequest {
			temboLog.ErrorLogging(
				"failed preparing request [ "+httpRequest.URL.Path+" ], requested from "+httpRequest.RemoteAddr+":",
				errorPrepareRequest.Message,
			)
		}
		defer waitGroup.Done()
	}()
	waitGroup.Wait()

	if isErrorConnect {
		return virest.Connection{}, errorConnect, isErrorConnect
	}
	if isErrorPrepareRequest {
		result.Close()
		return virest.Connection{}, errorPrepareRequest, isErrorPrepareRequest
	}

	return result, virest.Error{}, false
}
