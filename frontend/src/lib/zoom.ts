// Pinch-to-zoom, double-tap zoom and panning for the image viewer, with
// pointer events (works on iOS WKWebView and Android WebView). The page
// itself never zooms: the viewer has touch-action: none.

const MIN = 1;
const MAX = 5;
const DOUBLE_TAP_SCALE = 2.5;
const DOUBLE_TAP_MS = 300;

type Pt = { x: number; y: number };

export function attachZoom(container: HTMLElement, img: HTMLImageElement): () => void {
  let scale = 1;
  let tx = 0;
  let ty = 0;
  const pointers = new Map<number, Pt>();
  let pinch: { dist: number; scale: number; mid: Pt; tx: number; ty: number } | null = null;
  let lastTap = { t: 0, x: 0, y: 0 };

  const apply = (animate = false) => {
    img.style.transition = animate ? "transform 180ms ease-out" : "none";
    img.style.transform = `translate3d(${tx}px, ${ty}px, 0) scale(${scale})`;
  };

  /** Keeps the scaled image covering the viewport where possible. */
  const clamp = () => {
    const c = container.getBoundingClientRect();
    const w = img.offsetWidth * scale;
    const h = img.offsetHeight * scale;
    const maxX = Math.max(0, (w - c.width) / 2);
    const maxY = Math.max(0, (h - c.height) / 2);
    tx = Math.min(maxX, Math.max(-maxX, tx));
    ty = Math.min(maxY, Math.max(-maxY, ty));
  };

  const center = (): Pt => {
    const c = container.getBoundingClientRect();
    return { x: c.left + c.width / 2, y: c.top + c.height / 2 };
  };

  const zoomAt = (target: number, p: Pt) => {
    const c = center();
    const k = target / scale;
    // Keep the point under the finger fixed while scaling.
    tx = (tx - (p.x - c.x)) * k + (p.x - c.x);
    ty = (ty - (p.y - c.y)) * k + (p.y - c.y);
    scale = target;
    if (scale <= MIN) {
      scale = MIN;
      tx = ty = 0;
    }
    clamp();
  };

  const down = (e: PointerEvent) => {
    container.setPointerCapture(e.pointerId);
    pointers.set(e.pointerId, { x: e.clientX, y: e.clientY });
    if (pointers.size === 2) {
      const [a, b] = [...pointers.values()];
      pinch = {
        dist: Math.hypot(a.x - b.x, a.y - b.y),
        scale,
        mid: { x: (a.x + b.x) / 2, y: (a.y + b.y) / 2 },
        tx,
        ty,
      };
    }
  };

  const move = (e: PointerEvent) => {
    const prev = pointers.get(e.pointerId);
    if (!prev) return;
    const cur = { x: e.clientX, y: e.clientY };
    pointers.set(e.pointerId, cur);

    if (pointers.size === 2 && pinch) {
      const [a, b] = [...pointers.values()];
      const dist = Math.hypot(a.x - b.x, a.y - b.y);
      const mid = { x: (a.x + b.x) / 2, y: (a.y + b.y) / 2 };
      const target = Math.min(MAX, Math.max(MIN * 0.8, (pinch.scale * dist) / pinch.dist));
      scale = pinch.scale;
      tx = pinch.tx + (mid.x - pinch.mid.x);
      ty = pinch.ty + (mid.y - pinch.mid.y);
      zoomAt(Math.max(target, 0.8), mid);
      apply();
    } else if (pointers.size === 1 && scale > 1) {
      tx += cur.x - prev.x;
      ty += cur.y - prev.y;
      clamp();
      apply();
    }
  };

  const up = (e: PointerEvent) => {
    const start = pointers.get(e.pointerId);
    pointers.delete(e.pointerId);
    if (pointers.size < 2) pinch = null;
    if (pointers.size === 0 && scale < 1.02) {
      scale = 1;
      tx = ty = 0;
      apply(true);
    }
    // Double tap toggles zoom around the tapped point.
    if (start && pointers.size === 0 && e.type === "pointerup") {
      const now = Date.now();
      const near = Math.hypot(e.clientX - lastTap.x, e.clientY - lastTap.y) < 30;
      if (now - lastTap.t < DOUBLE_TAP_MS && near) {
        zoomAt(scale > 1 ? MIN : DOUBLE_TAP_SCALE, { x: e.clientX, y: e.clientY });
        apply(true);
        lastTap = { t: 0, x: 0, y: 0 };
      } else {
        lastTap = { t: now, x: e.clientX, y: e.clientY };
      }
    }
  };

  container.addEventListener("pointerdown", down);
  container.addEventListener("pointermove", move);
  container.addEventListener("pointerup", up);
  container.addEventListener("pointercancel", up);
  apply();

  return () => {
    container.removeEventListener("pointerdown", down);
    container.removeEventListener("pointermove", move);
    container.removeEventListener("pointerup", up);
    container.removeEventListener("pointercancel", up);
  };
}

/** Largest candidate of an <img> (srcset "url 1200w, ..."), else its src. */
export function largestSource(img: HTMLImageElement): string {
  let best = img.currentSrc || img.src;
  let bestW = 0;
  for (const part of (img.getAttribute("srcset") ?? "").split(",")) {
    const [url, w] = part.trim().split(/\s+/);
    const n = parseInt(w ?? "", 10);
    if (url && n > bestW) {
      best = url;
      bestW = n;
    }
  }
  return best;
}
