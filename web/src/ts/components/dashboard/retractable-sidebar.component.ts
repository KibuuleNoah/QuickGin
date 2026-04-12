import Alpine from "alpinejs";

// ── Retractable Sidebar ───────────────────────────────────────────────────────
// Replaces: components/DashboardLayoutI/RetractableSidebar.tsx
// Usage: <aside x-data="retractableSidebar">

Alpine.data("retractableSidebar", () => ({
  get expanded(): boolean {
    return (Alpine.store("dashboard") as any).sidebarExpanded;
  },

  get activeView(): string {
    return (Alpine.store("dashboard") as any).activeView;
  },

  toggle() {
    (Alpine.store("dashboard") as any).toggleSidebar();
  },

  isActive(label: string): boolean {
    return this.activeView === label;
  },

  navigate(label: string) {
    (Alpine.store("dashboard") as any).setView(label);
  },
}));
