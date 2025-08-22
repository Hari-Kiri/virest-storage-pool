package storagePool

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolList"
	"github.com/Hari-Kiri/virest-utilities/utils"
	"libvirt.org/go/libvirt"
)

func TestPoolListDetailActive(test *testing.T) {
	var incomingRequest poolList.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		fmt.Sprintf("/storage-pool/info?Option=%d&Inactive=%d", 0, 0),
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

	inactive, errorParseInactiveToUint, isErrorParseInactiveToUint := utils.StringToUint(incomingRequest.Inactive)
	if isErrorParseInactiveToUint {
		test.Fatalf("parsing inactive failed: %s", errorParseInactiveToUint)
	}

	poolList, errorGetPoolList, isErrorGetPoolList := poolConnection.PoolList(option, inactive)
	if isErrorGetPoolList {
		test.Errorf("get pool list test failed: %s", errorGetPoolList.Message)
	}
	poolInfolMarshaled, errorMarshalingPoolInfo := json.MarshalIndent(poolList, "", "  ")
	if errorMarshalingPoolInfo != nil {
		test.Errorf("marshaling result failed: %s", errorMarshalingPoolInfo.Error())
	}
	test.Logf("list storage pool: \n%s", string(poolInfolMarshaled))

	test.Cleanup(func() {
		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}
	})
}

func TestPoolListDetailInactive(test *testing.T) {
	var incomingRequest poolList.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		fmt.Sprintf("/storage-pool/info?Option=%d&Inactive=%d", 0, libvirt.STORAGE_XML_INACTIVE),
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

	inactive, errorParseInactiveToUint, isErrorParseInactiveToUint := utils.StringToUint(incomingRequest.Inactive)
	if isErrorParseInactiveToUint {
		test.Fatalf("parsing inactive failed: %s", errorParseInactiveToUint)
	}

	poolList, errorGetPoolList, isErrorGetPoolList := poolConnection.PoolList(option, inactive)
	if isErrorGetPoolList {
		test.Errorf("get pool list test failed: %s", errorGetPoolList.Message)
	}
	poolInfolMarshaled, errorMarshalingPoolInfo := json.MarshalIndent(poolList, "", "  ")
	if errorMarshalingPoolInfo != nil {
		test.Errorf("marshaling result failed: %s", errorMarshalingPoolInfo.Error())
	}
	test.Logf("list storage pool: \n%s", string(poolInfolMarshaled))

	test.Cleanup(func() {
		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}
	})
}
