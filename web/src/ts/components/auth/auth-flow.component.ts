import Alpine from "alpinejs";
import { createNavTree, getRemainingCooldown } from "../../lib/helpers";

// ── Auth Flow ─────────────────────────────────────────────────────────────────
// Replaces: pages/Auth.tsx + AuthStepIndicator + AuthProvider
// Usage: <div x-data="authFlow">

Alpine.data("authFlow", () => ({
  _navTree: createNavTree("0", "moxie-nav-tree"),
  stepIndex: 0,

  get step(): number {
    return this.stepIndex;
  },

  get stepMeta() {
    return [
      { label: "Get Started Now", hint: "Enter your email or phone to begin" },
      { label: "Verify", hint: "We sent you a 6-digit code" },
      { label: "Your Profile", hint: "Almost there — just a few details" },
    ];
  },

  get currentMeta() {
    return this.stepMeta[this.stepIndex] ?? this.stepMeta[0];
  },

  next() {
    console.log(this._navTree.tree);
    const nextStep = this._navTree.tree.length;
    this._navTree.append(String(nextStep));
    this.stepIndex = nextStep;
    console.log(this._navTree.tree);
  },

  back() {
    if (this._navTree.canPop()) {
      this._navTree.pop();
      this.stepIndex = Math.max(0, this.stepIndex - 1);
    }
  },

  init() {
    this.stepIndex = this._navTree.tree.length - 1;
    // Sync with browser back button
    window.addEventListener("popstate", () => {
      this.stepIndex = Math.max(0, this.stepIndex - 1);
    });
  },
}));
