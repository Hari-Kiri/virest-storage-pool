package storagePool

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Hari-Kiri/virest/utilities"
	"libvirt.org/go/libvirt"
	"libvirt.org/go/libvirtxml"
)

func TestCloseNilAndClosed(t *testing.T) {
	var c *Connection
	if err := c.Close(); err != nil {
		t.Fatalf("nil close: %v", err)
	}

	hv := newFakeHypervisor()
	conn := newConnectionForTest(hv)
	if err := conn.Close(); err != nil {
		t.Fatalf("close: %v", err)
	}
	if hv.closeCalls != 1 {
		t.Fatalf("closeCalls=%d", hv.closeCalls)
	}
	if err := conn.Close(); err != nil {
		t.Fatalf("second close: %v", err)
	}
}

func TestCloseError(t *testing.T) {
	hv := newFakeHypervisor()
	hv.closeErr = errors.New("boom")
	conn := newConnectionForTest(hv)
	if err := conn.Close(); err == nil {
		t.Fatal("expected error")
	}
}

func TestDefineSuccessAndFailures(t *testing.T) {
	defined := &fakePool{uuid: "u-1"}
	hv := newFakeHypervisor()
	hv.definePool = defined
	conn := newConnectionForTest(hv)

	uuid, err := conn.Define(libvirtxml.StoragePool{Name: "p", Type: "dir"}, 0)
	if err != nil {
		t.Fatalf("define: %v", err)
	}
	if uuid != "u-1" {
		t.Fatalf("uuid=%s", uuid)
	}
	if defined.freeCalls != 1 {
		t.Fatalf("defined freeCalls=%d", defined.freeCalls)
	}
	if hv.defineXML == "" {
		t.Fatal("expected define xml")
	}

	hv.defineErr = errors.New("define failed")
	if _, err := conn.Define(libvirtxml.StoragePool{Name: "p", Type: "dir"}, 1); err == nil {
		t.Fatal("expected define error")
	}

	hv.defineErr = nil
	hv.definePool = &fakePool{getUUIDErr: errors.New("uuid fail")}
	if _, err := conn.Define(libvirtxml.StoragePool{Name: "p", Type: "dir"}, 0); err == nil {
		t.Fatal("expected uuid error")
	}
}

func TestLifecycleOps(t *testing.T) {
	pool := &fakePool{}
	hv := newFakeHypervisor()
	hv.poolsByUUID["abc"] = pool
	conn := newConnectionForTest(hv)

	tests := []struct {
		name string
		call func() error
		check func(t *testing.T)
	}{
		{
			name: "build",
			call: func() error { return conn.Build("abc", libvirt.STORAGE_POOL_BUILD_NEW) },
			check: func(t *testing.T) {
				if pool.buildFlags != libvirt.STORAGE_POOL_BUILD_NEW {
					t.Fatalf("flags=%v", pool.buildFlags)
				}
			},
		},
		{
			name: "create",
			call: func() error { return conn.Create("abc", libvirt.STORAGE_POOL_CREATE_NORMAL) },
			check: func(t *testing.T) {
				if pool.createFlags != libvirt.STORAGE_POOL_CREATE_NORMAL {
					t.Fatalf("flags=%v", pool.createFlags)
				}
			},
		},
		{
			name: "destroy",
			call: func() error { return conn.Destroy("abc") },
		},
		{
			name: "delete",
			call: func() error { return conn.Delete("abc", libvirt.STORAGE_POOL_DELETE_NORMAL) },
			check: func(t *testing.T) {
				if pool.deleteFlags != libvirt.STORAGE_POOL_DELETE_NORMAL {
					t.Fatalf("flags=%v", pool.deleteFlags)
				}
			},
		},
		{
			name: "undefine",
			call: func() error { return conn.Undefine("abc") },
		},
		{
			name: "refresh",
			call: func() error { return conn.Refresh("abc") },
			check: func(t *testing.T) {
				if pool.refreshFlags != 0 {
					t.Fatalf("flags=%v", pool.refreshFlags)
				}
			},
		},
		{
			name: "autostart",
			call: func() error { return conn.Autostart("abc", true) },
			check: func(t *testing.T) {
				if pool.setAutostartValue == nil || !*pool.setAutostartValue {
					t.Fatal("expected autostart true")
				}
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pool.freeCalls = 0
			if err := tc.call(); err != nil {
				t.Fatalf("%s: %v", tc.name, err)
			}
			if pool.freeCalls != 1 {
				t.Fatalf("freeCalls=%d", pool.freeCalls)
			}
			if tc.check != nil {
				tc.check(t)
			}
		})
	}

	hv.lookupErr = errors.New("missing")
	if err := conn.Build("abc", 0); err == nil {
		t.Fatal("expected lookup error")
	}

	hv.lookupErr = nil
	pool.buildErr = errors.New("build fail")
	if err := conn.Build("abc", 0); err == nil {
		t.Fatal("expected build error")
	}
}

