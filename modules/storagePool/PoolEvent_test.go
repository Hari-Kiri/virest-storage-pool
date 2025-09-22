package storagePool

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolEvent"
	"libvirt.org/go/libvirt"
)

func TestPoolEvent(test *testing.T) {
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

	done := make(chan bool)
	go func() {
		var incomingRequest poolEvent.Request
		poolConnection, errorHttpRequestPrecondition, isErrorHttpRequestPrecondition := helperTestCreateRestApiConnection(
			test,
			"/home/hari/virest-storage-pool/.env",
			210000, // circa 2023 OWASP recommendation for PBKDF2-HMAC-SHA512 iterations
			64,
			fmt.Sprintf("/storage-pool/event?Uuid=%s&Types=%d&Timeout=%d", poolUuid, 0, -1),
			http.MethodGet,
			[]byte{},
			&incomingRequest,
		)
		if isErrorHttpRequestPrecondition {
			test.Fatalf("http request precondition failed: %s", errorHttpRequestPrecondition)
		}

		poolEvent, errorGetStoragePoolEvent, isErrorGetStoragePoolEvent := poolConnection.PoolEvent(requestBodyData.Uuid, types)
	}()
}
