package utilities

import (
	"errors"
	"testing"

	"libvirt.org/go/libvirt"
)

func TestWrapPositiveNegative(t *testing.T) {
	if Wrap("op", nil) != nil {
		t.Fatal("nil should stay nil")
	}
	err := Wrap("op", errors.New("x"))
	if err == nil || err.Error() != "op: x" {
		t.Fatalf("got %v", err)
	}
}

func TestAsLibvirtErrorAndCode(t *testing.T) {
	lv := libvirt.Error{Code: libvirt.ERR_NO_STORAGE_POOL, Message: "gone"}
	got, ok := AsLibvirtError(lv)
	if !ok || got.Code != libvirt.ERR_NO_STORAGE_POOL {
		t.Fatalf("AsLibvirtError: %+v ok=%v", got, ok)
	}
	code, ok := Code(lv)
	if !ok || code != libvirt.ERR_NO_STORAGE_POOL {
		t.Fatalf("Code: %v ok=%v", code, ok)
	}

	if _, ok := AsLibvirtError(errors.New("plain")); ok {
		t.Fatal("plain should not unwrap")
	}
	if _, ok := Code(errors.New("plain")); ok {
		t.Fatal("plain Code should fail")
	}
	wrapped := Wrap("lookup", lv)
	code, ok = Code(wrapped)
	if !ok || code != libvirt.ERR_NO_STORAGE_POOL {
		t.Fatalf("wrapped Code: %v ok=%v", code, ok)
	}
}

func TestWrapConnect(t *testing.T) {
	_, err := wrapConnect(nil, errors.New("dial"), "connect")
	if err == nil {
		t.Fatal("expected dial error")
	}
	_, err = wrapConnect(nil, nil, "connect")
	if err == nil {
		t.Fatal("expected nil connection error")
	}
}
