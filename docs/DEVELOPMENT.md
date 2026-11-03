# Development notes

## Toolchains

- **Go** 1.26+ and **Node** 24 (`nvm use`). Vite 8 needs Node ≥ 20.19.
- **wails3** `v3.1.0` — keep the CLI, `go.mod` and `@wailsio/runtime` on the
  same version. Do not use nightlies: mobile support is experimental.
  v3.1.0's experimental "Wake" build system (`wails.yaml`, enabled with
  `WAILS_EXP_USE_WAKE=1`) is not used: the migration is one-way and the
  project relies on custom Taskfile tasks (`device:ios`, `store:ios`).
  The Android project follows the v3.1.0 templates' AGP 9 and SDK 36 (Google
  Play requires target SDK 36 since 31 August 2026) but not their dev-server
  live reload.
- **iOS**: macOS with full Xcode *and an iOS Simulator runtime*
  (`xcodebuild -downloadPlatform iOS`, ~8 GB). `xcrun simctl list runtimes`
  must list an iOS runtime or `ios:run` fails with `Invalid runtime`.
- **Android**: minSdk 24 (Play Protect requires it), SDK platform 36, build-tools 36, NDK `26.3.11579264`, an
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

## Platform quirks (Wails v3.1.0)

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

- **Android back button.** `WebView.canGoBack()` ignores the fragment
  navigations the app uses (`#/noticia/123`), so back used to close the app.
  `MainActivity.onBackPressed` asks the page first via
  `window.pcnHandleBack()` (closes the image viewer, then the article, then
  walks back through views) and only leaves the app from the front page.
  The frontend keeps its own stack of visited hashes (`navStack`).
- **Lifecycle.** The Android WebView neither fires `visibilitychange` nor
  receives events while paused. `lifecycle.go` listens to the native
  `ActivityStopped/Resumed` (iOS `DidEnterBackground/DidBecomeActive`) events
  and emits `app:foreground` with the milliseconds spent away, measured in Go.
- **R8.** Release builds use `minifyEnabled` + `shrinkResources`. Go resolves
  `WailsBridge` methods by name through JNI and the page calls
  `@JavascriptInterface` methods, so `proguard-rules.pro` keeps those; always
  smoke-test a release APK (`wails3 task android:package`) after touching the
  Java side or the rules.
- **Emulator debuggability.** The Google APIs emulator image is `userdebug`,
  so every WebView (also release builds) is inspectable there; on real
  devices release builds are not.

## App icon

`build/appicon.png` (1024², source of the generated iOS AppIcon set) is a
white Fraunces "P" with the tella dot on the brand blue, rendered from the
app's own font. Android has legacy/round mipmaps plus an adaptive icon
(`mipmap-anydpi-v26`, monochrome layer for themed icons). The website's site
icon is the municipal coat of arms and is intentionally not used:
pontecesures.net is not the Concello, and stores reject apps that look
official.

## Running on an iPhone

```bash
wails3 task device:ios            # build, sign, install and launch (USB or same Wi-Fi)
```

Requirements: Xcode with the Apple Developer account added, Developer Mode on
the iPhone, a development provisioning profile at `bin/dev.mobileprovision`
(gitignored; Xcode's "iOS Team Provisioning Profile" or one downloaded from
developer.apple.com) and the matching "Apple Development" identity in the
keychain. `scripts/ios_signing.py` picks the identity the profile allows and
derives the entitlements (the account's App ID prefix, `3W5LP84889`, differs
from the team id `MA8KEHHSX2`).

Lessons learned (Wails v3.0.0–v3.1.0, Xcode 27, iOS 27):

- **iOS 27 requires the UIScene lifecycle** for apps linked against the
  iOS 27 SDK; Wails still uses the legacy app-delegate window, so the app
  quits immediately (exit code 0, no logs). `build/ios/Taskfile.yml` links
  with `-platform_version ios <min> 26.0` (`IOS_LINKED_SDK`) as a
  workaround. Upstream: wailsapp/wails#6377, #6194. Remove the workaround once
  Wails adopts UIScene, and before Apple requires the iOS 27 SDK for uploads.
