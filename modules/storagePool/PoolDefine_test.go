package storagePool

import (
	"os"
	"strings"
	"testing"

	"github.com/Hari-Kiri/virest-utilities/utils"
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

type filesystemDiskDevice struct {
	path           string
	format         string
	partitionTable string
}

var filesystemDiskDeviceValue = filesystemDiskDevice{
	path:           "/dev/vdb",
	format:         "raw",
	partitionTable: "gpt",
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

func (poolConnection *poolConnection) helperGetStoragePoolFilesystemStruct(test *testing.T, diskDevice filesystemDiskDevice) (libvirtxml.StoragePool, virest.Error, bool) {
	test.Helper()

	errorCreateNewSinglePrimaryPartitionOnInternalDiskDevice,
		isErrorCreateNewSinglePrimaryPartitionOnInternalDiskDevice := utils.CreateNewSinglePrimaryPartitionOnInternalDiskDevice(
		diskDevice.path,
		diskDevice.format,
		diskDevice.partitionTable,
	)
	if isErrorCreateNewSinglePrimaryPartitionOnInternalDiskDevice {
		test.Fail()
		return libvirtxml.StoragePool{}, errorCreateNewSinglePrimaryPartitionOnInternalDiskDevice, true
	}

	sysBlockDir, errorReadContentSysBlockDir := os.ReadDir("/sys/block")
	if errorReadContentSysBlockDir != nil {
		panic(errorReadContentSysBlockDir)
	}

	var (
		link         string
		errorGetLink error
	)
	diskDevicePathSplit := strings.Split(diskDevice.path, "/")
	blockDevice := diskDevicePathSplit[len(diskDevicePathSplit)-1]
	for i := 0; i < len(sysBlockDir); i++ {
		if sysBlockDir[i].Name() == blockDevice {
			link, errorGetLink = os.Readlink("/sys/block/" + sysBlockDir[i].Name())
		}
	}
	if errorGetLink != nil {
		panic(errorGetLink)
	}

	blockDeviceDir, errorGetBlockDeviceDir := os.ReadDir("/sys/" + link[2:])
	if errorGetBlockDeviceDir != nil {
		panic(errorGetBlockDeviceDir)
	}
	var partitionPath string
	for i := 0; i < len(blockDeviceDir); i++ {
		if strings.HasPrefix(blockDeviceDir[i].Name(), blockDevice) {
			partitionPath = "/dev/" + blockDeviceDir[i].Name()
		}
	}

	return libvirtxml.StoragePool{
		Type: "fs",
		Name: "unit-test-pool-filesystem",
		Source: &libvirtxml.StoragePoolSource{
			Device: []libvirtxml.StoragePoolSourceDevice{
				{
					Path: partitionPath,
				},
			},
			Format: &libvirtxml.StoragePoolSourceFormat{
				Type: "xfs",
			},
		},
		Target: &libvirtxml.StoragePoolTarget{
			Path: "/mnt/unit-test-pool-filesystem",
		},
	}, virest.Error{}, false
}

func helperDepleteDevicePartition(test *testing.T, diskDevice filesystemDiskDevice) (virest.Error, bool) {
	test.Helper()

	errorDeletePrimaryPartition, isErrorDeletePartition := utils.DepleteDevicePartition(diskDevice.path, diskDevice.format)
	if isErrorDeletePartition {
		test.Fail()
		return errorDeletePrimaryPartition, true
	}

	return virest.Error{}, false
}

func TestPoolDefineNoOption(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := poolConnection.helperTestPoolDefine(test, storagePoolDirectory, 0)
	if isErrorPoolDefine {
		test.Errorf("defining pool test failed: %s", errorPoolDefine.Message)
	}

	test.Cleanup(func() {
		if errorPoolUndefine, isErrorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		if errorCloseConnection, isErrorCloseConnection := poolConnection.helperTestCloseConnection(test); isErrorCloseConnection {
			test.Errorf("connection close() failed: %s", errorCloseConnection.Message)
		}
	})
}

func TestPoolDefineOptionValidate(test *testing.T) {
	poolConnection, errorGetPoolConnection, isErrorGetPoolConnection := helperTestConnection(test)
	if isErrorGetPoolConnection {
		test.Fatalf("connecting to host storage pool failed: %s", errorGetPoolConnection.Message)
	}

	storagePoolFilesystem, errorGetStoragePoolFilesystemStruct, isErrorGetStoragePoolFilesystemStruct := poolConnection.helperGetStoragePoolFilesystemStruct(
		test,
		filesystemDiskDeviceValue,
	)
	if isErrorGetStoragePoolFilesystemStruct {
		test.Fatalf("get storage pool filesystem struct failed: %s", errorGetStoragePoolFilesystemStruct.Message)
	}

	poolDefine, errorPoolDefine, isErrorPoolDefine := poolConnection.PoolDefine(storagePoolFilesystem, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Errorf("defining pool test failed: %s", errorPoolDefine.Message)
	}

	test.Cleanup(func() {
		if errorPoolUndefine, isErrorPoolUndefine := poolConnection.helperTestPoolUndefine(test, poolDefine.Uuid); isErrorPoolUndefine {
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
}
