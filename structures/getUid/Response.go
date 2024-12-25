package getUid

import "libvirt.org/go/libvirt"

type Response struct {
	Response bool          `json:"response"`
	Code     int           `json:"code"`
	Data     Uid           `json:"data"`
	Error    libvirt.Error `json:"error"`
}

type Uid struct {
	Uid int `json:"uid"`
}