- **Do not use the generated Xcode project** (`wails3 task ios:xcode`) for
  device builds: it disables signing, names the executable `wailsapp` while
  `PRODUCT_NAME` is the display name, does not force-load `WailsAppDelegate`
  and its Go build phase needs `go` in Xcode's PATH. `device:ios` uses
  `ios:package IOS_PLATFORM=device` instead, which also compiles the icons
  and (added here) the launch storyboard.
- **Signing.** Xcode looped on "Revoke certificate" because the Apple WWDR G3
  intermediate certificate was missing from the keychain, so no development
  certificate could be valid. Installing `AppleWWDRCAG3.cer` from
  apple.com/certificateauthority fixes it. A certificate downloaded from the
  developer portal never contains the private key.
- **Logs.** Release builds print nothing; `log collect --device-udid` needs
  root *and* a USB connection. Install, launch and `devicectl device copy
  from --domain-type appDataContainer` all work over Wi-Fi.

## App Store / TestFlight builds

```bash
wails3 task store:ios             # signed .ipa at bin/pontecesures-net-mobile.ipa
wails3 task store:ios:upload      # validate + upload with an App Store Connect API key
```

Requirements (one-off, developer.apple.com):

- An "Apple Distribution" certificate with its private key in the keychain.
  The CSR can be made with `openssl req` instead of Keychain Access; combine
  the downloaded `.cer` with the key into a `.p12` and `security import` it.
- An explicit App ID `net.pontecesures.app` (team prefix `MA8KEHHSX2`; App
  Store Connect does not accept wildcard App IDs). The description field
  only allows letters, digits and spaces.
- An **iOS** "App Store Connect" provisioning profile for that App ID and
  certificate at `bin/store.mobileprovision` (gitignored). Note that a macOS
  profile downloads as `.provisionprofile`; iOS ones are `.mobileprovision`.

What the build adds compared with `device:ios`:

- Distribution entitlements from the profile (`get-task-allow` false,
  `beta-reports-active` for TestFlight) via `scripts/ios_signing.py`.
- A new version per build: the patch number is bumped by default
  (`BUMP=minor|major|none`, or `APP_VERSION=X.Y.Z`). `scripts/version.py`
  keeps `build/config.yml`, both iOS plists, Android's `versionName` (and
  `versionCode`, +1), `appVersion` in `main.go` and `APP_VERSION` in
  `frontend/src/main.ts` in sync; commit the change after the build.
- A fresh `CFBundleVersion` per build (`YYMMDD.HHMM`, override with
  `BUILD_NUMBER=`).
- The `DT*` keys Xcode normally writes (SDK, platform and Xcode versions),
  which App Store Connect validates; `ios:package` adds them to device builds.
- `build/ios/PrivacyInfo.xcprivacy`: no tracking, no collected data, and the
  required-reason APIs the Go runtime links (`stat` family → file timestamps
  C617.1, `mach_absolute_time` → system boot time 35F9.1). Re-check with
  `nm -u bin/pontecesures-net-mobile` after adding dependencies.

`Info.plist` declares iPhone only (`UIDeviceFamily` 1) and
`ITSAppUsesNonExemptEncryption` false (HTTPS only).

The Wails runtime always links its optional mobile features (location,
camera torch, image picker, motion, biometrics). App Store Connect flagged
only location (ITMS-90683 on build 261031.1814), so `Info.plist` has a
`NSLocationWhenInUseUsageDescription` saying the app does not use it. If
another key is reported, add it the same way (an honest text; the app never
requests those permissions).

Uploading: `store:ios:upload` needs `ASC_KEY_ID`, `ASC_ISSUER_ID` and
`~/.appstoreconnect/private_keys/AuthKey_<ASC_KEY_ID>.p8`; alternatively drag
the `.ipa` into Apple's Transporter app.

## Google Play builds

```bash
wails3 task store:android         # signed .aab at bin/pontecesures-net-mobile.aab
```

