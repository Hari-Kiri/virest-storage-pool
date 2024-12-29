package poolDestroy

import "github.com/Hari-Kiri/virest-utilities/utils/structures/virest"

type Response struct {
	Response bool         `json:"response"`
	Code     int          `json:"code"`
	Data     Uuid         `json:"data"`
	Error    virest.Error `json:"error"`
}

type Uuid struct {
	Uuid string `json:"uuid"`
}
