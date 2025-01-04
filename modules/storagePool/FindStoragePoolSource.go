package storagePool

import (
	"fmt"

	"github.com/Hari-Kiri/virest-storage-pool/structures/findStoragePoolSources"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/libvirtxml"
	"github.com/Hari-Kiri/virest-utilities/utils/structures/virest"
	"libvirt.org/go/libvirt"
)

// Talks to a storage backend and attempts to auto-discover the set of available storage pool sources. e.g.
// For iSCSI this would be a set of iSCSI targets. For NFS this would be a list of exported paths. The srcSpec
// (optional for some storage pool types, e.g. local ones) is an instance of the storage pool's source element
// specifying where to look for the pools.
//
// srcSpec is not required for some types (e.g., those querying local storage resources only)
func FindStoragePoolSource(connection virest.Connection, pooltype string, srcSpec libvirtxml.Source) (findStoragePoolSources.Sources, virest.Error, bool) {
	var (
		srcSpecXml  string
		virestError virest.Error
		isError     bool
	)

	srcSpecXml, virestError, isError = srcSpec.Marshal()
	if isError {
		virestError.Message = fmt.Sprintf("failed to parse source spec of the storage pool: %s", virestError.Message)
		return findStoragePoolSources.Sources{}, virestError, isError
	}

	// extra flags; not used yet, so callers should always pass 0
	// https://libvirt.org/html/libvirt-libvirt-storage.html#virConnectFindStoragePoolSources
	discoverStoragePoolSources, errorDiscoverStoragePoolSources := connection.FindStoragePoolSources(pooltype, srcSpecXml, 0)
	virestError.Error, isError = errorDiscoverStoragePoolSources.(libvirt.Error)
	if isError {
		virestError.Message = fmt.Sprintf("failed to find potential storage pool sources: %s", virestError.Message)
		return findStoragePoolSources.Sources{}, virestError, isError
	}

	var result libvirtxml.Sources
	virestError, isError = result.Unmarshal(discoverStoragePoolSources)
	if isError {
		virestError.Message = fmt.Sprintf("failed parse discovered storage pool sources: %s", virestError.Message)
		return findStoragePoolSources.Sources{}, virestError, isError
	}

	return findStoragePoolSources.Sources{Sources: result}, virestError, false
}
