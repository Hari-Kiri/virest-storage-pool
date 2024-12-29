package poolList

import (
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

type Response struct {
	Response bool         `json:"response"`
	Code     int          `json:"code"`
	Data     []Data       `json:"data"`
	Error    virest.Error `json:"error"`
}

type Data struct {
	Uuid       string                     `json:"uuid"`
	Name       string                     `json:"name"`
	State      libvirt.StoragePoolState   `json:"state"`
	Autostart  bool                       `json:"autostart"`
	Persistent bool                       `json:"persistent"`
	Capacity   libvirtxml.StoragePoolSize `json:"capacity"`
	Allocation libvirtxml.StoragePoolSize `json:"allocation"`
	Available  libvirtxml.StoragePoolSize `json:"available"`
}
