package storagePool

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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
func (poolConnection *poolConnection) PoolEvent(poolUuid string, types uint) (poolEvent.Event, virest.Error, bool) {
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

	storagePoolObject, errorGetStoragePoolObject = poolConnection.LookupStoragePoolByUUIDString(poolUuid)
	virestError.Error, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed get storage pool object: %s", virestError.Message)
		return result, virestError, isError
	}
	defer storagePoolObject.Free()

	if types == 0 {
		var storagePoolEventLifecycleCallbackResult = make(chan libvirt.StoragePoolEventLifecycle)

		callbackId, errorGetCallbackId = poolConnection.StoragePoolEventLifecycleRegister(storagePoolObject, func(
			c *libvirt.Connect,
			n *libvirt.StoragePool,
			event *libvirt.StoragePoolEventLifecycle,
		) {
			storagePoolEventLifecycleCallbackResult <- *event
		})
		defer poolConnection.storagePoolEventDeregister(callbackId)

		result.Timestamp = time.Now().Unix()
		result.TimestampNano = time.Now().UnixNano()
		result.EventLifecycle = <-storagePoolEventLifecycleCallbackResult
		virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
	}

	if types == 1 {
		var storagePoolEventGenericCallbackResult = make(chan int)

		callbackId, errorGetCallbackId = poolConnection.StoragePoolEventRefreshRegister(storagePoolObject, func(
			c *libvirt.Connect,
			n *libvirt.StoragePool,
		) {
			storagePoolEventGenericCallbackResult <- 1
		})
		defer poolConnection.storagePoolEventDeregister(callbackId)

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
func (poolConnection *poolConnection) PoolEventTimeout(poolUuid string, httpResponseWriter http.ResponseWriter, types uint, timeout int) (virest.Error, bool) {
	var (
		eventStructure            poolEvent.Event
		storagePoolObject         *libvirt.StoragePool
		errorGetStoragePoolObject error
		virestError               virest.Error
		isError                   bool
	)

	eventStructure.EventRefresh = 0
	eventStructure.EventLifecycle = libvirt.StoragePoolEventLifecycle{
		Event: 6,
	}
	eventStructure.EventId = -1

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
		return virestError, isError
	}

	storagePoolObject, errorGetStoragePoolObject = poolConnection.LookupStoragePoolByUUIDString(poolUuid)
	virestError.Error, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		return virestError, isError
	}
	defer storagePoolObject.Free()

	if types == 0 {
		virestError, isError = poolConnection.poolEventLifecycleTimeout(httpResponseWriter, storagePoolObject, &eventStructure, timeout)
	}

	if types == 1 {
		virestError, isError = poolConnection.poolEventRefreshTimeout(httpResponseWriter, storagePoolObject, &eventStructure, timeout)
	}

	return virestError, isError
}

func (poolConnection *poolConnection) poolEventLifecycleTimeout(httpResponseWriter http.ResponseWriter, storagePoolObject *libvirt.StoragePool, eventStructure *poolEvent.Event, timeout int) (virest.Error, bool) {
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
	var freq int
	if timeout == 0 {
		freq = -1
	}
	if timeout >= 1 {
		freq = timeout * 1000
	}
	addTimeout, errorAddTimeout := libvirt.EventAddTimeout(
		freq,
		func(timer int) {
			usedCallbackId <- callbackId
			writeEventStreamEnd(httpResponseWriter)
		},
	)
	virestError.Error, isError = errorAddTimeout.(libvirt.Error)
	if isError {
		return virestError, true
	}

	callbackId, errorGetCallbackId = poolConnection.StoragePoolEventLifecycleRegister(
		storagePoolObject,
		func(c *libvirt.Connect, n *libvirt.StoragePool, event *libvirt.StoragePoolEventLifecycle) {
			eventStructure.EventId = callbackId
			writeEventStreamLifecycle(httpResponseWriter, eventStructure, virestError, event)
		},
	)
	virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
	if isError {
		return virestError, true
	}
	poolConnection.storagePoolEventDeregister(<-usedCallbackId)

	errorEventRemoveTimeout := libvirt.EventRemoveTimeout(addTimeout)
	virestError.Error, isError = errorEventRemoveTimeout.(libvirt.Error)
	if isError {
		return virestError, true
	}

	return virestError, false
}

