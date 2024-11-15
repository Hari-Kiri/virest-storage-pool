package poolList

import (
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

type Response struct {
	Response bool          `json:"response"`
	Code     int           `json:"code"`
	Data     []Data        `json:"data"`
	Error    libvirt.Error `json:"error"`
}

type Data struct {
	Config     libvirtxml.StoragePool   `json:"config"`
	State      libvirt.StoragePoolState `json:"state"` // Reference https://libvirt.org/html/libvirt-libvirt-storage.html#virStoragePoolState
	Autostart  bool                     `json:"autostart"`
	Persistent bool                     `json:"persistent"`
}
