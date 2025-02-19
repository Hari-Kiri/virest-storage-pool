package storagePool

import (
	"fmt"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolEvent"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

// Init pool event probe and wait for result. Probing will be done until selected storage
// pool events type occur. Please registering default event implementation using
//
//	libvirt.EventRegisterDefaultImpl()
//
// then run iteration of the event loop using
//
//	libvirt.EventRunDefaultImpl()
//
// inside goroutine in main package before initiate pool event probe using this function.
//
// types:
//   - 0 = lifecycle
//   - 1 = refresh
func PoolEvent(connection virest.Connection, poolUuid string, types uint) (poolEvent.Event, virest.Error, bool) {
	var (
		result                                        poolEvent.Event
		storagePoolObject                             *libvirt.StoragePool
		callbackId                                    int
		errorGetStoragePoolObject, errorGetCallbackId error
		virestError                                   virest.Error
		isError                                       bool
	)

	result.EventRefresh = 0
	result.EventLifecycle = libvirt.StoragePoolEventLifecycle{
		Event: 6,
	}

	if types < 0 {
		virestError.Error = libvirt.Error{
			Code:    libvirt.ERR_STORAGE_PROBE_FAILED,
			Domain:  libvirt.FROM_EVENT,
			Message: fmt.Sprintf("no event type: %d", types),
			Level:   libvirt.ERR_ERROR,
		}
		isError = true
		return result, virestError, isError
	}

	if types > 1 {
		virestError.Error = libvirt.Error{
			Code:    libvirt.ERR_STORAGE_PROBE_FAILED,
			Domain:  libvirt.FROM_EVENT,
			Message: fmt.Sprintf("no event type: %d", types),
			Level:   libvirt.ERR_ERROR,
		}
		isError = true
		return result, virestError, isError
	}

	storagePoolObject, errorGetStoragePoolObject = connection.LookupStoragePoolByUUIDString(poolUuid)
	virestError.Error, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed get storage pool object: %s", virestError.Message)
		return result, virestError, isError
	}
	defer storagePoolObject.Free()

	if types == 0 {
		var storagePoolEventLifecycleCallbackResult = make(chan libvirt.StoragePoolEventLifecycle)

		callbackId, errorGetCallbackId = connection.StoragePoolEventLifecycleRegister(storagePoolObject, func(
			c *libvirt.Connect,
			n *libvirt.StoragePool,
			event *libvirt.StoragePoolEventLifecycle,
		) {
			storagePoolEventLifecycleCallbackResult <- *event
		})
		defer StoragePoolEventDeregister(connection, callbackId)

		result.EventLifecycle = <-storagePoolEventLifecycleCallbackResult
		virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
	}

	if types == 1 {
		var storagePoolEventGenericCallbackResult = make(chan int)

		callbackId, errorGetCallbackId = connection.StoragePoolEventRefreshRegister(storagePoolObject, func(
			c *libvirt.Connect,
			n *libvirt.StoragePool,
		) {
			storagePoolEventGenericCallbackResult <- 1
		})
		defer StoragePoolEventDeregister(connection, callbackId)

		result.EventRefresh = <-storagePoolEventGenericCallbackResult
		virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
	}

	return result, virestError, isError
}

func StoragePoolEventDeregister(connection virest.Connection, callbackId int) {
	var (
		virestError virest.Error
		isError     bool
	)

	virestError.Error, isError = connection.StoragePoolEventDeregister(callbackId).(libvirt.Error)
	if isError {
		temboLog.ErrorLogging("failed to deregister:", virestError.Message)
		return
	}

	temboLog.InfoLogging("deregister storage pool event callback with id:", callbackId)
}
