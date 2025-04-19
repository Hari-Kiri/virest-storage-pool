package storagePool

import (
	"testing"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
)

func (poolConnection *poolConnection) helperTestPoolUndefine(test *testing.T, poolUuid string) (virest.Error, bool) {
	test.Helper()

	errorPoolUndefine, isErrorPoolUndefine := poolConnection.PoolUndefine(poolUuid)
	if isErrorPoolUndefine {
		test.Fail()
		return errorPoolUndefine, isErrorPoolUndefine
	}

	return virest.Error{}, false
}
