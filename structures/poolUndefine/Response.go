package poolUndefine

import "libvirt.org/go/libvirt"

type Response struct {
	Response bool          `json:"response"`
	Code     int           `json:"code"`
	Data     struct{}      `json:"data"`
	Error    libvirt.Error `json:"error"`
}
