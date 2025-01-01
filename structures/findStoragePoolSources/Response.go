package findStoragePoolSources

import (
	"github.com/Hari-Kiri/virest-storage-pool/structures"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
)

type Response struct {
	Response bool               `json:"response"`
	Code     int                `json:"code"`
	Data     structures.Sources `json:"data"`
	Error    virest.Error       `json:"error"`
}
