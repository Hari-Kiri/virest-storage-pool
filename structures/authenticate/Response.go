package authenticate

import "github.com/Hari-Kiri/virest-utilities/utils/structures/virest"

type Response struct {
	Response bool         `json:"response"`
	Code     int          `json:"code"`
	Data     Token        `json:"data"`
	Error    virest.Error `json:"error"`
}

type Token struct {
	Token string `json:"token"`
}
