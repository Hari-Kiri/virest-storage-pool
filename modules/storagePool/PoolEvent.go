package storagePool

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolEvent"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

// Init pool event probe and wait for result. Probing will be done until selected storage
// pool events type occur. Please registering default event implementation in main package.
//
//	errorEventRegisterDefaultImpl := libvirt.EventRegisterDefaultImpl()
//	if errorEventRegisterDefaultImpl != nil {
//		temboLog.FatalLogging("failed registers a default event implementation based on the poll() system call:", errorEventRegisterDefaultImpl)
//	}
//
// then run iteration of the event loop inside goroutine in main package before initiate pool
// event probe using this function.
//
//	go func() {
//		for {
//			errorEventRunDefaultImpl := libvirt.EventRunDefaultImpl()
//			if errorEventRunDefaultImpl != nil {
//				temboLog.FatalLogging("failed starting run one iteration of the event loop:", errorEventRunDefaultImpl)
//			}
//		}
//	}()
//
// types:
//   - 0 = storage pool event lifecycle
//   - 1 = storage pool event refresh
//
// timeout:
//   - -1 = wait until selected storage pool events type occur
//   - 0 = loop until timeout or interrupt, rather than one-shot
//   - more than 0 = timeout seconds
func PoolEvent(connection virest.Connection, poolUuid string, httpResponseWriter http.ResponseWriter, httpRequest *http.Request,
	types uint, timeout int) (poolEvent.Event, virest.Error, bool) {
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

	if types == 0 && timeout == -1 {
		var storagePoolEventLifecycleCallbackResult = make(chan libvirt.StoragePoolEventLifecycle)

		callbackId, errorGetCallbackId = connection.StoragePoolEventLifecycleRegister(storagePoolObject, func(
			c *libvirt.Connect,
			n *libvirt.StoragePool,
			event *libvirt.StoragePoolEventLifecycle,
		) {
			storagePoolEventLifecycleCallbackResult <- *event
		})
		defer StoragePoolEventDeregister(connection, callbackId)

		result.Timestamp = time.Now().Unix()
		result.TimestampNano = time.Now().UnixNano()
		result.EventLifecycle = <-storagePoolEventLifecycleCallbackResult
		virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
	}

	if types == 1 && timeout == -1 {
		var storagePoolEventGenericCallbackResult = make(chan int)

		callbackId, errorGetCallbackId = connection.StoragePoolEventRefreshRegister(storagePoolObject, func(
			c *libvirt.Connect,
			n *libvirt.StoragePool,
		) {
			storagePoolEventGenericCallbackResult <- 1
		})
		defer StoragePoolEventDeregister(connection, callbackId)

		result.Timestamp = time.Now().Unix()
		result.TimestampNano = time.Now().UnixNano()
		result.EventRefresh = <-storagePoolEventGenericCallbackResult
		virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
	}

	if types == 0 && timeout >= 1 {
		httpResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		httpResponseWriter.Header().Set("Access-Control-Expose-Headers", "Content-Type")
		httpResponseWriter.Header().Set("Content-Type", "text/event-stream")
		httpResponseWriter.Header().Set("Cache-Control", "no-cache")
		httpResponseWriter.Header().Set("Connection", "keep-alive")

		createEventStreamLifecycle(
			httpResponseWriter,
			&result,
			virestError,
			nil,
		)

		var (
			addTimeout      int
			errorAddTimeout error
		)
		addTimeout, errorAddTimeout = libvirt.EventAddTimeout(timeout*1000, func(timer int) {
			StoragePoolEventDeregister(connection, callbackId)
			createEventStreamEnd(httpResponseWriter)
			errorLog("failed to remove event timeout callback", libvirt.EventRemoveTimeout(addTimeout))
		})
		errorLog("failed to add event timeout", errorAddTimeout)

		httpConnection := httpRequest.Context()
		for {
			select {
			case <-httpConnection.Done():
				break
			default:
				var storagePoolEventLifecycleCallbackResult = make(chan libvirt.StoragePoolEventLifecycle)
				callbackId, errorGetCallbackId = connection.StoragePoolEventLifecycleRegister(storagePoolObject, func(
					c *libvirt.Connect,
					n *libvirt.StoragePool,
					event *libvirt.StoragePoolEventLifecycle,
				) {
					defer StoragePoolEventDeregister(connection, callbackId)
					storagePoolEventLifecycleCallbackResult <- *event
				})
				virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
				createEventStreamLifecycle(
					httpResponseWriter,
					&result,
					virestError,
					storagePoolEventLifecycleCallbackResult,
				)
			}
		}
	}

	if types == 1 && timeout >= 1 {
		httpResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
		httpResponseWriter.Header().Set("Access-Control-Expose-Headers", "Content-Type")
		httpResponseWriter.Header().Set("Content-Type", "text/event-stream")
		httpResponseWriter.Header().Set("Cache-Control", "no-cache")
		httpResponseWriter.Header().Set("Connection", "keep-alive")

		createEventStreamRefresh(
			httpResponseWriter,
			&result,
			virestError,
			nil,
		)

		var (
			addTimeout      int
			errorAddTimeout error
		)
		addTimeout, errorAddTimeout = libvirt.EventAddTimeout(timeout*1000, func(timer int) {
			StoragePoolEventDeregister(connection, callbackId)
			createEventStreamEnd(httpResponseWriter)
			errorLog("failed to remove event timeout callback", libvirt.EventRemoveTimeout(addTimeout))
		})
		errorLog("failed to add event timeout", errorAddTimeout)

		httpConnection := httpRequest.Context()
		for {
			select {
			case <-httpConnection.Done():
				break
			default:
				var storagePoolEventGenericCallbackResult = make(chan int)
				callbackId, errorGetCallbackId = connection.StoragePoolEventRefreshRegister(storagePoolObject, func(
					c *libvirt.Connect,
					n *libvirt.StoragePool,
				) {
					defer StoragePoolEventDeregister(connection, callbackId)
					storagePoolEventGenericCallbackResult <- 1
				})
				virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
				createEventStreamRefresh(
					httpResponseWriter,
					&result,
					virestError,
					storagePoolEventGenericCallbackResult,
				)
			}
		}
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

func createEventStreamLifecycle(httpResponseWriter http.ResponseWriter, result *poolEvent.Event, virestError virest.Error,
	storagePoolEventLifecycleCallbackResult chan libvirt.StoragePoolEventLifecycle) {
	result.Timestamp = time.Now().Unix()
	result.TimestampNano = time.Now().UnixNano()

	if storagePoolEventLifecycleCallbackResult != nil {
		result.EventLifecycle = <-storagePoolEventLifecycleCallbackResult
	}

	var httpBody poolEvent.Response
	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data = *result
	httpBody.Error = virestError

	var responseBuffer bytes.Buffer
	json.NewEncoder(&responseBuffer).Encode(&httpBody)

	response := []byte("event: lifecycle\n")
	response = append(response, append([]byte("data: "), responseBuffer.Bytes()...)...)
	response = append(response, []byte("\n\n")...)
	httpResponseWriter.Write(response)
	httpResponseWriter.(http.Flusher).Flush()
}

func createEventStreamRefresh(httpResponseWriter http.ResponseWriter, result *poolEvent.Event, virestError virest.Error,
	storagePoolEventGenericCallbackResult chan int) {
	result.Timestamp = time.Now().Unix()
	result.TimestampNano = time.Now().UnixNano()

	if storagePoolEventGenericCallbackResult != nil {
		result.EventRefresh = <-storagePoolEventGenericCallbackResult
	}

	var httpBody poolEvent.Response
	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data = *result
	httpBody.Error = virestError

	var responseBuffer bytes.Buffer
	json.NewEncoder(&responseBuffer).Encode(&httpBody)

	response := []byte("event: refresh\n")
	response = append(response, append([]byte("data: "), responseBuffer.Bytes()...)...)
	response = append(response, []byte("\n\n")...)
	httpResponseWriter.Write(response)
	httpResponseWriter.(http.Flusher).Flush()
}

func createEventStreamEnd(httpResponseWriter http.ResponseWriter) {
	response := []byte("event: end\n")
	response = append(response, []byte("data: null")...)
	response = append(response, []byte("\n\n")...)
	httpResponseWriter.Write(response)
	httpResponseWriter.(http.Flusher).Flush()
}

func errorLog(message string, error any) {
	if error != nil {
		temboLog.ErrorLogging(message, error.(libvirt.Error))
	}
}
