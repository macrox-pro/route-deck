# Prompt packets

Copy a packet into the chat. Fill every heading. Empty headings are not allowed.

| Packet | When |
| --- | --- |
| `feature.md` | New behavior |
| `bugfix.md` | Defect with current vs expected |
| `refactor.md` | Structure only, no behavior change |
| `test-gen.md` | Tests for existing code |
| `adr.md` | Decision record |
| `handoff.md` | End of a Loop |
| `plan-mode.md` | Before multi-file edits |

Fields match `docs/templates/TASK.md`.

## Anti-example

```text
Can you clean this up and make proxies work?
```

That prompt is forbidden. It has no Goal, no paths, and no acceptance criteria.
