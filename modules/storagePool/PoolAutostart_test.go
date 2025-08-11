package storagePool

import "testing"

func TestPoolAutostartTrue(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolDirectory, 0)
	if isErrorPoolDefine {
		test.Fatalf("defining pool test failed: %s", errorPoolDefine.Message)
	}

	test.Cleanup(func() {
		if errorPoolUndefine, isErrorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}
	})

	errorPoolAutostart, isErrorPoolAutostart := poolConnection.PoolAutostart(poolUuid, true)
	if isErrorPoolAutostart {
		test.Errorf("turn on pool autostart test failed: %s", errorPoolAutostart.Message)
	}
}

func TestPoolAutostartFalse(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolDirectory, 0)
	if isErrorPoolDefine {
		test.Fatalf("defining pool test failed: %s", errorPoolDefine.Message)
	}

	errorPoolAutostart, isErrorPoolAutostart := poolConnection.PoolAutostart(poolUuid, true)
	if isErrorPoolAutostart {
		test.Fatalf("turn on pool autostart test failed: %s", errorPoolAutostart.Message)
	}

	test.Cleanup(func() {
		if errorPoolUndefine, isErrorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}
	})

	errorPoolDisableAutostart, isErrorPoolDisableAutostart := poolConnection.PoolAutostart(poolUuid, false)
	if isErrorPoolDisableAutostart {
		test.Errorf("turn off pool autostart test failed: %s", errorPoolDisableAutostart.Message)
	}
}
