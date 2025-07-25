// Urgency model extensions and utilities
// These extend the generated models with frontend-specific functionality

import {
  UrgencyResponse,
  UrgencyLevel,
  UrgencyStatus
} from '../generated/urgency';

// Type aliases for cleaner imports
export type Urgency = UrgencyResponse;

// Utility functions
export const getUrgencyLevelColor = (level: UrgencyLevel): string => {
  switch (level) {
    case UrgencyLevel.Low:
      return 'green';
    case UrgencyLevel.Medium:
      return 'yellow';
    case UrgencyLevel.High:
      return 'orange';
    case UrgencyLevel.Critical:
      return 'red';
    default:
      return 'gray';
  }
};

export const getUrgencyStatusColor = (status: UrgencyStatus): string => {
  switch (status) {
    case UrgencyStatus.Open:
      return 'blue';
    case UrgencyStatus.InProgress:
      return 'orange';
    case UrgencyStatus.Resolved:
      return 'green';
    case UrgencyStatus.Closed:
      return 'gray';
    default:
      return 'gray';
  }
};
