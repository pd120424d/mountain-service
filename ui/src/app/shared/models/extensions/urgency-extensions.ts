// Urgency model extensions and utilities
// These extend the generated models with frontend-specific functionality

import {
  UrgencyResponse,
  UrgencyLevel,
  UrgencyStatus
} from '../generated/urgency';

// Type aliases for cleaner imports
export type Urgency = UrgencyResponse;

// Extended urgency interface with computed properties
export interface UrgencyWithDisplayName extends UrgencyResponse {
  displayName: string;
}

// Utility functions

// Create display name from firstName and lastName
export const createUrgencyDisplayName = (urgency: UrgencyResponse): string => {
  if (!urgency.firstName && !urgency.lastName) {
    return 'Unknown';
  }

  const firstName = urgency.firstName?.trim() || '';
  const lastName = urgency.lastName?.trim() || '';

  if (firstName && lastName) {
    return `${firstName} ${lastName}`;
  } else if (firstName) {
    return firstName;
  } else if (lastName) {
    return lastName;
  }

  return 'Unknown';
};

// Create urgency with display name
export const withDisplayName = (urgency: UrgencyResponse): UrgencyWithDisplayName => ({
  ...urgency,
  displayName: createUrgencyDisplayName(urgency)
});

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
