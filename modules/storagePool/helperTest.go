package storagePool

import (
	"bytes"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/Hari-Kiri/virest-utilities/utils/auth"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/pbkdf2"
	"libvirt.org/go/libvirt"
)

type helperTest struct {
	test *testing.T
}

func helperTestCreateRestApiConnection[RequestStructure utils.RequestStructure](test *testing.T, envFilePath string, passwordHashIteration, keylen int,
	requestUrl, requestMethod string, requestBody []byte, requestStructure *RequestStructure) (poolConnection, virest.Error, bool) {
	test.Helper()

	helperTest := helperTest{
		test: test,
	}
	errorProcessEnvFile, isErrorProcessEnvFile := helperTest.helperTestProcessEnvFile(envFilePath)
	if isErrorProcessEnvFile {
		test.Fail()
		return poolConnection{}, errorProcessEnvFile, true
	}

	jsonWebToken, errorAuthenticate, isErrorAuthenticate := helperTest.helperTestAuthenticate(passwordHashIteration, keylen)
	if isErrorAuthenticate {
		test.Fail()
		return poolConnection{}, errorAuthenticate, true
	}

	httpRequest, errorCreateHttpRequest, isErrorCreateHttpRequest := helperTest.helperTestCreateHttpRequest(
		requestUrl,
		requestMethod,
		requestBody,
		jsonWebToken,
	)
	if isErrorCreateHttpRequest {
		test.Fail()
		return poolConnection{}, errorCreateHttpRequest, true
	}

	return HttpRequestPrecondition(
		httpRequest,
		requestMethod,
		requestStructure,
		os.Getenv("VIREST_STORAGE_POOL_APPLICATION_NAME"),
		jwt.SigningMethodHS512,
		[]byte(os.Getenv("VIREST_STORAGE_POOL_APPLICATION_JWT_SIGNATURE_KEY")),
	)
}

func (helper helperTest) helperTestProcessEnvFile(envFilePath string) (virest.Error, bool) {
	helper.test.Helper()

	readEnvFile, errorReadEnvFile := os.ReadFile(envFilePath)
	if errorReadEnvFile != nil {
		helper.test.Fail()
		return virest.Error{Error: libvirt.Error{
			Code:    libvirt.ERR_INTERNAL_ERROR,
			Domain:  libvirt.FROM_TEST,
			Message: errorReadEnvFile.Error(),
			Level:   2,
		}}, true
	}

	rows := strings.Split(string(readEnvFile), "\n")
	for i := 0; i < len(rows); i++ {
		if len(rows[i]) == 0 {
			continue
		}
		if len(rows[i]) >= 1 && string(rows[i][0]) == "#" {
			continue
		}
		columns := strings.Split(rows[i], "=")
		os.Setenv(columns[0], columns[1])
	}

	return virest.Error{}, false
}

func (helper helperTest) helperTestAuthenticate(iterations, keyLen int) (string, virest.Error, bool) {
	helper.test.Helper()

	salt128Byte := make([]byte, 128)
	_, errorReadSalt := rand.Read(salt128Byte)
	if errorReadSalt != nil {
		helper.test.Fail()
		return "", virest.Error{Error: libvirt.Error{
			Code:    libvirt.ERR_INTERNAL_ERROR,
			Domain:  libvirt.FROM_TEST,
			Message: errorReadSalt.Error(),
			Level:   2,
		}}, true
	}
	salt128ByteBase64 := base64.StdEncoding.EncodeToString(salt128Byte)

	userKey := pbkdf2.Key(
		[]byte(os.Getenv("VIREST_STORAGE_POOL_APPLICATION_BA_USER")),
		[]byte(salt128ByteBase64),
		iterations,
		keyLen,
		sha512.New,
	)
	passwordKey := pbkdf2.Key(
		[]byte(os.Getenv("VIREST_STORAGE_POOL_APPLICATION_BA_PASSWORD")),
		[]byte(salt128ByteBase64),
		iterations,
		keyLen,
		sha512.New,
	)

	userKeyBase64 := base64.StdEncoding.EncodeToString(userKey)
	passwordKeyBase64 := base64.StdEncoding.EncodeToString(passwordKey)

	userForBearerToken := fmt.Sprintf("$%d$%s$%s", iterations, salt128ByteBase64, userKeyBase64)
	passwordForBearerToken := fmt.Sprintf("$%d$%s$%s", iterations, salt128ByteBase64, passwordKeyBase64)

	request := &http.Request{
		Header: map[string][]string{
			"Authorization": {fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString(
				fmt.Appendf(nil, "%s:%s", userForBearerToken, passwordForBearerToken),
			))},
		},
	}

	jwtLifetimeSeconds, errorParseJwtLifetimeSeconds, isErrorParseJwtLifetimeSeconds := utils.StringToUint32(
		os.Getenv("VIREST_STORAGE_POOL_APPLICATION_JWT_LIFETIME_SECONDS"),
	)
	if isErrorParseJwtLifetimeSeconds {
		helper.test.Fail()
		return "", errorParseJwtLifetimeSeconds, true
	}

	result, errorBasicAuth, isErrorBasicAuth := auth.BasicAuth(
		request,
		os.Getenv("VIREST_STORAGE_POOL_APPLICATION_BA_USER"),
		os.Getenv("VIREST_STORAGE_POOL_APPLICATION_BA_PASSWORD"),
		os.Getenv("VIREST_STORAGE_POOL_APPLICATION_NAME"),
		time.Second*time.Duration(jwtLifetimeSeconds),
		jwt.SigningMethodHS512,
		[]byte(os.Getenv("VIREST_STORAGE_POOL_APPLICATION_JWT_SIGNATURE_KEY")),
	)
	if isErrorBasicAuth {
		helper.test.Fail()
		return "", errorBasicAuth, true
	}

	return result, virest.Error{}, false
}

func (helper helperTest) helperTestCreateHttpRequest(requestUrl, requestMethod string, requestBody []byte, jwt string) (*http.Request, virest.Error, bool) {
	helper.test.Helper()

	urlParsed, errorParsingUrl := url.Parse(requestUrl)
	if errorParsingUrl != nil {
		helper.test.Fail()
		return nil, virest.Error{Error: libvirt.Error{
			Code:    libvirt.ERR_INTERNAL_ERROR,
			Domain:  libvirt.FROM_TEST,
			Message: errorParsingUrl.Error(),
			Level:   2,
		}}, true
	}

	// Create an io.Reader from the string
	bytesReader := bytes.NewReader(requestBody)
	ioReadCloser := io.NopCloser(bytesReader)
	defer ioReadCloser.Close()

	header := map[string][]string{
		"Accept":          {"*/*"},
		"Accept-Encoding": {"gzip", "deflate", "br"},
		"Authorization":   {fmt.Sprintf("Bearer %s", jwt)},
		"Connection":      {"keep-alive"},
		"Content-Length":  {"299"},
		"Content-Type":    {"application/json"},
		"Hypervisor-Uri":  {"qemu:///system"},
		"User-Agent":      {""},
	}

	return &http.Request{
		Method:     requestMethod,
		URL:        urlParsed,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     header,
		Body:       ioReadCloser,
	}, virest.Error{}, false
}
