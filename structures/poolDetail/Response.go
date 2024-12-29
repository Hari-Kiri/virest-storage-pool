package poolDetail

import (
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirtxml"
)

type Response struct {
	Response bool         `json:"response"`
	Code     int          `json:"code"`
	Data     Detail       `json:"data"`
	Error    virest.Error `json:"error"`
}

type Detail struct {
	libvirtxml.StoragePool
}
