package storagePool

import (
	"testing"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

func (poolConnection *poolConnection) helperTestPoolDestroy(test *testing.T, poolUuid string) (virest.Error, bool) {
	test.Helper()

	errorPoolDestroy, isErrorPoolDestroy := poolConnection.PoolDestroy(poolUuid)
	if isErrorPoolDestroy {
		test.Fail()
		return errorPoolDestroy, isErrorPoolDestroy
	}

	return virest.Error{}, false
}

func TestPoolDestroy(test *testing.T) {
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
		if errorPoolDelete, isErrorPoolDelete := poolConnection.PoolDelete(poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL); isErrorPoolDelete {
			test.Errorf("deleting pool test failed: %s", errorPoolDelete.Message)
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

	errorPoolDestroy, isErrorPoolDestroy := poolConnection.PoolDestroy(poolUuid)
	if isErrorPoolDestroy {
		test.Errorf("pool destroy test failed: %s", errorPoolDestroy.Message)
	}
}
