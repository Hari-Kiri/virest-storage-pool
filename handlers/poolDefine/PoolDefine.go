package poolDefine

import (
	"net/http"
	"sync"

	goVirtQemuConnector "github.com/Hari-Kiri/govirt-qemu-connector"
	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/modules/utils"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDefine"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

func PoolDefine(responseWriter http.ResponseWriter, request *http.Request) {
	var (
		qemuConnection  *libvirt.Connect
		requestBodyData libvirtxml.StoragePool
		httpBody        poolDefine.Response
		waitGroup       sync.WaitGroup
		libvirtError    libvirt.Error
		isError         bool
	)

	waitGroup.Add(2)
	go func() {
		// Connect to qemu hypervisor
		if !isError {
			var errorConnectToQemu error
			qemuConnection, errorConnectToQemu = goVirtQemuConnector.ConnectToLocalSystem()
			libvirtError, isError = errorConnectToQemu.(libvirt.Error)
		}
		if isError && libvirtError.Code != 11 {
			temboLog.ErrorLogging(
				"failed connect to hypervisor [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
				libvirtError.Message,
			)
		}
		defer waitGroup.Done()
	}()
	go func() {
		// Prepare request
		if !isError {
			libvirtError, isError = utils.PrepareRequest(request, http.MethodPost, &requestBodyData)
		}
		if isError && libvirtError.Code == 11 {
			temboLog.ErrorLogging(
				"failed preparing request [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
				libvirtError.Message,
			)
		}
		defer waitGroup.Done()
	}()
	waitGroup.Wait()
	// ConnectToLocalSystem().close() should be used to release the resources after the connection is no longer needed
	defer qemuConnection.Close()
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirtError.Code)
		httpBody.Error = libvirtError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		return
	}

	// Convert request body to libvirt xml
	libvirtXml, errorGetLibvirtXml := requestBodyData.Marshal()
	libvirtError, isError = errorGetLibvirtXml.(libvirt.Error)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirtError.Code)
		httpBody.Error = libvirtError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed to create pool config xml [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			libvirtError.Message,
		)
		return
	}

	// Define pool
	definePool, errorDefinePool := qemuConnection.StoragePoolDefineXML(libvirtXml, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	libvirtError, isError = errorDefinePool.(libvirt.Error)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirtError.Code)
		httpBody.Error = libvirtError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed to define new pool [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			libvirtError.Message,
		)
		return
	}

	// Get defined pool UUID
	definedPoolUuid, errorGetDefinedPoolUuid := definePool.GetUUIDString()
	libvirtError, isError = errorGetDefinedPoolUuid.(libvirt.Error)
	if isError {
		httpBody.Response = false
		httpBody.Code = utils.HttpErrorCode(libvirtError.Code)
		httpBody.Error = libvirtError
		utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
		temboLog.ErrorLogging(
			"failed to get defined pool UUID [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
			libvirtError.Message,
		)
		return
	}

	// Http ok response
	httpBody.Response = true
	httpBody.Code = http.StatusCreated
	httpBody.Message.Uuid = definedPoolUuid
	utils.JsonResponseBuilder(httpBody, responseWriter, httpBody.Code)
	temboLog.InfoLogging("new pool defined with uuid:", definedPoolUuid, "[", request.URL.Path, "]")

	go func() {
		// Free storage pool object
		libvirtError, isError = definePool.Free().(libvirt.Error)
		if isError {
			temboLog.ErrorLogging(
				"failed to free storage pool object [ "+request.URL.Path+" ], requested from "+request.RemoteAddr+":",
				libvirtError.Message,
			)
			return
		}
	}()
}
