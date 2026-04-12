import Alpine from "alpinejs";
import { validateIdentifier } from "../../lib/helpers";
import { api } from "../../lib/api";

// ── Auth Entry ────────────────────────────────────────────────────────────────
// Replaces: components/Auth/AuthEntry.tsx
// Usage: <div x-data="authEntry" @next="...">

Alpine.data("authEntry", () => ({
  mode: "email" as "email" | "mobile",
  identifierErr: "" as string,
  loading: false as boolean,

  get identifier(): string {
    return (Alpine.store("auth") as any).identifier;
  },

  setIdentifier(val: string) {
    (Alpine.store("auth") as any).setIdentifier(val);
  },

  get canSubmit(): boolean {
    return this.identifier.trim().length > 0 && this.identifierErr.length === 0;
  },

  switchMode(m: "email" | "mobile") {
    this.mode = m;
    this.setIdentifier("");
    (Alpine.store("auth") as any).setAuthWith(m);
    this.identifierErr = "";
  },

  validateOnBlur() {
    const auth = Alpine.store("auth") as any;
    if (!validateIdentifier(this.identifier, auth.authWith)) {
      this.identifierErr =
        auth.authWith === "mobile"
          ? "Invalid mobile number format!!, check again!"
          : "Invalid email address format!!, check again!";
    } else {
      this.identifierErr = "";
    }
  },

  async submit() {
    if (!this.canSubmit) return;
    this.loading = true;

    const res = await api.post("/user/", {
      name: "hshsysydy",
      identifier: this.identifier,
    });
    this.loading = false;

    if (!res.ok) {
      (Alpine.store("ui") as any).notify(
        res.error ?? "Something went wrong",
        "error",
      );
      return;
    }

    console.log(res);

    this.$dispatch("auth-next");
  },
}));
