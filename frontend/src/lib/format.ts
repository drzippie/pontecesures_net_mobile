const rtf = new Intl.RelativeTimeFormat("es", { numeric: "auto" });
const dateFmt = new Intl.DateTimeFormat("es-ES", { day: "numeric", month: "short", year: "numeric" });
const longFmt = new Intl.DateTimeFormat("es-ES", {
  weekday: "long",
  day: "numeric",
  month: "long",
  year: "numeric",
  hour: "2-digit",
  minute: "2-digit",
});
const numFmt = new Intl.NumberFormat("es-ES");

/**
 * WordPress returns local site time (Europe/Madrid) without a zone; the
 * audience is local so parsing it as device-local time is good enough.
 */
function parse(iso: string): Date {
  return new Date(iso);
}

/** "hace 3 horas", "ayer" or "31 oct 2026". */
export function relativeDate(iso: string, nowText: string): string {
  const d = parse(iso);
  if (isNaN(d.getTime())) return "";
  const diffMin = Math.round((d.getTime() - Date.now()) / 60000);
  if (diffMin > -1) return nowText;
  if (diffMin > -60) return rtf.format(diffMin, "minute");
  if (diffMin > -60 * 24) return rtf.format(Math.round(diffMin / 60), "hour");
  if (diffMin > -60 * 24 * 7) return rtf.format(Math.round(diffMin / 1440), "day");
  return dateFmt.format(d);
}

/** "sábado, 31 de octubre de 2026, 08:46" */
export function longDate(iso: string): string {
  const d = parse(iso);
  return isNaN(d.getTime()) ? "" : longFmt.format(d);
}

export function number(n: number): string {
  return numFmt.format(n);
}
