package storagePool

import (
	"context"
	"errors"
	"fmt"
	"time"

	"libvirt.org/go/libvirt"
)

var (
	errEmitRequired         = errors.New("emit callback is required")
	errUnsupportedEventKind = errors.New("unsupported event kind")
)

// WaitEvent blocks until the selected storage-pool event occurs.
//
// Callers must register a default libvirt event implementation and run the
// event loop (EventRegisterDefaultImpl + EventRunDefaultImpl) in the process.
func (c *Connection) WaitEvent(uuid string, kind EventKind) (ev Event, err error) {
	if kind > EventRefresh {
		return Event{}, wrap("wait event", fmt.Errorf("%w: %d", errUnsupportedEventKind, kind))
	}

	pool, err := c.lookupPool(uuid)
	if err != nil {
		return Event{}, err
	}
	defer finishFree(pool, &err)

	switch kind {
	case EventLifecycle:
		ev, err = c.waitLifecycleEvent(pool)
		return ev, err
	case EventRefresh:
		ev, err = c.waitRefreshEvent(pool)
		return ev, err
	}
	return Event{}, wrap("wait event", fmt.Errorf("%w: %d", errUnsupportedEventKind, kind))
}

func (c *Connection) waitLifecycleEvent(pool poolHandle) (ev Event, err error) {
	ch := make(chan libvirt.StoragePoolEventLifecycle, 1)
	callbackID, regErr := c.hv.StoragePoolEventLifecycleRegister(pool, func(
		_ *libvirt.Connect,
		_ *libvirt.StoragePool,
		event *libvirt.StoragePoolEventLifecycle,
	) {
		if event == nil {
			return
		}
		ch <- *event
	})
	if regErr != nil {
		return Event{}, wrap("register lifecycle event", regErr)
	}
	defer finishDeregister(c.hv, callbackID, &err)

	lifecycle := <-ch
	now := time.Now()
	return Event{
		EventLifecycle: lifecycle,
		Timestamp:      now.Unix(),
		TimestampNano:  now.UnixNano(),
		EventId:        callbackID,
	}, nil
}

func (c *Connection) waitRefreshEvent(pool poolHandle) (ev Event, err error) {
	ch := make(chan int, 1)
	callbackID, regErr := c.hv.StoragePoolEventRefreshRegister(pool, func(
		_ *libvirt.Connect,
		_ *libvirt.StoragePool,
	) {
		ch <- 1
	})
	if regErr != nil {
		return Event{}, wrap("register refresh event", regErr)
	}
	defer finishDeregister(c.hv, callbackID, &err)

	refresh := <-ch
	now := time.Now()
	return Event{
		EventRefresh:  refresh,
		Timestamp:     now.Unix(),
		TimestampNano: now.UnixNano(),
		EventId:       callbackID,
	}, nil
}

// StreamEvents registers for storage-pool events and invokes emit for each one
// until ctx is cancelled or emit returns an error.
//
// Callers must run the process-wide libvirt event loop.
func (c *Connection) StreamEvents(ctx context.Context, uuid string, kind EventKind, emit func(Event) error) (err error) {
	if kind > EventRefresh {
		return wrap("stream events", fmt.Errorf("%w: %d", errUnsupportedEventKind, kind))
	}
	if emit == nil {
		return errEmitRequired
	}

	pool, err := c.lookupPool(uuid)
	if err != nil {
		return err
	}
	defer finishFree(pool, &err)

	events := make(chan Event, 8)
	errCh := make(chan error, 1)

	callbackID, regErr := c.registerEventStream(pool, kind, events)
	if regErr != nil {
		return wrap("register event stream", regErr)
	}
	defer finishDeregister(c.hv, callbackID, &err)

	go func() {
		<-ctx.Done()
		select {
		case errCh <- ctx.Err():
		default:
		}
	}()

	return pumpEvents(events, errCh, callbackID, emit)
}

func (c *Connection) registerEventStream(
	pool poolHandle,
	kind EventKind,
	events chan Event,
) (int, error) {
	switch kind {
	case EventLifecycle:
		return c.hv.StoragePoolEventLifecycleRegister(pool, func(
			_ *libvirt.Connect,
			_ *libvirt.StoragePool,
			event *libvirt.StoragePoolEventLifecycle,
		) {
			if event == nil {
				return
			}
			now := time.Now()
			select {
			case events <- Event{
				EventLifecycle: *event,
				Timestamp:      now.Unix(),
				TimestampNano:  now.UnixNano(),
			}:
			default:
			}
		})
	case EventRefresh:
		return c.hv.StoragePoolEventRefreshRegister(pool, func(
			_ *libvirt.Connect,
			_ *libvirt.StoragePool,
		) {
			now := time.Now()
			select {
			case events <- Event{
				EventRefresh:  1,
				Timestamp:     now.Unix(),
				TimestampNano: now.UnixNano(),
			}:
			default:
			}
		})
	}
	return 0, fmt.Errorf("%w: %d", errUnsupportedEventKind, kind)
}

func pumpEvents(
	events <-chan Event,
	errCh <-chan error,
	callbackID int,
	emit func(Event) error,
) error {
	for {
		select {
		case streamErr := <-errCh:
			if errors.Is(streamErr, context.Canceled) || errors.Is(streamErr, context.DeadlineExceeded) {
				return nil
			}
			return streamErr
		case ev := <-events:
			ev.EventId = callbackID
			emitErr := emit(ev)
			if emitErr != nil {
				return wrap("emit event", emitErr)
			}
		}
	}
}
