import Alpine from "alpinejs";
import persist from "@alpinejs/persist";

export interface AuthTokens {
  accessToken: string;
  refreshToken: string;
  atExpires: string;
  rtExpires: string;
}

export type AuthWith = "mobile" | "email";

// ── Helpers ──────────────────────────────────────────────────────────────────

export function getStoredTokens(): AuthTokens | null {
  const raw = localStorage.getItem("auth-creds");
  return raw ? JSON.parse(raw) : null;
}

export function storeTokens(tokens: AuthTokens): void {
  localStorage.setItem("auth-creds", JSON.stringify(tokens));
}

export function clearTokens(): void {
  localStorage.removeItem("auth-creds");
}

// ── Store ─────────────────────────────────────────────────────────────────────
// Replaces: AuthCtx + useAuth hook + axiosInstance token logic

Alpine.store("auth", {
  // Auth flow state
  identifier: "" as string,
  maskedIdentifier: "no***********3@gmail.com",
  authWith: "email" as AuthWith,
  otpExpiryDate: "" as string,

  // Session state (persisted)
  tokens: getStoredTokens(),

  get isAuthenticated(): boolean {
    return !!this.tokens?.accessToken;
  },

  setIdentifier(val: string) {
    this.identifier = val;
  },

  setAuthWith(val: AuthWith) {
    this.authWith = val;
  },

  setOtpExpiry(isoDate: string) {
    this.otpExpiryDate = isoDate;
    localStorage.setItem("otp_expiry_date", isoDate);
  },

  login(tokens: AuthTokens) {
    this.tokens = tokens;
    storeTokens(tokens);
  },

  logout() {
    this.tokens = null;
    clearTokens();
    window.location.href = "/auth";
  },
});
