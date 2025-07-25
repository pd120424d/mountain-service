// Test for the main models index file to ensure all exports are working correctly
import {
  // Generated types
  Employee,
  EmployeeCreateRequest,
  EmployeeUpdateRequest,
  EmployeeLogin,
  TokenResponse,
  ErrorResponse,
  MessageResponse,
  UrgencyResponse,
  UrgencyCreateRequest,
  UrgencyUpdateRequest,
  GeneratedUrgencyLevel,
  GeneratedUrgencyStatus,
  ActivityResponse,
  ActivityCreateRequest,
  ActivityListResponse,
  ActivityStatsResponse,
  ActivityType,
  ActivityLevel,
  
  // Employee extensions
  MedicRole,
  TechnicalRole,
  AdministratorRole,
  EmployeeRole,
  EmployeeWithDisplayName,
  createDisplayName,
  isAdmin,
  isMedic,
  isTechnical,
  
  // Urgency extensions
  UrgencyLevel,
  Status,
  Urgency,
  mapGeneratedLevelToLegacy,
  mapLegacyLevelToGenerated,
  mapGeneratedStatusToLegacy,
  mapLegacyStatusToGenerated,
  getUrgencyLevelColor,
  getStatusColor,
  
  // Activity extensions
  Activity,
  getActivityLevelColor,
  getActivityTypeIcon,
  getActivityTypeDisplayName,
  isSystemActivity,
  isEmployeeActivity,
  isUrgencyActivity,
  isShiftActivity,
  isNotificationActivity
} from './index';

describe('Models Index', () => {
  describe('Generated Types Export', () => {
    it('should export Employee types', () => {
      expect(typeof MedicRole).toBe('string');
      expect(typeof TechnicalRole).toBe('string');
      expect(typeof AdministratorRole).toBe('string');
    });

    it('should export Urgency enums', () => {
      expect(GeneratedUrgencyLevel).toBeDefined();
      expect(GeneratedUrgencyStatus).toBeDefined();
    });

    it('should export Activity enums', () => {
      expect(ActivityType).toBeDefined();
      expect(ActivityLevel).toBeDefined();
    });
  });

  describe('Employee Extensions Export', () => {
    it('should export role constants', () => {
      expect(MedicRole).toBe('Medic');
      expect(TechnicalRole).toBe('Technical');
      expect(AdministratorRole).toBe('Administrator');
    });

    it('should export utility functions', () => {
      expect(typeof createDisplayName).toBe('function');
      expect(typeof isAdmin).toBe('function');
      expect(typeof isMedic).toBe('function');
      expect(typeof isTechnical).toBe('function');
    });

    it('should work with createDisplayName function', () => {
      const mockEmployee = {
        id: 1,
        firstName: 'John',
        lastName: 'Doe',
        email: 'john.doe@example.com',
        phone: '+1234567890',
        profileType: 'Medic',
        username: 'johndoe',
        gender: 'Male'
      } as Employee;

      const displayName = createDisplayName(mockEmployee);
      expect(displayName).toBe('John Doe');
    });
  });

  describe('Urgency Extensions Export', () => {
    it('should export legacy enums', () => {
      expect(UrgencyLevel.LOW).toBe('Low');
      expect(UrgencyLevel.MEDIUM).toBe('Medium');
      expect(UrgencyLevel.HIGH).toBe('High');
      expect(UrgencyLevel.CRITICAL).toBe('Critical');

      expect(Status.OPEN).toBe('Open');
      expect(Status.IN_PROGRESS).toBe('In Progress');
      expect(Status.RESOLVED).toBe('Resolved');
      expect(Status.CLOSED).toBe('Closed');
    });

    it('should export mapping functions', () => {
      expect(typeof mapGeneratedLevelToLegacy).toBe('function');
      expect(typeof mapLegacyLevelToGenerated).toBe('function');
      expect(typeof mapGeneratedStatusToLegacy).toBe('function');
      expect(typeof mapLegacyStatusToGenerated).toBe('function');
    });

    it('should export utility functions', () => {
      expect(typeof getUrgencyLevelColor).toBe('function');
      expect(typeof getStatusColor).toBe('function');
    });

    it('should work with utility functions', () => {
      expect(getUrgencyLevelColor(UrgencyLevel.LOW)).toBe('green');
      expect(getStatusColor(Status.OPEN)).toBe('blue');
    });
  });

  describe('Activity Extensions Export', () => {
    it('should export utility functions', () => {
      expect(typeof getActivityLevelColor).toBe('function');
      expect(typeof getActivityTypeIcon).toBe('function');
      expect(typeof getActivityTypeDisplayName).toBe('function');
      expect(typeof isSystemActivity).toBe('function');
      expect(typeof isEmployeeActivity).toBe('function');
      expect(typeof isUrgencyActivity).toBe('function');
      expect(typeof isShiftActivity).toBe('function');
      expect(typeof isNotificationActivity).toBe('function');
    });

    it('should work with utility functions', () => {
      expect(getActivityLevelColor(ActivityLevel.ActivityLevelInfo)).toBe('blue');
      expect(getActivityTypeIcon(ActivityType.ActivityEmployeeCreated)).toBe('person');
      expect(getActivityTypeDisplayName(ActivityType.ActivityEmployeeCreated)).toBe('Employee Created');
      expect(isSystemActivity(ActivityType.ActivitySystemReset)).toBe(true);
      expect(isEmployeeActivity(ActivityType.ActivityEmployeeCreated)).toBe(true);
    });
  });

  describe('Type Aliases', () => {
    it('should have correct type aliases', () => {
      // These tests mainly ensure the types compile correctly
      const mockEmployee: Employee = {
        id: 1,
        firstName: 'John',
        lastName: 'Doe',
        email: 'john.doe@example.com',
        phone: '+1234567890',
        profileType: 'Medic',
        username: 'johndoe',
        gender: 'Male'
      };

      const mockUrgency: Urgency = {
        id: 1,
        name: 'Test Urgency',
        description: 'Test Description',
        level: GeneratedUrgencyLevel.High,
        status: GeneratedUrgencyStatus.Open,
        location: 'Test Location',
        contactPhone: 'Test Contact',
        email: 'test@example.com',
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z'
      };

      const mockActivity: Activity = {
        id: 1,
        type: ActivityType.ActivityEmployeeCreated,
        level: ActivityLevel.ActivityLevelInfo,
        title: 'Test Activity',
        description: 'Test Activity Description',
        actorId: 1,
        actorName: 'Test Actor',
        createdAt: '2024-01-01T00:00:00Z',
        updatedAt: '2024-01-01T00:00:00Z'
      };

      expect(mockEmployee.id).toBe(1);
      expect(mockUrgency.id).toBe(1);
      expect(mockActivity.id).toBe(1);
    });
  });
});
