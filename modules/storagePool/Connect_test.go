package storagePool

import (
	"fmt"
	"testing"

	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

const hypervisorUri = "qemu:///system"

func TestConnection(test *testing.T) {
	virestConnection, errorConnect, isErrorConnect := utils.NewConnectWithAuth(hypervisorUri, nil, 0)
	if isErrorConnect {
		test.Fatalf("connection test failed: %s", errorConnect.Message)
	}

	result, errorResult := virestConnection.Close()
	if errorResult != nil {
		test.Fatalf("close() error: %s", errorResult.Error())
	}
	if result != 0 {
		test.Errorf("close() == %d, expected 0", result)
	}
}

func helperTestConnection(test *testing.T) (*poolConnection, virest.Error, bool) {
	test.Helper()

	virestConnection, errorConnect, isErrorConnect := utils.NewConnectWithAuth(hypervisorUri, nil, 0)
	if isErrorConnect {
		test.Fail()
		return &poolConnection{}, errorConnect, true
	}

	return &poolConnection{virestConnection.Connect}, virest.Error{}, false
}

func (poolConnection *poolConnection) helperTestCloseConnection(test *testing.T) (virest.Error, bool) {
	test.Helper()

	var (
		virestError virest.Error
		isError     bool
	)
	result, errorResult := poolConnection.Close()
	virestError.Error, isError = errorResult.(libvirt.Error)
	if isError {
		test.Fail()
		return virestError, isError
	}
	if result != 0 {
		test.Fail()
		return virest.Error{Error: libvirt.Error{
			Code:    libvirt.ERR_INTERNAL_ERROR,
			Domain:  libvirt.FROM_ACCESS,
			Message: fmt.Sprintf("close() == %d, expected 0", result),
			Level:   libvirt.ERR_WARNING,
		}}, true
	}

	return virest.Error{}, false
}
