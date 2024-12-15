package poolDetail

import (
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

type Response struct {
	Response bool                   `json:"response"`
	Code     int                    `json:"code"`
	Data     libvirtxml.StoragePool `json:"data"`
	Error    libvirt.Error          `json:"error"`
}
