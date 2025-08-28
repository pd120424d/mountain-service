// Activity model extensions and utilities
// CamelCase-only UI models and helpers

export interface Activity {
  id?: number;
  description?: string;
  employeeId?: number;
  urgencyId?: number;
  createdAt?: string;
  updatedAt?: string;
}

export interface ActivityCreatePayload {
  description: string;
  employeeId: number;
  urgencyId: number;
}

export const getActivityIcon = (activity: Activity): string => {
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

export const formatActivityDescription = (activity: Activity): string => {
  return activity.description || 'No description available';
};

export const getActivityDisplayTime = (activity: Activity): string => {
  if (!activity.createdAt) return '';

  const date = new Date(activity.createdAt);
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
