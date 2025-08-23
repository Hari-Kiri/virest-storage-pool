package storagePool

import (
	"fmt"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolInfo"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

// Get volatile information about the storage pool such as free space / usage summary
func (poolConnection *poolConnection) PoolInfo(uuid string) (poolInfo.Info, virest.Error, bool) {
	var (
		virestError virest.Error
		isError     bool
	)

	storagePoolObject, errorGetStoragePoolObject := poolConnection.LookupStoragePoolByUUIDString(uuid)
	virestError.Error, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed list storage pool: %s", virestError.Message)
		return poolInfo.Info{}, virestError, isError
	}

	storagePoolNameChannel := make(chan string)
	errorGetStoragePoolNameChannel := make(chan virest.Error)
	isErrorGetStoragePoolNameChannel := make(chan bool)
	go func() {
		var (
			virestError virest.Error
			isError     bool
		)

		errorGetStoragePoolRef := storagePoolObject.Ref()
		virestError.Error, isError = errorGetStoragePoolRef.(libvirt.Error)
		if isError {
			storagePoolName := ""

			storagePoolNameChannel <- storagePoolName
			errorGetStoragePoolNameChannel <- virestError
			isErrorGetStoragePoolNameChannel <- isError

			return
		}
		defer storagePoolObject.Free()

		storagePoolName, errorGetStoragePoolName := storagePoolObject.GetName()
		virestError.Error, isError = errorGetStoragePoolName.(libvirt.Error)
		if isError {
			storagePoolNameChannel <- storagePoolName
			errorGetStoragePoolNameChannel <- virestError
			isErrorGetStoragePoolNameChannel <- isError

			return
		}

		storagePoolNameChannel <- storagePoolName
		errorGetStoragePoolNameChannel <- virest.Error{}
		isErrorGetStoragePoolNameChannel <- false
	}()

	storagePoolStateChannel := make(chan libvirt.StoragePoolState)
	storagePoolCapacityChannel := make(chan uint64)
	storagePoolAllocationChannel := make(chan uint64)
	storagePoolAvailableChannel := make(chan uint64)
	errorGetStoragePoolInfoChannel := make(chan virest.Error)
	isErrorGetStoragePoolInfoChannel := make(chan bool)
	go func() {
		var (
			virestError virest.Error
			isError     bool
		)

		errorGetStoragePoolRef := storagePoolObject.Ref()
		virestError.Error, isError = errorGetStoragePoolRef.(libvirt.Error)
		if isError {
			storagePoolInfo := libvirt.StoragePoolInfo{}

			storagePoolStateChannel <- storagePoolInfo.State
			storagePoolCapacityChannel <- storagePoolInfo.Capacity
			storagePoolAllocationChannel <- storagePoolInfo.Allocation
			storagePoolAvailableChannel <- storagePoolInfo.Available
			errorGetStoragePoolInfoChannel <- virestError
			isErrorGetStoragePoolInfoChannel <- isError

			return
		}
		defer storagePoolObject.Free()

		storagePoolInfo, errorGetStoragePoolInfo := storagePoolObject.GetInfo()
		virestError.Error, isError = errorGetStoragePoolInfo.(libvirt.Error)
		if isError {
			storagePoolStateChannel <- storagePoolInfo.State
			storagePoolCapacityChannel <- storagePoolInfo.Capacity
			storagePoolAllocationChannel <- storagePoolInfo.Allocation
			storagePoolAvailableChannel <- storagePoolInfo.Available
			errorGetStoragePoolInfoChannel <- virestError
			isErrorGetStoragePoolInfoChannel <- isError

			return
		}

		storagePoolStateChannel <- storagePoolInfo.State
		storagePoolCapacityChannel <- storagePoolInfo.Capacity
		storagePoolAllocationChannel <- storagePoolInfo.Allocation
		storagePoolAvailableChannel <- storagePoolInfo.Available
		errorGetStoragePoolInfoChannel <- virest.Error{}
		isErrorGetStoragePoolInfoChannel <- false
	}()

	storagePoolAutostartChannel := make(chan bool)
	errorGetStoragePoolAutostartChannel := make(chan virest.Error)
	isErrorGetStoragePoolAutostartChannel := make(chan bool)
	go func() {
		var (
			virestError virest.Error
			isError     bool
		)

		errorGetStoragePoolRef := storagePoolObject.Ref()
		virestError.Error, isError = errorGetStoragePoolRef.(libvirt.Error)
		if isError {
			storagePoolAutostart := false

			storagePoolAutostartChannel <- storagePoolAutostart
			errorGetStoragePoolAutostartChannel <- virestError
			isErrorGetStoragePoolAutostartChannel <- isError

			return
		}
		defer storagePoolObject.Free()

		storagePoolAutostart, errorGetStoragePoolAutostart := storagePoolObject.GetAutostart()
		virestError.Error, isError = errorGetStoragePoolAutostart.(libvirt.Error)
		if isError {
			storagePoolAutostartChannel <- storagePoolAutostart
			errorGetStoragePoolAutostartChannel <- virestError
			isErrorGetStoragePoolAutostartChannel <- isError

			return
		}

		storagePoolAutostartChannel <- storagePoolAutostart
		errorGetStoragePoolAutostartChannel <- virest.Error{}
		isErrorGetStoragePoolAutostartChannel <- false
	}()

	storagePoolPersistentChannel := make(chan bool)
	errorStoragePoolPersistentChannel := make(chan virest.Error)
	isErrorStoragePoolPersistentChannel := make(chan bool)
	go func() {
		var (
			virestError virest.Error
			isError     bool
		)

		errorGetStoragePoolRef := storagePoolObject.Ref()
		virestError.Error, isError = errorGetStoragePoolRef.(libvirt.Error)
		if isError {
			storagePoolPersistent := false

			storagePoolPersistentChannel <- storagePoolPersistent
			errorStoragePoolPersistentChannel <- virestError
			isErrorStoragePoolPersistentChannel <- isError

			return
		}
		defer storagePoolObject.Free()

		storagePoolPersistent, errorGetStoragePoolPersistent := storagePoolObject.IsPersistent()
		virestError.Error, isError = errorGetStoragePoolPersistent.(libvirt.Error)
		if isError {
			storagePoolPersistentChannel <- storagePoolPersistent
			errorStoragePoolPersistentChannel <- virestError
			isErrorStoragePoolPersistentChannel <- isError

			return
		}

		storagePoolPersistentChannel <- storagePoolPersistent
		errorStoragePoolPersistentChannel <- virest.Error{}
		isErrorStoragePoolPersistentChannel <- false
	}()

	storagePoolName := <-storagePoolNameChannel
	errorGetStoragePoolName := <-errorGetStoragePoolNameChannel
	isErrorGetStoragePoolName := <-isErrorGetStoragePoolNameChannel

	storagePoolState := <-storagePoolStateChannel
	storagePoolCapacity := <-storagePoolCapacityChannel
	storagePoolAllocation := <-storagePoolAllocationChannel
	storagePoolAvailable := <-storagePoolAvailableChannel
	errorGetStoragePoolInfo := <-errorGetStoragePoolInfoChannel
	isErrorGetStoragePoolInfo := <-isErrorGetStoragePoolInfoChannel

	storagePoolAutostart := <-storagePoolAutostartChannel
	errorGetStoragePoolAutostart := <-errorGetStoragePoolAutostartChannel
	isErrorGetStoragePoolAutostart := <-isErrorGetStoragePoolAutostartChannel

	storagePoolPersistent := <-storagePoolPersistentChannel
	errorStoragePoolPersistent := <-errorStoragePoolPersistentChannel
	isErrorStoragePoolPersistent := <-isErrorStoragePoolPersistentChannel

	if isErrorGetStoragePoolName {
		return poolInfo.Info{}, errorGetStoragePoolName, isErrorGetStoragePoolName
	}

	if isErrorGetStoragePoolInfo {
		return poolInfo.Info{}, errorGetStoragePoolInfo, isErrorGetStoragePoolInfo
	}

	if isErrorGetStoragePoolAutostart {
		return poolInfo.Info{}, errorGetStoragePoolAutostart, isErrorGetStoragePoolAutostart
	}

	if isErrorStoragePoolPersistent {
		return poolInfo.Info{}, errorStoragePoolPersistent, isErrorStoragePoolPersistent
	}

	return poolInfo.Info{
		Uuid:       uuid,
		Name:       storagePoolName,
		State:      storagePoolState,
		Capacity:   storagePoolCapacity,
		Allocation: storagePoolAllocation,
		Available:  storagePoolAvailable,
		Autostart:  storagePoolAutostart,
		Persistent: storagePoolPersistent,
	}, virest.Error{}, false
}
