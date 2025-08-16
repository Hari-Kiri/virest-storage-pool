package storagePool

import (
	"testing"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

func (poolConnection *poolConnection) helperTestPoolCreate(test *testing.T, poolUuid string, option libvirt.StoragePoolCreateFlags) (virest.Error, bool) {
	test.Helper()

	errorPoolCreate, isErrorPoolCreate := poolConnection.PoolCreate(poolUuid, option)
	if isErrorPoolCreate {
		test.Fail()
		return errorPoolCreate, isErrorPoolCreate
	}

	return virest.Error{}, false
}

func TestPoolCreateActionStartingPool(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	errorPoolBuild, isErrorPoolBuild := poolConnection.helperTestPoolBuild(test, poolUuid, 0)
	if isErrorPoolBuild {
		test.Fatalf("pool build failed: %s", errorPoolBuild.Message)
	}

	test.Cleanup(func() {
		if errorPoolDestroy, isErrorPoolDestroy := poolConnection.helperTestPoolDestroy(test, poolUuid); isErrorPoolDestroy {
			test.Errorf("pool destroy failed: %s", errorPoolDestroy.Message)
		}

		if errorPoolDelete, isErrorPoolDelete := poolConnection.helperTestPoolDelete(test, poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL); isErrorPoolDelete {
			test.Errorf("pool delete failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}
	})

	errorPoolCreate, isErrorPoolCreate := poolConnection.PoolCreate(poolUuid, libvirt.STORAGE_POOL_CREATE_NORMAL)
	if isErrorPoolCreate {
		test.Errorf("starting pool test failed: %s", errorPoolCreate.Message)
	}
}

func TestPoolCreateActionBuildCreateAndStartingPool(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	test.Cleanup(func() {
		if errorPoolDestroy, isErrorPoolDestroy := poolConnection.helperTestPoolDestroy(test, poolUuid); isErrorPoolDestroy {
			test.Errorf("pool destroy failed: %s", errorPoolDestroy.Message)
		}

		if errorPoolDelete, isErrorPoolDelete := poolConnection.helperTestPoolDelete(test, poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL); isErrorPoolDelete {
			test.Errorf("pool delete failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}
	})

	errorPoolCreate, isErrorPoolCreate := poolConnection.PoolCreate(poolUuid, libvirt.STORAGE_POOL_CREATE_WITH_BUILD)
	if isErrorPoolCreate {
		test.Errorf("build, create and starting pool test failed: %s", errorPoolCreate.Message)
	}
}

func TestPoolCreateActionBuildCreateAndStartingPoolNoOverwriteDataInDirectory(test *testing.T) {
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

	test.Cleanup(func() {
		if errorPoolDestroy, isErrorPoolDestroy := poolConnection.helperTestPoolDestroy(test, poolUuid); isErrorPoolDestroy {
			test.Errorf("pool destroy failed: %s", errorPoolDestroy.Message)
		}

		if errorPoolDelete, isErrorPoolDelete := poolConnection.helperTestPoolDelete(test, poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL); isErrorPoolDelete {
			test.Errorf("pool delete failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}

		if errorDeletePartition, isErrorDeletePartition := helperDepleteDevicePartition(test, filesystemDiskDeviceValue); isErrorDeletePartition {
			test.Errorf("delete primary partition failed: %s", errorDeletePartition.Message)
		}
	})

	errorPoolCreate, isErrorPoolCreate := poolConnection.PoolCreate(poolUuid, libvirt.STORAGE_POOL_CREATE_WITH_BUILD_NO_OVERWRITE)
	if isErrorPoolCreate {
		test.Errorf("build, create and starting pool then no overwriting data inside directory test failed: %s", errorPoolCreate.Message)
	}
}

func TestPoolCreateActionBuildCreateAndStartingPoolOverwriteDataInDirectory(test *testing.T) {
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

	test.Cleanup(func() {
		if errorPoolDestroy, isErrorPoolDestroy := poolConnection.helperTestPoolDestroy(test, poolUuid); isErrorPoolDestroy {
			test.Errorf("pool destroy failed: %s", errorPoolDestroy.Message)
		}

		if errorPoolDelete, isErrorPoolDelete := poolConnection.helperTestPoolDelete(test, poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL); isErrorPoolDelete {
			test.Errorf("pool delete failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}

		if errorDeletePartition, isErrorDeletePartition := helperDepleteDevicePartition(test, filesystemDiskDeviceValue); isErrorDeletePartition {
			test.Errorf("delete primary partition failed: %s", errorDeletePartition.Message)
		}
	})

	errorPoolCreate, isErrorPoolCreate := poolConnection.PoolCreate(poolUuid, libvirt.STORAGE_POOL_CREATE_WITH_BUILD_OVERWRITE)
	if isErrorPoolCreate {
		test.Errorf("build, create and starting pool then overwriting data inside directory test failed: %s", errorPoolCreate.Message)
	}
}
