import { getActivityIcon, formatActivityDescription, getActivityDisplayTime, Activity } from './activity-extensions';

describe('activity-extensions', () => {
  it('getActivityIcon returns expected icons based on description keywords', () => {
    const cases: Array<{desc: string, icon: string}> = [
      { desc: 'New employee created', icon: 'person' },
      { desc: 'Urgency assigned', icon: 'warning' },
      { desc: 'Shift scheduled', icon: 'schedule' },
      { desc: 'Notification sent', icon: 'notifications' },
      // Avoid overlapping keyword 'user' that would map to 'person' first
      { desc: 'Account login', icon: 'login' },
      { desc: 'Miscellaneous', icon: 'info' },
    ];

    for (const c of cases) {
      expect(getActivityIcon({ description: c.desc } as Activity)).toBe(c.icon);
    }
  });

  it('formatActivityDescription returns description or fallback', () => {
    expect(formatActivityDescription({ description: 'Hello' })).toBe('Hello');
    expect(formatActivityDescription({} as Activity)).toBe('No description available');
  });

  it('getActivityDisplayTime formats relative time correctly', () => {
    const now = new Date();

    const justNow: Activity = { createdAt: new Date(now.getTime() - 30 * 1000).toISOString() };
    expect(getActivityDisplayTime(justNow)).toBe('Just now');

    const tenMinutes: Activity = { createdAt: new Date(now.getTime() - 10 * 60 * 1000).toISOString() };
    expect(getActivityDisplayTime(tenMinutes)).toContain('10 minute');

    const fiveHours: Activity = { createdAt: new Date(now.getTime() - 5 * 60 * 60 * 1000).toISOString() };
    expect(getActivityDisplayTime(fiveHours)).toContain('5 hour');

    const twoDays: Activity = { createdAt: new Date(now.getTime() - 2 * 24 * 60 * 60 * 1000).toISOString() };
    expect(getActivityDisplayTime(twoDays)).toContain('2 day');

    const eightDays: Activity = { createdAt: new Date(now.getTime() - 8 * 24 * 60 * 60 * 1000).toISOString() };
    expect(getActivityDisplayTime(eightDays)).toMatch(/\d{1,2}\/\d{1,2}\/\d{2,4}/);
  });
});

