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

// Init pool event probe and wait for result. This method will blocking until the probing task is completed. Probing will be done until selected storage
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
		defer storagePoolEventDeregister(connection, callbackId)

		result.Timestamp = time.Now().Unix()
		result.TimestampNano = time.Now().UnixNano()
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
		defer storagePoolEventDeregister(connection, callbackId)

		result.Timestamp = time.Now().Unix()
		result.TimestampNano = time.Now().UnixNano()
		result.EventRefresh = <-storagePoolEventGenericCallbackResult
		virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
	}

	return result, virestError, isError
}

// Init pool event probe with timeout or loop until interupt. This method will blocking and create server sent event protocol
// until the probing task is completed. Unlike PoolEvent(), probing task will be completed when it has timeout or interupt.
// Please registering default event implementation in main package.
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
//   - 0 = loop until interrupt
//   - more than 0 = timeout seconds
func PoolEventTimeout(connection virest.Connection, poolUuid string, httpResponseWriter http.ResponseWriter, httpRequest *http.Request, types uint, timeout int) {
	var (
		result                    poolEvent.Event
		storagePoolObject         *libvirt.StoragePool
		errorGetStoragePoolObject error
		virestError               virest.Error
		isError                   bool
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
	}
	if types > 1 {
		virestError.Error = libvirt.Error{
			Code:    libvirt.ERR_STORAGE_PROBE_FAILED,
			Domain:  libvirt.FROM_EVENT,
			Message: fmt.Sprintf("no event type: %d", types),
			Level:   libvirt.ERR_ERROR,
		}
		isError = true
	}
	if timeout < 0 {
		virestError.Error = libvirt.Error{
			Code:    libvirt.ERR_STORAGE_PROBE_FAILED,
			Domain:  libvirt.FROM_EVENT,
			Message: fmt.Sprintf("timeout must be positive integer: %d", timeout),
			Level:   libvirt.ERR_ERROR,
		}
		isError = true
	}
	if isError {
		temboLog.ErrorLogging("failed probing event:", virestError.Message)
		return
	}

	storagePoolObject, errorGetStoragePoolObject = connection.LookupStoragePoolByUUIDString(poolUuid)
	virestError.Error, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		temboLog.ErrorLogging("failed get storage pool object:", virestError.Message)
		return
	}
	defer storagePoolObject.Free()

	if types == 0 && timeout == 0 {
		poolEventLifecycleLoop(httpResponseWriter, httpRequest, connection, storagePoolObject, &result)
		return
	}

	if types == 1 && timeout == 0 {
		poolEventRefreshLoop(httpResponseWriter, httpRequest, connection, storagePoolObject, &result)
		return
	}

	if types == 0 && timeout >= 1 {
		poolEventLifecycleTimeout(httpResponseWriter, httpRequest, connection, storagePoolObject, &result, timeout)
		return
	}

	if types == 1 && timeout >= 1 {
		poolEventRefreshTimeout(httpResponseWriter, httpRequest, connection, storagePoolObject, &result, timeout)
		return
	}
}

func poolEventLifecycleLoop(httpResponseWriter http.ResponseWriter, httpRequest *http.Request, connection virest.Connection, storagePoolObject *libvirt.StoragePool, eventStructure *poolEvent.Event) {
	var (
		callbackId         int
		errorGetCallbackId error
		virestError        virest.Error
		isError            bool
	)

	httpResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	httpResponseWriter.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	httpResponseWriter.Header().Set("Content-Type", "text/event-stream")
	httpResponseWriter.Header().Set("Cache-Control", "no-cache")
	httpResponseWriter.Header().Set("Connection", "keep-alive")

	writeEventStreamLifecycle(httpResponseWriter, eventStructure, virestError, nil)

	httpConnection := httpRequest.Context()
	usedCallbackId := make(chan int)
	callbackId, errorGetCallbackId = connection.StoragePoolEventLifecycleRegister(storagePoolObject, func(
		c *libvirt.Connect,
		n *libvirt.StoragePool,
		event *libvirt.StoragePoolEventLifecycle,
	) {
		select {
		case <-httpConnection.Done():
			usedCallbackId <- callbackId
		default:
			writeEventStreamLifecycle(httpResponseWriter, eventStructure, virestError, event)
		}
	})
	virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
	if isError {
		temboLog.ErrorLogging("failed to probing pool event:", errorGetCallbackId)
		return
	}

	storagePoolEventDeregister(connection, <-usedCallbackId)
}

