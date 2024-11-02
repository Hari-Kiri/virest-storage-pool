package utils

import "libvirt.org/go/libvirt"

// This function should be called first to get a connection to the Hypervisor and xen store.
//
// If name is NULL, if the LIBVIRT_DEFAULT_URI environment variable is set, then it will be used.
// Otherwise if the client configuration file has the "uri_default" parameter set, then it will be used.
// Finally probing will be done to determine a suitable default driver to activate.
// This involves trying each hypervisor in turn until one successfully opens.
//
// If connecting to an unprivileged hypervisor driver which requires the libvirtd daemon to be active,
// it will automatically be launched if not already running. This can be prevented by setting the environment variable LIBVIRT_AUTOSTART=0
//
// URIs are documented at https://libvirt.org/uri.html
//
// Close() should be used to release the resources after the connection is no longer needed.
func NewConnect(uri string) (*libvirt.Connect, libvirt.Error, bool) {
	result, errorResult := libvirt.NewConnect(uri)
	libvirtError, isError := errorResult.(libvirt.Error)
	return result, libvirtError, isError
}
