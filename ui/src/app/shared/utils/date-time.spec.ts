import { formatISODateLocal, toISOStringUTC, parseBackendISO } from './date-time';

describe('date-time utility functions', () => {
  it('formatISODateLocal returns YYYY-MM-DD for local date', () => {
    const d = new Date(2025, 0, 9, 3, 4, 5); // Jan 9, 2025 local
    const s = formatISODateLocal(d);
    // Check basic structure and that it starts with the local year-month-day
    const parts = s.split('-');
    expect(parts.length).toBe(3);
    expect(parts[0]).toBe(String(d.getFullYear()));
    expect(parts[1].length).toBe(2);
    expect(parts[2].length).toBe(2);
  });

  it('toISOStringUTC returns ISO string with Z', () => {
    const d = new Date(Date.UTC(2025, 0, 9, 0, 0, 0));
    const s = toISOStringUTC(d);
    expect(s.endsWith('Z')).toBeTrue();
  });

  it('parseBackendISO returns a Date object', () => {
    const iso = new Date(Date.UTC(2025, 0, 9, 0, 0, 0)).toISOString();
    const parsed = parseBackendISO(iso);
    expect(parsed instanceof Date).toBeTrue();
    expect(parsed.toISOString()).toBe(iso);
  });
});

