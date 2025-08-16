package storagePool

import (
	"testing"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

func TestPoolBuildFromScratch(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	test.Cleanup(func() {
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

	errorPoolBuild, isErrorPoolBuild := poolConnection.PoolBuild(poolUuid, libvirt.STORAGE_POOL_BUILD_NEW)
	if isErrorPoolBuild {
		test.Errorf("build pool from scratch test failed: %s", errorPoolBuild.Message)
	}
}

// Unsupported in libvirt version 10.0.0 and virsh (the libvirt command line interface) still not implement this method.
func TestPoolBuildRepairOrReinitilize(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	test.Cleanup(func() {
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

	errorPoolBuildRepairOrReinitialize, isErrorPoolBuildRepairOrReinitialize := poolConnection.PoolBuild(poolUuid, libvirt.STORAGE_POOL_BUILD_REPAIR)
	if isErrorPoolBuildRepairOrReinitialize {
		test.Errorf("repair or reinitilize pool test failed: %s", errorPoolBuildRepairOrReinitialize.Message)
	}
}

// Unsupported in libvirt version 10.0.0 and virsh (the libvirt command line interface) still not implement this method.
func TestPoolBuildExtendExistingPool(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	test.Cleanup(func() {
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

	errorPoolBuildExtendExistingPool, isErrorPoolBuildExtendExistingPool := poolConnection.PoolBuild(poolUuid, libvirt.STORAGE_POOL_BUILD_RESIZE)
	if isErrorPoolBuildExtendExistingPool {
		test.Errorf("extend existing pool test failed: %s", errorPoolBuildExtendExistingPool.Message)
	}
}

// Currently only filesystem pool accepts flags VIR_STORAGE_POOL_BUILD_OVERWRITE (option 4).
func TestPoolBuildFilesystemNotOverwriteExistingPool(test *testing.T) {
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

	errorPoolBuildFilesystemNotOverwriteExistingPool, isErrorPoolBuildFilesystemNotOverwriteExistingPool := poolConnection.PoolBuild(
		poolUuid,
		libvirt.STORAGE_POOL_BUILD_NO_OVERWRITE,
	)
	if isErrorPoolBuildFilesystemNotOverwriteExistingPool {
		test.Errorf("build not overwrite existing filesystem pool test failed: %s", errorPoolBuildFilesystemNotOverwriteExistingPool.Message)
	}
}

// Currently only filesystem pool accepts flags VIR_STORAGE_POOL_BUILD_NO_OVERWRITE (option 8).
func TestPoolBuildFilesystemOverwriteData(test *testing.T) {
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
		if errorPoolDelete, isErrorPoolDelete := poolConnection.helperTestPoolDelete(test, poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL); isErrorPoolDelete {
			test.Errorf("pool delete failed: %s", errorPoolDelete.Message)
		}

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

	errorPoolBuildFilesystemOverwriteData, isErrorPoolBuildFilesystemOverwriteData := poolConnection.PoolBuild(poolUuid, libvirt.STORAGE_POOL_BUILD_OVERWRITE)
	if isErrorPoolBuildFilesystemOverwriteData {
		test.Errorf("build and overwrite existing filesystem pool test failed: %s", errorPoolBuildFilesystemOverwriteData.Message)
	}
}

func (poolConnection *poolConnection) helperTestPoolBuild(test *testing.T, poolUuid string, option libvirt.StoragePoolBuildFlags) (virest.Error, bool) {
	test.Helper()

	errorPoolBuild, isErrorPoolBuild := poolConnection.PoolBuild(poolUuid, option)
	if isErrorPoolBuild {
		test.Fail()
		return errorPoolBuild, isErrorPoolBuild
	}

	return virest.Error{}, false
}
