// Display preferences. Go (SettingsService) is the source of truth; a copy is
// kept in localStorage only to apply them synchronously at startup, before
// the first paint, so the app does not flash the wrong theme or language.
import type { Settings } from "../../bindings/github.com/drzippie/pontecesures_net_mobile/internal/settings/models";
import { detectLang, isLang, type Lang } from "../i18n";

export type Theme = "auto" | "light" | "dark";
export type TextSize = "normal" | "large" | "xlarge";

export interface Prefs {
  language: Lang;
  theme: Theme;
  textSize: TextSize;
}

const KEY = "pcn.settings";

const THEMES: Theme[] = ["auto", "light", "dark"];
const SIZES: TextSize[] = ["normal", "large", "xlarge"];

/** Fills empty or unknown values with the defaults. */
export function normalize(s: Partial<Settings> | null | undefined): Prefs {
  return {
    language: isLang(s?.language) ? s.language : detectLang(),
    theme: THEMES.includes(s?.theme as Theme) ? (s!.theme as Theme) : "auto",
    textSize: SIZES.includes(s?.textSize as TextSize) ? (s!.textSize as TextSize) : "normal",
  };
}

export function loadLocal(): Partial<Settings> | null {
  try {
    const raw = localStorage.getItem(KEY);
    return raw ? (JSON.parse(raw) as Partial<Settings>) : null;
  } catch {
    return null;
  }
}

export function saveLocal(p: Prefs) {
  try {
    localStorage.setItem(KEY, JSON.stringify(p));
  } catch {
    /* storage unavailable: Go still has them */
  }
}

/** Applies preferences to the document (read by style.css). */
export function applyPrefs(p: Prefs) {
  const html = document.documentElement;
  html.lang = p.language;
  html.dataset.theme = p.theme;
  html.dataset.textSize = p.textSize;
}
