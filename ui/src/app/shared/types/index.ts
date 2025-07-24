// Common types and utilities
// Add shared types that are not generated from backend contracts

// Example: Frontend-specific types
export type LoadingState = 'idle' | 'loading' | 'success' | 'error';

export interface ApiResponse<T> {
  data: T;
  loading: boolean;
  error?: string;
}