The bundle holds production `libwails.so` builds for arm64-v8a and x86_64
and is signed with the **upload key**; with Play App Signing, Google
re-signs it with the app signing key it keeps. Versioning works as in
`store:ios` (patch bump by default, `versionCode` +1); commit afterwards.
Upload the `.aab` in Play Console (Test and release > the track > Create new
release).

Upload keystore (one-off; keep a backup outside this machine, a lost upload
key can be reset from Play Console but it takes days):

```bash
mkdir -p ~/.config/pontecesures-net-mobile
keytool -genkeypair -keystore ~/.config/pontecesures-net-mobile/upload-keystore.jks \
  -storetype PKCS12 -alias upload -keyalg RSA -keysize 4096 -validity 10000 \
  -dname "CN=Antonio Cortés (DrZippie), O=Pontecesures.net, C=ES"
security add-generic-password -s pontecesures-android-upload -a upload -w
```

`store:android` reads the password from `ANDROID_KEYSTORE_PASSWORD` or that
keychain item, and refuses a bundle signed with the debug keystore.

Notes:

- `build/android/Taskfile.yml` passes `PRODUCTION` from `build` to
  `compile:go:shared`. Without it (Wails v3.1.0 template) `bundle:fat` and
  `package` shipped a *debug* `libwails.so` for arm64, the architecture
  real phones use. Check with `go version -m` on the extracted library.
- The manifest removes the permissions of Wails features the app never
  calls (location, camera, biometrics, notifications, foreground service)
  with `tools:node="remove"`, so Play Console asks for no permission
  declarations. The app keeps INTERNET, ACCESS_NETWORK_STATE and VIBRATE.
- Test on a phone with `wails3 task android:run:device` (debug build, USB
  debugging on; `adb devices` must list the phone).

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

## Third-party content in articles

