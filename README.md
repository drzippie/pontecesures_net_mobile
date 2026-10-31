# Pontecesures.net Mobile

Mobile app (iOS & Android) for [pontecesures.net](https://pontecesures.net), the
local news site of Pontecesures (Galicia) since 2005. Built with
[Wails v3](https://wails.io) — a Go backend plus a web frontend rendered in the
native WebView — so the same codebase also runs on desktop.

The app UI is in **Spanish**; news content comes from the site as-is (in
Galician).

## Features

- Front page with infinite scroll and section filter chips
- Section list with the website's section illustrations
- Full-text search across ~30,000 articles
- Article view with sanitised content; links to other articles open in-app,
  external links open in the system browser
- Native share sheet on mobile (copies the link on desktop)
- Offline reading of anything already loaded (file cache)
- Light & dark themes using the website's design tokens

User accounts (optional login and registration) are planned for a later phase
— see [docs/ROADMAP.md](docs/ROADMAP.md).

## Stack

| Layer | Tech |
| --- | --- |
| Shell | Wails v3.0.0 (mobile support is **experimental** upstream) |
| Backend | Go 1.26+, WordPress REST API client in `internal/wp` |
| Frontend | Vite + TypeScript, Alpine.js, Tailwind CSS v4 |
| Server side | Optional WordPress plugin in `wordpress-plugin/pcn-mobile` |

## Getting started

Prerequisites: Go 1.26+, Node 24 (see `.nvmrc`), the `wails3` CLI and the
mobile toolchains you need.

```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@v3.0.0
wails3 doctor
```

| Command | What it does |
| --- | --- |
| `wails3 dev` | Desktop window (phone-sized) with hot reload |
| `wails3 task ios:run` | Build and run on the iOS Simulator (macOS + Xcode with an iOS runtime) |
| `wails3 task android:run` | Build and run on the Android Emulator (SDK 35, NDK 26.3, JDK 21) |
| `wails3 task build:server && WAILS_SERVER_PORT=8765 bin/pontecesures-net-mobile-server` | Serve the full app over HTTP for quick checks in a browser |
| `go test ./internal/...` | Backend unit tests |

See [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) for toolchain setup and
architecture notes.

## Project layout

```
main.go               App bootstrap (window, iOS/Android options)
contentservice.go     Service bound to the frontend (posts, sections, share…)
internal/wp/          WordPress client, HTML sanitising, image cascade, cache
frontend/             Vite + Alpine.js + Tailwind UI
  src/i18n/es.ts      All UI strings (Spanish)
  bindings/           Generated Go → TS bindings (do not edit)
build/                Wails platform projects (iOS, Android, desktop)
wordpress-plugin/     Optional pcn/v1 REST endpoints for the site
docs/                 Roadmap and development notes
```

## License

Code released under the [MIT License](LICENSE). News content and images shown
in the app belong to pontecesures.net and their authors and are not covered by
this license.
