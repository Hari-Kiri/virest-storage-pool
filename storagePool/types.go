package storagePool

import (
	"encoding/xml"

	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

// Info is volatile information about a storage pool.
type Info struct {
	Name       string                   `json:"name"`
	Uuid       string                   `json:"uuid"`
	State      libvirt.StoragePoolState `json:"state"`
	Persistent bool                     `json:"persistent"`
	Autostart  bool                     `json:"autostart"`
	Capacity   uint64                   `json:"capacity"`
	Allocation uint64                   `json:"allocation"`
	Available  uint64                   `json:"available"`
}

// Detail is the full storage pool XML description.
type Detail struct {
	libvirtxml.StoragePool
}

// PoolSummary is a condensed view of a storage pool used by List.
type PoolSummary struct {
	Uuid       string                     `json:"uuid"`
	Name       string                     `json:"name"`
	State      libvirt.StoragePoolState   `json:"state"`
	Autostart  bool                       `json:"autostart"`
	Persistent bool                       `json:"persistent"`
	Capacity   libvirtxml.StoragePoolSize `json:"capacity"`
	Allocation libvirtxml.StoragePoolSize `json:"allocation"`
	Available  libvirtxml.StoragePoolSize `json:"available"`
}

// EventKind selects which storage-pool event to wait for or stream.
type EventKind uint

const (
	// EventLifecycle waits for storage pool lifecycle events.
	EventLifecycle EventKind = 0
	// EventRefresh waits for storage pool refresh events.
	EventRefresh EventKind = 1
)

// Event is a storage pool lifecycle or refresh notification.
type Event struct {
	EventLifecycle libvirt.StoragePoolEventLifecycle `json:"eventLifecycle"`
	EventRefresh   int                               `json:"eventRefresh"`
	Timestamp      int64                             `json:"timestamp"`
	TimestampNano  int64                             `json:"timestampNano"`
	EventId        int                               `json:"eventId"`
}

// SourceSpec is an optional storage-pool source element used by FindSources.
type SourceSpec struct {
	XMLName   xml.Name  `xml:"source" json:"-"`
	Host      Host      `xml:"host" json:"host"`
	Initiator Initiator `xml:"initiator" json:"initiator"`
}

// Host is a source host attribute pair.
type Host struct {
	Name string `xml:"name,attr" json:"name"`
	Port int    `xml:"port,attr" json:"port"`
}

// Initiator holds an iSCSI IQN.
type Initiator struct {
	Iqn Iqn `xml:"iqn" json:"iqn"`
}

// Iqn is an iSCSI qualified name.
type Iqn struct {
	Name string `xml:"name,attr" json:"name"`
}

// Sources is the discovered set of storage pool sources.
type Sources struct {
	XMLName xml.Name                       `xml:"sources" json:"-"`
	Source  []libvirtxml.StoragePoolSource `xml:"source" json:"source"`
}

// Capabilities describes supported storage pool types and formats.
type Capabilities struct {
	XMLName xml.Name         `xml:"storagepoolCapabilities" json:"-"`
	Pool    []CapabilityPool `xml:"pool" json:"pool"`
}

// CapabilityPool is one pool type entry in Capabilities.
type CapabilityPool struct {
	Type        string      `xml:"type,attr" json:"type"`
	Supported   string      `xml:"supported,attr" json:"supported"`
	VolOptions  VolOptions  `xml:"volOptions" json:"volOptions"`
	PoolOptions PoolOptions `xml:"poolOptions" json:"poolOptions"`
}

// VolOptions describes volume format options for a pool type.
type VolOptions struct {
	DefaultFormat DefaultFormat `xml:"defaultFormat" json:"defaultFormat"`
	Enum          Enum          `xml:"enum" json:"enum"`
}

// PoolOptions describes pool format options for a pool type.
type PoolOptions struct {
	DefaultFormat DefaultFormat `xml:"defaultFormat" json:"defaultFormat"`
	Enum          Enum          `xml:"enum" json:"enum"`
}

// DefaultFormat is a default format type attribute.
type DefaultFormat struct {
	Type string `xml:"type,attr" json:"type"`
}

// Enum lists allowed format values.
type Enum struct {
	Name  string   `xml:"name,attr" json:"name"`
	Value []string `xml:"value" json:"value"`
}
