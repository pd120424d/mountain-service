import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { provideRouter } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { ToastrService } from 'ngx-toastr';

import { UrgencyFormComponent } from './urgency-form.component';

describe('UrgencyFormComponent', () => {
  let component: UrgencyFormComponent;
  let fixture: ComponentFixture<UrgencyFormComponent>;

  beforeEach(async () => {
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
        { provide: ToastrService, useValue: toastrSpy }
      ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(UrgencyFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should initialize form with default values', () => {
    expect(component.urgencyForm).toBeDefined();
    expect(component.urgencyForm.get('name')?.value).toBe('');
    expect(component.urgencyForm.get('email')?.value).toBe('');
    expect(component.urgencyForm.get('contactPhone')?.value).toBe('');
    expect(component.urgencyForm.get('location')?.value).toBe('');
    expect(component.urgencyForm.get('description')?.value).toBe('');
    expect(component.urgencyForm.get('level')?.value).toBe('Medium');
  });

  it('should validate required fields', () => {
    const form = component.urgencyForm;

    // Form should be invalid when empty
    expect(form.valid).toBeFalsy();

    // Check individual field validations
    expect(form.get('name')?.hasError('required')).toBeTruthy();
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
});
