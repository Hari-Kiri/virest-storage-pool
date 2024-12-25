package poolDelete

import "libvirt.org/go/libvirt"

type Request struct {
	Uuid   string                         `json:"uuid"`
	Option libvirt.StoragePoolDeleteFlags `json:"option"`
}
