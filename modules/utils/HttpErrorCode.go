package utils

import "libvirt.org/go/libvirt"

// Convert the Libvirt Error Number to HTTP Status Code
func HttpErrorCode(libvirtErrorNumber libvirt.ErrorNumber) int {
	result := 500

	if libvirtErrorNumber == libvirt.ERR_OK {
		result = 200
	}

	if libvirtErrorNumber == libvirt.ERR_NO_MEMORY {
		result = 507
	}

	if libvirtErrorNumber == libvirt.ERR_NO_SUPPORT {
		result = 501
	}

	if libvirtErrorNumber == libvirt.ERR_UNKNOWN_HOST {
		result = 504
	}

	if libvirtErrorNumber == libvirt.ERR_NO_CONNECT {
		result = 504
	}

	if libvirtErrorNumber == libvirt.ERR_INVALID_CONN {
		result = 400
	}

	if libvirtErrorNumber == libvirt.ERR_INVALID_DOMAIN {
		result = 400
	}

	if libvirtErrorNumber == libvirt.ERR_INVALID_ARG {
		result = 400
	}

	if libvirtErrorNumber == libvirt.ERR_GET_FAILED {
		result = 405
	}

	if libvirtErrorNumber == libvirt.ERR_POST_FAILED {
		result = 405
	}

	if libvirtErrorNumber == libvirt.ERR_HTTP_ERROR {
		result = 405
	}

	if libvirtErrorNumber == libvirt.ERR_XML_ERROR {
		result = 406
	}

	if libvirtErrorNumber == libvirt.ERR_NO_STORAGE_POOL {
		result = 404
	}

	if libvirtErrorNumber == libvirt.ERR_NO_STORAGE_VOL {
		result = 404
	}

	if libvirtErrorNumber == libvirt.ERR_XML_INVALID_SCHEMA {
		result = 406
	}

	return result
}