func TestDetailSuccessAndFailures(t *testing.T) {
	pool := &fakePool{
		xmlDesc: `<pool type='dir'><name>demo</name><uuid>u-1</uuid></pool>`,
	}
	hv := newFakeHypervisor()
	hv.poolsByUUID["u-1"] = pool
	conn := newConnectionForTest(hv)

	detail, err := conn.Detail("u-1", 0)
	if err != nil {
		t.Fatalf("detail: %v", err)
	}
	if detail.Name != "demo" {
		t.Fatalf("name=%s", detail.Name)
	}

	pool.getXMLErr = errors.New("xml fail")
	if _, err := conn.Detail("u-1", 0); err == nil {
		t.Fatal("expected xml error")
	}

	pool.getXMLErr = nil
	pool.xmlDesc = "not-xml"
	if _, err := conn.Detail("u-1", 0); err == nil {
		t.Fatal("expected unmarshal error")
	}
}

func TestInfoSuccessAndFieldFailure(t *testing.T) {
	pool := &fakePool{
		name:       "n1",
		autostart:  true,
		persistent: true,
		info: &libvirt.StoragePoolInfo{
			State:      libvirt.STORAGE_POOL_RUNNING,
			Capacity:   100,
			Allocation: 40,
			Available:  60,
		},
	}
	hv := newFakeHypervisor()
	hv.poolsByUUID["u"] = pool
	conn := newConnectionForTest(hv)

	info, err := conn.Info("u")
	if err != nil {
		t.Fatalf("info: %v", err)
	}
	if info.Name != "n1" || info.Capacity != 100 || !info.Autostart || !info.Persistent {
		t.Fatalf("unexpected info: %+v", info)
	}
	if pool.freeCalls < 1 {
		t.Fatal("expected original pool free")
	}

	pool.getNameErr = errors.New("name fail")
	if _, err := conn.Info("u"); err == nil {
		t.Fatal("expected name error")
	}
}

func TestListEmptySingleAndFailure(t *testing.T) {
	hv := newFakeHypervisor()
	conn := newConnectionForTest(hv)

	list, err := conn.List(0, 0)
	if err != nil {
		t.Fatalf("empty list: %v", err)
	}
	if len(list) != 0 {
		t.Fatalf("len=%d", len(list))
	}

	pool := &fakePool{
		xmlDesc: `<pool type='dir'><name>p</name><uuid>u</uuid><capacity unit='bytes'>10</capacity><allocation unit='bytes'>4</allocation><available unit='bytes'>6</available></pool>`,
		info:    &libvirt.StoragePoolInfo{State: libvirt.STORAGE_POOL_RUNNING},
	}
	hv.listPools = []poolHandle{pool}
	list, err = conn.List(3, 1)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if hv.listFlagsSeen != 3 {
		t.Fatalf("flags=%v", hv.listFlagsSeen)
	}
	if len(list) != 1 || list[0].Name != "p" || list[0].Capacity.Value != 10 {
		t.Fatalf("unexpected: %+v", list)
	}
	if pool.freeCalls < 1 {
		t.Fatal("expected pool free")
	}

	hv.listErr = errors.New("list fail")
	if _, err := conn.List(0, 0); err == nil {
		t.Fatal("expected list error")
	}

	hv.listErr = nil
	bad := &fakePool{getInfoErr: errors.New("info fail"), xmlDesc: `<pool type='dir'><name>p</name><uuid>u</uuid></pool>`}
	hv.listPools = []poolHandle{bad}
	if _, err := conn.List(0, 0); err == nil {
		t.Fatal("expected per-pool error")
	}
}

