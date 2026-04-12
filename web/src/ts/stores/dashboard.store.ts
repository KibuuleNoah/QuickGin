import Alpine from "alpinejs";

export interface NavLink {
  path: string;
  label: string;
  icon: string; // SVG string or icon name — rendered with x-html
  handleOnClick?: string; // Optional JS expression string
}

// ── Store ─────────────────────────────────────────────────────────────────────
// Replaces: DashboardLayoutICtx + viewNavTree state in UserDashboard

const BASE_VIEWS = ["Home", "Settings", "Analytics", "Review Content"];

Alpine.store("dashboard", {
  activeView: "Home" as string,
  viewNavTree: ["Home"] as string[],
  sidebarExpanded: false as boolean,
  mobileNavOpen: false as boolean,

  get canGoBack(): boolean {
    return this.viewNavTree.length > 1;
  },

  setView(view: string) {
    // If navigating to a base view reset the tree
    if (BASE_VIEWS.includes(view) && this.viewNavTree.length !== 1) {
      this.viewNavTree = [view];
    } else if (!this.viewNavTree.includes(view)) {
      this.viewNavTree = [...this.viewNavTree, view];
    }

    this.activeView = view;
    this.mobileNavOpen = false;
  },

  goBack() {
    if (this.viewNavTree.length > 1) {
      const tree = [...this.viewNavTree];
      tree.pop();
      this.viewNavTree = tree;
      this.activeView = tree[tree.length - 1];
    }
  },

  toggleSidebar() {
    this.sidebarExpanded = !this.sidebarExpanded;
  },

  toggleMobileNav() {
    this.mobileNavOpen = !this.mobileNavOpen;
  },
});
