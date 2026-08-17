# AI Tool Configuration

## Cursor (this repo)

| Artifact | Path | Role |
|----------|------|------|
| Always-on rule | `.cursor/rules/ai-doc.mdc` | Force INDEX routing + post-change doc update |
| Skill | `.cursor/skills/maintain-ai-doc/SKILL.md` | Read ≤3 / update-one map |
| Stop hook | `.cursor/hooks/check-ai-doc-sync.sh` via `.cursor/hooks.json` | Nudge if impl changed without ai-doc |
| Go standards | `.cursor/rules/golang-standards.mdc` | Senior Go + go-libvirt style for `storagePool/`, `utilities/`, `cmd/`, `examples/` |
| Library bar | `.cursor/rules/mission-critical-library.mdc` | HA / production bar (always on) for library + cmd + examples |
| Libvirt skill | `.cursor/skills/go-libvirt-library/SKILL.md` | Library/cmd/examples: perf, N+1, deps, loops, lint, tests |

## Protocol (all agents)

1. Start at [`INDEX.md`](INDEX.md) only.
2. Open ≤3 linked docs.
3. After changing `storagepool/`, `cmd/`, or `examples/`, update one primary page (CATALOG update map).
4. Do not read entire `docs/ai-doc/`.

## Other tools (Copilot / Claude Code / Codex / Gemini / Cline / Roo / Windsurf)

Point project instructions / AGENTS / `.clinerules` / equivalent at:

- Primary: `docs/ai-doc/INDEX.md`
- Sync rule: same as Cursor protocol above

No separate per-vendor files required if the tool honors repo Cursor rules or a root instruction file that links INDEX.

## Related Docs

- [INDEX.md](INDEX.md)
- [CATALOG.md](CATALOG.md)
- [CODING_STANDARD.md](CODING_STANDARD.md)
