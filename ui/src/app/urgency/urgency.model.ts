
// Enums matching backend
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

export interface Urgency {
  id: number;
  name: string;
  email: string;
  contactPhone: string;
  location: string;
  description: string;
  level: UrgencyLevel;
  status: Status;
  createdAt: string;
  updatedAt: string;
}

export interface UrgencyCreateRequest {
  name: string;
  email: string;
  contactPhone: string;
  location: string;
  description: string;
  level: UrgencyLevel;
}

export interface UrgencyUpdateRequest {
  name?: string;
  email?: string;
  contactPhone?: string;
  location?: string;
  description?: string;
  level?: UrgencyLevel;
  status?: Status;
}