# Roadmap

## Phase 0 — Setup ✅
- Wails v3.0.0 project (vanilla-ts template) with Alpine.js + Tailwind v4
- App identity: `net.pontecesures.app`, display name "Pontecesures"

## Phase 1 — Public content (current)
- [x] Front page, sections, search, article view
- [x] Image cascade and orientation rules from the website theme
- [x] Offline cache, internal link handling, share
- [x] Verify on iOS Simulator (iOS 18.3) and Android Emulator (API 35)
- [ ] App icon and launch screen
- [ ] Pull-to-refresh gesture
- [ ] Fetch the illustration mapping from `pcn/v1/illustrations` once deployed

## Phase 2 — Store release
- [ ] Signing (Xcode automatic signing; Android upload keystore)
- [ ] TestFlight and Google Play internal testing
- [ ] Privacy policy and store listings (Spanish)
- [ ] CI builds (GitHub Actions)

## Phase 3 — Optional user accounts (last)
- [ ] `pcn/v1/auth/*` endpoints in the WordPress plugin
- [ ] Login, registration, password reset, profile
- [ ] In-app account deletion (App Store requirement)
- [ ] Token stored only in Go via `application.Mobile.SecureSet`

## Later ideas
- Push notifications (APNs/FCM — not provided by Wails)
- Comments from the app, favourites, deep links
