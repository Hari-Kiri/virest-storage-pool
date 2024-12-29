package getGid

import (
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
)

type Response struct {
	Response bool         `json:"response"`
	Code     int          `json:"code"`
	Data     Gid          `json:"data"`
	Error    virest.Error `json:"error"`
}

type Gid struct {
	Gid int `json:"gid"`
}
