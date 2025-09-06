package storagePool

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDetail"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"libvirt.org/go/libvirt"
)

func TestPoolDetail(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	errorPoolCreate, isErrorPoolCreate := helper.helperTestPoolCreate(poolUuid, libvirt.STORAGE_POOL_CREATE_WITH_BUILD)
	if isErrorPoolCreate {
		test.Fatalf("creating pool test failed: %s", errorPoolCreate.Message)
	}

	var incomingRequest poolDetail.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		fmt.Sprintf("/storage-pool/detail?Uuid=%s&Option=%d", poolUuid, 0),
		string(http.MethodGet),
		[]byte{},
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition)
	}

	option, errorParseOptionToUint, isErrorParseOptionToUint := utils.StringToUint(incomingRequest.Option)
	if isErrorParseOptionToUint {
		test.Fatalf("parsing option failed: %s", errorParseOptionToUint)
	}

	poolDetail, errorGetPoolDetail, isErrorGetPoolDetail := poolConnection.PoolDetail(incomingRequest.Uuid, option)
	if isErrorGetPoolDetail {
		test.Errorf("get pool detail test failed: %s", errorGetPoolDetail.Message)
	}
	poolDetailMarshaled, errorMarshalingPoolDetail := json.MarshalIndent(poolDetail, "", "  ")
	if errorMarshalingPoolDetail != nil {
		test.Errorf("marshaling result failed: %s", errorMarshalingPoolDetail.Error())
	}
	test.Logf("storage pool %s detail: \n%s", incomingRequest.Uuid, string(poolDetailMarshaled))

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
	})
}

func TestPoolDetailInactiveState(test *testing.T) {
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

	var incomingRequest poolDetail.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		fmt.Sprintf("/storage-pool/detail?Uuid=%s&Option=%d", poolUuid, libvirt.STORAGE_XML_INACTIVE),
		string(http.MethodGet),
		[]byte{},
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition)
	}

	option, errorParseOptionToUint, isErrorParseOptionToUint := utils.StringToUint(incomingRequest.Option)
	if isErrorParseOptionToUint {
		test.Fatalf("parsing option failed: %s", errorParseOptionToUint)
	}

	poolDetail, errorGetPoolDetail, isErrorGetPoolDetail := poolConnection.PoolDetail(incomingRequest.Uuid, option)
	if isErrorGetPoolDetail {
		test.Errorf("get pool detail test failed: %s", errorGetPoolDetail.Message)
	}
	poolDetailMarshaled, errorMarshalingPoolDetail := json.MarshalIndent(poolDetail, "", "  ")
	if errorMarshalingPoolDetail != nil {
		test.Errorf("marshaling result failed: %s", errorMarshalingPoolDetail.Error())
	}
	test.Logf("storage pool %s detail: \n%s", incomingRequest.Uuid, string(poolDetailMarshaled))

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

func TestPoolDetailWrongUuid(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	errorPoolCreate, isErrorPoolCreate := helper.helperTestPoolCreate(poolUuid, libvirt.STORAGE_POOL_CREATE_WITH_BUILD)
	if isErrorPoolCreate {
		test.Fatalf("creating pool test failed: %s", errorPoolCreate.Message)
	}

	var incomingRequest poolDetail.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		fmt.Sprintf("/storage-pool/detail?Uuid=%s&Option=%d", "ca3e11bd-b840-4412-9f44-89b8d09b9f4e", 0),
		string(http.MethodGet),
		[]byte{},
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition)
	}

	option, errorParseOptionToUint, isErrorParseOptionToUint := utils.StringToUint(incomingRequest.Option)
	if isErrorParseOptionToUint {
		test.Fatalf("parsing option failed: %s", errorParseOptionToUint)
	}

	_, errorGetPoolDetail, isErrorGetPoolDetail := poolConnection.PoolDetail(incomingRequest.Uuid, option)
	if isErrorGetPoolDetail {
		test.Logf("get pool detail test failed: %s", errorGetPoolDetail.Message)
	}
	if !isErrorGetPoolDetail {
		test.Error("this test must be failed, but it won't")
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
	})
}

func TestPoolDetailUuidNotValid(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("pool define failed: %s", errorPoolDefine.Message)
	}

	errorPoolCreate, isErrorPoolCreate := helper.helperTestPoolCreate(poolUuid, libvirt.STORAGE_POOL_CREATE_WITH_BUILD)
	if isErrorPoolCreate {
		test.Fatalf("creating pool test failed: %s", errorPoolCreate.Message)
	}

	var incomingRequest poolDetail.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		fmt.Sprintf("/storage-pool/detail?Uuid=%s&Option=%d", "ca3e11bd-b840-4412-9f44-", 0),
		string(http.MethodGet),
		[]byte{},
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition)
	}

	option, errorParseOptionToUint, isErrorParseOptionToUint := utils.StringToUint(incomingRequest.Option)
	if isErrorParseOptionToUint {
		test.Fatalf("parsing option failed: %s", errorParseOptionToUint)
	}

	_, errorGetPoolDetail, isErrorGetPoolDetail := poolConnection.PoolDetail(incomingRequest.Uuid, option)
	if isErrorGetPoolDetail {
		test.Logf("get pool detail test failed: %s", errorGetPoolDetail.Message)
	}
	if !isErrorGetPoolDetail {
		test.Error("this test must be failed, but it won't")
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
	})
}

func TestPoolDetailInactiveStateWrongUuid(test *testing.T) {
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

	var incomingRequest poolDetail.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		fmt.Sprintf("/storage-pool/detail?Uuid=%s&Option=%d", "ca3e11bd-b840-4412-9f44-89b8d09b9f4e", libvirt.STORAGE_XML_INACTIVE),
		string(http.MethodGet),
		[]byte{},
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition)
	}

	option, errorParseOptionToUint, isErrorParseOptionToUint := utils.StringToUint(incomingRequest.Option)
	if isErrorParseOptionToUint {
		test.Fatalf("parsing option failed: %s", errorParseOptionToUint)
	}

	_, errorGetPoolDetail, isErrorGetPoolDetail := poolConnection.PoolDetail(incomingRequest.Uuid, option)
	if isErrorGetPoolDetail {
		test.Logf("get pool detail test failed: %s", errorGetPoolDetail.Message)
	}
	if !isErrorGetPoolDetail {
		test.Error("this test must be failed, but it won't")
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

func TestPoolDetailInactiveStateUuidNotValid(test *testing.T) {
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

	var incomingRequest poolDetail.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		fmt.Sprintf("/storage-pool/detail?Uuid=%s&Option=%d", "ca3e11bd-b840-4412-9f44-", libvirt.STORAGE_XML_INACTIVE),
		string(http.MethodGet),
		[]byte{},
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition)
	}

	option, errorParseOptionToUint, isErrorParseOptionToUint := utils.StringToUint(incomingRequest.Option)
	if isErrorParseOptionToUint {
		test.Fatalf("parsing option failed: %s", errorParseOptionToUint)
	}

	_, errorGetPoolDetail, isErrorGetPoolDetail := poolConnection.PoolDetail(incomingRequest.Uuid, option)
	if isErrorGetPoolDetail {
		test.Logf("get pool detail test failed: %s", errorGetPoolDetail.Message)
	}
	if !isErrorGetPoolDetail {
		test.Error("this test must be failed, but it won't")
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
