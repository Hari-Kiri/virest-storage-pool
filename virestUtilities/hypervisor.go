package virestUtilities

import (
	"fmt"

	"libvirt.org/go/libvirt"
)

// Dial hooks (overridable in unit tests).
var (
	newConnect              = libvirt.NewConnect
	newConnectReadOnly      = libvirt.NewConnectReadOnly
	newConnectWithAuth      = libvirt.NewConnectWithAuth
	newConnectWithAuthDefault = libvirt.NewConnectWithAuthDefault
)

// Connect opens a connection to the hypervisor at uri.
// See https://libvirt.org/uri.html for URI formats.
func Connect(uri string) (Connection, error) {
	conn, err := newConnect(uri)
	return wrapConnect(conn, err, "connect")
}

// ConnectReadOnly opens a read-only connection to the hypervisor at uri.
func ConnectReadOnly(uri string) (Connection, error) {
	conn, err := newConnectReadOnly(uri)
	return wrapConnect(conn, err, "connect read-only")
}

// ConnectWithAuth opens a connection, invoking auth when credentials are required.
func ConnectWithAuth(uri string, auth *libvirt.ConnectAuth, flags libvirt.ConnectFlags) (Connection, error) {
	conn, err := newConnectWithAuth(uri, auth, flags)
	return wrapConnect(conn, err, "connect with auth")
}

// ConnectWithAuthDefault opens a connection using libvirt's default auth callback.
func ConnectWithAuthDefault(uri string, flags libvirt.ConnectFlags) (Connection, error) {
	conn, err := newConnectWithAuthDefault(uri, flags)
	return wrapConnect(conn, err, "connect with auth default")
}

func wrapConnect(conn *libvirt.Connect, err error, op string) (Connection, error) {
	if err != nil {
		return Connection{}, Wrap(op, err)
	}
	if conn == nil {
		return Connection{}, fmt.Errorf("%s: nil libvirt connection", op)
	}
	return Connection{Connect: conn}, nil
}
