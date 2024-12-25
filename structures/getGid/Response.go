package getGid

import "libvirt.org/go/libvirt"

type Response struct {
	Response bool          `json:"response"`
	Code     int           `json:"code"`
	Data     Gid           `json:"data"`
	Error    libvirt.Error `json:"error"`
}

type Gid struct {
	Gid int `json:"gid"`
}
