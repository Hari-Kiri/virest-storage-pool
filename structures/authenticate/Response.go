package authenticate

import "libvirt.org/go/libvirt"

type Response struct {
	Response bool          `json:"response"`
	Code     int           `json:"code"`
	Data     Token         `json:"data"`
	Error    libvirt.Error `json:"error"`
}

type Token struct {
	Token string `json:"token"`
}
