package poolInfo

import (
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

type Response struct {
	Response bool          `json:"response"`
	Code     int           `json:"code"`
	Data     Info          `json:"data"`
	Error    libvirt.Error `json:"error"`
}

type Info struct {
	Name       string                     `json:"name"`
	Uuid       string                     `json:"uuid"`
	State      libvirt.StoragePoolState   `json:"state"`
	Persistent bool                       `json:"persistent"`
	Autostart  bool                       `json:"autostart"`
	Capacity   libvirtxml.StoragePoolSize `json:"capacity"`
	Allocation libvirtxml.StoragePoolSize `json:"allocation"`
	Available  libvirtxml.StoragePoolSize `json:"available"`
}
