// Locale-aware formatting. Formatters are cached per locale.

const cache = new Map<string, unknown>();
function memo<T>(key: string, make: () => T): T {
  if (!cache.has(key)) cache.set(key, make());
  return cache.get(key) as T;
}

const rtf = (locale: string) =>
  memo(`rtf:${locale}`, () => new Intl.RelativeTimeFormat(locale, { numeric: "auto" }));
const shortDateFmt = (locale: string) =>
  memo(`short:${locale}`, () => new Intl.DateTimeFormat(locale, { day: "numeric", month: "short", year: "numeric" }));
const longDateFmt = (locale: string) =>
  memo(
    `long:${locale}`,
    () =>
      new Intl.DateTimeFormat(locale, {
        weekday: "long",
        day: "numeric",
        month: "long",
        year: "numeric",
        hour: "2-digit",
        minute: "2-digit",
      }),
  );
const numFmt = (locale: string) => memo(`num:${locale}`, () => new Intl.NumberFormat(locale));

/**
 * WordPress returns local site time (Europe/Madrid) without a zone; the
 * audience is local so parsing it as device-local time is good enough.
 */
function parse(iso: string): Date {
  return new Date(iso);
}

/** "hace 3 horas" / "hai 3 horas", "ayer" / "onte" or "31 oct 2026". */
export function relativeDate(iso: string, locale: string, nowText: string): string {
  const d = parse(iso);
  if (isNaN(d.getTime())) return "";
  const diffMin = Math.round((d.getTime() - Date.now()) / 60000);
  if (diffMin > -1) return nowText;
  if (diffMin > -60) return rtf(locale).format(diffMin, "minute");
  if (diffMin > -60 * 24) return rtf(locale).format(Math.round(diffMin / 60), "hour");
  if (diffMin > -60 * 24 * 7) return rtf(locale).format(Math.round(diffMin / 1440), "day");
  return shortDateFmt(locale).format(d);
}

/** "sábado, 31 de octubre de 2026, 08:46" */
export function longDate(iso: string, locale: string): string {
  const d = parse(iso);
  return isNaN(d.getTime()) ? "" : longDateFmt(locale).format(d);
}

export function number(n: number, locale: string): string {
  return numFmt(locale).format(n);
}
