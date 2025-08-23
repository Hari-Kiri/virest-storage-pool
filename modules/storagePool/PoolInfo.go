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
	storagePoolStateChannel := make(chan libvirt.StoragePoolState)
	storagePoolCapacityChannel := make(chan uint64)
	storagePoolAllocationChannel := make(chan uint64)
	storagePoolAvailableChannel := make(chan uint64)
	storagePoolAutostartChannel := make(chan bool)
	storagePoolPersistentChannel := make(chan bool)
	virestErrorChannel := make(chan virest.Error, 3)
	isErrorChannel := make(chan bool, 3)
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
			virestErrorChannel <- virestError
			isErrorChannel <- isError

			return
		}
		defer storagePoolObject.Free()

		storagePoolName, errorGetStoragePoolName := storagePoolObject.GetName()
		virestError.Error, isError = errorGetStoragePoolName.(libvirt.Error)
		if isError {
			storagePoolNameChannel <- storagePoolName
			virestErrorChannel <- virestError
			isErrorChannel <- isError

			return
		}

		storagePoolNameChannel <- storagePoolName
		virestErrorChannel <- virest.Error{}
		isErrorChannel <- false
	}()

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
			virestErrorChannel <- virestError
			isErrorChannel <- isError

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
			virestErrorChannel <- virestError
			isErrorChannel <- isError

			return
		}

		storagePoolStateChannel <- storagePoolInfo.State
		storagePoolCapacityChannel <- storagePoolInfo.Capacity
		storagePoolAllocationChannel <- storagePoolInfo.Allocation
		storagePoolAvailableChannel <- storagePoolInfo.Available
		virestErrorChannel <- virest.Error{}
		isErrorChannel <- false
	}()

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
			virestErrorChannel <- virestError
			isErrorChannel <- isError

			return
		}
		defer storagePoolObject.Free()

		storagePoolAutostart, errorGetStoragePoolAutostart := storagePoolObject.GetAutostart()
		virestError.Error, isError = errorGetStoragePoolAutostart.(libvirt.Error)
		if isError {
			storagePoolAutostartChannel <- storagePoolAutostart
			virestErrorChannel <- virestError
			isErrorChannel <- isError

			return
		}

		storagePoolAutostartChannel <- storagePoolAutostart
		virestErrorChannel <- virest.Error{}
		isErrorChannel <- false
	}()

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
			virestErrorChannel <- virestError
			isErrorChannel <- isError

			return
		}
		defer storagePoolObject.Free()

		storagePoolPersistent, errorGetStoragePoolPersistent := storagePoolObject.IsPersistent()
		virestError.Error, isError = errorGetStoragePoolPersistent.(libvirt.Error)
		if isError {
			storagePoolPersistentChannel <- storagePoolPersistent
			virestErrorChannel <- virestError
			isErrorChannel <- isError

			return
		}

		storagePoolPersistentChannel <- storagePoolPersistent
		virestErrorChannel <- virest.Error{}
		isErrorChannel <- false
	}()

	storagePoolName := <-storagePoolNameChannel
	storagePoolState := <-storagePoolStateChannel
	storagePoolCapacity := <-storagePoolCapacityChannel
	storagePoolAllocation := <-storagePoolAllocationChannel
	storagePoolAvailable := <-storagePoolAvailableChannel
	storagePoolAutostart := <-storagePoolAutostartChannel
	storagePoolPersistent := <-storagePoolPersistentChannel
	for i := 0; i < cap(virestErrorChannel); i++ {
		virestError = <-virestErrorChannel
		isError = <-isErrorChannel
		if isError {
			break
		}
	}
	if isError {
		return poolInfo.Info{}, virestError, isError
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
