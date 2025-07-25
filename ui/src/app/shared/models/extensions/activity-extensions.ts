// Activity model extensions and utilities
// These extend the generated models with frontend-specific functionality

import {
  GithubComPd120424DMountainServiceApiContractsActivityV1ActivityResponse as ActivityResponse,
  GithubComPd120424DMountainServiceApiContractsActivityV1ActivityType as ActivityType,
  GithubComPd120424DMountainServiceApiContractsActivityV1ActivityLevel as ActivityLevel
} from '../generated/activity';

// Type aliases for cleaner imports
export type Activity = ActivityResponse;

// Utility functions
export const getActivityLevelColor = (level: ActivityLevel): string => {
  switch (level) {
    case ActivityLevel.ActivityLevelInfo:
      return 'blue';
    case ActivityLevel.ActivityLevelWarning:
      return 'yellow';
    case ActivityLevel.ActivityLevelError:
      return 'orange';
    case ActivityLevel.ActivityLevelCritical:
      return 'red';
    default:
      return 'gray';
  }
};

export const getActivityTypeIcon = (type: ActivityType): string => {
  switch (type) {
    case ActivityType.ActivityEmployeeCreated:
    case ActivityType.ActivityEmployeeUpdated:
    case ActivityType.ActivityEmployeeDeleted:
      return 'person';
    case ActivityType.ActivityEmployeeLogin:
      return 'login';
    case ActivityType.ActivityShiftAssigned:
    case ActivityType.ActivityShiftRemoved:
      return 'schedule';
    case ActivityType.ActivityUrgencyCreated:
    case ActivityType.ActivityUrgencyUpdated:
    case ActivityType.ActivityUrgencyDeleted:
      return 'warning';
    case ActivityType.ActivityEmergencyAssigned:
    case ActivityType.ActivityEmergencyAccepted:
    case ActivityType.ActivityEmergencyDeclined:
      return 'emergency';
    case ActivityType.ActivityNotificationSent:
    case ActivityType.ActivityNotificationFailed:
      return 'notifications';
    case ActivityType.ActivitySystemReset:
      return 'refresh';
    default:
      return 'info';
  }
};

export const getActivityTypeDisplayName = (type: ActivityType): string => {
  switch (type) {
    case ActivityType.ActivityEmployeeCreated:
      return 'Employee Created';
    case ActivityType.ActivityEmployeeUpdated:
      return 'Employee Updated';
    case ActivityType.ActivityEmployeeDeleted:
      return 'Employee Deleted';
    case ActivityType.ActivityEmployeeLogin:
      return 'Employee Login';
    case ActivityType.ActivityShiftAssigned:
      return 'Shift Assigned';
    case ActivityType.ActivityShiftRemoved:
      return 'Shift Removed';
    case ActivityType.ActivityUrgencyCreated:
      return 'Urgency Created';
    case ActivityType.ActivityUrgencyUpdated:
      return 'Urgency Updated';
    case ActivityType.ActivityUrgencyDeleted:
      return 'Urgency Deleted';
    case ActivityType.ActivityEmergencyAssigned:
      return 'Emergency Assigned';
    case ActivityType.ActivityEmergencyAccepted:
      return 'Emergency Accepted';
    case ActivityType.ActivityEmergencyDeclined:
      return 'Emergency Declined';
    case ActivityType.ActivityNotificationSent:
      return 'Notification Sent';
    case ActivityType.ActivityNotificationFailed:
      return 'Notification Failed';
    case ActivityType.ActivitySystemReset:
      return 'System Reset';
    default:
      return (type as string).replace(/_/g, ' ').replace(/\b\w/g, (l: string) => l.toUpperCase());
  }
};

export const isSystemActivity = (type: ActivityType): boolean => {
  return type === ActivityType.ActivitySystemReset;
};

export const isEmployeeActivity = (type: ActivityType): boolean => {
  return [
    ActivityType.ActivityEmployeeCreated,
    ActivityType.ActivityEmployeeUpdated,
    ActivityType.ActivityEmployeeDeleted,
    ActivityType.ActivityEmployeeLogin
  ].includes(type);
};

export const isUrgencyActivity = (type: ActivityType): boolean => {
  return [
    ActivityType.ActivityUrgencyCreated,
    ActivityType.ActivityUrgencyUpdated,
    ActivityType.ActivityUrgencyDeleted,
    ActivityType.ActivityEmergencyAssigned,
    ActivityType.ActivityEmergencyAccepted,
    ActivityType.ActivityEmergencyDeclined
  ].includes(type);
};

export const isShiftActivity = (type: ActivityType): boolean => {
  return [
    ActivityType.ActivityShiftAssigned,
    ActivityType.ActivityShiftRemoved
  ].includes(type);
};

export const isNotificationActivity = (type: ActivityType): boolean => {
  return [
    ActivityType.ActivityNotificationSent,
    ActivityType.ActivityNotificationFailed
  ].includes(type);
};
