import Alpine from "alpinejs";
import intersect from "@alpinejs/intersect";
import "@fontsource-variable/fraunces";
import "@fontsource-variable/source-sans-3";
import "./style.css";

import { ContentService } from "../bindings/github.com/drzippie/pontecesures_net_mobile";
import type { Category, Post, PostSummary } from "../bindings/github.com/drzippie/pontecesures_net_mobile/internal/wp/models";
import { es, fmt } from "./i18n/es";
import { longDate, number, relativeDate } from "./lib/format";
import { href, isInternalPostLink, parseHash, type Route } from "./lib/router";

const APP_VERSION = "0.1.0";
const SITE = "https://pontecesures.net";

/** State of one paginated list (a section feed or the search results). */
interface ListState {
  items: PostSummary[];
  page: number; // last page loaded
  totalPages: number;
  total: number;
  loading: boolean;
  error: string;
  stale: boolean;
  scrollY: number;
}

function emptyList(): ListState {
  return { items: [], page: 0, totalPages: 1, total: 0, loading: false, error: "", stale: false, scrollY: 0 };
}

function errorText(err: unknown): string {
  const msg = String((err as any)?.message ?? err ?? "");
  if (/not found/i.test(msg)) return es.error.notFound;
  if (/dial|timeout|no such host|network|connection|deadline/i.test(msg)) return es.error.network;
  return es.error.generic;
}

Alpine.plugin(intersect);

Alpine.data("app", () => ({
  t: es,
  fmt,
  number,
  version: APP_VERSION,

  route: { name: "home" } as Route,
  /** Last non-article route, shown underneath the article overlay. */
  base: { name: "home" } as Route,
  /** Number of in-app navigations, to know whether "back" can pop history. */
  depth: 0,

  categories: [] as Category[],
  feeds: {} as Record<number, ListState>,
  search: { ...emptyList(), q: "", ran: "" },
  searchTimer: undefined as number | undefined,

  article: { id: 0, data: null as Post | PostSummary | null, full: false, loading: false, error: "" },
  toast: "",
  toastTimer: undefined as number | undefined,

  init() {
    window.addEventListener("hashchange", () => this.onRoute(true));
    this.onRoute(false);
    this.loadCategories();
  },

  // ---------- navigation ----------

  go(r: Route) {
    location.hash = href(r);
  },

  back() {
    if (this.depth > 0) history.back();
    else this.go(this.base.name === "article" ? { name: "home" } : this.base);
  },

  onRoute(fromNav: boolean) {
    const prev = this.route;
    const r = parseHash(location.hash);
    if (fromNav) this.depth++;

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
    if (list && list.page === 0 && r.name !== "search") this.loadMore();
    // Coming back from an article the list never moved (the article is an
    // overlay), so only restore the position when switching between lists.
    if (prev.name !== "article") this.$nextTick(() => window.scrollTo(0, list?.scrollY ?? 0));
    if (r.name === "search") this.$nextTick(() => (this.$refs.searchInput as HTMLInputElement)?.focus());
  },

  get currentCategory(): number {
    return this.base.name === "section" ? this.base.id : 0;
  },

  get isListView(): boolean {
    return ["home", "section", "search"].includes(this.base.name);
  },

  get sectionName(): string {
    return this.categories.find((c) => c.id === this.currentCategory)?.name ?? "";
  },

  listFor(r: Route): (ListState & { q?: string }) | null {
    if (r.name === "search") return this.search;
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
    const isSearch = this.base.name === "search";
    const q = isSearch ? this.search.q.trim() : "";
    if (isSearch && q.length < 3) return;

    list.loading = true;
    list.error = "";
    try {
      const res = await ContentService.Posts(list.page + 1, isSearch ? 0 : this.currentCategory, q);
      // Ignore late answers to an outdated search.
      if (isSearch && q !== this.search.q.trim()) return;
      const seen = new Set(list.items.map((p) => p.id));
      list.items.push(...(res.items ?? []).filter((p) => !seen.has(p.id)));
      list.page = res.page;
      list.totalPages = res.totalPages || res.page;
      list.total = res.total;
      list.stale = res.stale;
      if (isSearch) this.search.ran = q;
    } catch (err) {
      list.error = errorText(err);
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

  clearSearch() {
    Object.assign(this.search, emptyList(), { q: "", ran: "" });
    (this.$refs.searchInput as HTMLInputElement)?.focus();
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
    this.$nextTick(() => {
      const el = this.$refs.articleScroll as HTMLElement | undefined;
      if (el) el.scrollTop = 0;
    });
    try {
      const post = await ContentService.Post(id);
      if (this.article.id === id) this.article = { id, data: post, full: true, loading: false, error: "" };
    } catch (err) {
      if (this.article.id === id) Object.assign(this.article, { loading: false, error: errorText(err) });
    }
  },

  /** Intercepts taps on links inside an article body. */
  async onContentClick(ev: MouseEvent) {
    const a = (ev.target as HTMLElement).closest("a");
    if (!a) return;
    const link = a.getAttribute("href") ?? "";
    ev.preventDefault();
    if (!link || link.startsWith("#")) return;
    const abs = new URL(link, SITE).toString();
    if (isInternalPostLink(abs)) {
      this.article = { id: -1, data: null, full: false, loading: true, error: "" };
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
    if (!native) this.showToast(es.article.linkCopied);
  },

  openExternal(link: string) {
    ContentService.OpenURL(link).catch(() => this.showToast(es.error.generic));
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
    return relativeDate(p.date, es.time.now);
  },

  longDate,

  mainCategory(p: PostSummary): string {
    const cats = (p.categories ?? []).filter((c) => c.slug !== "sincategoria");
    return (cats[0] ?? p.categories?.[0])?.name ?? "";
  },

  isTab(name: string): boolean {
    const n = this.base.name === "section" ? "home" : this.base.name;
    return n === name;
  },
}));

window.Alpine = Alpine;
Alpine.start();

declare global {
  interface Window {
    Alpine: typeof Alpine;
  }
}
