package storagePool

import (
	"testing"
)

func TestPoolCreateActionStartingPool(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolDirectory, 1)
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

	errorPoolCreate, isErrorPoolCreate := poolConnection.PoolCreate(poolUuid, 0)
	if isErrorPoolCreate {
		test.Errorf("starting pool test failed: %s", errorPoolCreate.Message)
	}
}

func TestPoolCreateActionBuildCreateAndStartingPool(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolDirectory, 1)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	test.Cleanup(func() {
		if errorPoolDestroy, isErrorPoolDestroy := poolConnection.helperTestPoolDestroy(test, poolUuid); isErrorPoolDestroy {
			test.Errorf("pool destroy failed: %s", errorPoolDestroy.Message)
		}

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

	errorPoolCreate, isErrorPoolCreate := poolConnection.PoolCreate(poolUuid, 1)
	if isErrorPoolCreate {
		test.Errorf("starting pool test failed: %s", errorPoolCreate.Message)
	}
}

func TestPoolCreateActionBuildCreateAndStartingPoolOverwriteDataInDirectory(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolDirectory, 1)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	test.Cleanup(func() {
		if errorPoolDestroy, isErrorPoolDestroy := poolConnection.helperTestPoolDestroy(test, poolUuid); isErrorPoolDestroy {
			test.Errorf("pool destroy failed: %s", errorPoolDestroy.Message)
		}

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

	errorPoolCreate, isErrorPoolCreate := poolConnection.PoolCreate(poolUuid, 2)
	if isErrorPoolCreate {
		test.Errorf("starting pool test failed: %s", errorPoolCreate.Message)
	}
}
