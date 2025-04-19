package storagePool

import (
	"testing"

	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

func (poolConnection *poolConnection) helperTestPoolDelete(test *testing.T, poolUuid string, option libvirt.StoragePoolDeleteFlags) (virest.Error, bool) {
	test.Helper()

	errorPoolDelete, isErrorPoolDelete := poolConnection.PoolDelete(poolUuid, option)
	if isErrorPoolDelete {
		test.Fail()
		return errorPoolDelete, isErrorPoolDelete
	}

	return virest.Error{}, false
}
