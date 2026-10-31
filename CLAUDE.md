# CLAUDE.md

Guidance for Claude Code when working in this repository.

## Language rules

- **Everything that goes to GitHub is in English**: commit messages, PR titles
  and descriptions, issues, review comments, code comments and docs in this
  repo.
- **App UI text is Spanish** and lives only in `frontend/src/i18n/es.ts`.
- Post content comes from pontecesures.net in Galician; do not translate it.
  Section names are shown exactly as the site names them.

## Project

Wails v3 (v3.0.0, mobile support experimental) app for iOS/Android that shows
the public content of pontecesures.net (WordPress). See `README.md`,
`docs/DEVELOPMENT.md` and `docs/ROADMAP.md`.

- Go backend: `main.go`, `contentservice.go`, `weatherservice.go`,
  `internal/wp/` (WordPress), `internal/meteo/` (MeteoGalicia),
  `internal/cache/` (offline file cache)
- Frontend: `frontend/` — Vite + TypeScript + **Alpine.js + Tailwind CSS v4**
  (do not introduce React/Vue/Svelte)
- Optional WordPress plugin: `wordpress-plugin/pcn-mobile` (namespace
  `pcn/v1`). Not deployed; never deploy or activate it without the owner.

## Conventions

- All HTTP calls to WordPress happen in Go, never from the WebView. Auth
  tokens (future) must stay in Go / secure storage.
- pontecesures.net has no featured images: use the image cascade in
  `internal/wp/images.go` (body image → section illustration). Keep it in sync
  with the theme's `inc/imaxes.php`.
- Sanitise any HTML rendered with `x-html` in Go (`wp.SanitizeHTML`).
- Weather codes follow MeteoGalicia's official annex (see
  `docs/DEVELOPMENT.md`), not the website widget's table.
- Platform-specific Go code goes behind `//go:build ios` / `//go:build android`
  (`ios` implies `darwin`, `android` implies `linux`). Prefer
  `application.Mobile` / `application.System.IsMobile()` in shared code.
- Pin Wails to `v3.0.0` (CLI, `go.mod`, `@wailsio/runtime`). No nightlies.
- Do not edit generated files: `frontend/bindings/`, `build/ios/xcode/`,
  `build/android/gen/`.

## Commands

```bash
nvm use                                   # Node 24
cd frontend && npm install && npm run build && cd ..
go test ./internal/...                    # unit tests
wails3 dev                                # desktop, phone-sized window
wails3 task ios:run                       # iOS Simulator
wails3 task android:run                   # Android Emulator
wails3 generate bindings -clean=true -ts -i
```

Quick visual check without a simulator: `wails3 task build:server`, run
`WAILS_SERVER_PORT=8765 bin/pontecesures-net-mobile-server` and open
`http://127.0.0.1:8765/` in a mobile-sized browser window.
