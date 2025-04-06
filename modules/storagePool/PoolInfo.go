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

	var (
		result    poolInfo.Info
		waitGroup sync.WaitGroup
	)
	result.Uuid = uuid
	waitGroup.Add(4)
	go func() {
		defer waitGroup.Done()

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

		result.Name = storagePoolName
	}()
	go func() {
		defer waitGroup.Done()

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

		result.State = storagePoolInfo.State
		result.Capacity = storagePoolInfo.Capacity
		result.Allocation = storagePoolInfo.Allocation
		result.Available = storagePoolInfo.Available
	}()
	go func() {
		defer waitGroup.Done()

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

		result.Autostart = storagePoolAutostart
	}()
	go func() {
		defer waitGroup.Done()

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

		result.Persistent = storagePoolPersistent
	}()
	waitGroup.Wait()

	return result, virestError, false
}
