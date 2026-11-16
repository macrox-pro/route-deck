# MACROX RouteDeck

macOS desktop app that assigns per-app network policy (`Proxy`, `Direct`, `Block`) by driving the official ProxyBridge CLI. It does not implement a VPN.

## Read first

- [CONCEPTION.md](CONCEPTION.md) — product bounds
- [AGENTS.md](AGENTS.md) — how agents must work

## Develop

```bash
wails3 dev
wails3 task build
```

Go 1.26+, Wails v3.0.0. Frontend: `frontend/` (Vite + React + TypeScript).

Lint (Go): `golangci-lint run ./...` using `.golangci.yml`.

## Identity

| Item | Value |
| --- | --- |
| Binary | `route-deck` |
| Bundle ID (target) | `pro.macrox.routedeck` |
| Icon | [docs/branding/route-deck-icon.png](docs/branding/route-deck-icon.png) |

Platform v1 is macOS 13+ only.
