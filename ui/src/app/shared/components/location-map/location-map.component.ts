import { Component, OnInit, OnDestroy, Input, Output, EventEmitter, ElementRef, ViewChild, AfterViewInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { LocationCoordinates, EnhancedLocation } from '../../models';

// We'll load Leaflet dynamically to avoid build issues
declare var L: any;

@Component({
  selector: 'app-location-map',
  standalone: true,
  imports: [CommonModule, FormsModule, TranslateModule],
  templateUrl: './location-map.component.html',
  styleUrls: ['./location-map.component.css']
})
export class LocationMapComponent implements OnInit, AfterViewInit, OnDestroy {
  @ViewChild('mapContainer', { static: true }) mapContainer!: ElementRef;
  
  @Input() initialLocation?: EnhancedLocation;
  @Input() height: string = '400px';
  @Input() width: string = '100%';
  @Input() enableLocationSearch: boolean = true;
  @Input() enableCurrentLocation: boolean = true;
  @Input() readonly: boolean = false;
  
  @Output() locationSelected = new EventEmitter<EnhancedLocation>();
  @Output() coordinatesChanged = new EventEmitter<LocationCoordinates>();
  
  private map: any;
  private marker: any;
  private leafletLoaded = false;
  
  // Default center for mountain regions (can be customized)
  private defaultCenter: LocationCoordinates = {
    latitude: 44.0165, // Belgrade, Serbia as default
    longitude: 21.0059
  };
  
  searchText: string = '';
  currentLocation?: EnhancedLocation;
  isLoadingLocation = false;
  mapError?: string;

  ngOnInit(): void {
    this.currentLocation = this.initialLocation;
    this.loadLeaflet();
  }

  ngAfterViewInit(): void {
    // Initialize map after Leaflet is loaded
    if (this.leafletLoaded) {
      this.initializeMap();
    }
  }

  ngOnDestroy(): void {
    if (this.map) {
      this.map.remove();
    }
  }

  private async loadLeaflet(): Promise<void> {
    try {
      // Load Leaflet CSS
      if (!document.querySelector('link[href*="leaflet.css"]')) {
        const link = document.createElement('link');
        link.rel = 'stylesheet';
        link.href = 'https://unpkg.com/leaflet@1.9.4/dist/leaflet.css';
        document.head.appendChild(link);
      }

      // Load Leaflet JS
      if (typeof L === 'undefined') {
        await this.loadScript('https://unpkg.com/leaflet@1.9.4/dist/leaflet.js');
      }

      this.leafletLoaded = true;
      
      // Initialize map if view is ready
      if (this.mapContainer) {
        this.initializeMap();
      }
    } catch (error) {
      console.error('Failed to load Leaflet:', error);
      this.mapError = 'Failed to load map library';
    }
  }

  private loadScript(src: string): Promise<void> {
    return new Promise((resolve, reject) => {
      const script = document.createElement('script');
      script.src = src;
      script.onload = () => resolve();
      script.onerror = () => reject(new Error(`Failed to load script: ${src}`));
      document.head.appendChild(script);
    });
  }

  private initializeMap(): void {
    if (!this.leafletLoaded || this.map) {
      return;
    }

    try {
      // Determine initial center
      const center = this.currentLocation?.coordinates || this.defaultCenter;
      
      // Initialize map
      this.map = L.map(this.mapContainer.nativeElement).setView(
        [center.latitude, center.longitude], 
        this.currentLocation?.coordinates ? 15 : 10
      );

      // Add OpenStreetMap tiles
      L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '© OpenStreetMap contributors',
        maxZoom: 19
      }).addTo(this.map);

      // Add marker if location exists
      if (this.currentLocation?.coordinates) {
        this.addMarker(this.currentLocation.coordinates);
      }

      // Add click handler if not readonly
      if (!this.readonly) {
        this.map.on('click', (e: any) => {
          this.onMapClick(e.latlng);
        });
      }

    } catch (error) {
      console.error('Failed to initialize map:', error);
      this.mapError = 'Failed to initialize map';
    }
  }

  private addMarker(coordinates: LocationCoordinates): void {
    if (this.marker) {
      this.map.removeLayer(this.marker);
    }

    this.marker = L.marker([coordinates.latitude, coordinates.longitude], {
      draggable: !this.readonly
    }).addTo(this.map);

    if (!this.readonly) {
      this.marker.on('dragend', (e: any) => {
        const position = e.target.getLatLng();
        this.onCoordinatesChanged({
          latitude: position.lat,
          longitude: position.lng
        });
      });
    }
  }

  private onMapClick(latlng: any): void {
    if (this.readonly) return;

    const coordinates: LocationCoordinates = {
      latitude: latlng.lat,
      longitude: latlng.lng
    };

    this.onCoordinatesChanged(coordinates);
  }

  private onCoordinatesChanged(coordinates: LocationCoordinates): void {
    // Update marker position
    this.addMarker(coordinates);

    // Update current location
    this.currentLocation = {
      text: this.currentLocation?.text || `${coordinates.latitude.toFixed(6)}, ${coordinates.longitude.toFixed(6)}`,
      coordinates,
      source: 'map'
    };

    // Emit events
    this.coordinatesChanged.emit(coordinates);
    this.locationSelected.emit(this.currentLocation);
  }

  // Public methods for external control
  public setLocation(location: EnhancedLocation): void {
    this.currentLocation = location;
    
    if (location.coordinates && this.map) {
      this.map.setView([location.coordinates.latitude, location.coordinates.longitude], 15);
      this.addMarker(location.coordinates);
    }
  }

  public getCurrentLocation(): void {
    if (!this.enableCurrentLocation || this.readonly) return;

    this.isLoadingLocation = true;
    
    if (navigator.geolocation) {
      navigator.geolocation.getCurrentPosition(
        (position) => {
          const coordinates: LocationCoordinates = {
            latitude: position.coords.latitude,
            longitude: position.coords.longitude,
            accuracy: position.coords.accuracy
          };

          this.onCoordinatesChanged(coordinates);
          this.isLoadingLocation = false;
        },
        (error) => {
          console.error('Geolocation error:', error);
          this.isLoadingLocation = false;
        },
        {
          enableHighAccuracy: true,
          timeout: 10000,
          maximumAge: 60000
        }
      );
    } else {
      this.isLoadingLocation = false;
    }
  }

  public searchLocation(): void {
    if (!this.enableLocationSearch || !this.searchText.trim()) return;

    // Simple geocoding using Nominatim (OpenStreetMap's geocoding service)
    // In production, you might want to use a more robust geocoding service
    const query = encodeURIComponent(this.searchText.trim());
    const url = `https://nominatim.openstreetmap.org/search?format=json&q=${query}&limit=1`;

    fetch(url)
      .then(response => response.json())
      .then(data => {
        if (data && data.length > 0) {
          const result = data[0];
          const coordinates: LocationCoordinates = {
            latitude: parseFloat(result.lat),
            longitude: parseFloat(result.lon)
          };

          this.currentLocation = {
            text: result.display_name || this.searchText,
            coordinates,
            source: 'manual'
          };

          if (this.map) {
            this.map.setView([coordinates.latitude, coordinates.longitude], 15);
            this.addMarker(coordinates);
          }

          this.locationSelected.emit(this.currentLocation);
        }
      })
      .catch(error => {
        console.error('Geocoding error:', error);
      });
  }

  public clearLocation(): void {
    if (this.readonly) return;

    this.currentLocation = undefined;
    this.searchText = '';
    
    if (this.marker) {
      this.map.removeLayer(this.marker);
      this.marker = null;
    }

    this.locationSelected.emit(undefined as any);
  }
}
