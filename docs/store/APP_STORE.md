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
| Privacy policy URL | `https://pontecesures.net/colaborar/privacidade-app/` (to be published, text below) |
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

iPhone 6.9" (1320 × 2868), in this order, from `docs/store/screenshots-6.9/`.
App Store Connect scales them down for smaller iPhones. Taken on the iOS
Simulator (iPhone 16 Pro Max) with the status bar overridden to 9:41.

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
- Notes for the reviewer (English):

```
Pontecesures.net is the local news site of Pontecesures (Galicia, Spain), run by the developer since 2005. This app is its official mobile client.

The app is not a wrapped website: content is fetched from the site's public WordPress REST API and rendered natively in the app's own interface (front page, sections, search, related articles). It adds features the website does not have: saved articles and offline reading, a full weather forecast from MeteoGalicia (the Galician public weather service), an image viewer with zoom, and Spanish/Galician interface languages.

No account or login is needed. The app does not collect personal data, has no ads and no tracking. Articles are published in Galician.
```

## Privacy policy (to publish on pontecesures.net)

Publish as a page under "Colaborar" with the slug `privacidade-app`, then
use its URL above. Spanish text (the store's primary language):

```
Política de privacidad de la aplicación Pontecesures.net

Última actualización: 31 de octubre de 2026

La aplicación Pontecesures.net para iOS y Android (en adelante, «la aplicación») es un proyecto de Antonio Cortés (DrZippie), responsable de pontecesures.net. Contacto: admin@pontecesures.net.

Qué datos recogemos
La aplicación no recoge, almacena ni comparte datos personales. No necesita cuenta, no tiene publicidad, no usa herramientas de analítica ni de seguimiento y no incluye servicios de terceros que lo hagan.

Qué se guarda en tu dispositivo
Para que puedas leer sin conexión, la aplicación guarda en tu propio dispositivo una copia de las noticias que abres, las noticias que marcas como guardadas y tus preferencias (idioma, apariencia y tamaño del texto). Estos datos no salen de tu dispositivo y se borran al desinstalar la aplicación.

Conexiones que hace la aplicación
• pontecesures.net, para descargar las noticias públicas de la web.
• MeteoGalicia (Xunta de Galicia), para descargar la previsión del tiempo de Pontecesures.
Como en cualquier conexión a internet, estos servidores reciben la dirección IP de tu dispositivo para poder responder. La aplicación no les envía ningún otro dato sobre ti.

Enlaces externos
Los enlaces a otras webs se abren en el navegador de tu dispositivo y se rigen por las políticas de esas webs.

Cambios
Si en el futuro la aplicación permite crear una cuenta, actualizaremos esta política antes de publicar esa versión e indicaremos qué datos se tratan y cómo eliminarlos.
```
