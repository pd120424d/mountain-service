// Activity model extensions and utilities
// These extend the generated models with frontend-specific functionality

import {
  ActivityResponse,
  ActivityType,
  ActivityLevel
} from '../generated/activity';

// Type aliases for cleaner imports
export type Activity = ActivityResponse;

// Utility functions
export const getActivityLevelColor = (level: ActivityLevel): string => {
  switch (level) {
    case ActivityLevel.Info:
      return 'blue';
    case ActivityLevel.Warning:
      return 'yellow';
    case ActivityLevel.Error:
      return 'orange';
    case ActivityLevel.Critical:
      return 'red';
    default:
      return 'gray';
  }
};

export const getActivityTypeIcon = (type: ActivityType): string => {
  switch (type) {
    case ActivityType.EmployeeCreated:
    case ActivityType.EmployeeUpdated:
    case ActivityType.EmployeeDeleted:
      return 'person';
    case ActivityType.EmployeeLogin:
      return 'login';
    case ActivityType.ShiftAssigned:
    case ActivityType.ShiftRemoved:
      return 'schedule';
    case ActivityType.UrgencyCreated:
    case ActivityType.UrgencyUpdated:
    case ActivityType.UrgencyDeleted:
      return 'warning';
    case ActivityType.EmergencyAssigned:
    case ActivityType.EmergencyAccepted:
    case ActivityType.EmergencyDeclined:
      return 'emergency';
    case ActivityType.NotificationSent:
    case ActivityType.NotificationFailed:
      return 'notifications';
    case ActivityType.SystemReset:
      return 'refresh';
    default:
      return 'info';
  }
};

export const getActivityTypeDisplayName = (type: ActivityType): string => {
  switch (type) {
    case ActivityType.EmployeeCreated:
      return 'Employee Created';
    case ActivityType.EmployeeUpdated:
      return 'Employee Updated';
    case ActivityType.EmployeeDeleted:
      return 'Employee Deleted';
    case ActivityType.EmployeeLogin:
      return 'Employee Login';
    case ActivityType.ShiftAssigned:
      return 'Shift Assigned';
    case ActivityType.ShiftRemoved:
      return 'Shift Removed';
    case ActivityType.UrgencyCreated:
      return 'Urgency Created';
    case ActivityType.UrgencyUpdated:
      return 'Urgency Updated';
    case ActivityType.UrgencyDeleted:
      return 'Urgency Deleted';
    case ActivityType.EmergencyAssigned:
      return 'Emergency Assigned';
    case ActivityType.EmergencyAccepted:
      return 'Emergency Accepted';
    case ActivityType.EmergencyDeclined:
      return 'Emergency Declined';
    case ActivityType.NotificationSent:
      return 'Notification Sent';
    case ActivityType.NotificationFailed:
      return 'Notification Failed';
    case ActivityType.SystemReset:
      return 'System Reset';
    default:
      return (type as string).replace(/_/g, ' ').replace(/\b\w/g, (l: string) => l.toUpperCase());
  }
};

export const isSystemActivity = (type: ActivityType): boolean => {
  return type === ActivityType.SystemReset;
};

export const isEmployeeActivity = (type: ActivityType): boolean => {
  return [
    ActivityType.EmployeeCreated,
    ActivityType.EmployeeUpdated,
    ActivityType.EmployeeDeleted,
    ActivityType.EmployeeLogin
  ].includes(type);
};

export const isUrgencyActivity = (type: ActivityType): boolean => {
  return [
    ActivityType.UrgencyCreated,
    ActivityType.UrgencyUpdated,
    ActivityType.UrgencyDeleted,
    ActivityType.EmergencyAssigned,
    ActivityType.EmergencyAccepted,
    ActivityType.EmergencyDeclined
  ].includes(type);
};

export const isShiftActivity = (type: ActivityType): boolean => {
  return [
    ActivityType.ShiftAssigned,
    ActivityType.ShiftRemoved
  ].includes(type);
};

export const isNotificationActivity = (type: ActivityType): boolean => {
  return [
    ActivityType.NotificationSent,
    ActivityType.NotificationFailed
  ].includes(type);
};
