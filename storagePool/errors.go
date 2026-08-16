package storagePool

import "github.com/Hari-Kiri/virest/virestUtilities"

// wrap annotates err with an operation name via virestUtilities.
var wrap = virestUtilities.Wrap

// finishFree frees p; if free fails and *err is still nil, *err becomes the free error.
func finishFree(p poolHandle, err *error) {
	freeErr := p.Free()
	if freeErr == nil || *err != nil {
		return
	}
	*err = wrap("free pool", freeErr)
}

// finishDeregister deregisters callbackID; if that fails and *err is still nil, *err becomes it.
func finishDeregister(hv hypervisor, callbackID int, err *error) {
	deregErr := hv.StoragePoolEventDeregister(callbackID)
	if deregErr == nil || *err != nil {
		return
	}
	*err = wrap("deregister event", deregErr)
}

// freePoolsFrom frees pools[from:] and returns the first free error, if any.
func freePoolsFrom(pools []poolHandle, from int) error {
	n := len(pools)
	var first error
	for i := from; i < n; i++ {
		freeErr := pools[i].Free()
		if freeErr != nil && first == nil {
			first = wrap("free pool", freeErr)
		}
	}
	return first
}
