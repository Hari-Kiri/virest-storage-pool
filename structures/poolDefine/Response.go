package poolDefine

import "libvirt.org/go/libvirt"

type Response struct {
	Response bool          `json:"response"`
	Code     int           `json:"code"`
	Message  Uuid          `json:"message"`
	Error    libvirt.Error `json:"error"`
}

type Uuid struct {
	Uuid string `json:"uuid"`
}
