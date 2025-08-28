import {
  getUrgencyLevelColor,
  getUrgencyStatusColor,
  createUrgencyDisplayName,
  withDisplayName,
  parseLocationString,
  formatLocationForApi,
  LocationCoordinates,
  EnhancedLocation,
  hasAcceptedAssignment
} from './urgency-extensions';
import {
  UrgencyLevel,
  UrgencyStatus,
  UrgencyResponse
} from '../generated/urgency';

  describe('hasAcceptedAssignment', () => {
    it('returns false when assignedEmployeeId is undefined', () => {
      const u: UrgencyResponse = { id: 1, status: UrgencyStatus.Open } as any;
      expect(hasAcceptedAssignment(u as any)).toBe(false);
    });
    it('returns false when assignedEmployeeId is 0', () => {
      const u: UrgencyResponse = { id: 2, assignedEmployeeId: 0 } as any;
      expect(hasAcceptedAssignment(u as any)).toBe(false);
    });
    it('returns true when assignedEmployeeId > 0', () => {
      const u: UrgencyResponse = { id: 3, assignedEmployeeId: 7 } as any;
      expect(hasAcceptedAssignment(u as any)).toBe(true);
    });
  });

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

  describe('parseLocationString', () => {
    it('should parse legacy coordinate format (lat,lng|text)', () => {
      const locationString = '44.0165,21.0059|Belgrade, Serbia';
      const result = parseLocationString(locationString);

      expect(result).toBeTruthy();
      expect(result!.text).toBe('Belgrade, Serbia');
      expect(result!.coordinates?.latitude).toBe(44.0165);
      expect(result!.coordinates?.longitude).toBe(21.0059);
      expect(result!.source).toBe('map');
    });

    it('should parse backend coordinate format (N 43.401123 E 22.662756)', () => {
      const locationString = 'N 44.0165 E 21.0059';
      const result = parseLocationString(locationString);

      expect(result).toBeTruthy();
      expect(result!.text).toBe('44.016500, 21.005900');
      expect(result!.coordinates?.latitude).toBe(44.0165);
      expect(result!.coordinates?.longitude).toBe(21.0059);
      expect(result!.source).toBe('map');
    });

    it('should parse negative coordinates with S and W directions', () => {
      const locationString = 'S 44.0165 W 21.0059';
      const result = parseLocationString(locationString);

      expect(result).toBeTruthy();
      expect(result!.coordinates?.latitude).toBe(-44.0165);
      expect(result!.coordinates?.longitude).toBe(-21.0059);
    });

    it('should return text-only location for plain text', () => {
      const locationString = 'Belgrade, Serbia';
      const result = parseLocationString(locationString);

      expect(result).toBeTruthy();
      expect(result!.text).toBe('Belgrade, Serbia');
      expect(result!.coordinates).toBeUndefined();
      expect(result!.source).toBe('manual');
    });

    it('should return null for empty string', () => {
      const result = parseLocationString('');
      expect(result).toBeNull();
    });
  });

  describe('formatLocationForApi', () => {
    it('should format coordinates in backend format (N lat E lng)', () => {
      const location: EnhancedLocation = {
        text: 'Belgrade, Serbia',
        coordinates: {
          latitude: 44.0165,
          longitude: 21.0059
        },
        source: 'map'
      };

      const result = formatLocationForApi(location);
      expect(result).toBe('N 44.0165 E 21.0059');
    });

    it('should format negative coordinates with S and W directions', () => {
      const location: EnhancedLocation = {
        text: 'Southern Location',
        coordinates: {
          latitude: -44.0165,
          longitude: -21.0059
        },
        source: 'map'
      };

      const result = formatLocationForApi(location);
      expect(result).toBe('S 44.0165 W 21.0059');
    });

    it('should return text for location without coordinates', () => {
      const location: EnhancedLocation = {
        text: 'Belgrade, Serbia',
        source: 'manual'
      };

      const result = formatLocationForApi(location);
      expect(result).toBe('Belgrade, Serbia');
    });
  });
});
