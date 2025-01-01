package findStoragePoolSources

import (
	"libvirt.org/go/libvirtxml"
)

type Request struct {
	Type    string                       `json:"type"`
	SrcSpec libvirtxml.StoragePoolSource `json:"srcSpec"`
}
