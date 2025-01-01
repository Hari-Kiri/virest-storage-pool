package findStoragePoolSources

import (
	"github.com/Hari-Kiri/virest-storage-pool/structures"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
)

type Response struct {
	Response bool         `json:"response"`
	Code     int          `json:"code"`
	Data     Sources      `json:"data"`
	Error    virest.Error `json:"error"`
}

type Sources struct {
	structures.Sources `json:"sources"`
}
