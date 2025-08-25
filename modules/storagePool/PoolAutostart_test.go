package storagePool

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolAutostart"
	"libvirt.org/go/libvirt"
)

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
		test.Fatalf("turn on pool autostart test failed: %s", errorPoolAutostart.Message)
	}

	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition = helperTestCreateRestApiConnection(
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
		test.Fatalf("turn on pool autostart test failed: %s", errorPoolAutostart.Message)
	}

	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition = helperTestCreateRestApiConnection(
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
