package poolCreate

import "libvirt.org/go/libvirt"

type Request struct {
	Uuid   string                         `json:"uuid"`
	Option libvirt.StoragePoolCreateFlags `json:"option"`
}
