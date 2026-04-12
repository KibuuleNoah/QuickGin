// ── Time / OTP ────────────────────────────────────────────────────────────────

export interface CooldownResult {
  totalSeconds: number;
  formatted: string;
  isExpired: boolean;
}

export function getRemainingCooldown(expiresAt: string): CooldownResult {
  const expiryTime = new Date(expiresAt).getTime();
  const now = Date.now();
  const diff = expiryTime - now;

  if (diff <= 0) {
    return { totalSeconds: 0, formatted: "0:00", isExpired: true };
  }

  const totalSeconds = Math.floor(diff / 1000);
  const minutes = Math.floor(totalSeconds / 60);
  const seconds = totalSeconds % 60;

  return {
    totalSeconds,
    formatted: `${minutes}:${seconds.toString().padStart(2, "0")}`,
    isExpired: false,
  };
}

// ── PII masking ───────────────────────────────────────────────────────────────

// export function maskEmail(email: string): string {
//   const [local, domain] = email.split("@");
//   if (local.length <= 5) return email;
//   const start = local.slice(0, 2);
//   const end = local.slice(-4);
//   const masked = start + "x".repeat(local.length - 7) + end;
//   return `${masked}@${domain}`;
// }

export const maskEmail = (email: string): string => {
  const [localPart, domain] = email.split("@");

  // Guards against invalid formats
  if (!domain) return email;

  // Extract the first 2 characters
  const front = localPart.slice(0, 2);

  // Extract the last 1 character (only if the local part is long enough)
  const back = localPart.length > 2 ? localPart.slice(-1) : "";

  // Combine with 10 asterisks
  return `${front}${"*".repeat(10)}${back}@${domain}`;
};

export const saveMaskedEmail = (maskedEmail: string): void => {
  const now = new Date();

  // 24 hour * 60 min * 60 sec * 1000 ms
  const oneHourInMs = 24 * 60 * 60 * 1000;
  const expiryTime = now.getTime() + oneHourInMs;

  const item = {
    value: maskedEmail,
    expiry: expiryTime,
  };

  localStorage.setItem("user_masked_identifier", JSON.stringify(item));
};

export const getMaskedEmail = (): string | null => {
  const itemStr = localStorage.getItem("user_masked_identifier");
  if (!itemStr) return null;

  const item = JSON.parse(itemStr);
  const now = new Date().getTime();

  // If current time is past expiry, delete it
  if (now > item.expiry) {
    localStorage.removeItem("user_masked_email");
    return null;
  }

  return item.value;
};

// ── Greeting ──────────────────────────────────────────────────────────────────

export function greetUser(): string {
  const hour = new Date().getHours();
  if (hour >= 5 && hour < 12) return "Good Morning, ";
  if (hour >= 12 && hour < 17) return "Good Afternoon, ";
  if (hour >= 17 && hour < 21) return "Good Evening, ";
  return "Nights, ";
}

// ── Currency ──────────────────────────────────────────────────────────────────

export function formatUGX(value: number): string {
  return new Intl.NumberFormat("en-UG", {
    style: "currency",
    currency: "UGX",
    maximumFractionDigits: 0,
  }).format(value);
}

// ── Case conversion ───────────────────────────────────────────────────────────

export function camelToSnake(str: string): string {
  return str.replace(/[A-Z]/g, (l) => `_${l.toLowerCase()}`);
}

export function objectKeysToSnake<T extends Record<string, any>>(
  obj: T,
): Record<string, any> {
  return Object.entries(obj).reduce(
    (acc, [key, val]) => ({ ...acc, [camelToSnake(key)]: val }),
    {} as Record<string, any>,
  );
}

// ── Validation ────────────────────────────────────────────────────────────────

export function validateIdentifier(
  input: string,
  authWith: "mobile" | "email",
): boolean {
  input = input.trim();
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
  const phoneRegex = /^\+?[1-9]\d{1,14}$/;

  if (authWith === "mobile") {
    return phoneRegex.test(input.replace(/[\s\-\(\)]/g, ""));
  }
  return emailRegex.test(input);
}

// ── Nav tree (replaces useNavTree hook) ──────────────────────────────────────
// Used for the auth step wizard — push/pop with browser history integration

export function createNavTree(initialItem: string, storageKey: string) {
  const saved = localStorage.getItem(storageKey);
  let tree: string[] = saved ? JSON.parse(saved) : [initialItem];

  function getTree(): string[] {
    return tree;
  }

  function save() {
    localStorage.setItem(storageKey, JSON.stringify(tree));
  }

  function current(): string {
    return tree[tree.length - 1];
  }

  function append(item: string) {
    if (current() === item) return;
    tree = [...tree, item];
    window.history.pushState({ step: item }, "");
    save();
  }

  function pop() {
    if (tree.length > 1) {
      window.history.back();
    }
  }

  function handlePopState() {
    if (tree.length > 1) {
      tree = tree.slice(0, -1);
      save();
    }
  }

  window.addEventListener("popstate", handlePopState);

  return {
    current,
    append,
    pop,
    canPop: () => tree.length > 1,
    tree: getTree(),
  };
}
