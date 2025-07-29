import {
  getActivityLevelColor,
  getActivityTypeIcon,
  getActivityTypeDisplayName,
  isSystemActivity,
  isEmployeeActivity,
  isUrgencyActivity,
  isShiftActivity,
  isNotificationActivity
} from './activity-extensions';
import { ActivityLevel, ActivityType } from '../generated/activity';

describe('Activity Extensions', () => {
  describe('getActivityLevelColor', () => {
    it('should return correct colors for each level', () => {
      expect(getActivityLevelColor(ActivityLevel.Info)).toBe('blue');
      expect(getActivityLevelColor(ActivityLevel.Warning)).toBe('yellow');
      expect(getActivityLevelColor(ActivityLevel.Error)).toBe('orange');
      expect(getActivityLevelColor(ActivityLevel.Critical)).toBe('red');
    });

    it('should return gray for unknown level', () => {
      expect(getActivityLevelColor('unknown' as ActivityLevel)).toBe('gray');
    });
  });

  describe('getActivityTypeIcon', () => {
    it('should return correct icons for employee activities', () => {
      expect(getActivityTypeIcon(ActivityType.EmployeeCreated)).toBe('person');
      expect(getActivityTypeIcon(ActivityType.EmployeeUpdated)).toBe('person');
      expect(getActivityTypeIcon(ActivityType.EmployeeDeleted)).toBe('person');
      expect(getActivityTypeIcon(ActivityType.EmployeeLogin)).toBe('login');
    });

    it('should return correct icons for shift activities', () => {
      expect(getActivityTypeIcon(ActivityType.ShiftAssigned)).toBe('schedule');
      expect(getActivityTypeIcon(ActivityType.ShiftRemoved)).toBe('schedule');
    });

    it('should return correct icons for urgency activities', () => {
      expect(getActivityTypeIcon(ActivityType.UrgencyCreated)).toBe('warning');
      expect(getActivityTypeIcon(ActivityType.UrgencyUpdated)).toBe('warning');
      expect(getActivityTypeIcon(ActivityType.UrgencyDeleted)).toBe('warning');
      expect(getActivityTypeIcon(ActivityType.EmergencyAssigned)).toBe('emergency');
      expect(getActivityTypeIcon(ActivityType.EmergencyAccepted)).toBe('emergency');
      expect(getActivityTypeIcon(ActivityType.EmergencyDeclined)).toBe('emergency');
    });

    it('should return correct icons for notification activities', () => {
      expect(getActivityTypeIcon(ActivityType.NotificationSent)).toBe('notifications');
      expect(getActivityTypeIcon(ActivityType.NotificationFailed)).toBe('notifications');
    });

    it('should return correct icon for system activities', () => {
      expect(getActivityTypeIcon(ActivityType.SystemReset)).toBe('refresh');
    });

    it('should return info for unknown activity type', () => {
      expect(getActivityTypeIcon('unknown' as ActivityType)).toBe('info');
    });
  });

  describe('getActivityTypeDisplayName', () => {
    it('should return correct display names for employee activities', () => {
      expect(getActivityTypeDisplayName(ActivityType.EmployeeCreated)).toBe('Employee Created');
      expect(getActivityTypeDisplayName(ActivityType.EmployeeUpdated)).toBe('Employee Updated');
      expect(getActivityTypeDisplayName(ActivityType.EmployeeDeleted)).toBe('Employee Deleted');
      expect(getActivityTypeDisplayName(ActivityType.EmployeeLogin)).toBe('Employee Login');
    });

    it('should return correct display names for shift activities', () => {
      expect(getActivityTypeDisplayName(ActivityType.ShiftAssigned)).toBe('Shift Assigned');
      expect(getActivityTypeDisplayName(ActivityType.ShiftRemoved)).toBe('Shift Removed');
    });

    it('should return correct display names for urgency activities', () => {
      expect(getActivityTypeDisplayName(ActivityType.UrgencyCreated)).toBe('Urgency Created');
      expect(getActivityTypeDisplayName(ActivityType.UrgencyUpdated)).toBe('Urgency Updated');
      expect(getActivityTypeDisplayName(ActivityType.UrgencyDeleted)).toBe('Urgency Deleted');
      expect(getActivityTypeDisplayName(ActivityType.EmergencyAssigned)).toBe('Emergency Assigned');
      expect(getActivityTypeDisplayName(ActivityType.EmergencyAccepted)).toBe('Emergency Accepted');
      expect(getActivityTypeDisplayName(ActivityType.EmergencyDeclined)).toBe('Emergency Declined');
    });

    it('should return correct display names for notification activities', () => {
      expect(getActivityTypeDisplayName(ActivityType.NotificationSent)).toBe('Notification Sent');
      expect(getActivityTypeDisplayName(ActivityType.NotificationFailed)).toBe('Notification Failed');
    });

    it('should return correct display name for system activities', () => {
      expect(getActivityTypeDisplayName(ActivityType.SystemReset)).toBe('System Reset');
    });

    it('should format unknown activity types', () => {
      const result = getActivityTypeDisplayName('test_activity_type' as ActivityType);
      expect(result).toBe('Test Activity Type');
    });
  });

  describe('Activity type checkers', () => {
    it('should correctly identify system activities', () => {
      expect(isSystemActivity(ActivityType.SystemReset)).toBe(true);
      expect(isSystemActivity(ActivityType.EmployeeCreated)).toBe(false);
    });

    it('should correctly identify employee activities', () => {
      expect(isEmployeeActivity(ActivityType.EmployeeCreated)).toBe(true);
      expect(isEmployeeActivity(ActivityType.EmployeeUpdated)).toBe(true);
      expect(isEmployeeActivity(ActivityType.EmployeeDeleted)).toBe(true);
      expect(isEmployeeActivity(ActivityType.EmployeeLogin)).toBe(true);
      expect(isEmployeeActivity(ActivityType.ShiftAssigned)).toBe(false);
    });

    it('should correctly identify urgency activities', () => {
      expect(isUrgencyActivity(ActivityType.UrgencyCreated)).toBe(true);
      expect(isUrgencyActivity(ActivityType.UrgencyUpdated)).toBe(true);
      expect(isUrgencyActivity(ActivityType.UrgencyDeleted)).toBe(true);
      expect(isUrgencyActivity(ActivityType.EmergencyAssigned)).toBe(true);
      expect(isUrgencyActivity(ActivityType.EmergencyAccepted)).toBe(true);
      expect(isUrgencyActivity(ActivityType.EmergencyDeclined)).toBe(true);
      expect(isUrgencyActivity(ActivityType.EmployeeCreated)).toBe(false);
    });

    it('should correctly identify shift activities', () => {
      expect(isShiftActivity(ActivityType.ShiftAssigned)).toBe(true);
      expect(isShiftActivity(ActivityType.ShiftRemoved)).toBe(true);
      expect(isShiftActivity(ActivityType.EmployeeCreated)).toBe(false);
    });

    it('should correctly identify notification activities', () => {
      expect(isNotificationActivity(ActivityType.NotificationSent)).toBe(true);
      expect(isNotificationActivity(ActivityType.NotificationFailed)).toBe(true);
      expect(isNotificationActivity(ActivityType.EmployeeCreated)).toBe(false);
    });
  });
});
