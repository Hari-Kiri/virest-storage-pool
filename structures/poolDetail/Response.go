package poolDetail

import (
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

type Response struct {
	Response bool          `json:"response"`
	Code     int           `json:"code"`
	Data     Detail        `json:"data"`
	Error    libvirt.Error `json:"error"`
}

type Detail struct {
	libvirtxml.StoragePool
}
