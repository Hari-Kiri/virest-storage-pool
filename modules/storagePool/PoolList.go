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
		return nil, virestError, true
	}

	storagePoolDetailUuidChannel := make(chan string)
	storagePoolDetailNameChannel := make(chan string)
	storagePoolDetailCapacityChannel := make(chan libvirtxml.StoragePoolSize)
	storagePoolDetailAllocationChannel := make(chan libvirtxml.StoragePoolSize)
	storagePoolDetailAvailableChannel := make(chan libvirtxml.StoragePoolSize)
	errorGetStoragePoolDetailChannel := make(chan virest.Error)
	isErrorGetStoragePoolDetailChannel := make(chan bool)

	storagePoolInfoStateChannel := make(chan libvirt.StoragePoolState)
	errorGetStoragePoolInfoStateChannel := make(chan virest.Error)
	isErrorGetStoragePoolInfoStateChannel := make(chan bool)

	storagePoolAutostartChannel := make(chan bool)
	errorGetstoragePoolAutostartChannel := make(chan virest.Error)
	isErrorGetstoragePoolAutostartChannel := make(chan bool)

	storagePoolPersistentChannel := make(chan bool)
	errorGetstoragePoolPersistentChannel := make(chan virest.Error)
	isErrorGetstoragePoolPersistentChannel := make(chan bool)

	result := make([]poolList.Data, len(storagePools))
	for i := 0; i < len(storagePools); i++ {
		defer storagePools[i].Free()

		go func(storagePoolObject libvirt.StoragePool) {
			errorGetStoragePoolRef := storagePoolObject.Ref()
			if errorGetStoragePoolRef != nil {
				storagePoolDetail := libvirtxml.StoragePool{}
				isError := true

				storagePoolDetailUuidChannel <- storagePoolDetail.UUID
				storagePoolDetailNameChannel <- storagePoolDetail.Name
				storagePoolDetailCapacityChannel <- *storagePoolDetail.Capacity
				storagePoolDetailAllocationChannel <- *storagePoolDetail.Allocation
				storagePoolDetailAvailableChannel <- *storagePoolDetail.Available
				errorGetStoragePoolDetailChannel <- virest.Error{Error: errorGetStoragePoolRef.(libvirt.Error)}
				isErrorGetStoragePoolDetailChannel <- isError

				return
			}
			defer storagePoolObject.Free()

			storagePoolDetail, errorGetStoragePoolDetail, isError := getPoolDetail(storagePoolObject, libvirt.StorageXMLFlags(storageXmlFlags))
			if isError {
				storagePoolDetailUuidChannel <- storagePoolDetail.UUID
				storagePoolDetailNameChannel <- storagePoolDetail.Name
				storagePoolDetailCapacityChannel <- libvirtxml.StoragePoolSize{}
				storagePoolDetailAllocationChannel <- libvirtxml.StoragePoolSize{}
				storagePoolDetailAvailableChannel <- libvirtxml.StoragePoolSize{}
				errorGetStoragePoolDetailChannel <- virest.Error{Error: errorGetStoragePoolDetail}
				isErrorGetStoragePoolDetailChannel <- isError
			}

			storagePoolDetailUuidChannel <- storagePoolDetail.UUID
			storagePoolDetailNameChannel <- storagePoolDetail.Name
			storagePoolDetailCapacityChannel <- *storagePoolDetail.Capacity
			storagePoolDetailAllocationChannel <- *storagePoolDetail.Allocation
			storagePoolDetailAvailableChannel <- *storagePoolDetail.Available
			errorGetStoragePoolDetailChannel <- virest.Error{}
			isErrorGetStoragePoolDetailChannel <- isError
		}(storagePools[i])

		go func(storagePoolObject libvirt.StoragePool) {
			errorGetStoragePoolRef := storagePoolObject.Ref()
			if errorGetStoragePoolRef != nil {
				storagePoolInfo := &libvirt.StoragePoolInfo{}
				isError := true

				storagePoolInfoStateChannel <- storagePoolInfo.State
				errorGetStoragePoolInfoStateChannel <- virest.Error{Error: errorGetStoragePoolRef.(libvirt.Error)}
				isErrorGetStoragePoolInfoStateChannel <- isError

				return
			}
			defer storagePoolObject.Free()

			storagePoolInfo, errorGetStoragePoolInfo := storagePoolObject.GetInfo()
			if errorGetStoragePoolInfo != nil {
				isError := true

				storagePoolInfoStateChannel <- storagePoolInfo.State
				errorGetStoragePoolInfoStateChannel <- virest.Error{Error: errorGetStoragePoolInfo.(libvirt.Error)}
				isErrorGetStoragePoolInfoStateChannel <- isError

				return
			}

			storagePoolInfoStateChannel <- storagePoolInfo.State
			errorGetStoragePoolInfoStateChannel <- virest.Error{}
			isErrorGetStoragePoolInfoStateChannel <- isError
		}(storagePools[i])

		go func(storagePoolObject libvirt.StoragePool) {
			errorGetStoragePoolRef := storagePoolObject.Ref()
			if errorGetStoragePoolRef != nil {
				storagePoolAutostart := false
				isError := true

				storagePoolAutostartChannel <- storagePoolAutostart
				errorGetstoragePoolAutostartChannel <- virest.Error{Error: errorGetStoragePoolRef.(libvirt.Error)}
				isErrorGetstoragePoolAutostartChannel <- isError

				return
			}
			defer storagePoolObject.Free()

			storagePoolAutostart, errorGetStoragePoolAutostart := storagePoolObject.GetAutostart()
			if errorGetStoragePoolAutostart != nil {
				isError := true

				storagePoolAutostartChannel <- storagePoolAutostart
				errorGetstoragePoolAutostartChannel <- virest.Error{Error: errorGetStoragePoolAutostart.(libvirt.Error)}
				isErrorGetstoragePoolAutostartChannel <- isError

				return
			}

			storagePoolAutostartChannel <- storagePoolAutostart
			errorGetstoragePoolAutostartChannel <- virest.Error{}
			isErrorGetstoragePoolAutostartChannel <- isError
		}(storagePools[i])

		go func(storagePoolObject libvirt.StoragePool) {
			errorGetStoragePoolRef := storagePoolObject.Ref()
			if errorGetStoragePoolRef != nil {
				storagePoolPersistent := false
				isError := true

				storagePoolPersistentChannel <- storagePoolPersistent
				errorGetstoragePoolPersistentChannel <- virest.Error{Error: errorGetStoragePoolRef.(libvirt.Error)}
				isErrorGetstoragePoolPersistentChannel <- isError

				return
			}
			defer storagePoolObject.Free()

			storagePoolPersistent, errorGetStoragePoolPersistent := storagePoolObject.IsPersistent()
			if errorGetStoragePoolPersistent != nil {
				isError := true

				storagePoolPersistentChannel <- storagePoolPersistent
				errorGetstoragePoolPersistentChannel <- virest.Error{Error: errorGetStoragePoolPersistent.(libvirt.Error)}
				isErrorGetstoragePoolPersistentChannel <- isError

				return
			}

			storagePoolPersistentChannel <- storagePoolPersistent
			errorGetstoragePoolPersistentChannel <- virest.Error{}
			isErrorGetstoragePoolPersistentChannel <- isError
		}(storagePools[i])

		storagePoolUuid := <-storagePoolDetailUuidChannel
		storagePoolName := <-storagePoolDetailNameChannel
		storagePoolCapacity := <-storagePoolDetailCapacityChannel
		storagePoolAllocation := <-storagePoolDetailAllocationChannel
		storagePoolAvailable := <-storagePoolDetailAvailableChannel
		errorGetStoragePoolDetail := <-errorGetStoragePoolDetailChannel
		isErrorGetStoragePoolDetail := <-isErrorGetStoragePoolDetailChannel
		if isErrorGetStoragePoolDetail {
			virestError = errorGetStoragePoolDetail
			isError = isErrorGetStoragePoolDetail
			break
		}
		result[i].Uuid = storagePoolUuid
		result[i].Name = storagePoolName
		result[i].Capacity = storagePoolCapacity
		result[i].Allocation = storagePoolAllocation
		result[i].Available = storagePoolAvailable

		storagePoolInfoState := <-storagePoolInfoStateChannel
		errorGetStoragePoolInfoState := <-errorGetStoragePoolInfoStateChannel
		isErrorGetStoragePoolInfoState := <-isErrorGetStoragePoolInfoStateChannel
		if isErrorGetStoragePoolInfoState {
			virestError = errorGetStoragePoolInfoState
			isError = isErrorGetStoragePoolInfoState
			break
		}
		result[i].State = storagePoolInfoState

		storagePoolAutostart := <-storagePoolAutostartChannel
		errorGetStoragePoolAutostart := <-errorGetstoragePoolAutostartChannel
		isErrorGetStoragePoolAutostart := <-isErrorGetstoragePoolAutostartChannel
		if isErrorGetStoragePoolAutostart {
			virestError = errorGetStoragePoolAutostart
			isError = isErrorGetStoragePoolAutostart
			break
		}
		result[i].Autostart = storagePoolAutostart

		storagePoolPersistent := <-storagePoolPersistentChannel
		errorGetStoragePoolPersistent := <-errorGetstoragePoolPersistentChannel
		isErrorGetStoragePoolPersistent := <-isErrorGetstoragePoolPersistentChannel
		if isErrorGetStoragePoolPersistent {
			virestError = errorGetStoragePoolPersistent
			isError = isErrorGetStoragePoolPersistent
			break
		}
		result[i].Persistent = storagePoolPersistent
	}

	return result, virestError, isError
}
