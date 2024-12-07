package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/Hari-Kiri/virest-storage-pool/structures/poolDefine"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolList"
	"github.com/Hari-Kiri/virest-storage-pool/structures/poolUndefine"
	"libvirt.org/go/libvirt"
)

// Defined generic type constraint for request model structure.
type RequestStructure interface {
	poolList.Request | poolDefine.Request | poolUndefine.Request
}

// Check the expected HTTP request method and convert the JSON request body to structure if not HTTP GET method. HTTP GET method must be use
// query parameter.
//
// Notes for HTTP GET method:
//
// - Query parameter and structure field will be compared in case sensitive.
//
// - Every structure field data type must be string, so You must convert it to the right data type before You use it.
//
// - Untested for array query argument.
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
	if len(requestBody) == 0 && expectedRequestMethod != "GET" {
		return libvirt.Error{}, false
	}

	// Set GET parameter to structure. For now HTTP GET Method only have one query parameter "option"
	var errorSetStructure error
	if expectedRequestMethod == "GET" {
		errorSetStructure = setHttpGetStructure(httpRequest, structure)
	}
	if errorSetStructure != nil {
		return libvirt.Error{
			Code:    libvirt.ERR_INVALID_ARG,
			Domain:  libvirt.FROM_NET,
			Message: fmt.Sprintf("%s", errorSetStructure),
			Level:   libvirt.ERR_ERROR,
		}, true
	}

	// Parse JSON to model
	var errorUnmarshal error
	if expectedRequestMethod != "GET" {
		errorUnmarshal = json.Unmarshal(requestBody, structure)
	}
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

// Set HTTP GET method query parameter to structure. Query parameter and structure field will be compared in case sensitive.
// Every structure field data type must be string.
func setHttpGetStructure[Structure RequestStructure](httpRequest *http.Request, structure *Structure) error {
	var errorResult error

	keys := reflect.ValueOf(structure).Elem()
	for i := 0; i < keys.NumField(); i++ {
		reflect.ValueOf(structure).Elem().FieldByName(
			keys.Type().Field(i).Name,
		).SetString(
			httpRequest.URL.Query().Get(
				keys.Type().Field(i).Name,
			),
		)
	}

	return errorResult
}
