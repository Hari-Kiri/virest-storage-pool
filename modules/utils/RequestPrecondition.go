package utils

import (
	"net/http"
	"sync"

	goVirtQemuConnector "github.com/Hari-Kiri/govirt-qemu-connector"
	"github.com/Hari-Kiri/temboLog"
	"libvirt.org/go/libvirt"
)

// Connect to local qemu socket and check request.
func RequestPrecondition(httpRequest *http.Request, expectedRequestMethod string, structure any) (*libvirt.Connect, libvirt.Error, bool) {
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
			var errorConnectToQemu error
			result, errorConnectToQemu = goVirtQemuConnector.ConnectToLocalSystem()
			libvirtError, isError = errorConnectToQemu.(libvirt.Error)
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
