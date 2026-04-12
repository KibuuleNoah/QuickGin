import Alpine from "alpinejs";
import { api } from "../../lib/api";

// ── Auth Profile Setup ────────────────────────────────────────────────────────
// Replaces: components/Auth/AuthProfileSetup.tsx
// Usage: <div x-data="authProfile">

Alpine.data("authProfile", () => ({
  username: "" as string,
  name: "" as string,
  loading: false as boolean,

  get canSubmit(): boolean {
    return this.username.trim().length > 0 && this.name.trim().length > 0;
  },

  sanitizeUsername(val: string): string {
    return val.toLowerCase().replace(/\s/g, "_");
  },

  async submit() {
    if (!this.canSubmit) return;
    this.loading = true;

    const res = await api.post("/profile/setup", {
      name: this.name,
      username: this.username,
    });

    this.loading = false;

    if (!res.ok) {
      (Alpine.store("ui") as any).notify(
        res.error ?? "Profile setup failed",
        "error",
      );
      return;
    }

    this.$dispatch("auth-next");
  },
}));
