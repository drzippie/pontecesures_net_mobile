// Detects an Android System WebView too old for the UI's CSS. Tailwind CSS
// v4 relies on cascade layers, color-mix() and other features from
// Chromium 111; older WebViews ignore the stylesheet and show raw HTML.
// Phones normally update the WebView through Google Play, but some (Android
// Go, updates disabled, emulator images) keep the one they shipped with.

/** Oldest Chromium version whose WebView renders the app correctly. */
export const MIN_WEBVIEW = 111;

/** Play Store page of Android System WebView. */
export const WEBVIEW_STORE_URL = "https://play.google.com/store/apps/details?id=com.google.android.webview";

/** Returns the WebView's Chromium major version if it is too old, else null. */
export function outdatedWebView(ua: string = navigator.userAgent): number | null {
  if (!/\bAndroid\b/.test(ua)) return null;
  const m = /\bChrome\/(\d+)\./.exec(ua);
  if (!m) return null;
  const version = Number(m[1]);
  return version < MIN_WEBVIEW ? version : null;
}

/**
 * Replaces the page with a notice asking to update the WebView. Styled
 * inline with plain CSS that any WebView understands.
 */
export function showOutdatedWebView(
  text: { title: string; body: string; button: string },
  openStore: () => void,
): void {
  document.body.innerHTML = "";
  document.body.setAttribute(
    "style",
    "margin:0;min-height:100vh;display:flex;align-items:center;justify-content:center;" +
      "background:#1b3a4b;color:#fff;font-family:sans-serif;text-align:center;",
  );
  const box = document.createElement("div");
  box.setAttribute("style", "max-width:28rem;padding:2rem 1.5rem;");
  const title = document.createElement("h1");
  title.setAttribute("style", "font-size:1.5rem;line-height:1.25;margin:0 0 1rem;");
  title.textContent = text.title;
  const body = document.createElement("p");
  body.setAttribute("style", "font-size:1.05rem;line-height:1.5;margin:0 0 1.75rem;color:#dce6ea;");
  body.textContent = text.body;
  const button = document.createElement("button");
  button.setAttribute(
    "style",
    "font-size:1.05rem;font-weight:600;padding:.85rem 1.5rem;border:0;border-radius:999px;" +
      "background:#fff;color:#1b3a4b;",
  );
  button.textContent = text.button;
  button.addEventListener("click", openStore);
  box.append(title, body, button);
  document.body.append(box);
}
