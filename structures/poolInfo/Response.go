package poolInfo

import (
	"libvirt.org/go/libvirt"
)

type Response struct {
	Response bool          `json:"response"`
	Code     int           `json:"code"`
	Data     Info          `json:"data"`
	Error    libvirt.Error `json:"error"`
}

type Info struct {
	Name       string                   `json:"name"`
	Uuid       string                   `json:"uuid"`
	State      libvirt.StoragePoolState `json:"state"`
	Persistent bool                     `json:"persistent"`
	Autostart  bool                     `json:"autostart"`
	Capacity   uint64                   `json:"capacity"`
	Allocation uint64                   `json:"allocation"`
	Available  uint64                   `json:"available"`
}
