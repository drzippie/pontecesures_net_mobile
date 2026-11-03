# Google Play listing

Source of truth for the Google Play Console listing of **Pontecesures.net**
(package `net.pontecesures.app`). Default language Spanish (es-ES), with a
Galician (gl-ES) translation, which Google Play supports. Character limits
are Play Console's (in brackets). The iOS listing is in
[APP_STORE.md](APP_STORE.md).

The developer account is **personal and was created in 2024**, so production
access needs a **closed test with at least 12 opted-in testers for 14
consecutive days** first (Play Console > Test and release > Production shows
the requirement and the "Apply for production" button once it is met).

## Create app

| Field | Value |
| --- | --- |
| App name | `Pontecesures.net` |
| Default language | Spanish (Spain) – es-ES |
| App or game | App |
| Free or paid | Free (cannot be changed later) |

## Main store listing

### App name [30]

`Pontecesures.net` (same in both languages)

### Short description [80]

es-ES:

```
Noticias de Pontecesures desde 2005 y el tiempo de MeteoGalicia.
```

gl-ES:

```
Noticias de Pontecesures desde 2005 e o tempo de MeteoGalicia.
```

### Full description [4000]

es-ES:

```
Pontecesures.net es el diario digital de Pontecesures desde 2005. Con esta aplicación tienes todas sus noticias en el móvil: lo que pasa en el Concello, las fiestas, la cultura, el deporte y la vida del pueblo, tal y como se publican en la web.

NOTICIAS
• Portada con las últimas noticias y filtros por sección.
• Todas las secciones de la web, con sus ilustraciones.
• Buscador en más de 29.000 noticias, por relevancia o por fecha, también dentro de una sección.
• Noticias relacionadas al final de cada artículo.
• Visor de imágenes a pantalla completa con zoom, ideal para carteles y esquelas.

EL TIEMPO EN PONTECESURES
• Una línea discreta en la portada con el tiempo de hoy.
• Previsión por franjas (mañana, tarde y noche), próximos días, viento, lluvia e índice UV.
• Avisos meteorológicos activos.
• Datos oficiales de MeteoGalicia (Xunta de Galicia).

LEE A TU MANERA
• Guarda noticias para leerlas luego, también sin conexión.
• Lo que ya has abierto se puede volver a leer sin conexión.
• Tema claro, oscuro o automático.
• Tres tamaños de texto.
• Aplicación en español o en gallego. Las noticias se publican en gallego.

SIN REGISTRO Y SIN PUBLICIDAD
No necesitas cuenta para usar la aplicación, no hay anuncios y no recogemos datos personales ni hacemos seguimiento.

Pontecesures.net es un proyecto vecinal y abierto: si quieres colaborar, escribir o enviar una noticia, encontrarás cómo hacerlo en el apartado «Más».
```

gl-ES:

```
Pontecesures.net é o xornal dixital de Pontecesures desde 2005. Con esta aplicación tes todas as súas noticias no móbil: o que pasa no Concello, as festas, a cultura, o deporte e a vida da vila, tal e como se publican na web.

NOTICIAS
• Portada coas últimas noticias e filtros por sección.
• Todas as seccións da web, coas súas ilustracións.
• Buscador en máis de 29.000 noticias, por relevancia ou por data, tamén dentro dunha sección.
• Noticias relacionadas ao final de cada artigo.
• Visor de imaxes a pantalla completa con zoom, ideal para carteis e esquelas.

O TEMPO EN PONTECESURES
• Unha liña discreta na portada co tempo de hoxe.
• Previsión por franxas (mañá, tarde e noite), próximos días, vento, choiva e índice UV.
• Avisos meteorolóxicos activos.
• Datos oficiais de MeteoGalicia (Xunta de Galicia).

LE Á TÚA MANEIRA
• Garda noticias para lelas despois, tamén sen conexión.
• O que xa abriches pódese volver ler sen conexión.
• Tema claro, escuro ou automático.
• Tres tamaños de texto.
• Aplicación en galego ou en castelán. As noticias publícanse en galego.

SEN REXISTRO E SEN PUBLICIDADE
Non necesitas conta para usar a aplicación, non hai anuncios e non recollemos datos persoais nin facemos seguimento.

Pontecesures.net é un proxecto veciñal e aberto: se queres colaborar, escribir ou enviar unha noticia, atoparás como facelo no apartado «Máis».
```

### Graphics

All in `docs/store/google-play/`:

| Asset | File | Spec |
| --- | --- | --- |
| App icon | `icon-512.png` | 512 × 512, 32-bit PNG (made from `build/appicon.png`) |
| Feature graphic | `feature-graphic.png` | 1024 × 500, no alpha (source: `feature-graphic.html`) |
| Phone screenshots | `screenshots/01…05-*.jpg` | 1080 × 1920 (9:16), 2–8 required |

