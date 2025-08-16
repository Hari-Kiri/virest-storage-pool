package storagePool

import (
	"testing"

	"libvirt.org/go/libvirt"
)

func TestPoolDetail(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	errorPoolCreate, isErrorPoolCreate := poolConnection.helperTestPoolCreate(test, poolUuid, libvirt.STORAGE_POOL_CREATE_WITH_BUILD)
	if isErrorPoolCreate {
		test.Fatalf("creating pool test failed: %s", errorPoolCreate.Message)
	}

	test.Cleanup(func() {
		if errorPoolDestroy, isErrorPoolDestroy := poolConnection.PoolDestroy(poolUuid); isErrorPoolDestroy {
			test.Errorf("pool destroy test failed: %s", errorPoolDestroy.Message)
		}

		if errorPoolDelete, isErrorPoolDelete := poolConnection.PoolDelete(poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL); isErrorPoolDelete {
			test.Errorf("deleting pool test failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}
	})

	_, errorGetPoolDetail, isErrorGetPoolDetail := poolConnection.PoolDetail(poolUuid, 0)
	if isErrorGetPoolDetail {
		test.Errorf("get pool detail test failed: %s", errorGetPoolDetail.Message)
	}
}

func TestPoolDetailInactiveState(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	errorPoolCreate, isErrorPoolCreate := poolConnection.helperTestPoolCreate(test, poolUuid, libvirt.STORAGE_POOL_CREATE_WITH_BUILD)
	if isErrorPoolCreate {
		test.Fatalf("creating pool test failed: %s", errorPoolCreate.Message)
	}

	test.Cleanup(func() {
		if errorPoolDestroy, isErrorPoolDestroy := poolConnection.PoolDestroy(poolUuid); isErrorPoolDestroy {
			test.Errorf("pool destroy test failed: %s", errorPoolDestroy.Message)
		}

		if errorPoolDelete, isErrorPoolDelete := poolConnection.PoolDelete(poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL); isErrorPoolDelete {
			test.Errorf("deleting pool test failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}
	})

	_, errorGetPoolDetail, isErrorGetPoolDetail := poolConnection.PoolDetail(poolUuid, uint(libvirt.STORAGE_XML_INACTIVE))
	if isErrorGetPoolDetail {
		test.Errorf("get pool detail test failed: %s", errorGetPoolDetail.Message)
	}
}
