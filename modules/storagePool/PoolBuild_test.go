package storagePool

import (
	"testing"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

var storagePool = libvirtxml.StoragePool{
	Type: "dir",
	Name: "unit-test",
	Target: &libvirtxml.StoragePoolTarget{
		Path: "/home/dexip/unit-test",
		Permissions: &libvirtxml.StoragePoolTargetPermissions{
			Mode:  "0755",
			Owner: "1000",
			Group: "1000",
		},
	},
}

func TestPoolBuildFromScratch(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePool, 1)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	test.Cleanup(func() {
		if errorPoolDelete := poolConnection.helperTestPoolDelete(test, poolUuid, 0); errorPoolDelete != nil {
			test.Errorf("pool delete failed: %s", errorPoolDelete.Error())
		}

		if errorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolUuid); errorPoolUndefine != nil {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Error())
		}

		if errorCloseConnection := poolConnection.helperTestCloseConnection(test); errorCloseConnection != nil {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Error())
		}
	})

	errorPoolBuild, isErrorPoolBuild := poolConnection.PoolBuild(poolUuid, 0)
	if isErrorPoolBuild {
		test.Errorf("build pool from scratch test failed: %s", errorPoolBuild.Message)
	}
}

func (poolConnection *poolConnection) helperTestPoolDefine(test *testing.T, storagePool libvirtxml.StoragePool, option libvirt.StoragePoolDefineFlags) (string, virest.Error, bool) {
	test.Helper()

	result, errorPoolDefine, isErrorPoolDevine := poolConnection.PoolDefine(storagePool, option)
	if isErrorPoolDevine {
		test.Fail()
		return "", errorPoolDefine, isErrorPoolDevine
	}

	return result.Uuid, virest.Error{}, false
}

func (poolConnection *poolConnection) helperTestPoolDelete(test *testing.T, poolUuid string, option libvirt.StoragePoolDeleteFlags) error {
	test.Helper()

	errorPoolDelete, isErrorPoolDelete := poolConnection.PoolDelete(poolUuid, option)
	if isErrorPoolDelete {
		test.Fail()
		return errorPoolDelete.Error
	}

	return nil
}

func (poolConnection *poolConnection) helperTestPoolUndefine(test *testing.T, poolUuid string) error {
	test.Helper()

	errorPoolUndefine, isErrorPoolUndefine := poolConnection.PoolUndefine(poolUuid)
	if isErrorPoolUndefine {
		test.Fail()
		return errorPoolUndefine.Error
	}

	return nil
}
