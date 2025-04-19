package storagePool

import (
	"testing"
)

func TestPoolCapabilities(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	test.Cleanup(func() {
		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}
	})

	_, errorPoolCapabilities, isErrorPoolCapabilities := poolConnection.PoolCapabilities()
	if isErrorPoolCapabilities {
		test.Errorf("get storage pool capabilities test failed: %s", errorPoolCapabilities.Message)
	}
}
