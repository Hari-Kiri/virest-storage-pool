package utils

import (
	"net/http"
	"sync"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDefine"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolUndefine"
	"libvirt.org/go/libvirt"
)

// Defined generic type constraint for request model structure.
type requestStructure interface {
	poolDefine.Request | poolUndefine.Request
}

// Connect to local qemu socket and check request.
func RequestPrecondition[RequestStructure requestStructure](httpRequest *http.Request, expectedRequestMethod string, structure *RequestStructure) (*libvirt.Connect, libvirt.Error, bool) {
	var (
		result       *libvirt.Connect
		waitGroup    sync.WaitGroup
		libvirtError libvirt.Error
		isError      bool
	)

	waitGroup.Add(2)
	go func() {
		// Connect to qemu hypervisor
		if !isError {
			result, libvirtError, isError = NewConnect("qemu:///system")
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
			libvirtError, isError = CheckRequest(httpRequest, expectedRequestMethod, &structure)
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
