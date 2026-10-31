// All user-facing strings of the app (Spanish). Post content itself comes
// from pontecesures.net and stays in Galician. To add another language,
// create a sibling file with the same keys.
export const es = {
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
  },

  more: {
    title: "Más",
    about:
      "Pontecesures.net es el diario digital de Pontecesures desde 2005. Esta aplicación muestra sus noticias; el contenido se publica en gallego.",
    links: "Enlaces",
    website: "Abrir pontecesures.net",
    collaborate: "Colaborar",
    weather: "El tiempo en Pontecesures",
    rules: "Normas de participación",
    legal: "Aviso legal",
    version: "Versión {v}",
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

/** Replaces {placeholders} in a string. */
export function fmt(s: string, vars: Record<string, string | number>): string {
  return s.replace(/\{(\w+)\}/g, (_, k) => String(vars[k] ?? ""));
}
