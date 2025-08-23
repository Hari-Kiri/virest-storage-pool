package storagePool

import (
	"fmt"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolList"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

// Collect the list of storage pools, and allocate an array to store those objects.
// Normally, all storage pools are returned; however, flags can be used to filter the results for a smaller list of targeted pools.
// More about option UInteger [https://libvirt.org/html/libvirt-libvirt-storage.html#virConnectListAllStoragePoolsFlags].
func (poolConnection *poolConnection) PoolList(option uint, storageXmlFlags uint) ([]poolList.Data, virest.Error, bool) {
	var (
		virestError virest.Error
		isError     bool
	)

	storagePools, errorGetListOfStoragePool := poolConnection.ListAllStoragePools(libvirt.ConnectListAllStoragePoolsFlags(option))
	virestError.Error, isError = errorGetListOfStoragePool.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed list storage pool: %s", virestError.Message)
		return nil, virestError, isError
	}

	storagePoolDetailUuidChannel := make(chan string)
	storagePoolDetailNameChannel := make(chan string)
	storagePoolDetailCapacityChannel := make(chan libvirtxml.StoragePoolSize)
	storagePoolDetailAllocationChannel := make(chan libvirtxml.StoragePoolSize)
	storagePoolDetailAvailableChannel := make(chan libvirtxml.StoragePoolSize)
	storagePoolInfoStateChannel := make(chan libvirt.StoragePoolState)
	storagePoolAutostartChannel := make(chan bool)
	storagePoolPersistentChannel := make(chan bool)
	virestErrorChannel := make(chan virest.Error, 3)
	isErrorChannel := make(chan bool, 3)
	result := make([]poolList.Data, len(storagePools))
	for i := 0; i < len(storagePools); i++ {
		defer storagePools[i].Free()

		go func(storagePoolObject libvirt.StoragePool) {
			var (
				storagePoolDetail libvirtxml.StoragePool
				virestError       virest.Error
				isError           bool
			)

			errorGetStoragePoolRef := storagePoolObject.Ref()
			virestError.Error, isError = errorGetStoragePoolRef.(libvirt.Error)
			if isError {
				storagePoolDetailUuidChannel <- storagePoolDetail.UUID
				storagePoolDetailNameChannel <- storagePoolDetail.Name
				storagePoolDetailCapacityChannel <- libvirtxml.StoragePoolSize{}
				storagePoolDetailAllocationChannel <- libvirtxml.StoragePoolSize{}
				storagePoolDetailAvailableChannel <- libvirtxml.StoragePoolSize{}
				virestErrorChannel <- virestError
				isErrorChannel <- isError

				return
			}
			defer storagePoolObject.Free()

			storagePoolDetail, virestError.Error, isError = getPoolDetail(storagePoolObject, libvirt.StorageXMLFlags(storageXmlFlags))
			if isError {
				storagePoolDetailUuidChannel <- storagePoolDetail.UUID
				storagePoolDetailNameChannel <- storagePoolDetail.Name
				storagePoolDetailCapacityChannel <- libvirtxml.StoragePoolSize{}
				storagePoolDetailAllocationChannel <- libvirtxml.StoragePoolSize{}
				storagePoolDetailAvailableChannel <- libvirtxml.StoragePoolSize{}
				virestErrorChannel <- virestError
				isErrorChannel <- isError
			}

			storagePoolDetailUuidChannel <- storagePoolDetail.UUID
			storagePoolDetailNameChannel <- storagePoolDetail.Name
			storagePoolDetailCapacityChannel <- *storagePoolDetail.Capacity
			storagePoolDetailAllocationChannel <- *storagePoolDetail.Allocation
			storagePoolDetailAvailableChannel <- *storagePoolDetail.Available
			virestErrorChannel <- virest.Error{}
			isErrorChannel <- false
		}(storagePools[i])

		go func(storagePoolObject libvirt.StoragePool) {
			var (
				virestError virest.Error
				isError     bool
			)

			errorGetStoragePoolRef := storagePoolObject.Ref()
			virestError.Error, isError = errorGetStoragePoolRef.(libvirt.Error)
			if isError {
				storagePoolInfoStateChannel <- libvirt.StoragePoolInfo{}.State
				virestErrorChannel <- virestError
				isErrorChannel <- isError

				return
			}
			defer storagePoolObject.Free()

			storagePoolInfo, errorGetStoragePoolInfo := storagePoolObject.GetInfo()
			virestError.Error, isError = errorGetStoragePoolInfo.(libvirt.Error)
			if isError {
				storagePoolInfoStateChannel <- storagePoolInfo.State
				virestErrorChannel <- virestError
				isErrorChannel <- isError

				return
			}

			storagePoolInfoStateChannel <- storagePoolInfo.State
			virestErrorChannel <- virest.Error{}
			isErrorChannel <- false
		}(storagePools[i])

		go func(storagePoolObject libvirt.StoragePool) {
			var (
				virestError virest.Error
				isError     bool
			)

			errorGetStoragePoolRef := storagePoolObject.Ref()
			virestError.Error, isError = errorGetStoragePoolRef.(libvirt.Error)
			if isError {
				storagePoolAutostartChannel <- false
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
		}(storagePools[i])

		go func(storagePoolObject libvirt.StoragePool) {
			var (
				virestError virest.Error
				isError     bool
			)

			errorGetStoragePoolRef := storagePoolObject.Ref()
			virestError.Error, isError = errorGetStoragePoolRef.(libvirt.Error)
			if isError {
				storagePoolPersistentChannel <- false
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
		}(storagePools[i])

		storagePoolUuid := <-storagePoolDetailUuidChannel
		storagePoolName := <-storagePoolDetailNameChannel
		storagePoolCapacity := <-storagePoolDetailCapacityChannel
		storagePoolAllocation := <-storagePoolDetailAllocationChannel
		storagePoolAvailable := <-storagePoolDetailAvailableChannel
		storagePoolInfoState := <-storagePoolInfoStateChannel
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
			break
		}

		result[i].Uuid = storagePoolUuid
		result[i].Name = storagePoolName
		result[i].Capacity = storagePoolCapacity
		result[i].Allocation = storagePoolAllocation
		result[i].Available = storagePoolAvailable
		result[i].State = storagePoolInfoState
		result[i].Autostart = storagePoolAutostart
		result[i].Persistent = storagePoolPersistent
	}

	return result, virestError, isError
}
