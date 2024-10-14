package poolDefine

import (
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

type Request struct {
	PoolUuid    string                         `json:"poolUuid"`
	Option      libvirt.StoragePoolDefineFlags `json:"option"`
	StoragePool libvirtxml.StoragePool         `json:"storagePool"`
}
