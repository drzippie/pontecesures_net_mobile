// User-facing strings in Galician. Same keys as es.ts (enforced by the
// `Strings` type). Weather terms follow MeteoGalicia's official annexes;
// "Máis sobre isto" and "dende 2005" match the website.
import type { Strings } from "./es";

export const gl: Strings = {
  locale: "gl-ES",

  appName: "Pontecesures.net",
  since: "dende 2005",

  tabs: {
    home: "Portada",
    sections: "Seccións",
    search: "Buscar",
    more: "Máis",
  },

  feed: {
    all: "Todo",
    loading: "Cargando novas…",
    loadingMore: "Cargando máis…",
    end: "Non hai máis novas",
    empty: "Non hai novas nesta sección.",
    refresh: "Actualizar",
  },

  search: {
    placeholder: "Buscar novas",
    hint: "Escribe polo menos 3 letras para buscar en máis de 29.000 novas.",
    noResults: "Non atopamos nada para «{q}».",
    results: "{n} resultados",
    clear: "Borrar a busca",
  },

  sections: {
    title: "Seccións",
    posts: "{n} novas",
  },

  article: {
    back: "Volver",
    share: "Compartir",
    linkCopied: "Ligazón copiada",
    openWeb: "Ver en pontecesures.net",
    by: "Por {author}",
    loading: "Cargando a nova…",
    related: "Máis sobre isto",
    closeImage: "Pechar a imaxe",
    imageHint: "Belisca ou toca dúas veces para ampliar",
  },

  more: {
    title: "Máis",
    about: "Pontecesures.net é o xornal dixital de Pontecesures dende 2005. Esta aplicación amosa as súas novas.",
    settings: "Axustes",
    links: "Ligazóns",
    website: "Abrir pontecesures.net",
    collaborate: "Colaborar",
    weather: "O tempo en Pontecesures",
    rules: "Normas de participación",
    legal: "Aviso legal",
    version: "Versión {v}",
  },

  settings: {
    title: "Axustes",
    account: "Conta",
    guest: "Convidado",
    accountSoon: "Axiña poderás crear unha conta e iniciar sesión dende a aplicación.",
    language: "Idioma",
    languageHint: "Cambia os textos da aplicación. As novas publícanse en galego.",
    languages: { es: "Español", gl: "Galego" },
    appearance: "Aparencia",
    themes: { auto: "Automático", light: "Claro", dark: "Escuro" },
    themeHint: "«Automático» segue a configuración do sistema.",
    textSize: "Tamaño do texto",
    textSizes: { normal: "Normal", large: "Grande", xlarge: "Moi grande" },
    textPreview: "Así se verán as novas.",
    saveError: "Non se puideron gardar os axustes.",
  },

  weather: {
    title: "O tempo en Pontecesures",
    today: "Hoxe",
    more: "Ver previsión",
    periods: { morning: "Mañá", afternoon: "Tarde", night: "Noite" },
    rain: "Choiva",
    wind: "Vento",
    uv: "Índice UV",
    warning: "Aviso meteorolóxico activo",
    warningMore: "Ver avisos en MeteoGalicia",
    nextDays: "Próximos días",
    source: "Datos: MeteoGalicia (Xunta de Galicia)",
    updated: "Actualizado ás {time}",
    offline: "Sen conexión: amosando a última previsión gardada.",
    unavailable: "Previsión non dispoñible",
    loading: "Cargando a previsión…",

    sky: {
      101: "Despexado",
      102: "Nubes altas",
      103: "Nubes e claros",
      104: "Nubrado",
      105: "Moi nubrado",
      106: "Néboas",
      107: "Chuvascos",
      108: "Nubrado con chuvascos",
      109: "Chuvascos e neve",
      110: "Orballo",
      111: "Choiva",
      112: "Neve",
      113: "Treboada",
      114: "Brétema",
      115: "Bancos de néboa",
      116: "Nubes medias",
      117: "Choiva débil",
      118: "Chuvascos débiles",
      119: "Treboada con poucas nubes",
      120: "Auga neve",
      121: "Sarabia",
      122: "Anubrado con chuvascos débiles",
      123: "Anubrado con neve",
      124: "Calima",
      125: "Cuberto",
    },
    skyNight: {
      204: "Anubrado",
      205: "Cuberto",
      207: "Chuvascos",
      208: "Nubrado con chuvascos",
    },
    windCalm: "Calma",
    windVariable: "Variable",
    windForce: ["Frouxo", "Moderado", "Forte", "Moi forte"],
    windDir: ["do norte", "do nordés", "do leste", "do sueste", "do sur", "do sudoeste", "do oeste", "do noroeste"],
  },

  offline: "Sen conexión: amosando novas gardadas.",
  error: {
    generic: "Non se puido cargar o contido.",
    network: "Comproba a túa conexión a internet.",
    notFound: "Esta nova xa non existe.",
    retry: "Tentar de novo",
  },

  time: {
    now: "agora mesmo",
  },
};
