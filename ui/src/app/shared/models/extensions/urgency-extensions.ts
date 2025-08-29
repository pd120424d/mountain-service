// Urgency model extensions and utilities
// These extend the generated models with frontend-specific functionality

import {
  V1UrgencyResponse as UrgencyResponse,
  V1UrgencyCreateRequest as UrgencyCreateRequest,
  V1UrgencyUpdateRequest as UrgencyUpdateRequest,
  V1UrgencyLevel as UrgencyLevel,
  V1UrgencyStatus as UrgencyStatus
} from '../generated/urgency';

// Type aliases for cleaner imports
export type Urgency = UrgencyResponse;

export interface LocationCoordinates {
  latitude: number;
  longitude: number;
  accuracy?: number; // GPS accuracy in meters
}

export interface EnhancedLocation {
  text: string; // Original text location (for backward compatibility)
  coordinates?: LocationCoordinates; // Optional coordinates
  source?: 'manual' | 'gps' | 'map'; // How the location was obtained
}

export interface UrgencyCreateRequestWithCoordinates extends Omit<UrgencyCreateRequest, 'location'> {
  location: string; // Keep original for API compatibility
  enhancedLocation?: EnhancedLocation; // Additional location data
}

export interface UrgencyResponseWithCoordinates extends Omit<UrgencyResponse, 'location'> {
  location?: string; // Keep original for API compatibility
  enhancedLocation?: EnhancedLocation; // Additional location data
}

export interface UrgencyWithDisplayName extends UrgencyResponse {
  displayName: string;
}

export const parseLocationString = (locationString: string): EnhancedLocation | null => {
  if (!locationString) return null;

  // Try to parse coordinates from location string (format: "lat,lng|text" or just "text")
  const coordinatePattern = /^(-?\d+\.?\d*),(-?\d+\.?\d*)\|(.*)$/;
  const match = locationString.match(coordinatePattern);

  if (match) {
    const [, lat, lng, text] = match;
    return {
      text: text.trim(),
      coordinates: {
        latitude: parseFloat(lat),
        longitude: parseFloat(lng)
      },
      source: 'map'
    };
  }

  // Try to parse coordinates in backend format: 'N 43.401123 E 22.662756'
  const backendCoordinatePattern = /^([NS])\s+(-?\d+\.?\d*)\s+([EW])\s+(-?\d+\.?\d*)$/;
  const backendMatch = locationString.match(backendCoordinatePattern);

  if (backendMatch) {
    const [, latDirection, latValue, lngDirection, lngValue] = backendMatch;
    const latitude = (latDirection === 'S' ? -1 : 1) * parseFloat(latValue);
    const longitude = (lngDirection === 'W' ? -1 : 1) * parseFloat(lngValue);

    return {
      text: `${latitude.toFixed(6)}, ${longitude.toFixed(6)}`,
      coordinates: {
        latitude,
        longitude
      },
      source: 'map'
    };
  }

  // If no coordinates found, return as text-only location
  return {
    text: locationString,
    source: 'manual'
  };
};

export const formatLocationForApi = (enhancedLocation: EnhancedLocation): string => {
  if (enhancedLocation.coordinates) {
    // Format coordinates in the expected backend format: 'N 43.401123 E 22.662756'
    const lat = enhancedLocation.coordinates.latitude;
    const lng = enhancedLocation.coordinates.longitude;
    const latDirection = lat >= 0 ? 'N' : 'S';
    const lngDirection = lng >= 0 ? 'E' : 'W';
    return `${latDirection} ${Math.abs(lat)} ${lngDirection} ${Math.abs(lng)}`;
  }
  return enhancedLocation.text;
};

export const formatCoordinatesDisplay = (coordinates: LocationCoordinates): string => {
  return `${coordinates.latitude.toFixed(6)}, ${coordinates.longitude.toFixed(6)}`;
};

export const calculateDistance = (coord1: LocationCoordinates, coord2: LocationCoordinates): number => {
  const R = 6371; // Earth's radius in kilometers
  const dLat = (coord2.latitude - coord1.latitude) * Math.PI / 180;
  const dLon = (coord2.longitude - coord1.longitude) * Math.PI / 180;
  const a =
    Math.sin(dLat/2) * Math.sin(dLat/2) +
    Math.cos(coord1.latitude * Math.PI / 180) * Math.cos(coord2.latitude * Math.PI / 180) *
    Math.sin(dLon/2) * Math.sin(dLon/2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a));
  return R * c; // Distance in kilometers
};

export const isValidCoordinates = (coordinates: LocationCoordinates): boolean => {
  return coordinates.latitude >= -90 && coordinates.latitude <= 90 &&
         coordinates.longitude >= -180 && coordinates.longitude <= 180;
};

export const isInMountainRegion = (coordinates: LocationCoordinates): boolean => {
  // For now, this is a placeholder that accepts all coordinates
  return isValidCoordinates(coordinates);
};


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

}
export const hasAcceptedAssignment = (urgency: UrgencyResponse): boolean => {
  const assignedEmployeeId = (urgency as any)?.assignedEmployeeId as number | undefined;
  return typeof assignedEmployeeId === 'number' && assignedEmployeeId > 0;
};
