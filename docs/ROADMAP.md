# Roadmap

## Phase 0 — Setup ✅
- Wails v3.0.0 project (vanilla-ts template) with Alpine.js + Tailwind v4
- App identity: `net.pontecesures.app`, display name "Pontecesures"

## Phase 1 — Public content (current)
- [x] Front page, sections, search, article view
- [x] Image cascade and orientation rules from the website theme
- [x] Offline cache, internal link handling, share
- [x] Related posts under each article ("Más sobre esto")
- [x] Weather for Pontecesures (MeteoGalicia): front-page line + forecast view
- [x] Settings: language (Spanish/Galician), appearance, text size
- [x] Verify on iOS Simulator (iOS 18.3) and Android Emulator (API 35)
- [x] App icon and launch screen
- [x] Pull-to-refresh gesture and refresh on return to the app
- [x] Android back button / iOS swipe back navigate inside the app
- [x] Image viewer with zoom
- [x] Saved articles (offline)
- [x] Search by relevance and within a section
- [x] Accessibility pass (contrast, screen readers, touch targets, reduced motion)
- [x] Android release shrinking with R8 (8.8 MB APK)
- [x] Run on a physical iPhone (iOS 27) with development signing
- [ ] Verify the iOS edge-swipe back gesture on a device
- [ ] Drop the iOS SDK 26 link workaround once Wails supports UIScene (wailsapp/wails#6377)
- [ ] Fetch the illustration mapping from `pcn/v1/illustrations` once deployed
- [ ] Deploy `pcn-mobile` so related posts match the website exactly

## Phase 2 — Store release
- [x] iOS distribution signing and App Store build (`wails3 task store:ios`)
- [x] Android target SDK 36, trimmed permissions and Play build (`wails3 task store:android`)
- [ ] Android upload keystore
- [ ] TestFlight and Google Play internal testing
- [x] App Store listing draft and screenshots (`docs/store/APP_STORE.md`)
- [x] Publish the app privacy policy (https://pontecesures.net/privacidad-app/)
- [ ] CI builds (GitHub Actions)

## Phase 3 — Optional user accounts (last)
- [ ] `pcn/v1/auth/*` endpoints in the WordPress plugin
- [ ] Login, registration, password reset, profile (new "Cuenta" block at the top of Settings; hidden until then because App Review rejects "coming soon" features)
- [ ] In-app account deletion (App Store requirement)
- [ ] Token stored only in Go via `application.Mobile.SecureSet`

## Later ideas
- Push notifications (APNs/FCM — not provided by Wails)
- Comments from the app, favourites, deep links
