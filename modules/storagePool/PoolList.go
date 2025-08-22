package storagePool

import (
	"fmt"

	"github.com/Hari-Kiri/temboLog"
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

	result := make([]poolList.Data, len(storagePools))
	for i := 0; i < len(storagePools); i++ {
		defer storagePools[i].Free()

		storagePoolDetailUuidChannel := make(chan string)
		storagePoolDetailNameChannel := make(chan string)
		storagePoolDetailCapacityChannel := make(chan libvirtxml.StoragePoolSize)
		storagePoolDetailAllocationChannel := make(chan libvirtxml.StoragePoolSize)
		storagePoolDetailAvailableChannel := make(chan libvirtxml.StoragePoolSize)
		go func(storagePoolObject libvirt.StoragePool) {
			errorGetStoragePoolRef := storagePoolObject.Ref()
			if errorGetStoragePoolRef != nil {
				temboLog.ErrorLogging("error increase the reference count on the storage pool:", errorGetStoragePoolRef)
				return
			}
			defer storagePoolObject.Free()

			storagePoolDetail, errorGetStoragePoolDetail, isError := getPoolDetail(storagePoolObject, libvirt.StorageXMLFlags(storageXmlFlags))
			if isError {
				temboLog.ErrorLogging("failed get pool detail", errorGetStoragePoolDetail)
				return
			}

			storagePoolDetailUuidChannel <- storagePoolDetail.UUID
			storagePoolDetailNameChannel <- storagePoolDetail.Name
			storagePoolDetailCapacityChannel <- *storagePoolDetail.Capacity
			storagePoolDetailAllocationChannel <- *storagePoolDetail.Allocation
			storagePoolDetailAvailableChannel <- *storagePoolDetail.Available
		}(storagePools[i])

		storagePoolDetailStateChannel := make(chan libvirt.StoragePoolState)
		go func(storagePoolObject libvirt.StoragePool) {
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

			storagePoolDetailStateChannel <- storagePoolInfo.State
		}(storagePools[i])

		storagePoolDetailAutostartChannel := make(chan bool)
		go func(storagePoolObject libvirt.StoragePool) {
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

			storagePoolDetailAutostartChannel <- storagePoolAutostart
		}(storagePools[i])

		storagePoolDetailPersistentChannel := make(chan bool)
		go func(storagePoolObject libvirt.StoragePool) {
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

			storagePoolDetailPersistentChannel <- storagePoolPersistent
		}(storagePools[i])

		result[i].Uuid = <-storagePoolDetailUuidChannel
		result[i].Name = <-storagePoolDetailNameChannel
		result[i].Capacity = <-storagePoolDetailCapacityChannel
		result[i].Allocation = <-storagePoolDetailAllocationChannel
		result[i].Available = <-storagePoolDetailAvailableChannel
		result[i].State = <-storagePoolDetailStateChannel
		result[i].Autostart = <-storagePoolDetailAutostartChannel
		result[i].Persistent = <-storagePoolDetailPersistentChannel
	}

	return result, virestError, false
}
