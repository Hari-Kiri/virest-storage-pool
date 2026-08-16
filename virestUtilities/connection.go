package virestUtilities

import "libvirt.org/go/libvirt"

// Connection wraps a libvirt hypervisor connection.
type Connection struct {
	*libvirt.Connect
}
