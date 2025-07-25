import {
  UrgencyLevel,
  Status,
  type Urgency,
  mapGeneratedLevelToLegacy,
  mapLegacyLevelToGenerated,
  mapGeneratedStatusToLegacy,
  mapLegacyStatusToGenerated,
  getUrgencyLevelColor,
  getStatusColor
} from './urgency-extensions';
import {
  GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyLevel as GeneratedUrgencyLevel,
  GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyStatus as GeneratedUrgencyStatus
} from '../generated/urgency';

describe('Urgency Extensions', () => {
  describe('UrgencyLevel enum', () => {
    it('should have correct values', () => {
      expect(UrgencyLevel.LOW).toBe('Low');
      expect(UrgencyLevel.MEDIUM).toBe('Medium');
      expect(UrgencyLevel.HIGH).toBe('High');
      expect(UrgencyLevel.CRITICAL).toBe('Critical');
    });
  });

  describe('Status enum', () => {
    it('should have correct values', () => {
      expect(Status.OPEN).toBe('Open');
      expect(Status.IN_PROGRESS).toBe('In Progress');
      expect(Status.RESOLVED).toBe('Resolved');
      expect(Status.CLOSED).toBe('Closed');
    });
  });

  describe('mapGeneratedLevelToLegacy', () => {
    it('should map Low correctly', () => {
      expect(mapGeneratedLevelToLegacy(GeneratedUrgencyLevel.Low)).toBe(UrgencyLevel.LOW);
    });

    it('should map Medium correctly', () => {
      expect(mapGeneratedLevelToLegacy(GeneratedUrgencyLevel.Medium)).toBe(UrgencyLevel.MEDIUM);
    });

    it('should map High correctly', () => {
      expect(mapGeneratedLevelToLegacy(GeneratedUrgencyLevel.High)).toBe(UrgencyLevel.HIGH);
    });

    it('should map Critical correctly', () => {
      expect(mapGeneratedLevelToLegacy(GeneratedUrgencyLevel.Critical)).toBe(UrgencyLevel.CRITICAL);
    });

    it('should default to LOW for unknown values', () => {
      expect(mapGeneratedLevelToLegacy('Unknown' as any)).toBe(UrgencyLevel.LOW);
    });
  });

  describe('mapLegacyLevelToGenerated', () => {
    it('should map LOW correctly', () => {
      expect(mapLegacyLevelToGenerated(UrgencyLevel.LOW)).toBe(GeneratedUrgencyLevel.Low);
    });

    it('should map MEDIUM correctly', () => {
      expect(mapLegacyLevelToGenerated(UrgencyLevel.MEDIUM)).toBe(GeneratedUrgencyLevel.Medium);
    });

    it('should map HIGH correctly', () => {
      expect(mapLegacyLevelToGenerated(UrgencyLevel.HIGH)).toBe(GeneratedUrgencyLevel.High);
    });

    it('should map CRITICAL correctly', () => {
      expect(mapLegacyLevelToGenerated(UrgencyLevel.CRITICAL)).toBe(GeneratedUrgencyLevel.Critical);
    });

    it('should default to Low for unknown values', () => {
      expect(mapLegacyLevelToGenerated('Unknown' as any)).toBe(GeneratedUrgencyLevel.Low);
    });
  });

  describe('mapGeneratedStatusToLegacy', () => {
    it('should map Open correctly', () => {
      expect(mapGeneratedStatusToLegacy(GeneratedUrgencyStatus.Open)).toBe(Status.OPEN);
    });

    it('should map InProgress correctly', () => {
      expect(mapGeneratedStatusToLegacy(GeneratedUrgencyStatus.InProgress)).toBe(Status.IN_PROGRESS);
    });

    it('should map Resolved correctly', () => {
      expect(mapGeneratedStatusToLegacy(GeneratedUrgencyStatus.Resolved)).toBe(Status.RESOLVED);
    });

    it('should map Closed correctly', () => {
      expect(mapGeneratedStatusToLegacy(GeneratedUrgencyStatus.Closed)).toBe(Status.CLOSED);
    });

    it('should default to OPEN for unknown values', () => {
      expect(mapGeneratedStatusToLegacy('Unknown' as any)).toBe(Status.OPEN);
    });
  });

  describe('mapLegacyStatusToGenerated', () => {
    it('should map OPEN correctly', () => {
      expect(mapLegacyStatusToGenerated(Status.OPEN)).toBe(GeneratedUrgencyStatus.Open);
    });

    it('should map IN_PROGRESS correctly', () => {
      expect(mapLegacyStatusToGenerated(Status.IN_PROGRESS)).toBe(GeneratedUrgencyStatus.InProgress);
    });

    it('should map RESOLVED correctly', () => {
      expect(mapLegacyStatusToGenerated(Status.RESOLVED)).toBe(GeneratedUrgencyStatus.Resolved);
    });

    it('should map CLOSED correctly', () => {
      expect(mapLegacyStatusToGenerated(Status.CLOSED)).toBe(GeneratedUrgencyStatus.Closed);
    });

    it('should default to Open for unknown values', () => {
      expect(mapLegacyStatusToGenerated('Unknown' as any)).toBe(GeneratedUrgencyStatus.Open);
    });
  });

  describe('getUrgencyLevelColor', () => {
    it('should return green for low urgency', () => {
      expect(getUrgencyLevelColor(UrgencyLevel.LOW)).toBe('green');
      expect(getUrgencyLevelColor('low' as any)).toBe('green');
    });

    it('should return yellow for medium urgency', () => {
      expect(getUrgencyLevelColor(UrgencyLevel.MEDIUM)).toBe('yellow');
      expect(getUrgencyLevelColor('medium' as any)).toBe('yellow');
    });

    it('should return orange for high urgency', () => {
      expect(getUrgencyLevelColor(UrgencyLevel.HIGH)).toBe('orange');
      expect(getUrgencyLevelColor('high' as any)).toBe('orange');
    });

    it('should return red for critical urgency', () => {
      expect(getUrgencyLevelColor(UrgencyLevel.CRITICAL)).toBe('red');
      expect(getUrgencyLevelColor('critical' as any)).toBe('red');
    });

    it('should return gray for unknown urgency', () => {
      expect(getUrgencyLevelColor('unknown' as any)).toBe('gray');
    });
  });

  describe('getStatusColor', () => {
    it('should return blue for open status', () => {
      expect(getStatusColor(Status.OPEN)).toBe('blue');
      expect(getStatusColor('open' as any)).toBe('blue');
    });

    it('should return orange for in progress status', () => {
      expect(getStatusColor(Status.IN_PROGRESS)).toBe('orange');
      expect(getStatusColor('in_progress' as any)).toBe('orange');
    });

    it('should return green for resolved status', () => {
      expect(getStatusColor(Status.RESOLVED)).toBe('green');
      expect(getStatusColor('resolved' as any)).toBe('green');
    });

    it('should return gray for closed status', () => {
      expect(getStatusColor(Status.CLOSED)).toBe('gray');
      expect(getStatusColor('closed' as any)).toBe('gray');
    });

    it('should return gray for unknown status', () => {
      expect(getStatusColor('unknown' as any)).toBe('gray');
    });
  });
});
