import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { provideRouter, Router } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { ToastrService } from 'ngx-toastr';
import { of, throwError } from 'rxjs';

import { UrgencyFormComponent } from './urgency-form.component';
import { UrgencyService } from '../urgency.service';
import { UrgencyLevel } from '../../shared/models';

describe('UrgencyFormComponent', () => {
  let component: UrgencyFormComponent;
  let fixture: ComponentFixture<UrgencyFormComponent>;
  let urgencyService: jasmine.SpyObj<UrgencyService>;
  let router: jasmine.SpyObj<Router>;
  let toastrService: jasmine.SpyObj<ToastrService>;

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
});
