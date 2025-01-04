package findStoragePoolSources

import (
	"github.com/Hari-Kiri/virest-utilities/utils/structures/libvirtxml"
)

type Request struct {
	Type    string `json:"type"`
	SrcSpec Source `json:"srcSpec"`
}

type Source struct {
	Source libvirtxml.Source `json:"source"`
}
