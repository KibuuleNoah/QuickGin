import {
  getStoredTokens,
  storeTokens,
  clearTokens,
} from "../stores/auth.store";

const BASE_URL = "http://127.0.0.1:9000/v1";

export type RequestMethod =
  | "GET"
  | "POST"
  | "PUT"
  | "PATCH"
  | "DELETE"
  | "HEAD";

export interface ApiResponse<T = any> {
  data: T | null;
  status: number;
  ok: boolean;
  error?: Promise<string>;
}

// ── Core fetch wrapper ────────────────────────────────────────────────────────

async function fetchWithAuth(
  method: RequestMethod,
  path: string,
  body?: unknown,
  retry = false,
): Promise<Response> {
  const tokens = getStoredTokens();

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (tokens?.accessToken) {
    headers["Authorization"] = `Bearer ${tokens.accessToken}`;
  }

  const res = await fetch(`${BASE_URL}${path}`, {
    method,
    headers,
    body: body !== undefined ? JSON.stringify(body) : undefined,
  });

  // ── Auto-refresh on 401 ──
  if (res.status === 401 && !retry) {
    const stored = getStoredTokens();

    if (stored?.refreshToken) {
      try {
        const refreshRes = await fetch(`${BASE_URL}/auth/refresh`, {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ token: stored.refreshToken }),
        });

        if (refreshRes.ok) {
          const newTokens = await refreshRes.json();
          storeTokens(newTokens);
          // Retry once with the fresh token
          return fetchWithAuth(method, path, body, true);
        }
      } catch {
        // Refresh failed — clear session
      }
    }

    clearTokens();
    window.location.href = "/auth";
  }

  return res;
}

// ── Public API ────────────────────────────────────────────────────────────────

export const api = {
  async request<T = any>(
    method: RequestMethod,
    path: string,
    body?: unknown,
  ): Promise<ApiResponse<T>> {
    try {
      const res = await fetchWithAuth(method, path, body);
      const error = !res.ok
        ? res.json().then((err) => err?.message || "Something went Wrong!")
        : Promise.resolve("");

      const data =
        !error && res.status !== 204
          ? await res.json().catch(() => null)
          : null;

      return { data, error, status: res.status, ok: res.ok };
    } catch (err: any) {
      return {
        data: null,
        status: 0,
        ok: false,
        error: err?.message ?? "Network error",
      };
    }
  },

  get<T = any>(path: string) {
    return this.request<T>("GET", path);
  },

  post<T = any>(path: string, body: unknown) {
    return this.request<T>("POST", path, body);
  },

  put<T = any>(path: string, body: unknown) {
    return this.request<T>("PUT", path, body);
  },

  patch<T = any>(path: string, body: unknown) {
    return this.request<T>("PATCH", path, body);
  },

  delete<T = any>(path: string) {
    return this.request<T>("DELETE", path);
  },
};
