// Minimal hash router. Routes:
//   #/                  portada (all posts)
//   #/seccion/:id       posts of a section
//   #/secciones         section list
//   #/buscar            search
//   #/mas               about / links
//   #/tiempo            weather forecast
//   #/ajustes           settings
//   #/guardadas         saved articles
//   #/pagina/:id        static page of the website (collaborate, rules, legal)
//   #/noticia/:id       article (rendered as an overlay over the last list)

export type Route =
  | { name: "home" }
  | { name: "section"; id: number }
  | { name: "sections" }
  | { name: "search" }
  | { name: "more" }
  | { name: "weather" }
  | { name: "settings" }
  | { name: "saved" }
  | { name: "page"; id: number }
  | { name: "article"; id: number };

export function parseHash(hash: string): Route {
  const parts = hash.replace(/^#\/?/, "").split("/").filter(Boolean);
  const id = Number(parts[1]);
  switch (parts[0]) {
    case "seccion":
      return id > 0 ? { name: "section", id } : { name: "home" };
    case "secciones":
      return { name: "sections" };
    case "buscar":
      return { name: "search" };
    case "mas":
      return { name: "more" };
    case "tiempo":
      return { name: "weather" };
    case "ajustes":
      return { name: "settings" };
    case "guardadas":
      return { name: "saved" };
    case "pagina":
      return id > 0 ? { name: "page", id } : { name: "more" };
    case "noticia":
      return id > 0 ? { name: "article", id } : { name: "home" };
    default:
      return { name: "home" };
  }
}

export function href(r: Route): string {
  switch (r.name) {
    case "section":
      return `#/seccion/${r.id}`;
    case "sections":
      return "#/secciones";
    case "search":
      return "#/buscar";
    case "more":
      return "#/mas";
    case "weather":
      return "#/tiempo";
    case "settings":
      return "#/ajustes";
    case "saved":
      return "#/guardadas";
    case "page":
      return `#/pagina/${r.id}`;
    case "article":
      return `#/noticia/${r.id}`;
    default:
      return "#/";
  }
}

/** True for pontecesures.net article permalinks (/YYYY/MM/DD/slug/ or ?p=ID). */
export function isInternalPostLink(link: string): boolean {
  try {
    const u = new URL(link);
    if (!/^(www\.)?pontecesures\.net$/i.test(u.hostname)) return false;
    return /^\/\d{4}\/\d{2}\/\d{2}\/[^/]+\/?$/.test(u.pathname) || /^\d+$/.test(u.searchParams.get("p") ?? "");
  } catch {
    return false;
  }
}
