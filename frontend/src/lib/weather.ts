// Spanish texts for MeteoGalicia codes. Source: "Servizo JSON coa predición a
// curto prazo por concello", annex II (sky) and annex III (wind).
import type { Day } from "../../bindings/github.com/drzippie/pontecesures_net_mobile/internal/meteo/models";

const SKY_DAY: Record<number, string> = {
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
};

// Night codes (2xx) mirror the day ones except where the annex differs.
const SKY_NIGHT_OVERRIDES: Record<number, string> = {
  204: "Nublado",
  205: "Cubierto",
  207: "Chubascos",
  208: "Nublado con chubascos",
};

export function skyText(code: number | null | undefined): string {
  if (code == null) return "";
  if (code >= 200) return SKY_NIGHT_OVERRIDES[code] ?? SKY_DAY[code - 100] ?? "";
  return SKY_DAY[code] ?? "";
}

/**
 * Official MeteoGalicia SVG icon for a sky code, served by the app itself
 * (Go caches it on disk: internal/meteo/icons.go) so it works offline.
 */
export function skyIcon(code: number | null | undefined): string {
  return code == null ? "" : `/meteo/icons/${code}.svg`;
}

const WIND_DIR = ["del norte", "del nordeste", "del este", "del sudeste", "del sur", "del sudoeste", "del oeste", "del noroeste"];
const WIND_FORCE = ["Flojo", "Moderado", "Fuerte", "Muy fuerte"];

/** "Flojo del norte", "Calma", "Variable". */
export function windText(code: number | null | undefined): string {
  if (code == null) return "";
  if (code === 299) return "Calma";
  if (code === 300) return "Variable";
  if (code < 301 || code > 332) return "";
  const i = code - 301;
  return `${WIND_FORCE[Math.floor(i / 8)]} ${WIND_DIR[i % 8]}`;
}

/** Highest rain probability of the day's periods. */
export function rainMax(day: Day | null | undefined): number | null {
  const values = (day?.periods ?? []).map((p) => p.rain).filter((v): v is number => v != null);
  return values.length ? Math.max(...values) : null;
}

const weekday = new Intl.DateTimeFormat("es-ES", { weekday: "long" });
const shortDate = new Intl.DateTimeFormat("es-ES", { day: "numeric", month: "short" });

/**
 * "Hoy" or the weekday ("Domingo"). "Mañana" is avoided on purpose: it is
 * also the name of the morning time slot and would be ambiguous.
 */
export function dayLabel(date: string, today: string): string {
  const d = new Date(`${date}T12:00:00`);
  if (d.toDateString() === new Date().toDateString()) return today;
  const w = weekday.format(d);
  return w.charAt(0).toUpperCase() + w.slice(1);
}

export function dayDate(date: string): string {
  return shortDate.format(new Date(`${date}T12:00:00`));
}

export function hourMinute(iso: string): string {
  const d = new Date(iso);
  return isNaN(d.getTime()) ? "" : d.toLocaleTimeString("es-ES", { hour: "2-digit", minute: "2-digit" });
}
