import { LocalDatePipe, LocalDateTimePipe } from './date-time.pipe';

describe('LocalDatePipe', () => {
  let pipe: LocalDatePipe;

  beforeEach(() => {
    pipe = new LocalDatePipe();
  });

  it('returns empty string when value is undefined', () => {
    expect(pipe.transform(undefined)).toBe('');
  });

  it('returns empty string when value is null', () => {
    expect(pipe.transform(null as unknown as string)).toBe('');
  });

  it('formats a Date using Intl with provided options', () => {
    const d = new Date(2025, 0, 9, 3, 4, 5); // Local time: Jan 09 2025
    const options: Intl.DateTimeFormatOptions = { year: 'numeric', month: '2-digit', day: '2-digit' };
    const expected = new Intl.DateTimeFormat(undefined, options).format(d);
    expect(pipe.transform(d, options)).toBe(expected);
  });

  it('formats an ISO string using Intl with provided options', () => {
    // Use a stable, explicit local date via Date.UTC -> ISO string
    const iso = new Date(Date.UTC(2025, 0, 9, 0, 0, 0)).toISOString();
    const options: Intl.DateTimeFormatOptions = { year: 'numeric', month: 'long', day: '2-digit' };
    const expected = new Intl.DateTimeFormat(undefined, options).format(new Date(iso));
    expect(pipe.transform(iso, options)).toBe(expected);
  });

  it('formats using default options when none provided', () => {
    const d = new Date(2025, 5, 1, 12, 0, 0);
    const result = pipe.transform(d);
    expect(typeof result).toBe('string');
    expect(result.length).toBeGreaterThan(0);
  });

  it('returns original value string when invalid date input is provided', () => {
    const invalid = 'not-a-date';
    expect(pipe.transform(invalid)).toBe(invalid);
  });
});

describe('LocalDateTimePipe', () => {
  let pipe: LocalDateTimePipe;

  beforeEach(() => {
    pipe = new LocalDateTimePipe();
  });

  it('returns empty string when value is undefined', () => {
    expect(pipe.transform(undefined)).toBe('');
  });

  it('formats date-time with hours and minutes', () => {
    const d = new Date(2025, 6, 15, 14, 30, 0); // Local time: Jul 15 2025 14:30
    const options: Intl.DateTimeFormatOptions = { year: 'numeric', month: 'short', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false };
    const expected = new Intl.DateTimeFormat(undefined, options).format(d);
    expect(pipe.transform(d, options)).toBe(expected);
  });

  it('handles ISO string input for date-time formatting', () => {
    const iso = new Date(Date.UTC(2025, 6, 15, 14, 30, 0)).toISOString();
    const options: Intl.DateTimeFormatOptions = { year: 'numeric', month: 'short', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false };
    const expected = new Intl.DateTimeFormat(undefined, options).format(new Date(iso));
    expect(pipe.transform(iso, options)).toBe(expected);
  });

  it('returns original value string when invalid date input is provided', () => {
    const invalid = 'totally-invalid';
    expect(pipe.transform(invalid)).toBe(invalid);
  });
});

