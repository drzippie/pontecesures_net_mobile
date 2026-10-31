# Pontecesures.net Mobile API (WordPress plugin)

Small WordPress plugin that adds the REST endpoints the mobile app needs on
top of the core `wp/v2` API. Namespace: `pcn/v1`.

> **Status:** not installed on pontecesures.net yet. The app works without it:
> it ships a built-in copy of the illustration mapping and approximates
> related posts with the core REST API.

## Endpoints

| Method | Route | Auth | Description |
| --- | --- | --- | --- |
| GET | `/wp-json/pcn/v1/illustrations` | public | Section → illustration mapping used when a post has no image in its body. Reads the theme's `inc/imaxes.php` functions when available. |
| GET | `/wp-json/pcn/v1/posts/{id}/related?per_page=4` | public | IDs of the "Máis sobre isto" posts, from `pontecesures_novas_relacionadas()` (plugin `comments`). 404 if that function is missing; the app then uses its own approximation. Max `per_page` 8. |

Example response:

```json
{
  "version": "0.1.0",
  "fallback": "xeral",
  "urls": {
    "sociedade": "https://pontecesures.net/wp-content/themes/pontecesures/assets/img/categorias/sociedade-v2.jpg",
    "xeral": "https://pontecesures.net/wp-content/themes/pontecesures/assets/img/categorias/xeral.jpg"
  },
  "equivalents": { "concello": "politica", "pontecesures-cf": "deportes", "formacionemprego": "emprego" }
}
```

## Planned (last phase)

Optional user accounts — not implemented yet:

- `POST /pcn/v1/auth/register` — self-registration (nonce, rate limit, Akismet)
- `POST /pcn/v1/auth/login` — returns an opaque, revocable token (hashed in user meta, with expiry)
- `POST /pcn/v1/auth/logout`, `GET /pcn/v1/auth/me`
- `POST /pcn/v1/auth/password-reset`
- `DELETE /pcn/v1/auth/me` — account deletion (required by App Store guideline 5.1.1(v))

## Installation

Copy `pcn-mobile/` into `wp-content/plugins/` and activate it, or zip it:

```bash
cd wordpress-plugin && zip -r pcn-mobile.zip pcn-mobile
```
