package storagePool

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDelete"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

func (helper helperTest) helperTestPoolDelete(poolUuid string, option libvirt.StoragePoolDeleteFlags) (virest.Error, bool) {
	helper.test.Helper()

	var incomingRequest poolDelete.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		helper.test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/delete",
		http.MethodDelete,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"option\":%d}", poolUuid, option),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		helper.test.Fail()
		return errorHttpRequestPrecondition, isErrorHttpRequestPrecondition
	}
	defer poolConnection.Close()

	errorPoolDelete, isErrorPoolDelete := poolConnection.PoolDelete(incomingRequest.Uuid, incomingRequest.Option)
	if isErrorPoolDelete {
		helper.test.Fail()
		return errorPoolDelete, isErrorPoolDelete
	}

	return virest.Error{}, false
}

func TestPoolDeleteMetadaOnly(test *testing.T) {
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

	errorPoolDestroy, isErrorPoolDestroy := helper.helperTestPoolDestroy(poolUuid)
	if isErrorPoolDestroy {
		test.Fatalf("pool destroy failed: %s", errorPoolDestroy.Message)
	}

	var incomingRequest poolDelete.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/delete",
		http.MethodDelete,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"option\":%d}", poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolDelete, isErrorPoolDelete := poolConnection.PoolDelete(incomingRequest.Uuid, incomingRequest.Option)
	if isErrorPoolDelete {
		test.Errorf("deleting pool test failed: %s", errorPoolDelete.Message)
	}

	test.Cleanup(func() {
		if errorPoolUndefine, isErrorPoolUndefine := helper.helperTestPoolUndefine(poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}
	})
}

// Unsupported in libvirt version 10.0.0 and virsh (the libvirt command line interface) still not implement this method.
func TestPoolDeleteClearAllToZeros(test *testing.T) {
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

	errorPoolDestroy, isErrorPoolDestroy := helper.helperTestPoolDestroy(poolUuid)
	if isErrorPoolDestroy {
		test.Fatalf("pool destroy failed: %s", errorPoolDestroy.Message)
	}

	var incomingRequest poolDelete.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/delete",
		http.MethodDelete,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"option\":%d}", poolUuid, libvirt.STORAGE_POOL_DELETE_ZEROED),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolDelete, isErrorPoolDelete := poolConnection.PoolDelete(incomingRequest.Uuid, incomingRequest.Option)
	if isErrorPoolDelete {
		test.Logf("deleting pool and clear all data test failed: %s", errorPoolDelete.Message)
	}
	if !isErrorPoolDelete {
		test.Error("this test must be failed, but it won't")
	}

	test.Cleanup(func() {
		if errorPoolDelete, isErrorPoolDelete = helper.helperTestPoolDelete(poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL); isErrorPoolDelete {
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

func TestPoolDeleteMetadaOnlyWrongUuid(test *testing.T) {
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

	errorPoolDestroy, isErrorPoolDestroy := helper.helperTestPoolDestroy(poolUuid)
	if isErrorPoolDestroy {
		test.Fatalf("pool destroy failed: %s", errorPoolDestroy.Message)
	}

	var incomingRequest poolDelete.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/delete",
		http.MethodDelete,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"option\":%d}", "ca3e11bd-b840-4412-9f44-89b8d09b9f4e", libvirt.STORAGE_POOL_DELETE_NORMAL),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolDelete, isErrorPoolDelete := poolConnection.PoolDelete(incomingRequest.Uuid, incomingRequest.Option)
	if isErrorPoolDelete {
		test.Logf("deleting pool test failed: %s", errorPoolDelete.Message)
	}
	if !isErrorPoolDelete {
		test.Error("this test must be failed, but it won't")
	}

	test.Cleanup(func() {
		if errorPoolUndefine, isErrorPoolUndefine := helper.helperTestPoolUndefine(poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}
	})
}

func TestPoolDeleteMetadaOnlyUuidNotValid(test *testing.T) {
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

	errorPoolDestroy, isErrorPoolDestroy := helper.helperTestPoolDestroy(poolUuid)
	if isErrorPoolDestroy {
		test.Fatalf("pool destroy failed: %s", errorPoolDestroy.Message)
	}

	var incomingRequest poolDelete.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/delete",
		http.MethodDelete,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"option\":%d}", "ca3e11bd-b840-4412-9f44-", libvirt.STORAGE_POOL_DELETE_NORMAL),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolDelete, isErrorPoolDelete := poolConnection.PoolDelete(incomingRequest.Uuid, incomingRequest.Option)
	if isErrorPoolDelete {
		test.Logf("deleting pool test failed: %s", errorPoolDelete.Message)
	}
	if !isErrorPoolDelete {
		test.Error("this test must be failed, but it won't")
	}

	test.Cleanup(func() {
		if errorPoolUndefine, isErrorPoolUndefine := helper.helperTestPoolUndefine(poolUuid); isErrorPoolUndefine {
			test.Errorf("pool undefine failed: %s", errorPoolUndefine.Message)
		}

		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}
	})
}
