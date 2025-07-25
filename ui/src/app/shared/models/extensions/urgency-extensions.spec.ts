import {
  getUrgencyLevelColor,
  getUrgencyStatusColor
} from './urgency-extensions';
import {
  UrgencyLevel,
  UrgencyStatus
} from '../generated/urgency';

describe('Urgency Extensions', () => {
  describe('getUrgencyLevelColor', () => {
    it('should return green for low urgency', () => {
      expect(getUrgencyLevelColor(UrgencyLevel.Low)).toBe('green');
    });

    it('should return yellow for medium urgency', () => {
      expect(getUrgencyLevelColor(UrgencyLevel.Medium)).toBe('yellow');
    });

    it('should return orange for high urgency', () => {
      expect(getUrgencyLevelColor(UrgencyLevel.High)).toBe('orange');
    });

    it('should return red for critical urgency', () => {
      expect(getUrgencyLevelColor(UrgencyLevel.Critical)).toBe('red');
    });

    it('should return gray for unknown urgency', () => {
      expect(getUrgencyLevelColor('unknown' as any)).toBe('gray');
    });
  });

  describe('getUrgencyStatusColor', () => {
    it('should return blue for open status', () => {
      expect(getUrgencyStatusColor(UrgencyStatus.Open)).toBe('blue');
    });

    it('should return orange for in progress status', () => {
      expect(getUrgencyStatusColor(UrgencyStatus.InProgress)).toBe('orange');
    });

    it('should return green for resolved status', () => {
      expect(getUrgencyStatusColor(UrgencyStatus.Resolved)).toBe('green');
    });

    it('should return gray for closed status', () => {
      expect(getUrgencyStatusColor(UrgencyStatus.Closed)).toBe('gray');
    });

    it('should return gray for unknown status', () => {
      expect(getUrgencyStatusColor('unknown' as any)).toBe('gray');
    });
  });
});