func poolEventRefreshLoop(httpResponseWriter http.ResponseWriter, httpRequest *http.Request, connection virest.Connection, storagePoolObject *libvirt.StoragePool, eventStructure *poolEvent.Event) {
	var (
		callbackId         int
		errorGetCallbackId error
		virestError        virest.Error
		isError            bool
	)

	httpResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	httpResponseWriter.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	httpResponseWriter.Header().Set("Content-Type", "text/event-stream")
	httpResponseWriter.Header().Set("Cache-Control", "no-cache")
	httpResponseWriter.Header().Set("Connection", "keep-alive")

	writeEventStreamRefresh(httpResponseWriter, eventStructure, virestError, 0)

	httpConnection := httpRequest.Context()
	usedCallbackId := make(chan int)
	callbackId, errorGetCallbackId = connection.StoragePoolEventRefreshRegister(storagePoolObject, func(
		c *libvirt.Connect,
		n *libvirt.StoragePool,
	) {
		select {
		case <-httpConnection.Done():
			usedCallbackId <- callbackId
		default:
			writeEventStreamRefresh(httpResponseWriter, eventStructure, virestError, 1)
		}
	})
	virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
	if isError {
		temboLog.ErrorLogging("failed to probing pool event:", errorGetCallbackId)
		return
	}

	storagePoolEventDeregister(connection, <-usedCallbackId)
}

func poolEventLifecycleTimeout(httpResponseWriter http.ResponseWriter, httpRequest *http.Request, connection virest.Connection, storagePoolObject *libvirt.StoragePool, eventStructure *poolEvent.Event, timeout int) {
	var (
		callbackId         int
		errorGetCallbackId error
		virestError        virest.Error
		isError            bool
	)

	httpResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	httpResponseWriter.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	httpResponseWriter.Header().Set("Content-Type", "text/event-stream")
	httpResponseWriter.Header().Set("Cache-Control", "no-cache")
	httpResponseWriter.Header().Set("Connection", "keep-alive")

	writeEventStreamLifecycle(httpResponseWriter, eventStructure, virestError, nil)

	usedCallbackId := make(chan int)
	addTimeout, errorAddTimeout := libvirt.EventAddTimeout(timeout*1000, func(timer int) {
		usedCallbackId <- callbackId
		writeEventStreamEnd(httpResponseWriter)
	})
	virestError.Error, isError = errorAddTimeout.(libvirt.Error)
	if isError {
		temboLog.ErrorLogging("failed to probing pool event:", virestError.Message)
		return
	}

	httpConnection := httpRequest.Context()
	callbackId, errorGetCallbackId = connection.StoragePoolEventLifecycleRegister(storagePoolObject, func(
		c *libvirt.Connect,
		n *libvirt.StoragePool,
		event *libvirt.StoragePoolEventLifecycle,
	) {
		select {
		case <-httpConnection.Done():
			usedCallbackId <- callbackId
		default:
			writeEventStreamLifecycle(httpResponseWriter, eventStructure, virestError, event)
		}
	})
	virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
	if isError {
		temboLog.ErrorLogging("failed to probing pool event:", virestError.Message)
		return
	}
	storagePoolEventDeregister(connection, <-usedCallbackId)

	errorEventRemoveTimeout := libvirt.EventRemoveTimeout(addTimeout)
	virestError.Error, isError = errorEventRemoveTimeout.(libvirt.Error)
	if isError {
		temboLog.ErrorLogging("failed to remove event timeout callback", virestError.Message)
		return
	}
}

