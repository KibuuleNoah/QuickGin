import Alpine from "alpinejs";
import { formatUGX, greetUser } from "../../lib/helpers";
import type { NavLink } from "../../stores/dashboard.store";

// ── Dashboard Layout ──────────────────────────────────────────────────────────
// Replaces: layouts/DashboardLayoutI.tsx + DashboardLayoutProvider
// Usage: <div x-data="dashboardLayout">

Alpine.data("dashboardLayout", () => ({
  get activeView(): string {
    return (Alpine.store("dashboard") as any).activeView;
  },

  get sidebarExpanded(): boolean {
    return (Alpine.store("dashboard") as any).sidebarExpanded;
  },

  navigate(view: string) {
    (Alpine.store("dashboard") as any).setView(view);
  },

  goBack() {
    (Alpine.store("dashboard") as any).goBack();
  },

  get canGoBack(): boolean {
    return (Alpine.store("dashboard") as any).canGoBack;
  },

  // Utilities available in templates
  formatUGX,
  greetUser,
}));
