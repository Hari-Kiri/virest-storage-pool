package storagePool

import (
	"fmt"
	"testing"

	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
)

const hypervisorUri = "qemu:///system"

func TestConnection(test *testing.T) {
	virestConnection, errorConnect, isErrorConnect := utils.NewConnectWithAuth(hypervisorUri, nil, 0)
	if isErrorConnect {
		test.Errorf("connection test failed: %s", errorConnect.Message)
		return
	}

	result, errorResult := virestConnection.Close()
	if errorResult != nil {
		test.Fatalf("close() error: %s", errorResult.Error())
		return
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

func (poolConnection *poolConnection) helperTestCloseConnection(test *testing.T) error {
	test.Helper()

	result, errorResult := poolConnection.Close()
	if errorResult != nil {
		test.Fail()
		return errorResult
	}
	if result != 0 {
		test.Fail()
		return fmt.Errorf("close() == %d, expected 0", result)
	}

	return nil
}
