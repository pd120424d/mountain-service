import {
  type Activity,
  getActivityLevelColor,
  getActivityTypeIcon,
  getActivityTypeDisplayName,
  isSystemActivity,
  isEmployeeActivity,
  isUrgencyActivity,
  isShiftActivity,
  isNotificationActivity
} from './activity-extensions';
import {
  GithubComPd120424DMountainServiceApiContractsActivityV1ActivityType as ActivityType,
  GithubComPd120424DMountainServiceApiContractsActivityV1ActivityLevel as ActivityLevel
} from '../generated/activity';

describe('Activity Extensions', () => {
  describe('getActivityLevelColor', () => {
    it('should return blue for info level', () => {
      expect(getActivityLevelColor(ActivityLevel.ActivityLevelInfo)).toBe('blue');
    });

    it('should return yellow for warning level', () => {
      expect(getActivityLevelColor(ActivityLevel.ActivityLevelWarning)).toBe('yellow');
    });

    it('should return orange for error level', () => {
      expect(getActivityLevelColor(ActivityLevel.ActivityLevelError)).toBe('orange');
    });

    it('should return red for critical level', () => {
      expect(getActivityLevelColor(ActivityLevel.ActivityLevelCritical)).toBe('red');
    });

    it('should return gray for unknown level', () => {
      expect(getActivityLevelColor('Unknown' as any)).toBe('gray');
    });
  });

  describe('getActivityTypeIcon', () => {
    it('should return person icon for employee activities', () => {
      expect(getActivityTypeIcon(ActivityType.ActivityEmployeeCreated)).toBe('person');
      expect(getActivityTypeIcon(ActivityType.ActivityEmployeeUpdated)).toBe('person');
      expect(getActivityTypeIcon(ActivityType.ActivityEmployeeDeleted)).toBe('person');
    });

    it('should return login icon for employee login', () => {
      expect(getActivityTypeIcon(ActivityType.ActivityEmployeeLogin)).toBe('login');
    });

    it('should return schedule icon for shift activities', () => {
      expect(getActivityTypeIcon(ActivityType.ActivityShiftAssigned)).toBe('schedule');
      expect(getActivityTypeIcon(ActivityType.ActivityShiftRemoved)).toBe('schedule');
    });

    it('should return warning icon for urgency activities', () => {
      expect(getActivityTypeIcon(ActivityType.ActivityUrgencyCreated)).toBe('warning');
      expect(getActivityTypeIcon(ActivityType.ActivityUrgencyUpdated)).toBe('warning');
      expect(getActivityTypeIcon(ActivityType.ActivityUrgencyDeleted)).toBe('warning');
    });

    it('should return emergency icon for emergency activities', () => {
      expect(getActivityTypeIcon(ActivityType.ActivityEmergencyAssigned)).toBe('emergency');
      expect(getActivityTypeIcon(ActivityType.ActivityEmergencyAccepted)).toBe('emergency');
      expect(getActivityTypeIcon(ActivityType.ActivityEmergencyDeclined)).toBe('emergency');
    });

    it('should return notifications icon for notification activities', () => {
      expect(getActivityTypeIcon(ActivityType.ActivityNotificationSent)).toBe('notifications');
      expect(getActivityTypeIcon(ActivityType.ActivityNotificationFailed)).toBe('notifications');
    });

    it('should return refresh icon for system reset', () => {
      expect(getActivityTypeIcon(ActivityType.ActivitySystemReset)).toBe('refresh');
    });

    it('should return info icon for unknown activity type', () => {
      expect(getActivityTypeIcon('Unknown' as any)).toBe('info');
    });
  });

  describe('getActivityTypeDisplayName', () => {
    it('should return correct display names for employee activities', () => {
      expect(getActivityTypeDisplayName(ActivityType.ActivityEmployeeCreated)).toBe('Employee Created');
      expect(getActivityTypeDisplayName(ActivityType.ActivityEmployeeUpdated)).toBe('Employee Updated');
      expect(getActivityTypeDisplayName(ActivityType.ActivityEmployeeDeleted)).toBe('Employee Deleted');
      expect(getActivityTypeDisplayName(ActivityType.ActivityEmployeeLogin)).toBe('Employee Login');
    });

    it('should return correct display names for shift activities', () => {
      expect(getActivityTypeDisplayName(ActivityType.ActivityShiftAssigned)).toBe('Shift Assigned');
      expect(getActivityTypeDisplayName(ActivityType.ActivityShiftRemoved)).toBe('Shift Removed');
    });

    it('should return correct display names for urgency activities', () => {
      expect(getActivityTypeDisplayName(ActivityType.ActivityUrgencyCreated)).toBe('Urgency Created');
      expect(getActivityTypeDisplayName(ActivityType.ActivityUrgencyUpdated)).toBe('Urgency Updated');
      expect(getActivityTypeDisplayName(ActivityType.ActivityUrgencyDeleted)).toBe('Urgency Deleted');
    });

    it('should return correct display names for emergency activities', () => {
      expect(getActivityTypeDisplayName(ActivityType.ActivityEmergencyAssigned)).toBe('Emergency Assigned');
      expect(getActivityTypeDisplayName(ActivityType.ActivityEmergencyAccepted)).toBe('Emergency Accepted');
      expect(getActivityTypeDisplayName(ActivityType.ActivityEmergencyDeclined)).toBe('Emergency Declined');
    });

    it('should return correct display names for notification activities', () => {
      expect(getActivityTypeDisplayName(ActivityType.ActivityNotificationSent)).toBe('Notification Sent');
      expect(getActivityTypeDisplayName(ActivityType.ActivityNotificationFailed)).toBe('Notification Failed');
    });

    it('should return correct display name for system reset', () => {
      expect(getActivityTypeDisplayName(ActivityType.ActivitySystemReset)).toBe('System Reset');
    });

    it('should format unknown activity types', () => {
      const result = getActivityTypeDisplayName('activity_test_case' as any);
      expect(result).toBe('Activity Test Case');
    });
  });

  describe('isSystemActivity', () => {
    it('should return true for system reset', () => {
      expect(isSystemActivity(ActivityType.ActivitySystemReset)).toBe(true);
    });

    it('should return false for non-system activities', () => {
      expect(isSystemActivity(ActivityType.ActivityEmployeeCreated)).toBe(false);
      expect(isSystemActivity(ActivityType.ActivityUrgencyCreated)).toBe(false);
      expect(isSystemActivity(ActivityType.ActivityShiftAssigned)).toBe(false);
    });
  });

  describe('isEmployeeActivity', () => {
    it('should return true for employee activities', () => {
      expect(isEmployeeActivity(ActivityType.ActivityEmployeeCreated)).toBe(true);
      expect(isEmployeeActivity(ActivityType.ActivityEmployeeUpdated)).toBe(true);
      expect(isEmployeeActivity(ActivityType.ActivityEmployeeDeleted)).toBe(true);
      expect(isEmployeeActivity(ActivityType.ActivityEmployeeLogin)).toBe(true);
    });

    it('should return false for non-employee activities', () => {
      expect(isEmployeeActivity(ActivityType.ActivitySystemReset)).toBe(false);
      expect(isEmployeeActivity(ActivityType.ActivityUrgencyCreated)).toBe(false);
      expect(isEmployeeActivity(ActivityType.ActivityShiftAssigned)).toBe(false);
    });
  });

  describe('isUrgencyActivity', () => {
    it('should return true for urgency activities', () => {
      expect(isUrgencyActivity(ActivityType.ActivityUrgencyCreated)).toBe(true);
      expect(isUrgencyActivity(ActivityType.ActivityUrgencyUpdated)).toBe(true);
      expect(isUrgencyActivity(ActivityType.ActivityUrgencyDeleted)).toBe(true);
      expect(isUrgencyActivity(ActivityType.ActivityEmergencyAssigned)).toBe(true);
      expect(isUrgencyActivity(ActivityType.ActivityEmergencyAccepted)).toBe(true);
      expect(isUrgencyActivity(ActivityType.ActivityEmergencyDeclined)).toBe(true);
    });

    it('should return false for non-urgency activities', () => {
      expect(isUrgencyActivity(ActivityType.ActivitySystemReset)).toBe(false);
      expect(isUrgencyActivity(ActivityType.ActivityEmployeeCreated)).toBe(false);
      expect(isUrgencyActivity(ActivityType.ActivityShiftAssigned)).toBe(false);
    });
  });

  describe('isShiftActivity', () => {
    it('should return true for shift activities', () => {
      expect(isShiftActivity(ActivityType.ActivityShiftAssigned)).toBe(true);
      expect(isShiftActivity(ActivityType.ActivityShiftRemoved)).toBe(true);
    });

    it('should return false for non-shift activities', () => {
      expect(isShiftActivity(ActivityType.ActivitySystemReset)).toBe(false);
      expect(isShiftActivity(ActivityType.ActivityEmployeeCreated)).toBe(false);
      expect(isShiftActivity(ActivityType.ActivityUrgencyCreated)).toBe(false);
    });
  });

  describe('isNotificationActivity', () => {
    it('should return true for notification activities', () => {
      expect(isNotificationActivity(ActivityType.ActivityNotificationSent)).toBe(true);
      expect(isNotificationActivity(ActivityType.ActivityNotificationFailed)).toBe(true);
    });

    it('should return false for non-notification activities', () => {
      expect(isNotificationActivity(ActivityType.ActivitySystemReset)).toBe(false);
      expect(isNotificationActivity(ActivityType.ActivityEmployeeCreated)).toBe(false);
      expect(isNotificationActivity(ActivityType.ActivityUrgencyCreated)).toBe(false);
    });
  });
});
