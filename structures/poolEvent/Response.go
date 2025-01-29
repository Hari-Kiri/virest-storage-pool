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

	// Lifecycle event type default is 6, when occur a virStoragePoolEventLifecycleType (which
	// is less than 6) is emitted during storage pool lifecycle events. See more about Storage
	// Pool Event Type:
	// https://libvirt.org/html/libvirt-libvirt-storage.html#virStoragePoolEventLifecycleType
	EventLifecycle libvirt.StoragePoolEventLifecycle `json:"eventLifecycle"`

	// Refresh event default is 0, when occur it will change to 1.
	EventRefresh int `json:"eventRefresh"`
}
