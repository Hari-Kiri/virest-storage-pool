package poolEvent

import (
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

type Response struct {
	Response bool         `json:"response"`
	Code     int          `json:"code"`
	Data     Event        `json:"data"`
	Error    virest.Error `json:"error"`
}

type Event struct {
	EventRefresh   int                               `json:"eventRefresh"`
	EventLifecycle libvirt.StoragePoolEventLifecycle `json:"eventLifecycle"`
}
