# Development notes

## Toolchains

- **Go** 1.26+ and **Node** 24 (`nvm use`). Vite 8 needs Node ≥ 20.19.
- **wails3** `v3.0.0` — keep the CLI, `go.mod` and `@wailsio/runtime` on the
  same version. Do not use nightlies: mobile support is experimental.
- **iOS**: macOS with full Xcode *and an iOS Simulator runtime*
  (`xcodebuild -downloadPlatform iOS`, ~8 GB). `xcrun simctl list runtimes`
  must list an iOS runtime or `ios:run` fails with `Invalid runtime`.
- **Android**: SDK platform 35, build-tools 35, NDK `26.3.11579264`, an
  emulator image and JDK 21. Set `ANDROID_HOME` and `JAVA_HOME`.

Run `wails3 doctor` to check everything.

Reference setup on macOS (Apple Silicon):

```bash
# JDK 21 (Temurin tarball, no sudo needed)
export JAVA_HOME="$HOME/Library/Java/JavaVirtualMachines/jdk-21.0.13+9/Contents/Home"
export ANDROID_HOME="$HOME/Library/Android/sdk"
export ANDROID_SDK_ROOT="$ANDROID_HOME"
export PATH="$PATH:$ANDROID_HOME/platform-tools:$ANDROID_HOME/emulator:$ANDROID_HOME/cmdline-tools/latest/bin"

# SDK packages (sdkmanager from `brew install --cask android-commandlinetools`)
sdkmanager --sdk_root="$ANDROID_HOME" "cmdline-tools;latest" "platform-tools" \
  "platforms;android-35" "build-tools;35.0.0" "ndk;26.3.11579264" "emulator" \
  "system-images;android-35;google_apis;arm64-v8a"
avdmanager create avd --name wails --package "system-images;android-35;google_apis;arm64-v8a" --device pixel_7
```

## Platform quirks (Wails v3.0.0)

- **iOS executable name.** `CFBundleExecutable` in `build/ios/Info*.plist`
  must match the Taskfile binary name (`pontecesures-net-mobile`, with
  hyphens). `wails3 init` wrote it with underscores and the simulator refused
  to install the bundle.
- **iOS system bars.** The WebView is laid out *inside* the safe area, so the
  native window colour shows behind the status bar and home indicator. The
  window uses the fixed chrome colour (`chromeColour` in `main.go`,
  `--color-chrome` in CSS) and the status bar is set to light content.
- **Android insets.** Android 15 draws edge-to-edge, but the WebView only
  reports the top inset through `env(safe-area-inset-*)`. `main.ts` reads
  `ContentService.SafeArea()` (physical px ÷ `devicePixelRatio`) into
  `--safe-top`/`--safe-bottom`; the `pt-safe`/`pb-safe` utilities take the
  max of both.
- **Platform detection.** `System.IsAndroid()` from `@wailsio/runtime` returns
  `false` in the Android WebView (`window._wails.environment` is not set); use
  the user agent or ask Go (`application.System.IsMobile()`).
- **Phantom adb device.** If `adb devices` lists an `emulator-55xx` as
  `offline` while no emulator is running (another process listening on that
  port, e.g. a container), `android:run` skips booting the emulator and fails
  with `device offline`. Start the emulator yourself on a free port
  (`emulator -avd wails -port 5580`) and export `ANDROID_SERIAL=emulator-5580`.

## Debugging

- iOS: Safari → Develop → Simulator → the app. Logs: `wails3 task ios:logs:dev`.
- Android: `chrome://inspect`, or script it through the DevTools socket:
  `adb forward tcp:9333 localabstract:webview_devtools_remote_<pid>` and
  `curl localhost:9333/json`. Logs: `wails3 task android:logs`.
- Screenshots: `xcrun simctl io booted screenshot x.png`,
  `adb exec-out screencap -p > x.png`.

## Architecture

```
WebView (Alpine.js UI)  ──bindings──▶  ContentService (Go)  ──HTTPS──▶  pontecesures.net /wp-json
                                              │
                                              └── file cache (offline fallback)
```

- **All network access happens in Go.** The frontend only calls bound
  service methods (`frontend/bindings/...`). No CORS, one place for caching
  and, later, for auth tokens (which must never reach JavaScript).
- **Offline:** every successful API response is stored under
  `Mobile.StoragePath()/wp-cache` (desktop: user cache dir). If a request
  fails, the cached copy is returned with `stale: true` and the UI shows a
  banner. Entries older than 30 days are pruned on start.
- **HTML:** article bodies are sanitised with bluemonday (no scripts, styles
  or handlers; iframes only from YouTube/Vimeo/Facebook) and plain-http media
  is upgraded to https (iOS ATS).

## Images

pontecesures.net does **not** use featured images. The cascade, ported from
the theme (`wp-content/themes/pontecesures/inc/imaxes.php`), is:

1. First image in the post body (skipping gravatars, emoji and GIFs), using
   the `srcset` candidate closest to the needed width.
2. Otherwise the illustration of the post's section, with the theme's
   equivalences (e.g. `concello → politica`) and renamed files (`*-v2.jpg`),
   falling back to `xeral.jpg`.

Orientation follows the site: images with width ≥ 1.3 × height are cropped to
fill; everything else (posters, obituaries) is shown whole.

The mapping lives in `internal/wp/images.go`. The plugin exposes the same
data at `GET /wp-json/pcn/v1/illustrations` so the app can fetch it once the
plugin is deployed — keep both in sync with the theme until then.

## Frontend

- Hash router (`src/lib/router.ts`): `#/`, `#/seccion/:id`, `#/secciones`,
  `#/buscar`, `#/mas`, `#/noticia/:id`. The article is an overlay so the list
  keeps its scroll position.
- One root Alpine component (`app` in `src/main.ts`).
- Design tokens (colors, Fraunces + Source Sans 3) copied from the website
  theme into `src/style.css`; fonts are bundled for offline use.
- UI strings live in `src/i18n/es.ts` — never hard-code text in the HTML.

## Bindings

Bindings are regenerated by every build task. After changing a bound Go
method signature run:

```bash
wails3 generate bindings -clean=true -ts -i
```

`go build`/`go vet` of the root package need `frontend/dist` to exist
(`//go:embed`); run `npm run build` in `frontend/` first.
