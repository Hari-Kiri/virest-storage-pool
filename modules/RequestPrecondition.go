package modules

import (
	"net/http"
	"sync"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules/utils"
	"libvirt.org/go/libvirt"
)

// Connect to qemu hypervisor via SSH tunnel and check the expected HTTP request method and convert the JSON request body to structure if any.
// SSH tunnel work with Key-Based authentication. Please, create SSH Key on the host and copy it on the remote libvirt-daemon host
// ~/.ssh/authorized_keys.
func RequestPrecondition[RequestStructure utils.RequestStructure](httpRequest *http.Request, expectedRequestMethod string, connectionUri string, structure *RequestStructure) (*libvirt.Connect, libvirt.Error, bool) {
	var (
		result       *libvirt.Connect
		waitGroup    sync.WaitGroup
		libvirtError libvirt.Error
		isError      bool
	)

	waitGroup.Add(2)
	go func() {
		if !isError {
			result, libvirtError, isError = utils.NewConnectWithAuth(connectionUri, nil, 0)
		}
		if isError && libvirtError.Code != 11 {
			temboLog.ErrorLogging(
				"failed connect to hypervisor [ "+httpRequest.URL.Path+" ], requested from "+httpRequest.RemoteAddr+":",
				libvirtError.Message,
			)
		}
		defer waitGroup.Done()
	}()
	go func() {
		// Prepare request
		if !isError {
			libvirtError, isError = utils.CheckRequest(httpRequest, expectedRequestMethod, structure)
		}
		if isError && libvirtError.Code == 11 {
			temboLog.ErrorLogging(
				"failed preparing request [ "+httpRequest.URL.Path+" ], requested from "+httpRequest.RemoteAddr+":",
				libvirtError.Message,
			)
		}
		defer waitGroup.Done()
	}()
	waitGroup.Wait()

	if libvirtError.Code == 11 && result != nil {
		result.Close()
	}

	return result, libvirtError, isError
}
