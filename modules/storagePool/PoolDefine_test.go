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
		Path: "/home/dexip/unit-test-pool-directory",
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
				Path: "/dev/vdb1",
			},
		},
		Format: &libvirtxml.StoragePoolSourceFormat{
			Type: "xfs",
		},
	},
	Target: &libvirtxml.StoragePoolTarget{
		Path: "/mnt/unit-test-pool-filesystem",
	},
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
