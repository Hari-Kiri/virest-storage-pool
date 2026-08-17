package utilities

import (
	stderrors "errors"
	"fmt"

	"libvirt.org/go/libvirt"
)

// Wrap annotates err with an operation name. Nil input stays nil.
func Wrap(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", op, err)
}

// AsLibvirtError extracts a libvirt.Error from err when present.
func AsLibvirtError(err error) (libvirt.Error, bool) {
	var lv libvirt.Error
	if stderrors.As(err, &lv) {
		return lv, true
	}
	return libvirt.Error{}, false
}

// Code returns the libvirt error number when err unwraps to libvirt.Error.
func Code(err error) (libvirt.ErrorNumber, bool) {
	lv, ok := AsLibvirtError(err)
	if !ok {
		return 0, false
	}
	return lv.Code, true
}
