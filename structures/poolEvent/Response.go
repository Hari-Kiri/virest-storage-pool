package poolEvent

import (
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

type Response struct {
	Response bool         `json:"response"`
	Code     int          `json:"code"`
	Data     Info         `json:"data"`
	Error    virest.Error `json:"error"`
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
