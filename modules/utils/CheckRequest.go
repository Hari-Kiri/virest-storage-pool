package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDefine"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolUndefine"
	"libvirt.org/go/libvirt"
)

// Defined generic type constraint for request model structure.
type RequestStructure interface {
	poolDefine.Request | poolUndefine.Request
}

// Check the expected HTTP request method and convert the JSON request body to structure if any.
func CheckRequest[Structure RequestStructure](httpRequest *http.Request, expectedRequestMethod string, structure *Structure) (libvirt.Error, bool) {
	// Create libvirt error number
	var libvirtErrorNumber libvirt.ErrorNumber
	if expectedRequestMethod == "GET" {
		libvirtErrorNumber = libvirt.ERR_GET_FAILED
	}
	if expectedRequestMethod == "POST" {
		libvirtErrorNumber = libvirt.ERR_POST_FAILED
	}
	if expectedRequestMethod == "PUT" {
		libvirtErrorNumber = libvirt.ERR_HTTP_ERROR
	}
	if expectedRequestMethod == "PATCH" {
		libvirtErrorNumber = libvirt.ERR_HTTP_ERROR
	}
	if expectedRequestMethod == "DELETE" {
		libvirtErrorNumber = libvirt.ERR_HTTP_ERROR
	}

	// Check http method
	if httpRequest.Method != expectedRequestMethod {
		return libvirt.Error{
			Code:    libvirtErrorNumber,
			Domain:  libvirt.FROM_NET,
			Message: fmt.Sprintf("a HTTP %s command to failed", expectedRequestMethod),
			Level:   libvirt.ERR_ERROR,
		}, true
	}

	// Read request body
	requestBody, errorReadRequestBody := io.ReadAll(httpRequest.Body)
	if errorReadRequestBody != nil {
		return libvirt.Error{
			Code:    libvirt.ERR_HTTP_ERROR,
			Domain:  libvirt.FROM_NET,
			Message: fmt.Sprintf("%s", errorReadRequestBody),
			Level:   libvirt.ERR_ERROR,
		}, true
	}

	// Request body empty
	if len(requestBody) == 0 {
		return libvirt.Error{}, false
	}

	// Parse JSON to model
	errorUnmarshal := json.Unmarshal(requestBody, structure)
	if errorUnmarshal != nil {
		return libvirt.Error{
			Code:    libvirt.ERR_INTERNAL_ERROR,
			Domain:  libvirt.FROM_NET,
			Message: fmt.Sprintf("%s", errorUnmarshal),
			Level:   libvirt.ERR_ERROR,
		}, true
	}

	return libvirt.Error{}, false
}
