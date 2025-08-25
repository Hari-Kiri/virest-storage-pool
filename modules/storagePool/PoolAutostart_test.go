package storagePool

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolAutostart"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

func (helper helperTest) helperTestPoolAutostart(poolUuid string, autostart bool) (virest.Error, bool) {
	helper.test.Helper()

	var incomingRequest poolAutostart.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		helper.test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/define",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"autostart\":%t}", poolUuid, autostart),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		helper.test.Fail()
		return errorHttpRequestPrecondition, isErrorHttpRequestPrecondition
	}
	defer poolConnection.Close()

	errorPoolAutostart, isErrorPoolAutostart := poolConnection.PoolAutostart(incomingRequest.Uuid, incomingRequest.Autostart)
	if isErrorPoolAutostart {
		helper.test.Fail()
		return errorPoolAutostart, isErrorPoolAutostart
	}

	return virest.Error{}, false
}

func TestPoolAutostartTrue(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("defining pool test failed: %s", errorPoolDefine.Message)
	}

	var incomingRequest poolAutostart.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/autostart",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"autostart\":%t}", poolUuid, true),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolAutostart, isErrorPoolAutostart := poolConnection.PoolAutostart(incomingRequest.Uuid, incomingRequest.Autostart)
	if isErrorPoolAutostart {
		test.Errorf("turn on pool autostart test failed: %s", errorPoolAutostart.Message)
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

func TestPoolAutostartTrueTwice(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("defining pool test failed: %s", errorPoolDefine.Message)
	}

	errorPoolAutostart, isErrorPoolAutostart := helper.helperTestPoolAutostart(poolUuid, true)
	if isErrorPoolAutostart {
		test.Fatalf("turn on pool autostart test failed: %s", errorPoolAutostart.Message)
	}

	var incomingRequest poolAutostart.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/autostart",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"autostart\":%t}", poolUuid, true),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolAutostart, isErrorPoolAutostart = poolConnection.PoolAutostart(incomingRequest.Uuid, incomingRequest.Autostart)
	if isErrorPoolAutostart {
		test.Errorf("turn on pool autostart test failed: %s", errorPoolAutostart.Message)
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

func TestPoolAutostartFalse(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("defining pool test failed: %s", errorPoolDefine.Message)
	}

	errorPoolAutostart, isErrorPoolAutostart := helper.helperTestPoolAutostart(poolUuid, true)
	if isErrorPoolAutostart {
		test.Fatalf("turn on pool autostart test failed: %s", errorPoolAutostart.Message)
	}

	var incomingRequest poolAutostart.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/autostart",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"autostart\":%t}", poolUuid, false),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolAutostart, isErrorPoolAutostart = poolConnection.PoolAutostart(incomingRequest.Uuid, incomingRequest.Autostart)
	if isErrorPoolAutostart {
		test.Errorf("turn off pool autostart test failed: %s", errorPoolAutostart.Message)
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

func TestPoolAutostartFalseTwice(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("defining pool test failed: %s", errorPoolDefine.Message)
	}

	errorPoolAutostart, isErrorPoolAutostart := helper.helperTestPoolAutostart(poolUuid, true)
	if isErrorPoolAutostart {
		test.Fatalf("turn on pool autostart test failed: %s", errorPoolAutostart.Message)
	}

	errorPoolAutostart, isErrorPoolAutostart = helper.helperTestPoolAutostart(poolUuid, false)
	if isErrorPoolAutostart {
		test.Fatalf("turn on pool autostart test failed: %s", errorPoolAutostart.Message)
	}

	var incomingRequest poolAutostart.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/autostart",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"autostart\":%t}", poolUuid, false),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolAutostart, isErrorPoolAutostart = poolConnection.PoolAutostart(incomingRequest.Uuid, incomingRequest.Autostart)
	if isErrorPoolAutostart {
		test.Errorf("turn off pool autostart test failed: %s", errorPoolAutostart.Message)
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

func TestPoolAutostartWrongUuid(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("defining pool test failed: %s", errorPoolDefine.Message)
	}

	var incomingRequest poolAutostart.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/autostart",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"autostart\":%t}", "ca3e11bd-b840-4412-9f44-89b8d09b9f4e", true),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolAutostart, isErrorPoolAutostart := poolConnection.PoolAutostart(incomingRequest.Uuid, incomingRequest.Autostart)
	if isErrorPoolAutostart {
		test.Logf("turn on pool autostart test failed: %s", errorPoolAutostart.Message)
	}
	if !isErrorPoolAutostart {
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

func TestPoolAutostartUuidNotValid(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("defining pool test failed: %s", errorPoolDefine.Message)
	}

	var incomingRequest poolAutostart.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/autostart",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\",\"autostart\":%t}", "ca3e11bd-b840-4412-9f44-", true),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolAutostart, isErrorPoolAutostart := poolConnection.PoolAutostart(incomingRequest.Uuid, incomingRequest.Autostart)
	if isErrorPoolAutostart {
		test.Logf("turn on pool autostart test failed: %s", errorPoolAutostart.Message)
	}
	if !isErrorPoolAutostart {
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

func TestPoolAutostartBoolNotProvide(test *testing.T) {
	helper := helperTest{
		test: test,
	}

	poolUuid, errorPoolDefine, isErrorPoolDefine := helper.helperTestPoolDefine(storagePoolDirectory, libvirt.STORAGE_POOL_DEFINE_VALIDATE)
	if isErrorPoolDefine {
		test.Fatalf("defining pool test failed: %s", errorPoolDefine.Message)
	}

	poolInfo, errorGetPoolInfo, isErrorGetPoolInfo := helper.helperTestPoolInfo(poolUuid)
	if isErrorGetPoolInfo {
		test.Fatalf("get pool info test failed: %s", errorGetPoolInfo.Message)
	}
	test.Logf("define pool %s: %t", poolInfo.Name, poolInfo.Autostart)

	var incomingRequest poolAutostart.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/autostart",
		http.MethodPatch,
		fmt.Appendf(nil, "{\"uuid\":\"%s\"}", poolUuid),
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	errorPoolAutostart, isErrorPoolAutostart := poolConnection.PoolAutostart(incomingRequest.Uuid, incomingRequest.Autostart)
	if isErrorPoolAutostart {
		test.Errorf("turn on pool autostart test failed: %s", errorPoolAutostart.Message)
	}
	if !isErrorPoolAutostart {
		poolInfo, errorGetPoolInfo, isErrorGetPoolInfo = helper.helperTestPoolInfo(incomingRequest.Uuid)
	}
	if isErrorGetPoolInfo {
		test.Errorf("get pool info test failed: %s", errorGetPoolInfo.Message)
	}
	test.Logf("result pool autostart test of pool %s: %t", poolInfo.Name, poolInfo.Autostart)

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
