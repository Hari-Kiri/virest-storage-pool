package storagePool

import (
	"fmt"
	"sync"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolInfo"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

// Get volatile information about the storage pool such as free space / usage summary
func PoolInfo(connection virest.Connection, uuid string) (poolInfo.Info, virest.Error, bool) {
	var (
		virestError virest.Error
		isError     bool
	)

	storagePool, errorGetStoragePoolObject := connection.LookupStoragePoolByUUIDString(uuid)
	virestError.Error, isError = errorGetStoragePoolObject.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed list storage pool: %s", virestError.Message)
		return poolInfo.Info{}, virestError, true
	}

	var (
		result    poolInfo.Info
		waitGroup sync.WaitGroup
	)
	result.Uuid = uuid
	waitGroup.Add(4)
	go func() {
		defer waitGroup.Done()

		errorGetStoragePoolRef := storagePool.Ref()
		if errorGetStoragePoolRef != nil {
			temboLog.ErrorLogging("error increase the reference count on the storage pool:", errorGetStoragePoolRef)
			return
		}
		defer storagePool.Free()

		storagePoolName, errorGetStoragePoolName := storagePool.GetName()
		if errorGetStoragePoolName != nil {
			temboLog.ErrorLogging("failed get storage pool name", errorGetStoragePoolName)
			return
		}

		result.Name = storagePoolName
	}()
	go func() {
		defer waitGroup.Done()

		errorGetStoragePoolRef := storagePool.Ref()
		if errorGetStoragePoolRef != nil {
			temboLog.ErrorLogging("error increase the reference count on the storage pool:", errorGetStoragePoolRef)
			return
		}
		defer storagePool.Free()

		storagePoolInfo, errorGetStoragePoolInfo := storagePool.GetInfo()
		if errorGetStoragePoolInfo != nil {
			temboLog.ErrorLogging("failed get XML of pool", errorGetStoragePoolInfo)
			return
		}

		result.State = storagePoolInfo.State
		result.Capacity = storagePoolInfo.Capacity
		result.Allocation = storagePoolInfo.Allocation
		result.Available = storagePoolInfo.Available
	}()
	go func() {
		defer waitGroup.Done()

		errorGetStoragePoolRef := storagePool.Ref()
		if errorGetStoragePoolRef != nil {
			temboLog.ErrorLogging("error increase the reference count on the storage pool:", errorGetStoragePoolRef)
			return
		}
		defer storagePool.Free()

		storagePoolAutostart, errorGetStoragePoolAutostart := storagePool.GetAutostart()
		if errorGetStoragePoolAutostart != nil {
			temboLog.ErrorLogging("failed get XML of pool", errorGetStoragePoolAutostart)
			return
		}

		result.Autostart = storagePoolAutostart
	}()
	go func() {
		defer waitGroup.Done()

		errorGetStoragePoolRef := storagePool.Ref()
		if errorGetStoragePoolRef != nil {
			temboLog.ErrorLogging("error increase the reference count on the storage pool:", errorGetStoragePoolRef)
			return
		}
		defer storagePool.Free()

		storagePoolPersistent, errorGetStoragePoolPersistent := storagePool.IsPersistent()
		if errorGetStoragePoolPersistent != nil {
			temboLog.ErrorLogging("failed get XML of pool", errorGetStoragePoolPersistent)
			return
		}

		result.Persistent = storagePoolPersistent
	}()
	waitGroup.Wait()

	return result, virestError, false
}
