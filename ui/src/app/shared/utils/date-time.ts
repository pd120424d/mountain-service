// Centralized date-time helpers for UTC/local conversions

// Returns YYYY-MM-DD string using the user's local calendar day (no timezone info)
export function formatISODateLocal(date: Date): string {
  const y = date.getFullYear();
  const m = date.getMonth() + 1;
  const d = date.getDate();
  const pad = (n: number) => (n < 10 ? `0${n}` : `${n}`);
  return `${y}-${pad(m)}-${pad(d)}`;
}

// Ensures a backend-friendly ISO string with trailing Z (UTC) for full timestamps
export function toISOStringUTC(date: Date): string {
  return date.toISOString();
}

// Parse backend ISO string and return a local Date instance
export function parseBackendISO(iso: string): Date {
  return new Date(iso);
}

