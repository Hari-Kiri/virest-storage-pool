package poolDelete

import "libvirt.org/go/libvirt"

type Response struct {
	Response bool          `json:"response"`
	Code     int           `json:"code"`
	Error    libvirt.Error `json:"error"`
}
