package storagePool

import (
	"fmt"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolEvent"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

// Get the probing result. Run PoolEvent() before, to init event probe, or it will always
// return invalid event.
var PoolEventProbingResult poolEvent.Event

// Init pool event probe and get the result using 'PoolEventProbingResult' variable.
// Probing will be done until selected storage pool events (using types to define it)
// occur.
//
// types:
//   - 0 = lifecycle
//   - 1 = refresh
func PoolEvent(connection virest.Connection, poolUuid string, types uint) (virest.Error, bool) {
	var (
		storagePoolObject                             *libvirt.StoragePool
		callbackId                                    int
		errorGetStoragePoolObject, errorGetCallbackId error
		virestError                                   virest.Error
		isError                                       bool
	)

	eventRunDefaultImpl := true
	go func() {
		for eventRunDefaultImpl {
			errorEventRunDefaultImpl := libvirt.EventRunDefaultImpl()
			if errorEventRunDefaultImpl != nil {
				virestError.Error = libvirt.Error{
					Code:    libvirt.ERR_STORAGE_PROBE_FAILED,
					Domain:  libvirt.FROM_EVENT,
					Message: fmt.Sprintf("failed start EventRunDefaultImpl(): %s", errorEventRunDefaultImpl),
					Level:   libvirt.ERR_ERROR,
				}
				isError = true
				break
			}
		}
	}()
	if isError {
		return virestError, isError
	}

	if types < 0 {
		virestError.Error = libvirt.Error{
			Code:    libvirt.ERR_STORAGE_PROBE_FAILED,
			Domain:  libvirt.FROM_EVENT,
			Message: fmt.Sprintf("no event type: %d", types),
			Level:   libvirt.ERR_ERROR,
		}
		isError = true
		return virestError, isError
	}
	if types > 1 {
		virestError.Error = libvirt.Error{
			Code:    libvirt.ERR_STORAGE_PROBE_FAILED,
			Domain:  libvirt.FROM_EVENT,
			Message: fmt.Sprintf("no event type: %d", types),
			Level:   libvirt.ERR_ERROR,
		}
		isError = true
		return virestError, isError
	}

	storagePoolObject, errorGetStoragePoolObject = connection.LookupStoragePoolByUUIDString(poolUuid)
	virestError.Error, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed get storage pool object: %s", virestError.Message)
		return virestError, isError
	}
	defer storagePoolObject.Free()

	PoolEventProbingResult.EventRefresh = 0
	PoolEventProbingResult.EventLifecycle = libvirt.StoragePoolEventLifecycle{
		Event: 6,
	}

	if types == 0 {
		callbackId, errorGetCallbackId = connection.StoragePoolEventLifecycleRegister(storagePoolObject, func(
			connection *libvirt.Connect,
			storagePoolObject *libvirt.StoragePool,
			storagePoolEvent *libvirt.StoragePoolEventLifecycle,
		) {
			PoolEventProbingResult.EventLifecycle = *storagePoolEvent
			eventRunDefaultImpl = false
		})

		virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
		return virestError, isError
	}
	if types == 1 {
		callbackId, errorGetCallbackId = connection.StoragePoolEventRefreshRegister(storagePoolObject, func(
			connection *libvirt.Connect,
			storagePoolObject *libvirt.StoragePool,
		) {
			PoolEventProbingResult.EventRefresh = 1
			eventRunDefaultImpl = false
		})

		virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
		return virestError, isError
	}

	if !eventRunDefaultImpl {
		virestError.Error, isError = connection.StoragePoolEventDeregister(callbackId).(libvirt.Error)
	}
	if isError {
		return virestError, isError
	}

	return virestError, isError
}