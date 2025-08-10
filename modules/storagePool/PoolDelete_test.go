package storagePool

import (
	"testing"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

func (poolConnection *poolConnection) helperTestPoolDelete(test *testing.T, poolUuid string, option libvirt.StoragePoolDeleteFlags) (virest.Error, bool) {
	test.Helper()

	errorPoolDelete, isErrorPoolDelete := poolConnection.PoolDelete(poolUuid, option)
	if isErrorPoolDelete {
		test.Fail()
		return errorPoolDelete, isErrorPoolDelete
	}

	return virest.Error{}, false
}

func TestPoolDeleteMetadaOnly(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	errorPoolBuild, isErrorPoolBuild := poolConnection.helperTestPoolBuild(test, poolUuid, libvirt.STORAGE_POOL_BUILD_NEW)
	if isErrorPoolBuild {
		test.Fatalf("pool build failed: %s", errorPoolBuild.Message)
	}

	errorPoolCreate, isErrorPoolCreate := poolConnection.helperTestPoolCreate(test, poolUuid, libvirt.STORAGE_POOL_CREATE_NORMAL)
	if isErrorPoolCreate {
		test.Fatalf("creating pool test failed: %s", errorPoolCreate.Message)
	}

	errorPoolDestroy, isErrorPoolDestroy := poolConnection.helperTestPoolDestroy(test, poolUuid)
	if isErrorPoolDestroy {
		test.Fatalf("pool destroy failed: %s", errorPoolDestroy.Message)
	}

	test.Cleanup(func() {
		if errorPoolUndefine, isErrorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}
	})

	errorPoolDelete, isErrorPoolDelete := poolConnection.PoolDelete(poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL)
	if isErrorPoolDelete {
		test.Errorf("deleting pool test failed: %s", errorPoolDelete.Message)
	}
}

// Unsupported in libvirt version 10.0.0
func TestPoolDeleteClearAllToZeros(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	storagePoolFilesystem, errorGetStoragePoolFilesystemStruct, isErrorGetStoragePoolFilesystemStruct := poolConnection.helperGetStoragePoolFilesystemStruct(
		test,
		filesystemDiskDeviceValue,
	)
	if isErrorGetStoragePoolFilesystemStruct {
		test.Fatalf("get storage pool filesystem struct failed: %s", errorGetStoragePoolFilesystemStruct.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolFilesystem, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	errorPoolCreate, isErrorPoolCreate := poolConnection.helperTestPoolCreate(test, poolUuid, libvirt.STORAGE_POOL_CREATE_WITH_BUILD_OVERWRITE)
	if isErrorPoolCreate {
		test.Fatalf("build, create and starting pool then overwriting data inside directory failed: %s", errorPoolCreate.Message)
	}

	errorPoolDestroy, isErrorPoolDestroy := poolConnection.helperTestPoolDestroy(test, poolUuid)
	if isErrorPoolDestroy {
		test.Fatalf("pool destroy failed: %s", errorPoolDestroy.Message)
	}

	test.Cleanup(func() {
		if errorPoolUndefine, isErrorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}

		errorDeletePartition, isErrorDeletePartition := helperDepleteDevicePartition(test, filesystemDiskDeviceValue)
		if isErrorDeletePartition {
			test.Errorf("delete primary partition failed: %s", errorDeletePartition.Message)
		}
	})

	errorPoolDelete, isErrorPoolDelete := poolConnection.PoolDelete(poolUuid, libvirt.STORAGE_POOL_DELETE_ZEROED)
	if isErrorPoolDelete {
		test.Errorf("deleting pool test failed: %s", errorPoolDelete.Message)
	}
}
