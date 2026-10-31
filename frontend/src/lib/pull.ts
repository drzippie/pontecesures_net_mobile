// Pull-to-refresh for the page scroll (lists). Touch events are passive, so
// scrolling is never blocked; the gesture only counts when the page is at
// the very top and the finger moves down.

export const PULL_TRIGGER = 72; // px of (damped) pull needed to refresh
const MAX_PULL = 110;

export interface PullOptions {
  enabled: () => boolean;
  onPull: (distance: number) => void; // 0 when released/cancelled
  onRefresh: () => void;
}

export function attachPullToRefresh(o: PullOptions): () => void {
  let startY: number | null = null;
  let distance = 0;

  const start = (e: TouchEvent) => {
    startY = o.enabled() && window.scrollY <= 0 && e.touches.length === 1 ? e.touches[0].clientY : null;
    distance = 0;
  };
  const move = (e: TouchEvent) => {
    if (startY === null) return;
    const dy = e.touches[0].clientY - startY;
    if (dy <= 0 || window.scrollY > 0) {
      distance = 0;
      o.onPull(0);
      return;
    }
    // Damped so it feels elastic.
    distance = Math.min(MAX_PULL, dy * 0.5);
    o.onPull(distance);
  };
  const end = () => {
    if (startY !== null && distance >= PULL_TRIGGER) o.onRefresh();
    startY = null;
    distance = 0;
    o.onPull(0);
  };

  window.addEventListener("touchstart", start, { passive: true });
  window.addEventListener("touchmove", move, { passive: true });
  window.addEventListener("touchend", end, { passive: true });
  window.addEventListener("touchcancel", end, { passive: true });
  return () => {
    window.removeEventListener("touchstart", start);
    window.removeEventListener("touchmove", move);
    window.removeEventListener("touchend", end);
    window.removeEventListener("touchcancel", end);
  };
}
