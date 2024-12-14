package storagePool

import (
	"net/http"
	"sync"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"libvirt.org/go/libvirt"
)

// Connect to qemu hypervisor via SSH tunnel and check the expected HTTP request method and convert the JSON request body to structure if any.
// SSH tunnel work with Key-Based authentication. Please, create SSH Key on the host and copy it on the remote libvirt-daemon host
// ~/.ssh/authorized_keys.
//
// Notes for HTTP GET method:
//
// - Query parameter and structure field will be compared in case sensitive.
//
// - Every structure field data type must be string, so You must convert it to the right data type before You use it.
//
// - Untested for array query argument.
func RequestPrecondition[RequestStructure utils.RequestStructure](httpRequest *http.Request, expectedRequestMethod string, connectionUri string, structure *RequestStructure) (*libvirt.Connect, libvirt.Error, bool) {
	var (
		result                                          *libvirt.Connect
		waitGroup                                       sync.WaitGroup
		libvirtErrorConnect, libvirtErrorPrepareRequest libvirt.Error
		isErrorConnect, isErrorPrepareRequest           bool
	)

	waitGroup.Add(2)
	go func() {
		result, libvirtErrorConnect, isErrorConnect = utils.NewConnectWithAuth(connectionUri, nil, 0)
		if isErrorConnect {
			temboLog.ErrorLogging(
				"failed connect to hypervisor [ "+httpRequest.URL.Path+" ], requested from "+httpRequest.RemoteAddr+":",
				libvirtErrorConnect.Message,
			)
		}
		defer waitGroup.Done()
	}()
	go func() {
		// Prepare request
		libvirtErrorPrepareRequest, isErrorPrepareRequest = utils.CheckRequest(httpRequest, expectedRequestMethod, structure)
		if isErrorPrepareRequest {
			temboLog.ErrorLogging(
				"failed preparing request [ "+httpRequest.URL.Path+" ], requested from "+httpRequest.RemoteAddr+":",
				libvirtErrorPrepareRequest.Message,
			)
		}
		defer waitGroup.Done()
	}()
	waitGroup.Wait()

	if isErrorConnect {
		return nil, libvirtErrorConnect, isErrorConnect
	}
	if isErrorPrepareRequest {
		result.Close()
		return nil, libvirtErrorPrepareRequest, isErrorPrepareRequest
	}

	return result, libvirt.Error{}, false
}
