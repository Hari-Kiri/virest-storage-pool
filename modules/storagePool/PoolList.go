package storagePool

import (
	"fmt"
	"sync"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolList"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

// Collect the list of storage pools, and allocate an array to store those objects.
// Normally, all storage pools are returned; however, flags can be used to filter the results for a smaller list of targeted pools.
// More about option UInteger [https://libvirt.org/html/libvirt-libvirt-storage.html#virConnectListAllStoragePoolsFlags].
func PoolList(connection virest.Connection, option uint, storageXmlFlags uint) ([]poolList.Data, virest.Error, bool) {
	var (
		virestError virest.Error
		isError     bool
	)

	storagePools, errorGetListOfStoragePool := connection.ListAllStoragePools(libvirt.ConnectListAllStoragePoolsFlags(option))
	virestError.Error, isError = errorGetListOfStoragePool.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed list storage pool: %s", virestError.Message)
		return nil, virestError, true
	}

	var waitGroup sync.WaitGroup
	result := make([]poolList.Data, len(storagePools))
	waitGroup.Add(len(storagePools) * 4)
	for i := 0; i < len(storagePools); i++ {
		defer storagePools[i].Free()

		go func(index int) {
			defer waitGroup.Done()

			errorGetStoragePoolRef := storagePools[index].Ref()
			if errorGetStoragePoolRef != nil {
				temboLog.ErrorLogging("error increase the reference count on the storage pool:", errorGetStoragePoolRef)
				return
			}
			defer storagePools[index].Free()

			storagePoolDetail, errorGetStoragePoolDetail, isError := getPoolDetail(storagePools[index], libvirt.StorageXMLFlags(storageXmlFlags))
			if isError {
				temboLog.ErrorLogging("failed get pool detail", errorGetStoragePoolDetail)
				return
			}

			result[index].Uuid = storagePoolDetail.UUID
			result[index].Name = storagePoolDetail.Name
			result[index].Capacity = *storagePoolDetail.Capacity
			result[index].Allocation = *storagePoolDetail.Allocation
			result[index].Available = *storagePoolDetail.Available
		}(i)
		go func(index int) {
			defer waitGroup.Done()

			errorGetStoragePoolRef := storagePools[index].Ref()
			if errorGetStoragePoolRef != nil {
				temboLog.ErrorLogging("error increase the reference count on the storage pool:", errorGetStoragePoolRef)
				return
			}
			defer storagePools[index].Free()

			storagePoolInfo, errorGetStoragePoolInfo := storagePools[index].GetInfo()
			if errorGetStoragePoolInfo != nil {
				temboLog.ErrorLogging("failed get XML of pool", errorGetStoragePoolInfo)
				return
			}

			result[index].State = storagePoolInfo.State
		}(i)
		go func(index int) {
			defer waitGroup.Done()

			errorGetStoragePoolRef := storagePools[index].Ref()
			if errorGetStoragePoolRef != nil {
				temboLog.ErrorLogging("error increase the reference count on the storage pool:", errorGetStoragePoolRef)
				return
			}
			defer storagePools[index].Free()

			storagePoolAutostart, errorGetStoragePoolAutostart := storagePools[index].GetAutostart()
			if errorGetStoragePoolAutostart != nil {
				temboLog.ErrorLogging("failed get XML of pool", errorGetStoragePoolAutostart)
				return
			}

			result[index].Autostart = storagePoolAutostart
		}(i)
		go func(index int) {
			defer waitGroup.Done()

			errorGetStoragePoolRef := storagePools[index].Ref()
			if errorGetStoragePoolRef != nil {
				temboLog.ErrorLogging("error increase the reference count on the storage pool:", errorGetStoragePoolRef)
				return
			}
			defer storagePools[index].Free()

			storagePoolPersistent, errorGetStoragePoolPersistent := storagePools[index].IsPersistent()
			if errorGetStoragePoolPersistent != nil {
				temboLog.ErrorLogging("failed get XML of pool", errorGetStoragePoolPersistent)
				return
			}

			result[index].Persistent = storagePoolPersistent
		}(i)
	}
	waitGroup.Wait()

	return result, virestError, false
}
