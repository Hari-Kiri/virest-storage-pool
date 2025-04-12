package storagePool

import (
	"testing"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

var storagePoolDirectory = libvirtxml.StoragePool{
	Type: "dir",
	Name: "unit-test-pool-directory",
	Target: &libvirtxml.StoragePoolTarget{
		Path: "/home/dexip/unit-test",
		Permissions: &libvirtxml.StoragePoolTargetPermissions{
			Mode:  "0755",
			Owner: "1000",
			Group: "1000",
		},
	},
}

var storagePoolFilesystem = libvirtxml.StoragePool{
	Type: "fs",
	Name: "unit-test-pool-filesystem",
	Source: &libvirtxml.StoragePoolSource{
		Device: []libvirtxml.StoragePoolSourceDevice{
			{
				Path: "/dev/vdb",
			},
		},
	},
	Target: &libvirtxml.StoragePoolTarget{
		Path: "/mnt/unit-test-pool-filesystem",
	},
}

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

func (poolConnection *poolConnection) helperTestPoolDefine(test *testing.T, storagePool libvirtxml.StoragePool, option libvirt.StoragePoolDefineFlags) (string, virest.Error, bool) {
	test.Helper()

	result, errorPoolDefine, isErrorPoolDevine := poolConnection.PoolDefine(storagePool, option)
	if isErrorPoolDevine {
		test.Fail()
		return "", errorPoolDefine, isErrorPoolDevine
	}

	return result.Uuid, virest.Error{}, false
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
