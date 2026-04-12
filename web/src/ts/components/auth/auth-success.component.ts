import Alpine from "alpinejs";

// ── Auth Success ──────────────────────────────────────────────────────────────
// Replaces: components/Auth/AuthSuccess.tsx
// Usage: <div x-data="authSuccess">

Alpine.data("authSuccess", () => ({
  goToDashboard() {
    window.location.href = "/user";
  },
}));