func (poolConnection *poolConnection) poolEventRefreshTimeout(httpResponseWriter http.ResponseWriter, storagePoolObject *libvirt.StoragePool, eventStructure *poolEvent.Event, timeout int) (virest.Error, bool) {
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
	var freq int
	if timeout == 0 {
		freq = -1
	}
	if timeout >= 1 {
		freq = timeout * 1000
	}
	addTimeout, errorAddTimeout := libvirt.EventAddTimeout(
		freq,
		func(timer int) {
			usedCallbackId <- callbackId
			writeEventStreamEnd(httpResponseWriter)
		},
	)
	virestError.Error, isError = errorAddTimeout.(libvirt.Error)
	if isError {
		return virestError, true
	}

	callbackId, errorGetCallbackId = poolConnection.StoragePoolEventRefreshRegister(
		storagePoolObject,
		func(c *libvirt.Connect, n *libvirt.StoragePool) {
			eventStructure.EventId = callbackId
			writeEventStreamRefresh(httpResponseWriter, eventStructure, virestError, 1)
		},
	)
	virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
	if isError {
		return virestError, true
	}
	poolConnection.storagePoolEventDeregister(<-usedCallbackId)

	errorEventRemoveTimeout := libvirt.EventRemoveTimeout(addTimeout)
	virestError.Error, isError = errorEventRemoveTimeout.(libvirt.Error)
	if isError {
		return virestError, true
	}

	return virestError, false
}

func (poolConnection *poolConnection) storagePoolEventDeregister(callbackId int) {
	var (
		virestError virest.Error
		isError     bool
	)

	virestError.Error, isError = poolConnection.StoragePoolEventDeregister(callbackId).(libvirt.Error)
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

func (poolConnection *poolConnection) PoolEventStream(poolUuid string, types uint, timeout int, pipeWriter *io.PipeWriter) (virest.Error, bool) {
	defer pipeWriter.Close()
	var (
		storagePoolObject                             *libvirt.StoragePool
		callbackId                                    int
		errorGetStoragePoolObject, errorGetCallbackId error
		virestError                                   virest.Error
		isError                                       bool
	)

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

	storagePoolObject, errorGetStoragePoolObject = poolConnection.LookupStoragePoolByUUIDString(poolUuid)
	virestError.Error, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed get storage pool object: %s", virestError.Message)
		return virestError, isError
	}
	defer storagePoolObject.Free()

	if types == 0 {
		usedCallbackId := make(chan int)
		var freq int
		if timeout == 0 {
			freq = -1
		}
		if timeout >= 1 {
			freq = timeout * 1000
		}
		addTimeout, errorAddTimeout := libvirt.EventAddTimeout(
			freq,
			func(timer int) {
				usedCallbackId <- callbackId
			},
		)
		virestError.Error, isError = errorAddTimeout.(libvirt.Error)
		if isError {
			return virestError, true
		}

		callbackId, errorGetCallbackId = poolConnection.StoragePoolEventRefreshRegister(
			storagePoolObject,
			func(c *libvirt.Connect, n *libvirt.StoragePool) {
				bytesBuffer := new(bytes.Buffer)
				bytesBuffer.WriteString(fmt.Sprintf("callback id %d: restart\n", callbackId))
				bufio.NewWriter(pipeWriter)
			},
		)
		virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
		if isError {
			return virestError, true
		}
		poolConnection.storagePoolEventDeregister(<-usedCallbackId)

		errorEventRemoveTimeout := libvirt.EventRemoveTimeout(addTimeout)
		virestError.Error, isError = errorEventRemoveTimeout.(libvirt.Error)
		if isError {
			return virestError, true
		}
	}

	if types == 1 {
		usedCallbackId := make(chan int)
		var freq int
		if timeout == 0 {
			freq = -1
		}
		if timeout >= 1 {
			freq = timeout * 1000
		}
		addTimeout, errorAddTimeout := libvirt.EventAddTimeout(
			freq,
			func(timer int) {
				usedCallbackId <- callbackId
			},
		)
		virestError.Error, isError = errorAddTimeout.(libvirt.Error)
		if isError {
			return virestError, true
		}

		callbackId, errorGetCallbackId = poolConnection.StoragePoolEventRefreshRegister(
			storagePoolObject,
			func(c *libvirt.Connect, n *libvirt.StoragePool) {
				bytesBuffer := new(bytes.Buffer)
				bytesBuffer.WriteString(fmt.Sprintf("callback id %d: restart\n", callbackId))
				io.Copy(pipeWriter, bytesBuffer)
			},
		)
		virestError.Error, isError = errorGetCallbackId.(libvirt.Error)
		if isError {
			return virestError, true
		}
		poolConnection.storagePoolEventDeregister(<-usedCallbackId)

		errorEventRemoveTimeout := libvirt.EventRemoveTimeout(addTimeout)
		virestError.Error, isError = errorEventRemoveTimeout.(libvirt.Error)
		if isError {
			return virestError, true
		}
	}

	return virestError, isError
}
