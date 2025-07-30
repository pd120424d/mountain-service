import {
  getUrgencyLevelColor,
  getUrgencyStatusColor,
  createUrgencyDisplayName,
  withDisplayName
} from './urgency-extensions';
import {
  UrgencyLevel,
  UrgencyStatus,
  UrgencyResponse
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

  describe('createUrgencyDisplayName', () => {
    it('should create display name from first and last name', () => {
      const urgency: UrgencyResponse = {
        id: 1,
        firstName: 'John',
        lastName: 'Doe',
        level: UrgencyLevel.High,
        status: UrgencyStatus.Open,
        description: 'Test urgency'
      };
      expect(createUrgencyDisplayName(urgency)).toBe('John Doe');
    });

    it('should handle only first name', () => {
      const urgency: UrgencyResponse = {
        id: 1,
        firstName: 'John',
        lastName: '',
        level: UrgencyLevel.High,
        status: UrgencyStatus.Open,
        description: 'Test urgency'
      };
      expect(createUrgencyDisplayName(urgency)).toBe('John');
    });

    it('should handle only last name', () => {
      const urgency: UrgencyResponse = {
        id: 1,
        firstName: '',
        lastName: 'Doe',
        level: UrgencyLevel.High,
        status: UrgencyStatus.Open,
        description: 'Test urgency'
      };
      expect(createUrgencyDisplayName(urgency)).toBe('Doe');
    });

    it('should handle whitespace trimming', () => {
      const urgency: UrgencyResponse = {
        id: 1,
        firstName: '  John  ',
        lastName: '  Doe  ',
        level: UrgencyLevel.High,
        status: UrgencyStatus.Open,
        description: 'Test urgency'
      };
      expect(createUrgencyDisplayName(urgency)).toBe('John Doe');
    });

    it('should return Unknown for empty names', () => {
      const urgency: UrgencyResponse = {
        id: 1,
        firstName: '',
        lastName: '',
        level: UrgencyLevel.High,
        status: UrgencyStatus.Open,
        description: 'Test urgency'
      };
      expect(createUrgencyDisplayName(urgency)).toBe('Unknown');
    });

    it('should return Unknown for undefined names', () => {
      const urgency: UrgencyResponse = {
        id: 1,
        level: UrgencyLevel.High,
        status: UrgencyStatus.Open,
        description: 'Test urgency'
      };
      expect(createUrgencyDisplayName(urgency)).toBe('Unknown');
    });
  });

  describe('withDisplayName', () => {
    it('should add display name to urgency object', () => {
      const urgency: UrgencyResponse = {
        id: 1,
        firstName: 'John',
        lastName: 'Doe',
        level: UrgencyLevel.High,
        status: UrgencyStatus.Open,
        description: 'Test urgency'
      };
      const result = withDisplayName(urgency);
      expect(result.displayName).toBe('John Doe');
      expect(result.id).toBe(1);
      expect(result.firstName).toBe('John');
      expect(result.lastName).toBe('Doe');
    });
  });
});
