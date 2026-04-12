import Alpine from "alpinejs";

export type AlertType = "success" | "error" | "info";

export interface ToastMessage {
  id: number;
  message: string;
  type: AlertType;
}

// ── Store ─────────────────────────────────────────────────────────────────────
// Replaces: Alert.tsx + any global loading state

Alpine.store("ui", {
  loading: false as boolean,
  theme: (localStorage.getItem("theme") ?? "light") as "light" | "dark",
  toasts: [] as ToastMessage[],
  _toastId: 0 as number,

  toggleTheme() {
    this.theme = this.theme === "dark" ? "light" : "dark";
    localStorage.setItem("theme", this.theme);
    document.documentElement.classList.toggle("dark", this.theme === "dark");
  },

  notify(message: string, type: AlertType = "info", durationMs = 5500) {
    const id = ++this._toastId;
    this.toasts.push({ id, message, type });
    setTimeout(() => this.dismiss(id), durationMs);
  },

  dismiss(id: number) {
    this.toasts = this.toasts.filter((t) => t.id !== id);
  },
});

// Apply saved theme on boot
const saved = localStorage.getItem("theme") ?? "light";
document.documentElement.classList.toggle("dark", saved === "dark");
