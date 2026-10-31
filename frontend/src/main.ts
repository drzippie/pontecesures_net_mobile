import Alpine from "alpinejs";
import { Events } from "@wailsio/runtime";
import intersect from "@alpinejs/intersect";
import "@fontsource-variable/fraunces";
import "@fontsource-variable/source-sans-3";
import "./style.css";

import { ContentService, SettingsService, WeatherService } from "../bindings/github.com/drzippie/pontecesures_net_mobile";
import type { Forecast } from "../bindings/github.com/drzippie/pontecesures_net_mobile/internal/meteo/models";
import type { Category, Page, Post, PostSummary } from "../bindings/github.com/drzippie/pontecesures_net_mobile/internal/wp/models";
import { dictionaries, fmt, type Strings } from "./i18n";
import { longDate, number, relativeDate } from "./lib/format";
import { applyPrefs, loadLocal, normalize, saveLocal, type Prefs } from "./lib/prefs";
import { href, isInternalPostLink, parseHash, type Route } from "./lib/router";
import { dayDate, dayLabel, hourMinute, rainMax, skyIcon, skyText, windText } from "./lib/weather";
import { attachZoom, largestSource } from "./lib/zoom";
import { attachPullToRefresh, PULL_TRIGGER } from "./lib/pull";

/** After this long in the background, the app refreshes when it comes back. */
const RESUME_REFRESH_MS = 15 * 60 * 1000;

/** Re-fetch the forecast when the copy shown is older than this. */
const WEATHER_TTL_MS = 30 * 60 * 1000;

const APP_VERSION = "0.1.1";
const SITE = "https://pontecesures.net";

/** State of one paginated list (a section feed or the search results). */
interface ListState {
  items: PostSummary[];
  page: number; // last page loaded
  totalPages: number;
  total: number;
  loading: boolean;
  error: ErrorKind;
  stale: boolean;
  scrollY: number;
}

function emptyList(): ListState {
  return { items: [], page: 0, totalPages: 1, total: 0, loading: false, error: "", stale: false, scrollY: 0 };
}

/** Error kind, rendered with t.error[kind] so it follows the language. */
type ErrorKind = "" | "generic" | "network" | "notFound";

function errorKind(err: unknown): ErrorKind {
  const msg = String((err as any)?.message ?? err ?? "");
  if (/not found/i.test(msg)) return "notFound";
  if (/dial|timeout|no such host|network|connection|deadline/i.test(msg)) return "network";
  return "generic";
}

// Apply saved preferences before Alpine renders anything.
const initialPrefs = normalize(loadLocal());
applyPrefs(initialPrefs);

/**
 * Android draws the WebView edge-to-edge but does not report the system bar
 * insets through CSS env(); read them from Go (physical px) and expose them
 * as CSS variables used by the pt-safe/pb-safe utilities.
 */
async function applyAndroidSafeArea() {
  // System.IsAndroid() is not reliable here: window._wails.environment is
  // not populated in the Android WebView (Wails v3.0.0).
  if (!/\bAndroid\b/.test(navigator.userAgent)) return;
  try {
    const sa = await ContentService.SafeArea();
    const dpr = window.devicePixelRatio || 1;
    const root = document.documentElement.style;
    root.setProperty("--safe-top", `${sa.top / dpr}px`);
    root.setProperty("--safe-bottom", `${sa.bottom / dpr}px`);
  } catch {
    /* keep env() values */
  }
}
applyAndroidSafeArea();
window.addEventListener("resize", () => applyAndroidSafeArea());

Alpine.plugin(intersect);

