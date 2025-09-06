package storagePool

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolUndefine"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

func (helper helperTest) helperTestPoolUndefine(poolUuid string) (virest.Error, bool) {
	helper.test.Helper()

	var incomingRequest poolUndefine.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		helper.test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/undefine",
		http.MethodDelete,
		fmt.Appendf(nil, "{\"uuid\":\"%s\"}", poolUuid),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		helper.test.Fail()
		return errorHttpRequestPrecondition, isErrorHttpRequestPrecondition
	}
	defer poolConnection.Close()

	errorPoolUndefine, isErrorPoolUndefine := poolConnection.PoolUndefine(incomingRequest.Uuid)
	if isErrorPoolUndefine {
		helper.test.Fail()
		return errorPoolUndefine, isErrorPoolUndefine
	}

	return virest.Error{}, false
}

func TestPoolUndefine(test *testing.T) {
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
		test.Fatalf("pool destroy test failed: %s", errorPoolDestroy.Message)
	}

	errorPoolDelete, isErrorPoolDelete := helper.helperTestPoolDelete(poolUuid, libvirt.STORAGE_POOL_DELETE_NORMAL)
	if isErrorPoolDelete {
		test.Fatalf("deleting pool test failed: %s", errorPoolDelete.Message)
	}

	var incomingRequest poolUndefine.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/undefine",
		http.MethodDelete,
		fmt.Appendf(nil, "{\"uuid\":\"%s\"}", poolUuid),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition)
	}

	errorPoolUndefine, isErrorPoolUndefine := poolConnection.PoolUndefine(incomingRequest.Uuid)
	if isErrorPoolUndefine {
		test.Errorf("pool undefine test failed: %s", errorPoolUndefine.Message)
	}

	test.Cleanup(func() {
		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}

		if errorDeletePartition, isErrorDeletePartition := helper.helperDepleteDevicePartition(filesystemDiskDeviceValue); isErrorDeletePartition {
			test.Errorf("delete primary partition failed: %s", errorDeletePartition.Message)
		}
	})
}

func TestPoolUndefineWrongUuid(test *testing.T) {
	var incomingRequest poolUndefine.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/undefine",
		http.MethodDelete,
		fmt.Appendf(nil, "{\"uuid\":\"%s\"}", "ca3e11bd-b840-4412-9f44-89b8d09b9f4e"),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition)
	}

	errorPoolUndefine, isErrorPoolUndefine := poolConnection.PoolUndefine(incomingRequest.Uuid)
	if isErrorPoolUndefine {
		test.Logf("pool undefine test failed: %s", errorPoolUndefine.Message)
	}
	if !isErrorPoolUndefine {
		test.Error("this test must be failed, but it won't")
	}

	test.Cleanup(func() {
		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}
	})
}

func TestPoolUndefineInvalidUuid(test *testing.T) {
	var incomingRequest poolUndefine.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/undefine",
		http.MethodDelete,
		fmt.Appendf(nil, "{\"uuid\":\"%s\"}", "ca3e11bd-b840-4412-9f44-"),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition)
	}

	errorPoolUndefine, isErrorPoolUndefine := poolConnection.PoolUndefine(incomingRequest.Uuid)
	if isErrorPoolUndefine {
		test.Logf("pool undefine test failed: %s", errorPoolUndefine.Message)
	}
	if !isErrorPoolUndefine {
		test.Error("this test must be failed, but it won't")
	}

	test.Cleanup(func() {
		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}
	})
}
