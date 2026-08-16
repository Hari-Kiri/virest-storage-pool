# Business Process

Domain = **libvirt storage pool lifecycle** (not SaaS orders/payments).

## Primary: Create and run a pool

```text
Discover capabilities / sources (optional)
        ↓
Define (XML → UUID, inactive)
        ↓
Build (prepare backing storage; optional for some types)
        ↓
Create / start (active)
        ↓
Autostart configure (optional)
        ↓
Refresh volumes / Query info|detail|list
        ↓
Operate (events optional)
        ↓
Destroy (stop)
        ↓
Delete resources (optional) → Undefine definition
```

## Processes

### 1. Pool provisioning

| Step | Library | Notes |
|------|---------|-------|
| Define | `Connection.Define(model, flags)` → UUID | Inactive persistent definition |
| Build | `Build(uuid, flags)` | May be no-op / soft-fail for `dir` |
| Create | `Create(uuid, flags)` | Start pool |

### 2. Pool inspection

| Step | Library |
|------|---------|
| List | `List(listFlags, xmlFlags)` → `[]PoolSummary` |
| Info | `Info(uuid)` → volatile size/state |
| Detail | `Detail(uuid, xmlFlags)` → full XML model |

### 3. Pool discovery

| Step | Library |
|------|---------|
| Capabilities | `Capabilities()` |
| Find sources | `FindSources(poolType, SourceSpec)` |

### 4. Pool teardown

| Step | Library | Order |
|------|---------|-------|
| Destroy | `Destroy(uuid)` | Stop active |
| Delete | `Delete(uuid, flags)` | Remove backing resources |
| Undefine | `Undefine(uuid)` | Remove definition |

### 5. Autostart / refresh

| Step | Library |
|------|---------|
| Autostart | `Autostart(uuid, bool)` |
| Refresh | `Refresh(uuid)` |

### 6. Event observation

| Step | Library | Prerequisite |
|------|---------|--------------|
| Wait one | `WaitEvent(uuid, kind)` | Process event loop |
| Stream | `StreamEvents(ctx, uuid, kind, emit)` | Same |

### 7. Testserver lab session

```text
Basic Auth → JWT
        ↓
Bearer + Hypervisor-Uri (allowlisted)
        ↓
HTTP pool ops → same library flows above
```

## Related Docs

- [features/pool-lifecycle.md](features/pool-lifecycle.md)
- [DATA_FLOW.md](DATA_FLOW.md)
- [GLOSSARY.md](GLOSSARY.md)
- [INDEX.md](INDEX.md)
