package storagePool

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolCapabilities"
)

func TestPoolCapabilities(test *testing.T) {
	var incomingRequest poolCapabilities.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/capabilities",
		http.MethodGet,
		[]byte{},
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	storagePoolCapabilities, errorPoolCapabilities, isErrorPoolCapabilities := poolConnection.PoolCapabilities()
	if isErrorPoolCapabilities {
		test.Errorf("get storage pool capabilities test failed: %s", errorPoolCapabilities.Message)
	}
	storagePoolCapabilitiesMarshaled, errorMarshalingStoragePoolCapabilities := json.MarshalIndent(storagePoolCapabilities, "", "  ")
	if errorMarshalingStoragePoolCapabilities != nil {
		test.Errorf("marshaling result failed: %s", errorMarshalingStoragePoolCapabilities.Error())
	}
	test.Logf("storage pool capabilities: \n%s", string(storagePoolCapabilitiesMarshaled))

	test.Cleanup(func() {
		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}
	})
}
