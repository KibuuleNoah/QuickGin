import Alpine from "alpinejs";
import { getRemainingCooldown } from "../../lib/helpers";

// ── Auth OTP ──────────────────────────────────────────────────────────────────
// Replaces: components/Auth/AuthOTP.tsx
// Usage: <div x-data="authOtp">

Alpine.data("authOtp", () => ({
  digits: ["", "", "", "", "", ""] as string[],
  timeLeft: { totalSeconds: 0, formatted: "0:00", isExpired: true },
  _timer: null as ReturnType<typeof setInterval> | null,

  get filled(): boolean {
    return this.digits.every((d) => d !== "");
  },

  get otpValue(): string {
    return this.digits.join("");
  },

  init() {
    // Set expiry to 2 minutes from now (same as the React component's 12s demo — use real value in prod)
    const expiry = new Date(Date.now() + 120_000).toISOString();
    (Alpine.store("auth") as any).setOtpExpiry(expiry);
    this._startTimer(expiry);
  },

  destroy() {
    if (this._timer) clearInterval(this._timer);
  },

  _startTimer(expiry: string) {
    this.timeLeft = getRemainingCooldown(expiry);
    if (this.timeLeft.isExpired) return;

    this._timer = setInterval(() => {
      this.timeLeft = getRemainingCooldown(expiry);
      if (this.timeLeft.isExpired && this._timer) {
        clearInterval(this._timer);
        this._timer = null;
      }
    }, 1000);
  },

  handleKeydown(e: KeyboardEvent, index: number) {
    if (e.key === "Backspace") {
      if (this.digits[index]) {
        this.digits[index] = "";
      } else if (index > 0) {
        (document.getElementById(`otp-${index - 1}`) as HTMLInputElement)?.focus();
      }
      return;
    }

    if (/^[0-9]$/.test(e.key)) {
      e.preventDefault();
      this.digits[index] = e.key;
      if (index < this.digits.length - 1) {
        (document.getElementById(`otp-${index + 1}`) as HTMLInputElement)?.focus();
      }
    }
  },

  async resend() {
    if (!this.timeLeft.isExpired) return;
    this.digits = ["", "", "", "", "", ""];
    // TODO: call resend endpoint
    (Alpine.store("ui") as any).notify("OTP resent!", "info");
  },

  async verify() {
    if (!this.filled) return;
    // TODO: call verify endpoint with this.otpValue
    this.$dispatch("auth-next");
  },
}));
