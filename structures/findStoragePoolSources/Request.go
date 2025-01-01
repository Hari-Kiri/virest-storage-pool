package findStoragePoolSources

import (
	"github.com/Hari-Kiri/virest-storage-pool/structures"
)

type Request struct {
	Type    string `json:"type"`
	SrcSpec Source `json:"srcSpec"`
}

type Source struct {
	Source structures.Source `json:"source"`
}
