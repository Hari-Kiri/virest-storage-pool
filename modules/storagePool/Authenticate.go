package storagePool

import (
	"net/http"
	"time"

	"github.com/Hari-Kiri/virest-utilities/utils"
	"github.com/golang-jwt/jwt"
	"libvirt.org/go/libvirt"
)

// Authenticate and generate JWT token for user with basic auth. After authentication succeed, new JWT will be generated with issuer name
// from supplied argument 'applicationName' and valid until supplied argument 'jwtLifetimeDuration'.
func Authenticate(httpRequest *http.Request, applicationName string, jwtLifetimeDuration time.Duration, jwtSigningMethod *jwt.SigningMethodHMAC, jwtSignatureKey []byte) (string, libvirt.Error, bool) {
	username, password, ok := httpRequest.BasicAuth()
	if !ok {
		return "", libvirt.Error{
			Code:    libvirt.ERR_AUTH_FAILED,
			Domain:  libvirt.FROM_AUTH,
			Message: "basic authentication credential not found",
			Level:   2,
		}, true
	}

	succeed, errorBasicAuth := utils.BasicAuthVerification(username, password)

	if !succeed {
		return "", libvirt.Error{
			Code:    libvirt.ERR_AUTH_FAILED,
			Domain:  libvirt.FROM_AUTH,
			Message: errorBasicAuth.Error(),
			Level:   2,
		}, true
	}

	token := jwt.NewWithClaims(
		jwtSigningMethod,
		jwt.StandardClaims{
			Issuer:    applicationName,
			ExpiresAt: time.Now().Add(jwtLifetimeDuration).Unix(),
		},
	)

	signedToken, errorSigningToken := token.SignedString(jwtSignatureKey)
	if errorSigningToken != nil {
		return "", libvirt.Error{
			Code:    libvirt.ERR_AUTH_FAILED,
			Domain:  libvirt.FROM_AUTH,
			Message: errorSigningToken.Error(),
			Level:   2,
		}, true
	}

	return signedToken, libvirt.Error{}, false
}
