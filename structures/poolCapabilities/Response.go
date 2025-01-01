package poolCapabilities

import (
	"github.com/Hari-Kiri/virest-utilities/utils/structures/libvirtxml"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
)

type Response struct {
	Response bool                    `json:"response"`
	Code     int                     `json:"code"`
	Data     StoragepoolCapabilities `json:"data"`
	Error    virest.Error            `json:"error"`
}

type StoragepoolCapabilities struct {
	libvirtxml.StoragepoolCapabilities `json:"capabilities"`
}
