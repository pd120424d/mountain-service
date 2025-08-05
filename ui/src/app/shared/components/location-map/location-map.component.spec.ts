import { ComponentFixture, TestBed } from '@angular/core/testing';
import { FormsModule } from '@angular/forms';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { LocationMapComponent } from './location-map.component';
import { LocationCoordinates, EnhancedLocation } from '../../models';

// Mock Leaflet
const mockMap = {
  setView: jasmine.createSpy('setView').and.returnValue({}),
  on: jasmine.createSpy('on'),
  remove: jasmine.createSpy('remove'),
  removeLayer: jasmine.createSpy('removeLayer')
};

const mockMarker: any = {
  addTo: jasmine.createSpy('addTo'),
  on: jasmine.createSpy('on'),
  getLatLng: jasmine.createSpy('getLatLng').and.returnValue({ lat: 44.0165, lng: 21.0059 }),
  setLatLng: jasmine.createSpy('setLatLng'),
  remove: jasmine.createSpy('remove')
};

// Set up return values after declaration to avoid circular reference
mockMarker.addTo.and.returnValue(mockMarker);
mockMarker.on.and.returnValue(mockMarker);
mockMarker.setLatLng.and.returnValue(mockMarker);
mockMarker.remove.and.returnValue(mockMarker);

const mockTileLayer = {
  addTo: jasmine.createSpy('addTo')
};

const mockLeaflet = {
  map: jasmine.createSpy('map').and.returnValue(mockMap),
  tileLayer: jasmine.createSpy('tileLayer').and.returnValue(mockTileLayer),
  marker: jasmine.createSpy('marker').and.returnValue(mockMarker)
};

// Mock global L
(window as any).L = mockLeaflet;

describe('LocationMapComponent', () => {
  let component: LocationMapComponent;
  let fixture: ComponentFixture<LocationMapComponent>;
  let translateService: TranslateService;

  const mockCoordinates: LocationCoordinates = {
    latitude: 44.0165,
    longitude: 21.0059
  };

  const mockLocation: EnhancedLocation = {
    text: 'Belgrade, Serbia',
    coordinates: mockCoordinates,
    source: 'map'
  };

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [
        LocationMapComponent,
        FormsModule,
        TranslateModule.forRoot()
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(LocationMapComponent);
    component = fixture.componentInstance;
    translateService = TestBed.inject(TranslateService);

    // Mock the component's map and marker properties to prevent Leaflet initialization
    (component as any).map = mockMap;
    (component as any).marker = mockMarker;
    (component as any).leafletLoaded = true;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should initialize with default values', () => {
    expect(component.height).toBe('400px');
    expect(component.width).toBe('100%');
    expect(component.enableLocationSearch).toBe(true);
    expect(component.enableCurrentLocation).toBe(true);
    expect(component.readonly).toBe(false);
  });

  it('should accept initial location input', () => {
    component.initialLocation = mockLocation;
    component.ngOnInit();
    
    expect(component.currentLocation).toEqual(mockLocation);
  });

  it('should emit location when coordinates change', () => {
    spyOn(component.locationSelected, 'emit');
    spyOn(component.coordinatesChanged, 'emit');
    spyOn(component as any, 'addMarker'); // Prevent actual marker creation

    // Simulate private method call (would normally be called by map click)
    (component as any).onCoordinatesChanged(mockCoordinates);

    expect(component.coordinatesChanged.emit).toHaveBeenCalledWith(mockCoordinates);
    expect(component.locationSelected.emit).toHaveBeenCalled();
  });

  it('should set location programmatically', () => {
    // Mock the map being initialized
    (component as any).map = mockLeaflet.map();
    
    component.setLocation(mockLocation);
    
    expect(component.currentLocation).toEqual(mockLocation);
  });

  it('should clear location', () => {
    component.currentLocation = mockLocation;
    spyOn(component.locationSelected, 'emit');
    
    component.clearLocation();
    
    expect(component.currentLocation).toBeUndefined();
    expect(component.searchText).toBe('');
    expect(component.locationSelected.emit).toHaveBeenCalled();
  });

  it('should handle readonly mode', () => {
    component.readonly = true;
    
    // Should not allow clearing in readonly mode
    component.currentLocation = mockLocation;
    component.clearLocation();
    expect(component.currentLocation).toEqual(mockLocation);
    
    // Should not allow getting current location in readonly mode
    component.getCurrentLocation();
    expect(component.isLoadingLocation).toBe(false);
  });

  it('should handle search text input', () => {
    const searchText = 'Belgrade';
    component.searchText = searchText;
    
    expect(component.searchText).toBe(searchText);
  });

  it('should validate coordinates in enhanced location', () => {
    spyOn(component as any, 'addMarker'); // Prevent actual marker creation

    const validLocation: EnhancedLocation = {
      text: 'Valid Location',
      coordinates: { latitude: 45.0, longitude: 20.0 },
      source: 'map'
    };

    component.setLocation(validLocation);
    expect(component.currentLocation?.coordinates?.latitude).toBe(45.0);
    expect(component.currentLocation?.coordinates?.longitude).toBe(20.0);
  });

  it('should handle location without coordinates', () => {
    const textOnlyLocation: EnhancedLocation = {
      text: 'Text Only Location',
      source: 'manual'
    };
    
    component.setLocation(textOnlyLocation);
    expect(component.currentLocation?.text).toBe('Text Only Location');
    expect(component.currentLocation?.coordinates).toBeUndefined();
  });

  it('should emit events when location is selected', () => {
    spyOn(component.locationSelected, 'emit');
    spyOn(component as any, 'addMarker'); // Prevent actual marker creation

    component.setLocation(mockLocation);

    // Simulate location selection
    (component as any).onCoordinatesChanged(mockCoordinates);

    expect(component.locationSelected.emit).toHaveBeenCalled();
  });

  it('should handle geolocation success', () => {
    spyOn(component as any, 'addMarker'); // Prevent actual marker creation

    const mockPosition = {
      coords: {
        latitude: 44.0165,
        longitude: 21.0059,
        accuracy: 10
      }
    };

    spyOn(navigator.geolocation, 'getCurrentPosition').and.callFake((success: any) => {
      success(mockPosition);
    });

    component.getCurrentLocation();

    expect(component.isLoadingLocation).toBe(false);
  });

  it('should handle geolocation error', () => {
    const mockError = new Error('Geolocation failed');

    spyOn(navigator.geolocation, 'getCurrentPosition').and.callFake((success: any, error: any) => {
      error(mockError);
    });

    spyOn(console, 'error');
    
    component.getCurrentLocation();
    
    expect(component.isLoadingLocation).toBe(false);
    expect(console.error).toHaveBeenCalled();
  });

  it('should handle map initialization error', () => {
    spyOn(console, 'error');

    // Reset the mock to throw an error
    (component as any).leafletLoaded = true;
    (component as any).map = null; // Reset map
    mockLeaflet.map.and.throwError('Map init failed');

    try {
      (component as any).initializeMap();
    } catch (error) {
      // Expected to throw
    }

    expect(console.error).toHaveBeenCalled();
    expect(component.mapError).toBe('Failed to initialize map');

    // Reset the spy to its original behavior for subsequent tests
    mockLeaflet.map.and.returnValue(mockMap);
  });

  it('should cleanup map on destroy', () => {
    const mockMap = { remove: jasmine.createSpy('remove') };
    (component as any).map = mockMap;
    
    component.ngOnDestroy();
    
    expect(mockMap.remove).toHaveBeenCalled();
  });
});
