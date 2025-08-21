package storagePool

import (
	"fmt"
	"net/http"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolUndefine"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
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