func poolEventRefreshTimeout(httpResponseWriter http.ResponseWriter, httpRequest *http.Request, connection virest.Connection, storagePoolObject *libvirt.StoragePool, eventStructure *poolEvent.Event, timeout int) {
	var (
		callbackId         int
		errorGetCallbackId error
		virestError        virest.Error
		isError            bool
	)

	httpResponseWriter.Header().Set("Access-Control-Allow-Origin", "*")
	httpResponseWriter.Header().Set("Access-Control-Expose-Headers", "Content-Type")
	httpResponseWriter.Header().Set("Content-Type", "text/event-stream")
	httpResponseWriter.Header().Set("Cache-Control", "no-cache")
	httpResponseWriter.Header().Set("Connection", "keep-alive")

	writeEventStreamRefresh(httpResponseWriter, eventStructure, virestError, 0)

	usedCallbackId := make(chan int)
	addTimeout, errorAddTimeout := libvirt.EventAddTimeout(timeout*1000, func(timer int) {
		usedCallbackId <- callbackId
		writeEventStreamEnd(httpResponseWriter)
	})
	virestError.Error, isError = errorAddTimeout.(libvirt.Error)
	if isError {
		temboLog.ErrorLogging("failed to add event timeout", virestError.Message)
		return
	}

	httpConnection := httpRequest.Context()
	callbackId, errorGetCallbackId = connection.StoragePoolEventRefreshRegister(storagePoolObject, func(
		c *libvirt.Connect,
		n *libvirt.StoragePool,
	) {
		select {
		case <-httpConnection.Done():
			usedCallbackId <- callbackId
		default:
			writeEventStreamRefresh(httpResponseWriter, eventStructure, virestError, 1)
		}
	})
	virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
	if isError {
		temboLog.ErrorLogging("failed to probing pool event:", virestError.Message)
		return
	}
	storagePoolEventDeregister(connection, <-usedCallbackId)

	errorEventRemoveTimeout := libvirt.EventRemoveTimeout(addTimeout)
	virestError.Error, isError = errorEventRemoveTimeout.(libvirt.Error)
	if isError {
		temboLog.ErrorLogging("failed to remove event timeout callback", virestError.Message)
		return
	}
}

func storagePoolEventDeregister(connection virest.Connection, callbackId int) {
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

func writeEventStreamLifecycle(httpResponseWriter http.ResponseWriter, result *poolEvent.Event, virestError virest.Error, event *libvirt.StoragePoolEventLifecycle) {
	result.Timestamp = time.Now().Unix()
	result.TimestampNano = time.Now().UnixNano()

	if event != nil {
		result.EventLifecycle = *event
	}

	var httpBody poolEvent.Response
	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data = *result
	httpBody.Error = virestError

	var responseBuffer bytes.Buffer
	errorEncodeToJson := json.NewEncoder(&responseBuffer).Encode(&httpBody)
	if errorEncodeToJson != nil {
		temboLog.ErrorLogging("failed to write storage pool event lifecycle stream:", errorEncodeToJson)
		return
	}

	response := []byte("event: lifecycle\n")
	response = append(response, append([]byte("data: "), responseBuffer.Bytes()...)...)
	response = append(response, []byte("\n\n")...)
	httpResponseWriter.Write(response)
	httpResponseWriter.(http.Flusher).Flush()
}

func writeEventStreamRefresh(httpResponseWriter http.ResponseWriter, result *poolEvent.Event, virestError virest.Error, flag int) {
	result.Timestamp = time.Now().Unix()
	result.TimestampNano = time.Now().UnixNano()

	if flag != 0 {
		result.EventRefresh = flag
	}

	var httpBody poolEvent.Response
	httpBody.Response = true
	httpBody.Code = http.StatusOK
	httpBody.Data = *result
	httpBody.Error = virestError

	var responseBuffer bytes.Buffer
	errorEncodeToJson := json.NewEncoder(&responseBuffer).Encode(&httpBody)
	if errorEncodeToJson != nil {
		temboLog.ErrorLogging("failed to write storage pool event refresh stream:", errorEncodeToJson)
		return
	}

	response := []byte("event: refresh\n")
	response = append(response, append([]byte("data: "), responseBuffer.Bytes()...)...)
	response = append(response, []byte("\n\n")...)
	httpResponseWriter.Write(response)
	httpResponseWriter.(http.Flusher).Flush()
}

func writeEventStreamEnd(httpResponseWriter http.ResponseWriter) {
	response := []byte("event: end\n")
	response = append(response, []byte("data: null")...)
	response = append(response, []byte("\n\n")...)
	httpResponseWriter.Write(response)
	httpResponseWriter.(http.Flusher).Flush()
}
