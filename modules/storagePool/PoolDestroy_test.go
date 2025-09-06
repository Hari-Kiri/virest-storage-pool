package storagePool

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDestroy"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

func (helper helperTest) helperTestPoolDestroy(poolUuid string) (virest.Error, bool) {
	helper.test.Helper()

	var incomingRequest poolDestroy.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		helper.test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/destroy",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\"}", poolUuid),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		helper.test.Fail()
		return errorHttpRequestPrecondition, isErrorHttpRequestPrecondition
	}
	defer poolConnection.Close()

	errorPoolDestroy, isErrorPoolDestroy := poolConnection.PoolDestroy(incomingRequest.Uuid)
	if isErrorPoolDestroy {
		helper.test.Fail()
		return errorPoolDestroy, isErrorPoolDestroy
	}

	return virest.Error{}, false
}

func TestPoolDestroy(test *testing.T) {
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

	var incomingRequest poolDestroy.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/destroy",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\"}", poolUuid),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolDestroy, isErrorPoolDestroy := poolConnection.PoolDestroy(incomingRequest.Uuid)
	if isErrorPoolDestroy {
		test.Errorf("pool destroy test failed: %s", errorPoolDestroy.Message)
	}

	test.Cleanup(func() {
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

func TestPoolDestroyWrongUuid(test *testing.T) {
	var incomingRequest poolDestroy.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/destroy",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\"}", "ca3e11bd-b840-4412-9f44-89b8d09b9f4e"),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolDestroy, isErrorPoolDestroy := poolConnection.PoolDestroy(incomingRequest.Uuid)
	if isErrorPoolDestroy {
		test.Logf("pool destroy test failed: %s", errorPoolDestroy.Message)
	}
	if !isErrorPoolDestroy {
		test.Error("this test must be failed, but it won't")
	}

	test.Cleanup(func() {
		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}
	})
}

func TestPoolDestroyUuidNotValis(test *testing.T) {
	var incomingRequest poolDestroy.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/destroy",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\"}", "ca3e11bd-b840-4412-9f44-"),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolDestroy, isErrorPoolDestroy := poolConnection.PoolDestroy(incomingRequest.Uuid)
	if isErrorPoolDestroy {
		test.Logf("pool destroy test failed: %s", errorPoolDestroy.Message)
	}
	if !isErrorPoolDestroy {
		test.Error("this test must be failed, but it won't")
	}

	test.Cleanup(func() {
		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}
	})
}
