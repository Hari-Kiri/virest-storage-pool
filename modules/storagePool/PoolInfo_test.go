package storagePool

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolInfo"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

func (helper helperTest) helperTestPoolInfo(poolUuid string) (poolInfo.Info, virest.Error, bool) {
	helper.test.Helper()

	var incomingRequest poolInfo.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		helper.test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		fmt.Sprintf("/storage-pool/info?Uuid=%s", poolUuid),
		string(http.MethodGet),
		[]byte{},
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		helper.test.Fail()
		return poolInfo.Info{}, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition
	}
	defer poolConnection.Close()

	result, errorGetPoolInfo, isErrorGetPoolInfo := poolConnection.PoolInfo(incomingRequest.Uuid)
	if isErrorGetPoolInfo {
		helper.test.Fail()
		return poolInfo.Info{}, errorGetPoolInfo, isErrorGetPoolInfo
	}

	return result, virest.Error{}, false
}

func TestPoolInfo(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	storagePoolFilesystem, errorGetStoragePoolFilesystemStruct, isErrorGetStoragePoolFilesystemStruct := helper.helperGetStoragePoolFilesystemStruct(
		filesystemDiskDeviceValue,
	)
	if isErrorGetStoragePoolFilesystemStruct {
		test.Fatalf("get storage pool filesystem struct failed: %s", errorGetStoragePoolFilesystemStruct.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolFilesystem, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	errorPoolCreate, isErrorPoolCreate := helper.helperTestPoolCreate(poolUuid, libvirt.STORAGE_POOL_CREATE_WITH_BUILD_NO_OVERWRITE)
	if isErrorPoolCreate {
		test.Fatalf("creating pool test failed: %s", errorPoolCreate.Message)
	}

	var incomingRequest poolInfo.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		fmt.Sprintf("/storage-pool/info?Uuid=%s", poolUuid),
		string(http.MethodGet),
		[]byte{},
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition)
	}

	poolInfo, errorGetPoolInfo, isErrorGetPoolInfo := poolConnection.PoolInfo(incomingRequest.Uuid)
	if isErrorGetPoolInfo {
		test.Errorf("get pool info test failed: %s", errorGetPoolInfo.Message)
	}
	poolInfolMarshaled, errorMarshalingPoolInfo := json.MarshalIndent(poolInfo, "", "  ")
	if errorMarshalingPoolInfo != nil {
		test.Errorf("marshaling result failed: %s", errorMarshalingPoolInfo.Error())
	}
	test.Logf("storage pool %s info: \n%s", incomingRequest.Uuid, string(poolInfolMarshaled))

	test.Cleanup(func() {
		if errorPoolDestroy, isErrorPoolDestroy := helper.helperTestPoolDestroy(poolUuid); isErrorPoolDestroy {
			test.Errorf("pool destroy test failed: %s", errorPoolDestroy.Message)
		}

		if errorPoolDelete, isErrorPoolDelete := helper.helperTestPoolDelete(poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL); isErrorPoolDelete {
			test.Errorf("deleting pool test failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := helper.helperTestPoolUndefine(poolUuid); isErrorPoolUndefine {
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

func TestPoolInfoWrongUuid(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	storagePoolFilesystem, errorGetStoragePoolFilesystemStruct, isErrorGetStoragePoolFilesystemStruct := helper.helperGetStoragePoolFilesystemStruct(
		filesystemDiskDeviceValue,
	)
	if isErrorGetStoragePoolFilesystemStruct {
		test.Fatalf("get storage pool filesystem struct failed: %s", errorGetStoragePoolFilesystemStruct.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolFilesystem, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	errorPoolCreate, isErrorPoolCreate := helper.helperTestPoolCreate(poolUuid, libvirt.STORAGE_POOL_CREATE_WITH_BUILD_NO_OVERWRITE)
	if isErrorPoolCreate {
		test.Fatalf("creating pool test failed: %s", errorPoolCreate.Message)
	}

	var incomingRequest poolInfo.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		fmt.Sprintf("/storage-pool/info?Uuid=%s", "ca3e11bd-b840-4412-9f44-89b8d09b9f4e"),
		string(http.MethodGet),
		[]byte{},
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition)
	}

	poolInfo, errorGetPoolInfo, isErrorGetPoolInfo := poolConnection.PoolInfo(incomingRequest.Uuid)
	if isErrorGetPoolInfo {
		test.Logf("get pool info test failed: %s", errorGetPoolInfo.Message)
	}
	if !isErrorGetPoolInfo {
		test.Errorf("this test must be failed, but it won't: %v", poolInfo)
	}

	test.Cleanup(func() {
		if errorPoolDestroy, isErrorPoolDestroy := helper.helperTestPoolDestroy(poolUuid); isErrorPoolDestroy {
			test.Errorf("pool destroy test failed: %s", errorPoolDestroy.Message)
		}

		if errorPoolDelete, isErrorPoolDelete := helper.helperTestPoolDelete(poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL); isErrorPoolDelete {
			test.Errorf("deleting pool test failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := helper.helperTestPoolUndefine(poolUuid); isErrorPoolUndefine {
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

func TestPoolInfoUuidNotValid(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	storagePoolFilesystem, errorGetStoragePoolFilesystemStruct, isErrorGetStoragePoolFilesystemStruct := helper.helperGetStoragePoolFilesystemStruct(
		filesystemDiskDeviceValue,
	)
	if isErrorGetStoragePoolFilesystemStruct {
		test.Fatalf("get storage pool filesystem struct failed: %s", errorGetStoragePoolFilesystemStruct.Message)
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolFilesystem, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	errorPoolCreate, isErrorPoolCreate := helper.helperTestPoolCreate(poolUuid, libvirt.STORAGE_POOL_CREATE_WITH_BUILD_NO_OVERWRITE)
	if isErrorPoolCreate {
		test.Fatalf("creating pool test failed: %s", errorPoolCreate.Message)
	}

	var incomingRequest poolInfo.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		fmt.Sprintf("/storage-pool/info?Uuid=%s", "ca3e11bd-b840-4412-9f44-"),
		string(http.MethodGet),
		[]byte{},
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition)
	}

	poolInfo, errorGetPoolInfo, isErrorGetPoolInfo := poolConnection.PoolInfo(incomingRequest.Uuid)
	if isErrorGetPoolInfo {
		test.Logf("get pool info test failed: %s", errorGetPoolInfo.Message)
	}
	if !isErrorGetPoolInfo {
		test.Errorf("this test must be failed, but it won't: %v", poolInfo)
	}

	test.Cleanup(func() {
		if errorPoolDestroy, isErrorPoolDestroy := helper.helperTestPoolDestroy(poolUuid); isErrorPoolDestroy {
			test.Errorf("pool destroy test failed: %s", errorPoolDestroy.Message)
		}

		if errorPoolDelete, isErrorPoolDelete := helper.helperTestPoolDelete(poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL); isErrorPoolDelete {
			test.Errorf("deleting pool test failed: %s", errorPoolDelete.Message)
		}

		if errorPoolUndefine, isErrorPoolUndefine := helper.helperTestPoolUndefine(poolUuid); isErrorPoolUndefine {
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
