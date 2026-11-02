# App Store listing

Source of truth for the App Store Connect listing of **Pontecesures.net**
(Apple ID 6829266361, bundle `net.pontecesures.app`, SKU
`pontecesures-net-mobile`). The store copy is Spanish (primary language,
Spain); the App Store has no Galician localisation. Character limits are
App Store Connect's (in brackets).

## App information

| Field | Value |
| --- | --- |
| Name [30] | `Pontecesures.net` |
| Subtitle [30] | `Noticias de Pontecesures` |
| Primary category | News (Noticias) |
| Secondary category | Weather (Tiempo) |
| Content rights | The app shows third-party content: **yes**, and the developer has the rights (owner of pontecesures.net since 2005) |
| Age rating | No unrestricted web access (external links open in Safari), no user-generated content, no messaging or ads. Local news can mention accidents or crimes without graphic content: answer "Infrequent/Mild" to realistic violence only if the questionnaire asks |
| Privacy policy URL | `https://pontecesures.net/privacidad-app/` (published 2026-10-31) |
| Support URL | `https://pontecesures.net/colaborar/` |
| Marketing URL | `https://pontecesures.net/` |
| Copyright | `2026 Antonio Cortés (DrZippie)` |

## Version 1.0 page

### Promotional text [170] (editable without a new build)

```
Las noticias de Pontecesures desde 2005, ahora en tu móvil. Consulta el tiempo de MeteoGalicia y guarda noticias para leerlas sin conexión.
```

### Description [4000]

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

### Keywords [100] (comma separated, no spaces; the name is already indexed)

```
noticias,galicia,pontevedra,ulla,cesures,concello,tiempo,meteogalicia,periódico,diario,local
```

### What's new (version 1.0)

```
Primera versión de la aplicación de Pontecesures.net.
```

### Screenshots

In this order, from the folder matching the size App Store Connect asks
for (it lists the accepted sizes under each drop zone; in October 2026 the
iPhone slot asked for 6.1"/6.3" sizes):

- `docs/store/screenshots-6.3/` — 1206 × 2622 (iPhone 16 Pro simulator)
- `docs/store/screenshots-6.9/` — 1320 × 2868 (iPhone 16 Pro Max simulator)

Taken on the iOS Simulator with the status bar overridden to 9:41.

1. `01-portada.jpg` — front page with the weather line
2. `02-noticia.jpg` — article
3. `03-tiempo.jpg` — weather forecast
4. `04-secciones.jpg` — sections
5. `05-portada-oscuro.jpg` — dark theme

Avoid screenshots of tragic news (accidents, obituaries); retake if the
front page shows one.

## App privacy ("nutrition label")

Data collection: **No, we do not collect data from this app.** The app has
no analytics, no ads, no accounts and no third-party SDKs; it only
downloads public content from pontecesures.net and the forecast from
MeteoGalicia. Revisit when user accounts are added (name, email and user
ID linked to the user, for app functionality).

## App Review information

- Sign-in required: **No**.
- Contact: Antonio Cortés, antonio@antoniocortes.com.
- Notes for the reviewer (English). Version 1.0 got a Guideline 2.1
  "Information Needed" request (2026-11-02, new developer account); this
  text answers it and must also stay in the Notes field. Attach a screen
  recording made on a physical iPhone (see below) to the reply.

```
1. Screen recording
Attached: a recording made on an iPhone 15 Pro running the latest iOS. It starts by launching the app from the Home Screen and shows the typical flow: front page, opening an article, the full-screen image viewer, related articles, sections, search, the weather forecast, saving an article and reading it offline, and the settings (language, theme, text size). The app has no account registration or login, no user-generated content and no paid content or features.

2. Purpose and target audience
Pontecesures.net is the local news site of Pontecesures, a small town in Galicia (Spain), run by the developer since 2005. This app is its official mobile client. It is aimed at residents of Pontecesures and nearby towns, and at people from the area living elsewhere, who want to follow local news (town council, festivities, culture, sports, community life) and the local weather. It provides a faster and more readable way to follow the site than the mobile website: native navigation, search across about 30,000 articles, saved articles and offline reading, a full weather forecast for the town, and a Spanish or Galician interface.

3. Setup and access
No setup, account or login is required: all features are available as soon as the app opens.
- Front page: latest articles; the chips at the top filter by section. Pull down to refresh.
- Article: tap any article. Tap an image to open it full screen (pinch or double tap to zoom). Related articles appear at the end. The bookmark icon saves the article; the share icon opens the share sheet.
- Weather: tap the weather line at the top of the front page for the forecast (today by time of day, next days, warnings).
- Sections and search: "Secciones" and "Buscar" in the bottom tab bar. Search can be sorted by relevance or date and limited to a section.
- Saved articles: "Más" tab > "Guardadas"; they can be read without a connection (try airplane mode).
- Settings: "Más" tab > "Ajustes": language (Spanish/Galician), theme and text size.

4. External services
- pontecesures.net: the developer's own WordPress site. The app reads its public WordPress REST API (articles, sections, search) and loads article images from it.
- MeteoGalicia (servizos.meteogalicia.gal and www.meteogalicia.gal), the public weather service of the Galician Government: town forecast and weather icons, from its public open-data service.
- YouTube thumbnails (i.ytimg.com) when an article contains a video; the video itself opens in Safari or the YouTube app.
No authentication, payment, analytics, advertising or AI services are used.

5. Regional differences
None. The app works the same in every region. Its content is about Pontecesures and is published in Galician; the interface is in Spanish or Galician (Spanish by default unless the device prefers Galician).

6. Regulated industry / third-party material
The app does not operate in a regulated industry. The news content and images belong to pontecesures.net, which the developer owns and has run since 2005; the app and the site are by the same developer (Antonio Cortés). Weather data comes from MeteoGalicia's public open-data service and is credited in the app. Articles are shown as published on the site.
```

## Privacy policy

Published at https://pontecesures.net/privacidad-app/ (WordPress page
192198). Keep it in line with what the app does: it must mention every
server the app talks to (pontecesures.net and MeteoGalicia) and what is
stored on the device (cached content, saved articles, display
preferences). Update it before releasing user accounts.
