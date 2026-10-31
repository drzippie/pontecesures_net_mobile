// Language selection. UI strings live in es.ts (default) and gl.ts.
import { es, type Strings } from "./es";
import { gl } from "./gl";

export type Lang = "es" | "gl";
export type { Strings };

export const dictionaries: Record<Lang, Strings> = { es, gl };

/** Galician if the device prefers it, Spanish otherwise. */
export function detectLang(): Lang {
  const prefs = navigator.languages?.length ? navigator.languages : [navigator.language];
  return prefs.some((l) => /^gl\b/i.test(l ?? "")) ? "gl" : "es";
}

export function isLang(v: unknown): v is Lang {
  return v === "es" || v === "gl";
}

/** Replaces {placeholders} in a string. */
export function fmt(s: string, vars: Record<string, string | number>): string {
  return s.replace(/\{(\w+)\}/g, (_, k) => String(vars[k] ?? ""));
}
