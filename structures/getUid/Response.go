package getUid

import "github.com/Hari-Kiri/virest-utilities/utils/structures/virest"

type Response struct {
	Response bool         `json:"response"`
	Code     int          `json:"code"`
	Data     Uid          `json:"data"`
	Error    virest.Error `json:"error"`
}

type Uid struct {
	Uid int `json:"uid"`
}
