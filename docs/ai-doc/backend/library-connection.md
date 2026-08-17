# Library Connection

## Purpose

Open/close hypervisor connection; pool lookup helper.

## Responsibilities

- `Connect(uri)` → `*Connection` wrapping `hypervisor`
- `Close()` releases libvirt conn; nil-safe
- `lookupPool(uuid)` → `poolHandle` (caller must `Free`)
- `newConnectionForTest(hv)` test injection (`*_test.go` only)

## Inputs / Outputs

| Fn | In | Out |
|----|----|-----|
| Connect | URI | Connection or error (`wrap("connect", …)`; nil libvirt conn → `errNilLibvirtConnect`) |
| Close | — | `wrap("close", …)` |
| lookupPool | uuid | poolHandle or `errConnectionClosed` |

## Uses

- `utilities.ConnectWithAuth(uri, nil, 0)`
- Production adapter: `&libvirtHypervisor{conn: virestConn.Connect}`

## Important Notes

- Auth to libvirt itself is via `utilities` (`nil` creds, flags `0` today)
- Closed connection: `hv` set nil after Close
- Pool Free errors surfaced via `finishFree` when primary `err` is nil

## Related Docs

- [adapters.md](adapters.md)
- [integrations/libvirt.md](../integrations/libvirt.md)
- [INDEX.md](../INDEX.md)