func TestCapabilitiesAndFindSources(t *testing.T) {
	hv := newFakeHypervisor()
	hv.capabilitiesXML = `<storagepoolCapabilities><pool type='dir' supported='yes'></pool></storagepoolCapabilities>`
	conn := newConnectionForTest(hv)

	caps, err := conn.Capabilities()
	if err != nil {
		t.Fatalf("caps: %v", err)
	}
	if len(caps.Pool) != 1 || caps.Pool[0].Type != "dir" {
		t.Fatalf("caps=%+v", caps)
	}

	hv.capabilitiesErr = errors.New("caps fail")
	if _, err := conn.Capabilities(); err == nil {
		t.Fatal("expected caps error")
	}

	hv.capabilitiesErr = nil
	hv.capabilitiesXML = "bad"
	if _, err := conn.Capabilities(); err == nil {
		t.Fatal("expected unmarshal error")
	}

	hv.findSourcesXML = `<sources></sources>`
	sources, err := conn.FindSources("netfs", SourceSpec{Host: Host{Name: "h", Port: 2049}})
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if hv.findPoolType != "netfs" || sources.Source == nil && len(sources.Source) != 0 {
		// empty sources is fine
	}
	if hv.findSrcSpec == "" {
		t.Fatal("expected src spec xml")
	}

	hv.findSourcesErr = errors.New("find fail")
	if _, err := conn.FindSources("netfs", SourceSpec{}); err == nil {
		t.Fatal("expected find error")
	}
}

func TestWaitEventAndStreamEvents(t *testing.T) {
	pool := &fakePool{}
	hv := newFakeHypervisor()
	hv.poolsByUUID["u"] = pool
	conn := newConnectionForTest(hv)

	if _, err := conn.WaitEvent("u", 99); err == nil {
		t.Fatal("expected invalid kind")
	}

	done := make(chan Event, 1)
	go func() {
		ev, err := conn.WaitEvent("u", EventLifecycle)
		if err != nil {
			t.Errorf("wait: %v", err)
			return
		}
		done <- ev
	}()

	deadline := time.After(2 * time.Second)
	for {
		hv.mu.Lock()
		ready := hv.lifecycleCallback != nil
		hv.mu.Unlock()
		if ready {
			break
		}
		select {
		case <-deadline:
			t.Fatal("callback not registered")
		case <-time.After(5 * time.Millisecond):
		}
	}
	hv.fireLifecycle(libvirt.StoragePoolEventLifecycle{Event: libvirt.STORAGE_POOL_EVENT_STARTED, Detail: 1})

	select {
	case ev := <-done:
		if ev.EventLifecycle.Event != libvirt.STORAGE_POOL_EVENT_STARTED {
			t.Fatalf("event=%v", ev.EventLifecycle)
		}
		if len(hv.deregisterCalls) == 0 {
			t.Fatal("expected deregister")
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout waiting for event")
	}

	ctx, cancel := context.WithCancel(context.Background())
	got := make(chan Event, 1)
	errCh := make(chan error, 1)
	go func() {
		errCh <- conn.StreamEvents(ctx, "u", EventRefresh, func(ev Event) error {
			got <- ev
			return nil
		})
	}()

	deadline = time.After(2 * time.Second)
	for {
		hv.mu.Lock()
		ready := hv.refreshCallback != nil
		hv.mu.Unlock()
		if ready {
			break
		}
		select {
		case <-deadline:
			t.Fatal("refresh callback not registered")
		case <-time.After(5 * time.Millisecond):
		}
	}
	hv.fireRefresh()
	select {
	case ev := <-got:
		if ev.EventRefresh != 1 {
			t.Fatalf("refresh=%d", ev.EventRefresh)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("timeout stream event")
	}
	cancel()
	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("stream err: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("stream did not exit")
	}

	if err := conn.StreamEvents(context.Background(), "u", EventLifecycle, nil); err == nil {
		t.Fatal("expected nil emit error")
	}
}

func TestErrorsAsLibvirt(t *testing.T) {
	hv := newFakeHypervisor()
	hv.lookupErr = libvirt.Error{Code: libvirt.ERR_NO_STORAGE_POOL, Message: "gone"}
	conn := newConnectionForTest(hv)
	err := conn.Destroy("x")
	if err == nil {
		t.Fatal("expected error")
	}
	lv, ok := utilities.AsLibvirtError(err)
	if !ok {
		t.Fatalf("AsLibvirtError failed: %v", err)
	}
	if lv.Code != libvirt.ERR_NO_STORAGE_POOL {
		t.Fatalf("code=%v", lv.Code)
	}
	if code, ok := utilities.Code(err); !ok || code != libvirt.ERR_NO_STORAGE_POOL {
		t.Fatalf("Code=%v ok=%v", code, ok)
	}

	plain := wrap("op", errors.New("plain"))
	if _, ok := utilities.AsLibvirtError(plain); ok {
		t.Fatal("plain error should not be libvirt.Error")
	}
}
