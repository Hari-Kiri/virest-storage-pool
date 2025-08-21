package storagePool

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDefine"
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

func (helper helperTest) helperTestPoolDefine(storagePool libvirtxml.StoragePool, option libvirt.StoragePoolDefineFlags) (string, virest.Error, bool) {
	helper.test.Helper()

	requestData := poolDefine.Request{
		Option:      option,
		StoragePool: storagePool,
	}
	var (
		virestError virest.Error
		isError     bool
	)
	requestBody, errorMarshalStoragePool := json.Marshal(requestData)
	virestError.Error, isError = errorMarshalStoragePool.(libvirt.Error)
	if isError {
		helper.test.Fail()
		return "", virestError, isError
	}

	var incomingRequest poolDefine.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		helper.test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/define",
		http.MethodPost,
		requestBody,
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		helper.test.Fail()
		return "", errorHttpRequestPrecondition, isErrorHttpRequestPrecondition
	}
	defer poolConnection.Close()

	result, errorPoolDefine, isErrorPoolDevine := poolConnection.PoolDefine(incomingRequest.StoragePool, incomingRequest.Option)
	if isErrorPoolDevine {
		helper.test.Fail()
		return "", errorPoolDefine, isErrorPoolDevine
	}

	return result.Uuid, virest.Error{}, false
}

func (helper helperTest) helperGetStoragePoolFilesystemStruct(diskDevice filesystemDiskDevice) (libvirtxml.StoragePool, virest.Error, bool) {
	helper.test.Helper()

	errorCreateNewSinglePrimaryPartitionOnInternalDiskDevice,
		isErrorCreateNewSinglePrimaryPartitionOnInternalDiskDevice := utils.CreateNewSinglePrimaryPartitionOnInternalDiskDevice(
		diskDevice.path,
		diskDevice.format,
		diskDevice.partitionTable,
	)
	if isErrorCreateNewSinglePrimaryPartitionOnInternalDiskDevice {
		helper.test.Fail()
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

func (helper helperTest) helperDepleteDevicePartition(diskDevice filesystemDiskDevice) (virest.Error, bool) {
	helper.test.Helper()

	errorDeletePrimaryPartition, isErrorDeletePartition := utils.DepleteDevicePartition(diskDevice.path, diskDevice.format)
	if isErrorDeletePartition {
		helper.test.Fail()
		return errorDeletePrimaryPartition, true
	}

	return virest.Error{}, false
}

func TestPoolDefineNoOption(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	requestData := poolDefine.Request{
		Option:      0,
		StoragePool: storagePoolDirectory,
	}
	var (
		virestError virest.Error
		isError     bool
	)
	requestBody, errorMarshalStoragePool := json.Marshal(requestData)
	virestError.Error, isError = errorMarshalStoragePool.(libvirt.Error)
	if isError {
		test.Fatalf("marshaling request body failed: %s", virestError.Message)
	}

	var incomingRequest poolDefine.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/define",
		http.MethodPost,
		requestBody,
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	poolDefine, errorPoolDefine, isErrorPoolDefine := poolConnection.PoolDefine(incomingRequest.StoragePool, incomingRequest.Option)
	if isErrorPoolDefine {
		test.Errorf("defining pool test failed: %s", errorPoolDefine.Message)
	}

	test.Cleanup(func() {
		if errorPoolUndefine, isErrorPoolUndefine := helper.helperTestPoolUndefine(poolDefine.Uuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}
	})
}

func TestPoolDefineOptionValidate(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	storagePoolFilesystem, errorGetStoragePoolFilesystemStruct, isErrorGetStoragePoolFilesystemStruct := helper.helperGetStoragePoolFilesystemStruct(
		filesystemDiskDeviceValue,
	)
	if isErrorGetStoragePoolFilesystemStruct {
		test.Fatalf("get storage pool filesystem struct failed: %s", errorGetStoragePoolFilesystemStruct.Message)
	}

	requestData := poolDefine.Request{
		Option:      libvirt.STORAGE_POOL_DEFINE_VALIDATE,
		StoragePool: storagePoolFilesystem,
	}
	var (
		virestError virest.Error
		isError     bool
	)
	requestBody, errorMarshalStoragePool := json.Marshal(requestData)
	virestError.Error, isError = errorMarshalStoragePool.(libvirt.Error)
	if isError {
		test.Fatalf("marshaling request body failed: %s", virestError.Message)
	}

	var incomingRequest poolDefine.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/define",
		http.MethodPost,
		requestBody,
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	poolDefine, errorPoolDefine, isErrorPoolDefine := poolConnection.PoolDefine(incomingRequest.StoragePool, incomingRequest.Option)
	if isErrorPoolDefine {
		test.Errorf("defining pool test failed: %s", errorPoolDefine.Message)
	}

	test.Cleanup(func() {
		if errorPoolUndefine, isErrorPoolUndefine := helper.helperTestPoolUndefine(poolDefine.Uuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}

		if errorDeletePartition, isErrorDeletePartition := helper.helperDepleteDevicePartition(filesystemDiskDeviceValue); isErrorDeletePartition {
			test.Errorf("delete primary partition failed: %s", errorDeletePartition.Message)
		}
	})
}