Screenshots, in order (Android Emulator, Android 16, light theme, Spanish
UI, status bar in demo mode at 9:41):

1. `01-portada.jpg` — front page with the weather line
2. `02-noticia.jpg` — article
3. `03-tiempo.jpg` — weather forecast
4. `04-secciones.jpg` — sections
5. `05-portada-oscuro.jpg` — dark theme

The same screenshots serve the Galician listing (the UI language does not
have to match). No tablet screenshots: the app is phone-first and Play only
requires them to list the app as designed for tablets.

How they were taken: `adb shell wm size 1080x1920`, `wm density 420`,
SystemUI demo mode for the status bar, and the WebView's remote debugging
(`adb forward tcp:9333 localabstract:webview_devtools_remote_<pid>`,
`Runtime.evaluate` to set `location.hash`). Reload the page on each route:
the emulator's WebView sometimes keeps a stale paint of the taller header
with section chips when switching routes. Avoid tragic news on the front
page. Reset with `adb shell wm size reset` and `wm density reset`.

## Store settings

| Field | Value |
| --- | --- |
| App category | News & Magazines (Noticias y revistas) |
| Tags | Local news, News, Weather (choose the closest Play offers) |
| Email (public) | antonio@antoniocortes.com, or a pontecesures.net address |
| Website | `https://pontecesures.net/` |
| Phone | leave empty |
| External marketing | allowed |

## App content (Policy > App content)

| Section | Answer |
| --- | --- |
| Privacy policy | `https://pontecesures.net/privacidad-app/` |
| Ads | No, the app does not contain ads |
| App access | All functionality is available without special access |
| Content rating | Questionnaire below |
| Target audience | 18 and over (optionally also 16–17). Not under 13, which would bring in the Families policy |
| News apps | Yes, it is a news app (see below) |
| Data safety | Questionnaire below |
| Government app | No |
| Financial features | None |
| Health | None |
| Advertising ID | No (the app does not use it; nothing in the manifest declares `AD_ID`) |

### Content rating (IARC questionnaire)

- Email: the developer's.
- Category: **Reference, News, or Educational**.
- Violence, sexuality, language, controlled substances, crude humour,
  gambling: **No** — the app is not built around such content. Local news
  may mention accidents or crimes in a factual, non-graphic way.
- Users can interact or exchange content: **No** (no accounts, comments or
  messaging).
- Shares the user's location with others: **No**.
- Digital purchases: **No**.
- Unrestricted internet access (web browser or search engine): **No** —
  external links open in the system browser.

Expected rating: PEGI 3 / Everyone (IARC may give PEGI 12 for news).

### News app declaration

Answer "Yes". Google Play's News policy asks news apps to show who owns
the content and how to contact them; the website does (footer, legal notice
and "Colaborar" pages, also reachable in the app's "Más" tab). The app is
listed in News & Magazines, does not mask its source and has no
subscription.

### Data safety

- Does your app collect or share any of the required user data types?
  **No.**
  - The app downloads public content from pontecesures.net and
    MeteoGalicia; requests carry no personal data. Display preferences,
    saved articles and the offline cache stay on the device.
  - No analytics, ads, crash reporting or third-party SDKs.
- (If Play still asks:) data encrypted in transit: **Yes** (HTTPS only);
  deletion request: not applicable, no account.

Revisit both the data safety form and the privacy policy when user accounts
arrive (name, email and user ID, collected for app functionality, deletable
from the app).

## Testing tracks

1. **Internal testing** (no review, live in minutes): create an email list
   with the developer's Google account and a few more, upload the `.aab`
   from `wails3 task store:android`, roll out, and open the opt-in link on
   the phone. Up to 100 testers.
2. **Closed testing** (reviewed): add a Google Group (e.g.
   `probadores-pontecesures@googlegroups.com`) as testers so people can
   join without touching Play Console. Needs the store listing and all of
   "App content" complete. At least 12 testers must stay opted in for 14
   consecutive days; ask them to open the app now and then and to send
   feedback — Play asks about it when applying for production.
3. **Production**: "Apply for production" (questions about the test, the
   feedback and what changed), then a normal release review.

### Release notes [500] (first release)

es-ES:

```
Primera versión de la aplicación de Pontecesures.net.
```

gl-ES:

```
Primeira versión da aplicación de Pontecesures.net.
```

### Closed test feedback channel

Play shows testers a feedback URL or email; use
antonio@antoniocortes.com or `https://pontecesures.net/colaborar/`.
