package storagePool

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/Hari-Kiri/virest-storage-pool/structures/findStoragePoolSources"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/libvirtxml"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

// Change it to Your's.
const (
	netfsHost = "192.168.122.4"
	iscsiHost = "192.168.122.6"
)

func TestFindStoragePoolSourceNetfs(test *testing.T) {
	requestData := findStoragePoolSources.Request{
		Type: "netfs",
		SrcSpec: findStoragePoolSources.Source{
			Source: libvirtxml.Source{
				Host: libvirtxml.Host{
					Name: netfsHost,
					Port: 2049,
				},
			},
		},
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

	var incomingRequest findStoragePoolSources.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/find-storage-pool-sources",
		http.MethodPost,
		requestBody,
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	sources, errorFindStoragePoolSource, isErrorFindStoragePoolSource := poolConnection.FindStoragePoolSource(
		incomingRequest.Type,
		incomingRequest.SrcSpec.Source,
	)
	if isErrorFindStoragePoolSource {
		test.Errorf("find storage pool netfs source test failed: %s", errorFindStoragePoolSource.Message)
	}
	sourcesMarshaled, errorMarshalingSources := json.MarshalIndent(sources, "", "  ")
	if errorMarshalingSources != nil {
		test.Errorf("marshaling result failed: %s", errorMarshalingSources.Error())
	}
	test.Logf("found netfs storage source: \n%s", string(sourcesMarshaled))

	test.Cleanup(func() {
		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}
	})
}

func TestFindStoragePoolSourceIscsi(test *testing.T) {
	requestData := findStoragePoolSources.Request{
		Type: "iscsi",
		SrcSpec: findStoragePoolSources.Source{
			Source: libvirtxml.Source{
				Host: libvirtxml.Host{
					Name: iscsiHost,
					Port: 3260,
				},
			},
		},
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

	var incomingRequest findStoragePoolSources.Request
	poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
		test,
		"/home/hari/virest-storage-pool/.env",
		210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
		64,
		"/storage-pool/find-storage-pool-sources",
		http.MethodPost,
		requestBody,
		&incomingRequest,
	)
	if isErrorHttpRequestPrecondition {
		test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition.Message)
	}

	sources, errorFindStoragePoolSource, isErrorFindStoragePoolSource := poolConnection.FindStoragePoolSource(
		incomingRequest.Type,
		incomingRequest.SrcSpec.Source,
	)
	if isErrorFindStoragePoolSource {
		test.Errorf("find storage pool iscsi source test failed: %s", errorFindStoragePoolSource.Message)
	}
	sourcesMarshaled, errorMarshalingSources := json.MarshalIndent(sources, "", "  ")
	if errorMarshalingSources != nil {
		test.Errorf("marshaling result failed: %s", errorMarshalingSources.Error())
	}
	test.Logf("found iscsi storage source: \n%s", string(sourcesMarshaled))

	test.Cleanup(func() {
		connectionReference, errorClosingPoolConnection := poolConnection.Close()
		if errorClosingPoolConnection != nil {
			test.Errorf("closing pool connection %d failed: %s", connectionReference, errorClosingPoolConnection.Error())
		}
	})
}
