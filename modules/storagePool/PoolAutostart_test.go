package storagePool

import "testing"

const poolAutostartTestUuid = "850c4194-a85a-4a7c-9941-5b53c113ed5d"

func TestPoolAutostartTrue(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	test.Cleanup(func() {
		if errorCloseConnection := poolConnection.helperTestCloseConnection(test); errorCloseConnection != nil {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Error())
		}
	})

	errorPoolAutostart, isErrorPoolAutostart := poolConnection.PoolAutostart(poolAutostartTestUuid, true)
	if isErrorPoolAutostart {
		test.Errorf("turn on pool autostart failed: %s", errorPoolAutostart.Message)
	}
}

func TestPoolAutostartFalse(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	test.Cleanup(func() {
		if errorCloseConnection := poolConnection.helperTestCloseConnection(test); errorCloseConnection != nil {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Error())
		}
	})

	errorPoolAutostart, isErrorPoolAutostart := poolConnection.PoolAutostart(poolAutostartTestUuid, false)
	if isErrorPoolAutostart {
		test.Errorf("turn off pool autostart failed: %s", errorPoolAutostart.Message)
	}
}
