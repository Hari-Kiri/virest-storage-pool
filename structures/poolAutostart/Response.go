package poolAutostart

import "libvirt.org/go/libvirt"

type Response struct {
	Response bool          `json:"response"`
	Code     int           `json:"code"`
	Data     Uuid          `json:"data"`
	Error    libvirt.Error `json:"error"`
}

type Uuid struct {
	Uuid string `json:"uuid"`
}
