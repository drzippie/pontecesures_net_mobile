// User-facing strings in Spanish (default language). Every key must also
// exist in gl.ts (the `Strings` type enforces it). Post content itself comes
// from pontecesures.net and is not translated.
export const es = {
  /** BCP 47 locale for dates and numbers. */
  locale: "es-ES",

  appName: "Pontecesures.net",
  since: "desde 2005",

  tabs: {
    home: "Portada",
    sections: "Secciones",
    search: "Buscar",
    more: "Más",
  },

  feed: {
    all: "Todo",
    loading: "Cargando noticias…",
    loadingMore: "Cargando más…",
    end: "No hay más noticias",
    empty: "No hay noticias en esta sección.",
    refresh: "Actualizar",
  },

  search: {
    placeholder: "Buscar noticias",
    hint: "Escribe al menos 3 letras para buscar en más de 29.000 noticias.",
    noResults: "No hemos encontrado nada para «{q}».",
    results: "{n} resultados",
    clear: "Borrar búsqueda",
  },

  sections: {
    title: "Secciones",
    posts: "{n} noticias",
  },

  article: {
    back: "Volver",
    share: "Compartir",
    linkCopied: "Enlace copiado",
    openWeb: "Ver en pontecesures.net",
    by: "Por {author}",
    loading: "Cargando noticia…",
    related: "Más sobre esto",
    closeImage: "Cerrar imagen",
    imageHint: "Pellizca o toca dos veces para ampliar",
  },

  more: {
    title: "Más",
    about:
      "Pontecesures.net es el diario digital de Pontecesures desde 2005. Esta aplicación muestra sus noticias; el contenido se publica en gallego.",
    settings: "Ajustes",
    links: "Enlaces",
    website: "Abrir pontecesures.net",
    collaborate: "Colaborar",
    weather: "El tiempo en Pontecesures",
    rules: "Normas de participación",
    legal: "Aviso legal",
    version: "Versión {v}",
  },

  settings: {
    title: "Ajustes",
    account: "Cuenta",
    guest: "Invitado",
    accountSoon: "Próximamente podrás crear una cuenta e iniciar sesión desde la aplicación.",
    language: "Idioma",
    languageHint: "Cambia los textos de la aplicación. Las noticias se publican en gallego.",
    languages: { es: "Español", gl: "Galego" } as Record<string, string>,
    appearance: "Apariencia",
    themes: { auto: "Automático", light: "Claro", dark: "Oscuro" } as Record<string, string>,
    themeHint: "«Automático» sigue la configuración del sistema.",
    textSize: "Tamaño del texto",
    textSizes: { normal: "Normal", large: "Grande", xlarge: "Muy grande" } as Record<string, string>,
    textPreview: "Así se verán las noticias.",
    saveError: "No se han podido guardar los ajustes.",
  },

  weather: {
    title: "El tiempo en Pontecesures",
    today: "Hoy",
    more: "Ver previsión",
    periods: { morning: "Mañana", afternoon: "Tarde", night: "Noche" } as Record<string, string>,
    rain: "Lluvia",
    wind: "Viento",
    uv: "Índice UV",
    warning: "Aviso meteorológico activo",
    warningMore: "Ver avisos en MeteoGalicia",
    nextDays: "Próximos días",
    source: "Datos: MeteoGalicia (Xunta de Galicia)",
    updated: "Actualizado {time}",
    offline: "Sin conexión: mostrando la última previsión guardada.",
    unavailable: "Previsión no disponible",
    loading: "Cargando previsión…",

    // MeteoGalicia codes ("Servizo JSON coa predición a curto prazo por
    // concello", annex II for sky and annex III for wind).
    sky: {
      101: "Despejado",
      102: "Nubes altas",
      103: "Nubes y claros",
      104: "Nublado",
      105: "Muy nublado",
      106: "Nieblas",
      107: "Chubascos",
      108: "Nublado con chubascos",
      109: "Chubascos de nieve",
      110: "Llovizna",
      111: "Lluvia",
      112: "Nieve",
      113: "Tormenta",
      114: "Bruma",
      115: "Bancos de niebla",
      116: "Nubes medias",
      117: "Lluvia débil",
      118: "Chubascos débiles",
      119: "Tormenta con pocas nubes",
      120: "Aguanieve",
      121: "Granizo",
      122: "Nublado con chubascos débiles",
      123: "Nublado con nieve",
      124: "Calima",
      125: "Cubierto",
    } as Record<number, string>,
    /** Night codes (2xx) mirror the day ones except these. */
    skyNight: {
      204: "Nublado",
      205: "Cubierto",
      207: "Chubascos",
      208: "Nublado con chubascos",
    } as Record<number, string>,
    windCalm: "Calma",
    windVariable: "Variable",
    windForce: ["Flojo", "Moderado", "Fuerte", "Muy fuerte"],
    windDir: ["del norte", "del nordeste", "del este", "del sudeste", "del sur", "del sudoeste", "del oeste", "del noroeste"],
  },

  offline: "Sin conexión: mostrando noticias guardadas.",
  error: {
    generic: "No se ha podido cargar el contenido.",
    network: "Comprueba tu conexión a internet.",
    notFound: "Esta noticia ya no existe.",
    retry: "Reintentar",
  },

  time: {
    now: "ahora mismo",
  },
};

export type Strings = typeof es;