`wp.SanitizeHTML` keeps what the WebView contacts in line with the privacy
policy (https://pontecesures.net/privacidad-app/) and the "no data
collected" App Store label:

- Images from ad/tracking hosts (adnxs, Criteo, DoubleClick…) are dropped;
  some old posts copied from other media carry their tracking pixels.
- Embedded players become links opened outside the app. YouTube refuses to
  play embeds without a web Referer (error 153), which pages served from the
  `wails://` scheme cannot send; YouTube embeds become the video thumbnail
  with a play button (`a.pcn-video`), Vimeo and Facebook embeds a plain
  link.
- Images hot-linked from other news sites are kept (the policy mentions
  them).



Under each article the app shows "Más sobre esto" (4 posts), like the
website's "Máis sobre isto".

- **Preferred:** `GET /pcn/v1/posts/{id}/related` from the `pcn-mobile`
  plugin returns the IDs computed by the site itself
  (`pontecesures_novas_relacionadas()` in the `comments` plugin: tags weighted
  by rarity, recency bonus, same-category top-up).
- **Fallback** (plugin not installed — the route answers 404 and the client
  stops asking for 6 hours): `internal/wp/related.go` queries the post's
  specific tags (fewer than 1,500 posts) one by one from the rarest to the most
  common, newest first, then tops up with the latest posts of the same
  category. Posts tagged only with generic tags ("Pontecesures", "Concello")
  get recent posts of those tags, which is also what the website does.

## Weather

Source: MeteoGalicia (Xunta de Galicia), the same public endpoint the website
uses — `https://servizos.meteogalicia.gal/mgrss/predicion/jsonPredConcellos.action?idConc=36044`
(no API key; 36044 = Pontecesures). Spec and code tables:
[JSON_Pred_Concello_gl.pdf](https://www.meteogalicia.gal/datosred/infoweb/meteo/docs/rss/JSON_Pred_Concello_gl.pdf).

- `internal/meteo` parses 3–4 days with morning/afternoon/night slots. `-9999`
  ("not available") becomes `null`. Responses are cached for 30 minutes and
  served stale when offline.
- Go returns codes; Spanish texts live in `frontend/src/lib/weather.ts`
  (annex II sky codes 101–125 day / 201–225 night, annex III wind 299–332).
  Icons are MeteoGalicia's official SVGs
  (`https://www.meteogalicia.gal/web/assets/icons/svg/<code>.svg`). They are
  not shipped in the repo: `meteo.IconStore` serves them under
  `/meteo/icons/<code>.svg` through an asset-server middleware, downloading
  each one once into `<storage>/meteo-icons/` and serving it from disk
  afterwards. Every fresh forecast prefetches its icons in the background,
  so the weather view works offline before it is ever opened. Only annex II
  codes are accepted and responses must be real SVGs.
- UI: one discreet line on the front page (today's icon, max/min, max rain
  probability) that opens `#/tiempo` with today's slots, UV index, next days,
  active warnings and the attribution.
- The website widget (`class-weather-widget.php`) maps codes with its own
  table that does not match the official annex (e.g. it shows night code 201,
  "clear", as "Choiva escasa"). The app follows the official annex.

## Settings and languages

`#/ajustes` (from "Más") holds the display preferences (an account block
will go first once accounts exist; no placeholder before, App Review rejects
"coming soon" features):

| Setting | Values | Effect |
| --- | --- | --- |
| Language | `es`, `gl` | UI strings (`src/i18n/{es,gl}.ts`), dates/numbers via `Intl` with `es-ES`/`gl-ES`. Default: Galician if the device prefers it, else Spanish |
| Appearance | `auto`, `light`, `dark` | `<html data-theme>`; `auto` follows `prefers-color-scheme` |
| Text size | `normal`, `large`, `xlarge` | `<html data-text-size>` scales the root font size (100/112.5/125 %) |

- Go (`SettingsService`, `internal/settings`) is the source of truth: a
  validated `settings.json` in the app's private storage (desktop: user
  config dir), not in a cache directory.
- A copy is kept in `localStorage` only to apply the preferences
  synchronously before the first paint (no theme/language flash).
- Weather code texts are part of the dictionaries; the Galician ones are the
  official MeteoGalicia terms.
- `wails3 task android:run` uninstalls the app first, so saved settings are
  reset on every dev deploy.

## Reading features

- **Image viewer** (`src/lib/zoom.ts`): article images open full screen at
  their largest `srcset` size; pinch, double tap and pan with pointer events
  (`touch-action: none` on the viewer, the page itself never zooms).
- **Saved articles** (`internal/saved`): full posts stored as JSON in
  persistent storage (`dataPath("saved")`); `Post()` falls back to them when
  offline and refreshes them when the network copy loads. Remote images in
  saved posts depend on the WebView cache.
- **Search**: `orderby=relevance` (only with a search term) or date, plus a
  section filter.
- **Refresh**: pull to refresh on lists (`src/lib/pull.ts`); after 15+ minutes
  in the background the forecast reloads, and the list too if the reader is
  near the top.

## Frontend

- Hash router (`src/lib/router.ts`): `#/`, `#/seccion/:id`, `#/secciones`,
  `#/buscar`, `#/mas`, `#/tiempo`, `#/ajustes`, `#/guardadas`,
  `#/pagina/:id`, `#/noticia/:id`.
- Website pages shown in-app (`#/pagina/:id`, `ContentService.Page`, same
  sanitising and offline cache as posts): "Como colaborar" (10), "Normas de
  uso" (287) and "Legal" (487), referenced by WordPress page ID. The article is an overlay so the list
  keeps its scroll position.
- One root Alpine component (`app` in `src/main.ts`).
- Design tokens (colors, Fraunces + Source Sans 3) copied from the website
  theme into `src/style.css`; fonts are bundled for offline use.
- UI strings live in `src/i18n/es.ts` and `gl.ts` — never hard-code text in
  the HTML.

## Bindings

Bindings are regenerated by every build task. After changing a bound Go
method signature run:

```bash
wails3 generate bindings -clean=true -ts -i
```

`go build`/`go vet` of the root package need `frontend/dist` to exist
(`//go:embed`); run `npm run build` in `frontend/` first.