Alpine.data("app", () => ({
  prefs: initialPrefs as Prefs,
  t: dictionaries[initialPrefs.language] as Strings,
  fmt,
  version: APP_VERSION,

  route: { name: "home" } as Route,
  /** Last non-article route, shown underneath the article overlay. */
  base: { name: "home" } as Route,
  /**
   * Hashes visited in this session, newest last. Lets "Volver" know whether
   * there is an in-app page to go back to, also after the system back
   * gesture (Android back button, iOS edge swipe) moved through history.
   */
  navStack: [] as string[],

  categories: [] as Category[],
  feeds: {} as Record<number, ListState>,
  /** Saved articles (offline reading). */
  savedList: emptyList(),
  savedIds: [] as number[],
  search: { ...emptyList(), q: "", ran: "", category: 0, order: "relevance" as "relevance" | "date" },
  searchTimer: undefined as number | undefined,

  article: { id: 0, data: null as Post | PostSummary | null, full: false, loading: false, error: "" as ErrorKind },
  /** Static website page shown in #/pagina/:id. */
  page: { id: 0, data: null as Page | null, loading: false, error: "" as ErrorKind },
  /** "Más sobre esto" under the open article. */
  related: { id: 0, items: [] as PostSummary[], loading: false },
  weather: { data: null as Forecast | null, loading: false, error: false, loadedAt: 0 },
  /** Full-screen image viewer. */
  lightbox: { open: false, src: "", alt: "", hint: false },
  detachZoom: null as null | (() => void),

  /** Pull-to-refresh: current (damped) pull distance in px. */
  pull: 0,
  pullTrigger: PULL_TRIGGER,
  hiddenAt: 0,

  toast: "",
  toastTimer: undefined as number | undefined,

  init() {
    window.addEventListener("hashchange", () => this.onRoute());
    this.onRoute();
    // Called by the Android host (MainActivity.onBackPressed).
    window.pcnHandleBack = () => this.handleSystemBack();

    attachPullToRefresh({
      enabled: () => this.isListView && this.route.name !== "article" && !this.lightbox.open && !this.list.loading,
      onPull: (d) => (this.pull = d),
      onRefresh: () => {
        this.refresh();
        if (this.base.name === "home") this.loadWeather(true);
      },
    });

    // Native lifecycle (lifecycle.go): the WebView's visibilitychange does not
    // fire on Android, so rely on the app events, with visibilitychange as a
    // fallback for desktop/browsers.
    Events.On("app:foreground", (e) => this.onForeground(Number(e.data) || 0));
    document.addEventListener("visibilitychange", () => {
      if (document.visibilityState === "hidden") this.hiddenAt = Date.now();
      else if (this.hiddenAt) this.onForeground(Date.now() - this.hiddenAt);
    });
    this.loadCategories();
    this.loadWeather();
    this.syncSettings();
    this.loadSavedIds();
  },

  /**
   * Refresh stale content when the app returns from the background.
   * `away` comes from Go on mobile (lifecycle.go), or from visibilitychange
   * on desktop/browsers.
   */
  onForeground(away: number) {
    this.hiddenAt = 0;
    if (away < RESUME_REFRESH_MS) return;
    this.loadWeather(true);
    // Only reload the list if the reader is near the top: never move content
    // under someone who is reading.
    if (this.isListView && this.route.name !== "article" && window.scrollY < 200 && this.base.name !== "search") {
      this.refresh();
    }
  },

  // ---------- settings ----------

  /** Go is the source of truth; adopt its values or seed it from local ones. */
  async syncSettings() {
    try {
      const saved = await SettingsService.Get();
      if (saved.language || saved.theme || saved.textSize) {
        this.setPrefs(normalize(saved), false);
      } else if (loadLocal()) {
        await SettingsService.Save(this.prefs);
      }
    } catch {
      /* keep local preferences */
    }
  },

  setPrefs(p: Prefs, persist = true) {
    this.prefs = p;
    this.t = dictionaries[p.language];
    applyPrefs(p);
    saveLocal(p);
    if (persist) SettingsService.Save(p).catch(() => this.showToast(this.t.settings.saveError));
  },

  setPref<K extends keyof Prefs>(key: K, value: Prefs[K]) {
    if (this.prefs[key] === value) return;
    this.setPrefs({ ...this.prefs, [key]: value });
  },

  // ---------- navigation ----------

  go(r: Route) {
    location.hash = href(r);
  },

  back() {
    if (this.navStack.length > 1) history.back();
    else this.go(this.base.name === "article" ? { name: "home" } : this.base);
  },

  /**
   * System back (Android back button). Returns true when handled in-app;
   * false on the front page so the OS can leave the app.
   */
  handleSystemBack(): boolean {
    if (this.lightbox.open) {
      this.closeLightbox();
      return true;
    }
    if (this.route.name === "home" && this.navStack.length <= 1) return false;
    this.back();
    return true;
  },

  /** Keeps navStack in sync: a hash equal to the previous entry is a "back". */
  trackNavigation() {
    const h = location.hash || "#/";
    const st = this.navStack;
    if (st.length > 1 && st[st.length - 2] === h) st.pop();
    else if (st[st.length - 1] !== h) st.push(h);
  },

  onRoute() {
    const prev = this.route;
    const r = parseHash(location.hash);
    this.trackNavigation();

    // Remember the scroll position of the list we are leaving.
    if (prev.name !== "article" && r.name !== "article") {
      const leaving = this.listFor(prev);
      if (leaving) leaving.scrollY = window.scrollY;
    }

    this.route = r;
    if (r.name === "article") {
      this.openArticle(r.id);
      return;
    }
    this.base = r;
    this.article.id = 0;

    const list = this.listFor(r);
    // The saved list can change from any article: reload it on every visit.
    if (r.name === "saved") Object.assign(this.savedList, emptyList());
    if (list && list.page === 0 && r.name !== "search") this.loadMore();
    // Coming back from an article the list never moved (the article is an
    // overlay), so only restore the position when switching between lists.
    if (prev.name !== "article") this.$nextTick(() => window.scrollTo(0, list?.scrollY ?? 0));
    if (r.name === "search") this.$nextTick(() => (this.$refs.searchInput as HTMLInputElement)?.focus());
    if (r.name === "weather" || r.name === "home") this.loadWeather();
    if (r.name === "page") this.loadPage(r.id);
  },

  get currentCategory(): number {
    return this.base.name === "section" ? this.base.id : 0;
  },

  get isListView(): boolean {
    return ["home", "section", "search", "saved"].includes(this.base.name);
  },

  get sectionName(): string {
    return this.categories.find((c) => c.id === this.currentCategory)?.name ?? "";
  },

  listFor(r: Route): (ListState & { q?: string }) | null {
    if (r.name === "search") return this.search;
    if (r.name === "saved") return this.savedList;
    if (r.name !== "home" && r.name !== "section") return null;
    const id = r.name === "section" ? r.id : 0;
    if (!this.feeds[id]) this.feeds[id] = emptyList();
    return this.feeds[id];
  },

  get list() {
    return this.listFor(this.base) ?? emptyList();
  },

  // ---------- data ----------

  async loadCategories() {
    try {
      this.categories = (await ContentService.Categories()) ?? [];
    } catch {
      /* the chips simply stay hidden */
    }
  },

  async loadMore() {
    const list = this.list;
    if (list.loading || list.page >= list.totalPages) return;
    if (this.base.name === "saved") {
      list.items = (await ContentService.Saved()) ?? [];
      list.page = list.totalPages = 1;
      list.total = list.items.length;
      return;
    }
    const isSearch = this.base.name === "search";
    const q = isSearch ? this.search.q.trim() : "";
    if (isSearch && q.length < 3) return;

    list.loading = true;
    list.error = "";
    try {
      const category = isSearch ? this.search.category : this.currentCategory;
      const order = isSearch ? this.search.order : "date";
      const res = await ContentService.Posts(list.page + 1, category, q, order);
      // Ignore late answers to an outdated search (text, section or order changed).
      if (isSearch && (q !== this.search.q.trim() || category !== this.search.category || order !== this.search.order)) return;
      const seen = new Set(list.items.map((p) => p.id));
      list.items.push(...(res.items ?? []).filter((p) => !seen.has(p.id)));
      list.page = res.page;
      list.totalPages = res.totalPages || res.page;
      list.total = res.total;
      list.stale = res.stale;
      if (isSearch) this.search.ran = q;
    } catch (err) {
      list.error = errorKind(err);
    } finally {
      list.loading = false;
    }
    // If the sentinel is still on screen (short page), keep filling.
    this.$nextTick(() => {
      const s = this.$refs.sentinel as HTMLElement | undefined;
      if (!list.error && list.page < list.totalPages && s && s.getBoundingClientRect().top < window.innerHeight + 600) {
        this.loadMore();
      }
    });
  },

  async refresh() {
    const list = this.list;
    Object.assign(list, emptyList(), this.base.name === "search" ? { q: this.search.q } : {});
    window.scrollTo({ top: 0, behavior: "smooth" });
    await this.loadMore();
  },

  onSearchInput() {
    clearTimeout(this.searchTimer);
    this.searchTimer = window.setTimeout(() => {
      Object.assign(this.search, emptyList());
      this.search.ran = "";
      this.loadMore();
    }, 350);
  },

  /** Section or order changed: run the search again from page 1. */
  rerunSearch() {
    Object.assign(this.search, emptyList(), { ran: "" });
    this.loadMore();
  },

  clearSearch() {
    Object.assign(this.search, emptyList(), { q: "", ran: "" });
    (this.$refs.searchInput as HTMLInputElement)?.focus();
  },

  // ---------- weather ----------

  async loadWeather(force = false) {
    const w = this.weather;
    if (w.loading) return;
    if (!force && w.data && Date.now() - w.loadedAt < WEATHER_TTL_MS) return;
    w.loading = true;
    try {
      w.data = await WeatherService.Forecast();
      w.error = false;
      w.loadedAt = Date.now();
    } catch {
      w.error = true;
    } finally {
      w.loading = false;
    }
  },

  get today() {
    return this.weather.data?.days?.[0] ?? null;
  },

  skyIcon,
  rainMax,
  skyText(code: number | null | undefined): string {
    return skyText(code, this.t);
  },
  windText(code: number | null | undefined): string {
    return windText(code, this.t);
  },
  dayLabel(date: string): string {
    return dayLabel(date, this.t);
  },
  dayDate(date: string): string {
    return dayDate(date, this.t.locale);
  },
  hourMinute(iso: string): string {
    return hourMinute(iso, this.t.locale);
  },

  // ---------- article ----------

  showPost(p: PostSummary) {
    // Show what we already have immediately; the full body loads on route.
    this.article = { id: p.id, data: p, full: false, loading: true, error: "" };
    this.go({ name: "article", id: p.id });
  },

  async openArticle(id: number) {
    if (this.article.id !== id) {
      this.article = { id, data: null, full: false, loading: true, error: "" };
    }
    this.loadRelated(id);
    this.$nextTick(() => {
      const el = this.$refs.articleScroll as HTMLElement | undefined;
      if (el) el.scrollTop = 0;
      // Screen readers start reading at the article title.
      (this.$refs.articleTitle as HTMLElement | undefined)?.focus({ preventScroll: true });
    });
    try {
      const post = await ContentService.Post(id);
      if (this.article.id === id) this.article = { id, data: post, full: true, loading: false, error: "" };
    } catch (err) {
      if (this.article.id === id) Object.assign(this.article, { loading: false, error: errorKind(err) });
    }
  },

  async loadRelated(id: number) {
    if (this.related.id === id && (this.related.loading || this.related.items.length)) return;
    this.related = { id, items: [], loading: true };
    try {
      const items = (await ContentService.Related(id)) ?? [];
      if (this.related.id === id) this.related = { id, items, loading: false };
    } catch {
      // Related posts are optional: just leave the section out.
      if (this.related.id === id) this.related.loading = false;
    }
  },

  // ---------- image viewer ----------

  openLightbox(img: HTMLImageElement) {
    this.lightbox = { open: true, src: largestSource(img), alt: img.alt, hint: true };
    // The usage hint is only needed for a moment; it would cover the image.
    window.setTimeout(() => (this.lightbox.hint = false), 2500);
    this.$nextTick(() => {
      const box = this.$refs.lightbox as HTMLElement;
      const big = this.$refs.lightboxImg as HTMLImageElement;
      this.detachZoom?.();
      this.detachZoom = attachZoom(box, big);
      (this.$refs.lightboxClose as HTMLElement | undefined)?.focus();
    });
  },

  closeLightbox() {
    this.detachZoom?.();
    this.detachZoom = null;
    this.lightbox = { open: false, src: "", alt: "", hint: false };
  },

  // ---------- static pages ----------

  async loadPage(id: number) {
    if (this.page.id === id && (this.page.data || this.page.loading)) return;
    this.page = { id, data: null, loading: true, error: "" };
    try {
      const data = await ContentService.Page(id);
      if (this.page.id === id) this.page = { id, data, loading: false, error: "" };
    } catch (err) {
      if (this.page.id === id) Object.assign(this.page, { loading: false, error: errorKind(err) });
    }
  },

  // ---------- saved articles ----------

  async loadSavedIds() {
    try {
      this.savedIds = (await ContentService.SavedIDs()) ?? [];
    } catch {
      /* none */
    }
  },

  isSaved(id: number): boolean {
    return this.savedIds.includes(id);
  },

  async toggleSaved() {
    const id = this.article.data?.id;
    if (!id) return;
    try {
      if (this.isSaved(id)) {
        await ContentService.RemoveSaved(id);
        this.savedIds = this.savedIds.filter((x) => x !== id);
        this.showToast(this.t.article.removedToast);
      } else {
        await ContentService.SaveArticle(id);
        this.savedIds = [...this.savedIds, id];
        this.showToast(this.t.article.savedToast);
      }
    } catch {
      this.showToast(this.t.error.generic);
    }
  },

  /** Intercepts taps inside an article body: images open the viewer, links navigate. */
  async onContentClick(ev: MouseEvent) {
    const target = ev.target as HTMLElement;
    const img = target.closest("img");
    // Images open full screen, unless they are the content of a link.
    if (img && !target.closest("a")) {
      ev.preventDefault();
      this.openLightbox(img);
      return;
    }
    const a = target.closest("a");
    if (!a) return;
    // A linked image that points to an image file is also shown in the viewer.
    const inner = a.querySelector("img");
    if (inner && /\.(jpe?g|png|webp|gif)(\?|$)/i.test(a.getAttribute("href") ?? "")) {
      ev.preventDefault();
      this.openLightbox(inner);
      return;
    }
    const link = a.getAttribute("href") ?? "";
    ev.preventDefault();
    if (!link || link.startsWith("#")) return;
    const abs = new URL(link, SITE).toString();
    if (isInternalPostLink(abs)) {
      this.article = { id: -1, data: null, full: false, loading: true, error: "" as ErrorKind };
      try {
        const post = await ContentService.PostByURL(abs);
        this.article = { id: post.id, data: post, full: true, loading: false, error: "" };
        this.go({ name: "article", id: post.id });
        return;
      } catch {
        this.article.loading = false;
      }
    }
    this.openExternal(abs);
  },

  async share() {
    const p = this.article.data;
    if (!p) return;
    const native = await ContentService.Share(p.title, p.link).catch(() => false);
    if (!native) this.showToast(this.t.article.linkCopied);
  },

  openExternal(link: string) {
    ContentService.OpenURL(link).catch(() => this.showToast(this.t.error.generic));
  },

  openSite(path = "/") {
    this.openExternal(SITE + path);
  },

  showToast(msg: string) {
    this.toast = msg;
    clearTimeout(this.toastTimer);
    this.toastTimer = window.setTimeout(() => (this.toast = ""), 2500);
  },

  // ---------- view helpers ----------

  when(p: PostSummary): string {
    return relativeDate(p.date, this.t.locale, this.t.time.now);
  },

  longDate(iso: string): string {
    return longDate(iso, this.t.locale);
  },

  number(n: number): string {
    return number(n, this.t.locale);
  },

  mainCategory(p: PostSummary): string {
    const cats = (p.categories ?? []).filter((c) => c.slug !== "sincategoria");
    return (cats[0] ?? p.categories?.[0])?.name ?? "";
  },

  isTab(name: string): boolean {
    // Sections and the weather view are reached from the front page; the
    // settings from "More".
    const b = this.base.name;
    const n = b === "section" || b === "weather" ? "home" : ["settings", "saved", "page"].includes(b) ? "more" : b;
    return n === name;
  },
}));

window.Alpine = Alpine;
Alpine.start();

declare global {
  interface Window {
    Alpine: typeof Alpine;
    pcnHandleBack?: () => boolean;
  }
}
