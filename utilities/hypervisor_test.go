package utilities

import (
	"errors"
	"testing"

	"libvirt.org/go/libvirt"
)

func TestConnectWrappers(t *testing.T) {
	origConnect := newConnect
	origRO := newConnectReadOnly
	origAuth := newConnectWithAuth
	origAuthDef := newConnectWithAuthDefault
	t.Cleanup(func() {
		newConnect = origConnect
		newConnectReadOnly = origRO
		newConnectWithAuth = origAuth
		newConnectWithAuthDefault = origAuthDef
	})

	newConnect = func(uri string) (*libvirt.Connect, error) {
		return nil, errors.New("fail")
	}
	if _, err := Connect("u"); err == nil {
		t.Fatal("Connect expected error")
	}
	newConnect = func(uri string) (*libvirt.Connect, error) {
		return nil, nil
	}
	if _, err := Connect("u"); err == nil {
		t.Fatal("Connect expected nil conn error")
	}

	newConnectReadOnly = func(uri string) (*libvirt.Connect, error) {
		return nil, errors.New("fail")
	}
	if _, err := ConnectReadOnly("u"); err == nil {
		t.Fatal("ConnectReadOnly expected error")
	}

	newConnectWithAuth = func(uri string, auth *libvirt.ConnectAuth, flags libvirt.ConnectFlags) (*libvirt.Connect, error) {
		return nil, errors.New("fail")
	}
	if _, err := ConnectWithAuth("u", nil, 0); err == nil {
		t.Fatal("ConnectWithAuth expected error")
	}

	newConnectWithAuthDefault = func(uri string, flags libvirt.ConnectFlags) (*libvirt.Connect, error) {
		return nil, errors.New("fail")
	}
	if _, err := ConnectWithAuthDefault("u", 0); err == nil {
		t.Fatal("ConnectWithAuthDefault expected error")
	}
}
