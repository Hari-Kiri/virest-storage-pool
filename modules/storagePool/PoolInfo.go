package storagePool

import (
	"fmt"

	"github.com/Hari-Kiri/temboLog"
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
		return poolInfo.Info{}, virestError, true
	}

	storagePoolNameChannel := make(chan string)
	go func() {
		errorGetStoragePoolRef := storagePoolObject.Ref()
		if errorGetStoragePoolRef != nil {
			temboLog.ErrorLogging("error increase the reference count on the storage pool:", errorGetStoragePoolRef)
			return
		}
		defer storagePoolObject.Free()

		storagePoolName, errorGetStoragePoolName := storagePoolObject.GetName()
		if errorGetStoragePoolName != nil {
			temboLog.ErrorLogging("failed get storage pool name", errorGetStoragePoolName)
			return
		}

		storagePoolNameChannel <- storagePoolName
	}()

	storagePoolStateChannel := make(chan libvirt.StoragePoolState)
	storagePoolCapacityChannel := make(chan uint64)
	storagePoolAllocationChannel := make(chan uint64)
	storagePoolAvailableChannel := make(chan uint64)
	go func() {
		errorGetStoragePoolRef := storagePoolObject.Ref()
		if errorGetStoragePoolRef != nil {
			temboLog.ErrorLogging("error increase the reference count on the storage pool:", errorGetStoragePoolRef)
			return
		}
		defer storagePoolObject.Free()

		storagePoolInfo, errorGetStoragePoolInfo := storagePoolObject.GetInfo()
		if errorGetStoragePoolInfo != nil {
			temboLog.ErrorLogging("failed get XML of pool", errorGetStoragePoolInfo)
			return
		}

		storagePoolStateChannel <- storagePoolInfo.State
		storagePoolCapacityChannel <- storagePoolInfo.Capacity
		storagePoolAllocationChannel <- storagePoolInfo.Allocation
		storagePoolAvailableChannel <- storagePoolInfo.Available
	}()

	storagePoolAutostartChannel := make(chan bool)
	go func() {
		errorGetStoragePoolRef := storagePoolObject.Ref()
		if errorGetStoragePoolRef != nil {
			temboLog.ErrorLogging("error increase the reference count on the storage pool:", errorGetStoragePoolRef)
			return
		}
		defer storagePoolObject.Free()

		storagePoolAutostart, errorGetStoragePoolAutostart := storagePoolObject.GetAutostart()
		if errorGetStoragePoolAutostart != nil {
			temboLog.ErrorLogging("failed get XML of pool", errorGetStoragePoolAutostart)
			return
		}

		storagePoolAutostartChannel <- storagePoolAutostart
	}()

	storagePoolPersistentChannel := make(chan bool)
	go func() {
		errorGetStoragePoolRef := storagePoolObject.Ref()
		if errorGetStoragePoolRef != nil {
			temboLog.ErrorLogging("error increase the reference count on the storage pool:", errorGetStoragePoolRef)
			return
		}
		defer storagePoolObject.Free()

		storagePoolPersistent, errorGetStoragePoolPersistent := storagePoolObject.IsPersistent()
		if errorGetStoragePoolPersistent != nil {
			temboLog.ErrorLogging("failed get XML of pool", errorGetStoragePoolPersistent)
			return
		}

		storagePoolPersistentChannel <- storagePoolPersistent
	}()

	return poolInfo.Info{
		Uuid:       uuid,
		Name:       <-storagePoolNameChannel,
		State:      <-storagePoolStateChannel,
		Capacity:   <-storagePoolCapacityChannel,
		Allocation: <-storagePoolAllocationChannel,
		Available:  <-storagePoolAvailableChannel,
		Autostart:  <-storagePoolAutostartChannel,
		Persistent: <-storagePoolPersistentChannel,
	}, virestError, false
}
