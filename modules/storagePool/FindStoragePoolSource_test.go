package storagePool

import (
	"testing"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/libvirtxml"
)

// Change it to Your's.
const (
	netfsHost = "192.168.122.4"
	iscsiHost = "192.168.122.6"
)

func TestFindStoragePoolSourceNetfs(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	test.Cleanup(func() {
		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}
	})

	srcSpec := libvirtxml.Source{
		Host: libvirtxml.Host{
			Name: netfsHost,
			Port: 2049,
		},
	}
	_, errorFindStoragePoolSource, isErrorFindStoragePoolSource := poolConnection.FindStoragePoolSource("netfs", srcSpec)
	if isErrorFindStoragePoolSource {
		test.Errorf("find storage pool netfs source test failed: %s", errorFindStoragePoolSource.Message)
	}
}

func TestFindStoragePoolSourceIscsi(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	test.Cleanup(func() {
		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}
	})

	srcSpec := libvirtxml.Source{
		Host: libvirtxml.Host{
			Name: iscsiHost,
			Port: 3260,
		},
	}
	_, errorFindStoragePoolSource, isErrorFindStoragePoolSource := poolConnection.FindStoragePoolSource("iscsi", srcSpec)
	if isErrorFindStoragePoolSource {
		test.Errorf("find storage pool iscsi source test failed: %s", errorFindStoragePoolSource.Message)
	}
}
