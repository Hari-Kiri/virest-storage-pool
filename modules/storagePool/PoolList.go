package storagePool

import (
	"fmt"
	"sync"

	"github.com/Hari-Kiri/temboLog"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolList"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

// Collect the list of storage pools, and allocate an array to store those objects.
// Normally, all storage pools are returned; however, flags can be used to filter the results for a smaller list of targeted pools.
// More about option UInteger [https://libvirt.org/html/libvirt-libvirt-storage.html#virConnectListAllStoragePoolsFlags].
func PoolList(qemuConnection *libvirt.Connect, option libvirt.ConnectListAllStoragePoolsFlags, storageXmlFlags libvirt.StorageXMLFlags) ([]poolList.Data, libvirt.Error, bool) {
	var (
		waitGroup    sync.WaitGroup
		libvirtError libvirt.Error
		isError      bool
	)

	storagePools, errorGetListOfStoragePool := qemuConnection.ListAllStoragePools(option)
	libvirtError, isError = errorGetListOfStoragePool.(libvirt.Error)
	if isError {
		libvirtError.Message = fmt.Sprintf("failed list storage pool: %s", libvirtError.Message)
		return nil, libvirtError, true
	}

	result := make([]poolList.Data, len(storagePools))
	waitGroup.Add(len(storagePools) * 4)
	for i := 0; i < len(storagePools); i++ {
		go func(index int) {
			defer waitGroup.Done()

			errorGetStoragePoolRef := storagePools[index].Ref()
			if errorGetStoragePoolRef != nil {
				temboLog.ErrorLogging("error increase the reference count on the storage pool:", errorGetStoragePoolRef)
				return
			}
			defer storagePools[index].Free()

			storagePoolXml, errorGetStoragePoolXml := storagePools[index].GetXMLDesc(storageXmlFlags)
			if errorGetStoragePoolXml != nil {
				temboLog.ErrorLogging("failed get XML of pool", errorGetStoragePoolXml)
				return
			}

			var storagePool libvirtxml.StoragePool
			errorUnmarshallStoragePool := storagePool.Unmarshal(storagePoolXml)
			if errorUnmarshallStoragePool != nil {
				temboLog.ErrorLogging("failed Unmarshal XML of pool", errorUnmarshallStoragePool)
				return
			}

			result[index].Uuid = storagePool.UUID
			result[index].Name = storagePool.Name
			result[index].Capacity = *storagePool.Capacity
			result[index].Allocation = *storagePool.Allocation
			result[index].Available = *storagePool.Available
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

	return result, libvirtError, false
}