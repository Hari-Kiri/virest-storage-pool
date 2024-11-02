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

// This function should be called first to get a restricted connection to the library functionalities.
// The set of APIs usable are then restricted on the available methods to control the domains.
//
// See https://pkg.go.dev/github.com/Hari-Kiri/virest-storage-pool/modules/utils#NewConnect or
// https://libvirt.org/html/libvirt-libvirt-host.html#virConnectOpen for notes about environment variables
// which can have an effect on opening drivers and freeing the connection resources
//
// URIs are documented at https://libvirt.org/uri.html
func NewConnectReadOnly(uri string) (*libvirt.Connect, libvirt.Error, bool) {
	result, errorResult := libvirt.NewConnectReadOnly(uri)
	libvirtError, isError := errorResult.(libvirt.Error)
	return result, libvirtError, isError
}

// This function should be called first to get a connection to the Hypervisor. If necessary, authentication
// will be performed fetching credentials via the callback
//
// See https://pkg.go.dev/github.com/Hari-Kiri/virest-storage-pool/modules/utils#NewConnect or
// https://libvirt.org/html/libvirt-libvirt-host.html#virConnectOpen for notes about environment variables
// which can have an effect on opening drivers and freeing the connection resources
//
// URIs are documented at https://libvirt.org/uri.html
func NewConnectWithAuth(uri string, auth *libvirt.ConnectAuth, flags libvirt.ConnectFlags) (*libvirt.Connect, libvirt.Error, bool) {
	result, errorResult := libvirt.NewConnectWithAuth(uri, auth, flags)
	libvirtError, isError := errorResult.(libvirt.Error)
	return result, libvirtError, isError
}

// This function should be called first to get a connection to the Hypervisor. If necessary, authentication
// will be performed fetching credentials via the callback
//
// See https://pkg.go.dev/github.com/Hari-Kiri/virest-storage-pool/modules/utils#NewConnect or
// https://libvirt.org/html/libvirt-libvirt-host.html#virConnectOpen for notes about environment variables
// which can have an effect on opening drivers and freeing the connection resources
//
// URIs are documented at https://libvirt.org/uri.html
func NewConnectWithAuthDefault(uri string, flags libvirt.ConnectFlags) (*libvirt.Connect, libvirt.Error, bool) {
	result, errorResult := libvirt.NewConnectWithAuthDefault(uri, flags)
	libvirtError, isError := errorResult.(libvirt.Error)
	return result, libvirtError, isError
}
