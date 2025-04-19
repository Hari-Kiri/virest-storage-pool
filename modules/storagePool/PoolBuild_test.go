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

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolDirectory, 1)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	test.Cleanup(func() {
		if errorPoolDelete, isErrorPoolDelete := poolConnection.helperTestPoolDelete(test, poolUuid, 0); isErrorPoolDelete {
			test.Errorf("pool delete failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}
	})

	errorPoolBuild, isErrorPoolBuild := poolConnection.PoolBuild(poolUuid, 0)
	if isErrorPoolBuild {
		test.Errorf("build pool from scratch test failed: %s", errorPoolBuild.Message)
	}
}

// Unsupported in libvirt version 10.0.0
func TestPoolBuildRepairOrReinitilize(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolDirectory, 1)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	test.Cleanup(func() {
		if errorPoolDelete, isErrorPoolDelete := poolConnection.helperTestPoolDelete(test, poolUuid, 0); isErrorPoolDelete {
			test.Errorf("pool delete failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}
	})

	errorPoolBuildRepairOrReinitialize, isErrorPoolBuildRepairOrReinitialize := poolConnection.PoolBuild(poolUuid, 1)
	if isErrorPoolBuildRepairOrReinitialize {
		test.Errorf("repair or reinitilize pool test failed: %s", errorPoolBuildRepairOrReinitialize.Message)
	}
}

// Unsupported in libvirt version 10.0.0
func TestPoolBuildExtendExistingPool(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolDirectory, 1)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	test.Cleanup(func() {
		if errorPoolDelete, isErrorPoolDelete := poolConnection.helperTestPoolDelete(test, poolUuid, 0); isErrorPoolDelete {
			test.Errorf("pool delete failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}
	})

	errorPoolBuildExtendExistingPool, isErrorPoolBuildExtendExistingPool := poolConnection.PoolBuild(poolUuid, 2)
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

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolFilesystem, 1)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	test.Cleanup(func() {
		if errorPoolDelete, isErrorPoolDelete := poolConnection.helperTestPoolDelete(test, poolUuid, 0); isErrorPoolDelete {
			test.Errorf("pool delete failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}
	})

	errorPoolBuildFilesystemNotOverwriteExistingPool, isErrorPoolBuildFilesystemNotOverwriteExistingPool := poolConnection.PoolBuild(poolUuid, 4)
	if isErrorPoolBuildFilesystemNotOverwriteExistingPool {
		test.Errorf("overwrite existing filesystem pool test failed: %s", errorPoolBuildFilesystemNotOverwriteExistingPool.Message)
	}
}

// Currently only filesystem pool accepts flags VIR_STORAGE_POOL_BUILD_NO_OVERWRITE (option 8).
func TestPoolBuildFilesystemOverwriteData(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolFilesystem, 1)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	test.Cleanup(func() {
		if errorPoolDelete, isErrorPoolDelete := poolConnection.helperTestPoolDelete(test, poolUuid, 0); isErrorPoolDelete {
			test.Errorf("pool delete failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}
	})

	errorPoolBuildFilesystemOverwriteData, isErrorPoolBuildFilesystemOverwriteData := poolConnection.PoolBuild(poolUuid, 8)
	if isErrorPoolBuildFilesystemOverwriteData {
		test.Errorf("overwrite existing filesystem pool test failed: %s", errorPoolBuildFilesystemOverwriteData.Message)
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

func (poolConnection *poolConnection) helperTestPoolDelete(test *testing.T, poolUuid string, option libvirt.StoragePoolDeleteFlags) (virest.Error, bool) {
	test.Helper()

	errorPoolDelete, isErrorPoolDelete := poolConnection.PoolDelete(poolUuid, option)
	if isErrorPoolDelete {
		test.Fail()
		return errorPoolDelete, isErrorPoolDelete
	}

	return virest.Error{}, false
}

func (poolConnection *poolConnection) helperTestPoolUndefine(test *testing.T, poolUuid string) (virest.Error, bool) {
	test.Helper()

	errorPoolUndefine, isErrorPoolUndefine := poolConnection.PoolUndefine(poolUuid)
	if isErrorPoolUndefine {
		test.Fail()
		return errorPoolUndefine, isErrorPoolUndefine
	}

	return virest.Error{}, false
}
