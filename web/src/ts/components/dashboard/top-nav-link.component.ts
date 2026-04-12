import Alpine from "alpinejs";

// ── Top Nav Link ──────────────────────────────────────────────────────────────
// Replaces: components/DashboardLayoutI/TopNavLink.tsx
// Usage: <a x-data="topNavLink($el, link)">

Alpine.data("topNavLink", (label: string) => ({
  label,

  get isActive(): boolean {
    return (Alpine.store("dashboard") as any).activeView === this.label;
  },

  navigate() {
    (Alpine.store("dashboard") as any).setView(this.label);
  },
}));
