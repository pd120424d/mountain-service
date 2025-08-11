import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { provideRouter, Router } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { ToastrService } from 'ngx-toastr';
import { of, throwError } from 'rxjs';

import { UrgencyFormComponent } from './urgency-form.component';
import { UrgencyService } from '../urgency.service';
import { UrgencyLevel, EnhancedLocation, LocationCoordinates } from '../../shared/models';

describe('UrgencyFormComponent', () => {
  let component: UrgencyFormComponent;
  let fixture: ComponentFixture<UrgencyFormComponent>;
  let urgencyService: jasmine.SpyObj<UrgencyService>;
  let router: jasmine.SpyObj<Router>;
  let toastrService: jasmine.SpyObj<ToastrService>;

  // Mock data for map functionality
  let mockCoordinates: LocationCoordinates;
  let mockLocation: EnhancedLocation;

  beforeEach(async () => {
    const urgencyServiceSpy = jasmine.createSpyObj('UrgencyService', ['addUrgency']);
    const routerSpy = jasmine.createSpyObj('Router', ['navigate']);
    const toastrSpy = jasmine.createSpyObj('ToastrService', ['success', 'error', 'info', 'warning']);

    await TestBed.configureTestingModule({
      imports: [
        UrgencyFormComponent,
        TranslateModule.forRoot()
      ],
      providers: [
        provideHttpClient(),
        provideHttpClientTesting(),
        provideRouter([]),
        TranslateService,
        { provide: UrgencyService, useValue: urgencyServiceSpy },
        { provide: Router, useValue: routerSpy },
        { provide: ToastrService, useValue: toastrSpy }
      ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(UrgencyFormComponent);
    component = fixture.componentInstance;
    urgencyService = TestBed.inject(UrgencyService) as jasmine.SpyObj<UrgencyService>;
    router = TestBed.inject(Router) as jasmine.SpyObj<Router>;
    toastrService = TestBed.inject(ToastrService) as jasmine.SpyObj<ToastrService>;

    // Initialize mock data
    mockCoordinates = {
      latitude: 44.0165,
      longitude: 21.0059
    };

    mockLocation = {
      text: 'Belgrade, Serbia',
      coordinates: mockCoordinates,
      source: 'map'
    };

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should initialize form with default values', () => {
    expect(component.urgencyForm).toBeDefined();
    expect(component.urgencyForm.get('firstName')?.value).toBe('');
    expect(component.urgencyForm.get('lastName')?.value).toBe('');
    expect(component.urgencyForm.get('email')?.value).toBe('');
    expect(component.urgencyForm.get('contactPhone')?.value).toBe('');
    expect(component.urgencyForm.get('location')?.value).toBe('');
    expect(component.urgencyForm.get('description')?.value).toBe('');
    expect(component.urgencyForm.get('level')?.value).toBe('medium');
  });

  it('should validate required fields', () => {
    const form = component.urgencyForm;

    // Form should be invalid when empty
    expect(form.valid).toBeFalsy();

    // Check individual field validations
    expect(form.get('firstName')?.hasError('required')).toBeTruthy();
    expect(form.get('lastName')?.hasError('required')).toBeTruthy();
    expect(form.get('email')?.hasError('required')).toBeTruthy();
    expect(form.get('contactPhone')?.hasError('required')).toBeTruthy();
    expect(form.get('location')?.hasError('required')).toBeTruthy();
    expect(form.get('description')?.hasError('required')).toBeTruthy();
  });

  it('should validate email format', () => {
    const emailControl = component.urgencyForm.get('email');

    emailControl?.setValue('invalid-email');
    expect(emailControl?.hasError('email')).toBeTruthy();

    emailControl?.setValue('valid@email.com');
    expect(emailControl?.hasError('email')).toBeFalsy();
  });

  describe('onSubmit', () => {
    it('should submit form when valid', () => {
      // Fill form with valid data
      component.urgencyForm.patchValue({
        firstName: 'John',
        lastName: 'Doe',
        email: 'test@example.com',
        contactPhone: '1234567890',
        location: 'Test Location',
        description: 'This is a test emergency description',
        level: UrgencyLevel.High
      });

      const mockResponse = { id: 1, firstName: 'John', lastName: 'Doe' };
      urgencyService.addUrgency.and.returnValue(of(mockResponse));

      component.onSubmit();

      expect(urgencyService.addUrgency).toHaveBeenCalledWith(component.urgencyForm.value);
      expect(toastrService.success).toHaveBeenCalled();
      expect(router.navigate).toHaveBeenCalledWith(['/']);
    });

    it('should handle submission error', () => {
      // Fill form with valid data
      component.urgencyForm.patchValue({
        firstName: 'John',
        lastName: 'Doe',
        email: 'test@example.com',
        contactPhone: '1234567890',
        location: 'Test Location',
        description: 'This is a test emergency description',
        level: UrgencyLevel.High
      });

      urgencyService.addUrgency.and.returnValue(throwError(() => new Error('Server error')));

      component.onSubmit();

      expect(urgencyService.addUrgency).toHaveBeenCalled();
      expect(component.isSubmitting).toBe(false);
      expect(toastrService.success).not.toHaveBeenCalled();
      expect(router.navigate).not.toHaveBeenCalled();
    });

    it('should not submit when form is invalid', () => {
      component.onSubmit();

      expect(urgencyService.addUrgency).not.toHaveBeenCalled();
      expect(toastrService.success).not.toHaveBeenCalled();
      expect(router.navigate).not.toHaveBeenCalled();
    });

    it('should not submit when already submitting', () => {
      component.urgencyForm.patchValue({
        firstName: 'John',
        lastName: 'Doe',
        email: 'test@example.com',
        contactPhone: '1234567890',
        location: 'Test Location',
        description: 'This is a test emergency description',
        level: UrgencyLevel.High
      });

      component.isSubmitting = true;
      component.onSubmit();

      expect(urgencyService.addUrgency).not.toHaveBeenCalled();
    });
  });

  describe('cancel', () => {
    it('should navigate to home page', () => {
      component.cancel();
      expect(router.navigate).toHaveBeenCalledWith(['/']);
    });
  });

  describe('isFieldInvalid', () => {
    it('should return true for invalid touched field', () => {
      const firstNameControl = component.urgencyForm.get('firstName');
      firstNameControl?.markAsTouched();

      expect(component.isFieldInvalid('firstName')).toBe(true);
    });

    it('should return false for valid field', () => {
      const firstNameControl = component.urgencyForm.get('firstName');
      firstNameControl?.setValue('Valid Name');
      firstNameControl?.markAsTouched();

      expect(component.isFieldInvalid('firstName')).toBe(false);
    });

    it('should return false for invalid untouched field', () => {
      expect(component.isFieldInvalid('firstName')).toBe(false);
    });
  });

  describe('getFieldError', () => {
    it('should return required error message', () => {
      const firstNameControl = component.urgencyForm.get('firstName');
      firstNameControl?.markAsTouched();

      expect(component.getFieldError('firstName')).toBe('FIRSTNAME_FORM.REQUIRED');
    });

    it('should return email error message', () => {
      const emailControl = component.urgencyForm.get('email');
      emailControl?.setValue('invalid-email');
      emailControl?.markAsTouched();

      expect(component.getFieldError('email')).toBe('EMAIL_FORM.INVALID_EMAIL');
    });

    it('should return minlength error message', () => {
      const firstNameControl = component.urgencyForm.get('firstName');
      firstNameControl?.setValue('a');
      firstNameControl?.markAsTouched();

      expect(component.getFieldError('firstName')).toBe('FIRSTNAME_FORM.MIN_LENGTH');
    });

    it('should return empty string for valid field', () => {
      const firstNameControl = component.urgencyForm.get('firstName');
      firstNameControl?.setValue('Valid Name');

      expect(component.getFieldError('firstName')).toBe('');
    });
  });

  describe('Map functionality', () => {

    describe('toggleMap', () => {
      it('should toggle map visibility', () => {
        expect(component.showMap).toBe(false);

        component.toggleMap();
        expect(component.showMap).toBe(true);

        component.toggleMap();
        expect(component.showMap).toBe(false);
      });

      it('should parse existing location text when opening map', () => {
        component.urgencyForm.patchValue({ location: 'Test Location' });

        component.toggleMap();

        expect(component.currentLocation?.text).toBe('Test Location');
        expect(component.currentLocation?.source).toBe('manual');
      });
    });

    describe('onLocationSelected', () => {
      it('should update current location and form control', () => {
        component.onLocationSelected(mockLocation);

        expect(component.currentLocation).toEqual(mockLocation);
        expect(component.urgencyForm.get('location')?.value).toBe('Belgrade, Serbia');
      });

      it('should mark location field as touched', () => {
        component.onLocationSelected(mockLocation);

        const locationControl = component.urgencyForm.get('location');
        expect(locationControl?.touched).toBe(true);
      });
    });

    describe('onCoordinatesChanged', () => {
      it('should update coordinates in current location', () => {
        component.currentLocation = { text: 'Test', source: 'manual' };

        component.onCoordinatesChanged(mockCoordinates);

        expect(component.currentLocation.coordinates).toEqual(mockCoordinates);
        expect(component.currentLocation.source).toBe('map');
      });
    });

    describe('clearMapLocation', () => {
      it('should clear current location and form control', () => {
        component.currentLocation = mockLocation;
        component.urgencyForm.patchValue({ location: 'Test Location' });

        component.clearMapLocation();

        expect(component.currentLocation).toBeUndefined();
        expect(component.urgencyForm.get('location')?.value).toBe('');
      });
    });

    describe('hasCoordinates', () => {
      it('should return true when location has coordinates', () => {
        component.currentLocation = mockLocation;
        expect(component.hasCoordinates()).toBe(true);
      });

      it('should return false when location has no coordinates', () => {
        component.currentLocation = { text: 'Text only', source: 'manual' };
        expect(component.hasCoordinates()).toBe(false);
      });

      it('should return false when no current location', () => {
        component.currentLocation = undefined;
        expect(component.hasCoordinates()).toBe(false);
      });
    });
  });

  describe('Form submission with coordinates', () => {
    it('should format location with coordinates for API', () => {
      // Mock the service response
      urgencyService.addUrgency.and.returnValue(of({ id: 1, message: 'Success' }));

      // Set up valid form
      component.urgencyForm.patchValue({
        firstName: 'John',
        lastName: 'Doe',
        email: 'test@example.com',
        contactPhone: '1234567890',
        location: 'Belgrade, Serbia',
        description: 'This is a test emergency description',
        level: UrgencyLevel.High
      });

      // Set location with coordinates
      component.currentLocation = mockLocation;

      component.onSubmit();

      expect(urgencyService.addUrgency).toHaveBeenCalledWith(
        jasmine.objectContaining({
          location: 'N 44.0165 E 21.0059'
        })
      );
    });

    it('should use text location when no coordinates available', () => {
      // Mock the service response
      urgencyService.addUrgency.and.returnValue(of({ id: 1, message: 'Success' }));

      // Set up valid form
      component.urgencyForm.patchValue({
        firstName: 'John',
        lastName: 'Doe',
        email: 'test@example.com',
        contactPhone: '1234567890',
        location: 'Text only location',
        description: 'This is a test emergency description',
        level: UrgencyLevel.High
      });

      // No coordinates set
      component.currentLocation = undefined;

      component.onSubmit();

      expect(urgencyService.addUrgency).toHaveBeenCalledWith(
        jasmine.objectContaining({
          location: 'Text only location'
        })
      );
    });
  });
});
