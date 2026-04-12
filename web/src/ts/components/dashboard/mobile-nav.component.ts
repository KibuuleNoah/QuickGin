import Alpine from "alpinejs";

// ── Mobile Nav ────────────────────────────────────────────────────────────────
// Replaces: components/DashboardLayoutI/MobileNav.tsx + useScrollDirection
// Usage: <nav x-data="mobileNav">

Alpine.data("mobileNav", () => ({
  showMore: false as boolean,
  scrollDirection: "up" as "up" | "down",
  _prevScrollY: 0 as number,
  _scrollHandler: null as ((e: Event) => void) | null,

  get activeView(): string {
    return (Alpine.store("dashboard") as any).activeView;
  },

  get isHidden(): boolean {
    return this.scrollDirection === "down";
  },

  init() {
    this._scrollHandler = () => {
      const currentY = window.pageYOffset;
      if (Math.abs(currentY - this._prevScrollY) < 10) return;

      const direction = currentY > this._prevScrollY ? "down" : "up";
      if (direction !== this.scrollDirection) {
        this.scrollDirection = direction;
      }
      this._prevScrollY = currentY > 0 ? currentY : 0;
    };

    window.addEventListener("scroll", this._scrollHandler, { passive: true });
  },

  destroy() {
    if (this._scrollHandler) {
      window.removeEventListener("scroll", this._scrollHandler);
    }
  },

  navigate(label: string) {
    (Alpine.store("dashboard") as any).setView(label);
    this.showMore = false;
  },

  isActive(label: string): boolean {
    return this.activeView === label;
  },

  toggleMore() {
    this.showMore = !this.showMore;
  },
}));
