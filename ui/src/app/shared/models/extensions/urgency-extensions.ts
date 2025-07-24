// Urgency model extensions and utilities
// These extend the generated models with frontend-specific functionality

import { UrgencyResponse, UrgencyLevel as GeneratedUrgencyLevel, UrgencyStatus as GeneratedUrgencyStatus } from '../generated/urgency';

// Legacy enum mappings for backward compatibility
export enum UrgencyLevel {
  LOW = 'Low',
  MEDIUM = 'Medium',
  HIGH = 'High',
  CRITICAL = 'Critical'
}

export enum Status {
  OPEN = 'Open',
  IN_PROGRESS = 'In Progress',
  RESOLVED = 'Resolved',
  CLOSED = 'Closed'
}

// Type aliases for cleaner imports
export type Urgency = UrgencyResponse;

// Utility functions to convert between generated and legacy enums
export const mapGeneratedLevelToLegacy = (level: GeneratedUrgencyLevel): UrgencyLevel => {
  switch (level) {
    case GeneratedUrgencyLevel.Low:
      return UrgencyLevel.LOW;
    case GeneratedUrgencyLevel.Medium:
      return UrgencyLevel.MEDIUM;
    case GeneratedUrgencyLevel.High:
      return UrgencyLevel.HIGH;
    case GeneratedUrgencyLevel.Critical:
      return UrgencyLevel.CRITICAL;
    default:
      return UrgencyLevel.LOW;
  }
};

export const mapLegacyLevelToGenerated = (level: UrgencyLevel): GeneratedUrgencyLevel => {
  switch (level) {
    case UrgencyLevel.LOW:
      return GeneratedUrgencyLevel.Low;
    case UrgencyLevel.MEDIUM:
      return GeneratedUrgencyLevel.Medium;
    case UrgencyLevel.HIGH:
      return GeneratedUrgencyLevel.High;
    case UrgencyLevel.CRITICAL:
      return GeneratedUrgencyLevel.Critical;
    default:
      return GeneratedUrgencyLevel.Low;
  }
};

export const mapGeneratedStatusToLegacy = (status: GeneratedUrgencyStatus): Status => {
  switch (status) {
    case GeneratedUrgencyStatus.Open:
      return Status.OPEN;
    case GeneratedUrgencyStatus.InProgress:
      return Status.IN_PROGRESS;
    case GeneratedUrgencyStatus.Resolved:
      return Status.RESOLVED;
    case GeneratedUrgencyStatus.Closed:
      return Status.CLOSED;
    default:
      return Status.OPEN;
  }
};

export const mapLegacyStatusToGenerated = (status: Status): GeneratedUrgencyStatus => {
  switch (status) {
    case Status.OPEN:
      return GeneratedUrgencyStatus.Open;
    case Status.IN_PROGRESS:
      return GeneratedUrgencyStatus.InProgress;
    case Status.RESOLVED:
      return GeneratedUrgencyStatus.Resolved;
    case Status.CLOSED:
      return GeneratedUrgencyStatus.Closed;
    default:
      return GeneratedUrgencyStatus.Open;
  }
};

// Utility functions
export const getUrgencyLevelColor = (level: UrgencyLevel | GeneratedUrgencyLevel): string => {
  const normalizedLevel = typeof level === 'string' ? level.toLowerCase() : level;
  
  switch (normalizedLevel) {
    case 'low':
    case UrgencyLevel.LOW.toLowerCase():
      return 'green';
    case 'medium':
    case UrgencyLevel.MEDIUM.toLowerCase():
      return 'yellow';
    case 'high':
    case UrgencyLevel.HIGH.toLowerCase():
      return 'orange';
    case 'critical':
    case UrgencyLevel.CRITICAL.toLowerCase():
      return 'red';
    default:
      return 'gray';
  }
};

export const getStatusColor = (status: Status | GeneratedUrgencyStatus): string => {
  const normalizedStatus = typeof status === 'string' ? status.toLowerCase().replace(' ', '_') : status;
  
  switch (normalizedStatus) {
    case 'open':
    case Status.OPEN.toLowerCase():
      return 'blue';
    case 'in_progress':
    case Status.IN_PROGRESS.toLowerCase().replace(' ', '_'):
      return 'orange';
    case 'resolved':
    case Status.RESOLVED.toLowerCase():
      return 'green';
    case 'closed':
    case Status.CLOSED.toLowerCase():
      return 'gray';
    default:
      return 'gray';
  }
};
