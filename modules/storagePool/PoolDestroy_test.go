package storagePool

import (
	"testing"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
)

func (poolConnection *poolConnection) helperTestPoolDestroy(test *testing.T, poolUuid string) (virest.Error, bool) {
	test.Helper()

	errorPoolDestroy, isErrorPoolDestroy := poolConnection.PoolDestroy(poolUuid)
	if isErrorPoolDestroy {
		test.Fail()
		return errorPoolDestroy, isErrorPoolDestroy
	}

	return virest.Error{}, false
}
