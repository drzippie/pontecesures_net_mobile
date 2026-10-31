// Helpers for MeteoGalicia forecast data. Texts come from the active
// dictionary (src/i18n), which holds the official code tables translated.
import type { Day } from "../../bindings/github.com/drzippie/pontecesures_net_mobile/internal/meteo/models";
import type { Strings } from "../i18n";

export function skyText(code: number | null | undefined, t: Strings): string {
  if (code == null) return "";
  if (code >= 200) return t.weather.skyNight[code] ?? t.weather.sky[code - 100] ?? "";
  return t.weather.sky[code] ?? "";
}

/**
 * Official MeteoGalicia SVG icon for a sky code, served by the app itself
 * (Go caches it on disk: internal/meteo/icons.go) so it works offline.
 */
export function skyIcon(code: number | null | undefined): string {
  return code == null ? "" : `/meteo/icons/${code}.svg`;
}

/** "Flojo del norte" / "Frouxo do norte", "Calma", "Variable". */
export function windText(code: number | null | undefined, t: Strings): string {
  if (code == null) return "";
  if (code === 299) return t.weather.windCalm;
  if (code === 300) return t.weather.windVariable;
  if (code < 301 || code > 332) return "";
  const i = code - 301;
  return `${t.weather.windForce[Math.floor(i / 8)]} ${t.weather.windDir[i % 8]}`;
}

/** Highest rain probability of the day's periods. */
export function rainMax(day: Day | null | undefined): number | null {
  const values = (day?.periods ?? []).map((p) => p.rain).filter((v): v is number => v != null);
  return values.length ? Math.max(...values) : null;
}

/**
 * "Hoy"/"Hoxe" or the capitalised weekday. "Mañana"/"Mañá" is avoided on
 * purpose: it is also the name of the morning time slot.
 */
export function dayLabel(date: string, t: Strings): string {
  const d = new Date(`${date}T12:00:00`);
  if (d.toDateString() === new Date().toDateString()) return t.weather.today;
  const w = d.toLocaleDateString(t.locale, { weekday: "long" });
  return w.charAt(0).toUpperCase() + w.slice(1);
}

export function dayDate(date: string, locale: string): string {
  return new Date(`${date}T12:00:00`).toLocaleDateString(locale, { day: "numeric", month: "short" });
}

export function hourMinute(iso: string, locale: string): string {
  const d = new Date(iso);
  return isNaN(d.getTime()) ? "" : d.toLocaleTimeString(locale, { hour: "2-digit", minute: "2-digit" });
}
