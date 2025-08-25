// Activity model extensions and utilities
// These extend the generated models with frontend-specific functionality

import {
  ActivityResponse
} from '../generated/activity';

// Type aliases for cleaner imports
export type Activity = ActivityResponse;

export const getActivityIcon = (activity: ActivityResponse): string => {
  const description = activity.description?.toLowerCase() || '';

  if (description.includes('employee') || description.includes('user')) {
    return 'person';
  } else if (description.includes('urgency') || description.includes('emergency')) {
    return 'warning';
  } else if (description.includes('shift') || description.includes('schedule')) {
    return 'schedule';
  } else if (description.includes('notification') || description.includes('message')) {
    return 'notifications';
  } else if (description.includes('login') || description.includes('auth')) {
    return 'login';
  } else {
    return 'info';
  }
};

export const formatActivityDescription = (activity: ActivityResponse): string => {
  return activity.description || 'No description available';
};

export const getActivityDisplayTime = (activity: ActivityResponse): string => {
  if (!activity.created_at) return '';

  const date = new Date(activity.created_at);
  const now = new Date();
  const diffMs = now.getTime() - date.getTime();
  const diffMins = Math.floor(diffMs / (1000 * 60));
  const diffHours = Math.floor(diffMs / (1000 * 60 * 60));
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));

  if (diffMins < 1) {
    return 'Just now';
  } else if (diffMins < 60) {
    return `${diffMins} minute${diffMins === 1 ? '' : 's'} ago`;
  } else if (diffHours < 24) {
    return `${diffHours} hour${diffHours === 1 ? '' : 's'} ago`;
  } else if (diffDays < 7) {
    return `${diffDays} day${diffDays === 1 ? '' : 's'} ago`;
  } else {
    return date.toLocaleDateString();
  }
};
