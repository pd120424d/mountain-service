import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { ResolveModalComponent } from './resolve-modal.component';
import { sharedTestingProviders } from '../../test-utils/shared-test-imports';

describe('ResolveModalComponent', () => {
  let component: ResolveModalComponent;
  let fixture: ComponentFixture<ResolveModalComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ResolveModalComponent, ReactiveFormsModule],
      providers: [...sharedTestingProviders]
    }).compileComponents();

    fixture = TestBed.createComponent(ResolveModalComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should initialize with form controls', () => {
    expect(component.resolveForm.get('addActivity')).toBeTruthy();
    expect(component.resolveForm.get('activityDescription')).toBeTruthy();
  });

  it('should show activity field when checkbox is checked', () => {
    component.resolveForm.patchValue({ addActivity: true });
    expect(component.showActivityField).toBe(true);
  });

  it('should hide activity field when checkbox is unchecked', () => {
    component.resolveForm.patchValue({ addActivity: false });
    expect(component.showActivityField).toBe(false);
  });

  it('should emit confirmed result on confirm', () => {
    spyOn(component.result, 'emit');
    component.onConfirm();
    expect(component.result.emit).toHaveBeenCalledWith({ confirmed: true, activityDescription: undefined });
  });

  it('should emit cancelled result on cancel', () => {
    spyOn(component.result, 'emit');
    component.onCancel();
    expect(component.result.emit).toHaveBeenCalledWith({ confirmed: false });
  });

  it('should require activity description when checkbox is checked', () => {
    component.resolveForm.patchValue({ addActivity: true });
    const activityControl = component.resolveForm.get('activityDescription');
    expect(activityControl?.hasError('required')).toBe(true);
  });

  it('should not require activity description when checkbox is unchecked', () => {
    component.resolveForm.patchValue({ addActivity: false });
    const activityControl = component.resolveForm.get('activityDescription');
    expect(activityControl?.hasError('required')).toBe(false);
  });

  it('should validate and return error keys properly', () => {
    // Initially hidden field has no validators
    expect(component.isFieldInvalid('activityDescription')).toBe(false);
    expect(component.getFieldError('activityDescription')).toBe('');

    // Enable and require the field
    component.resolveForm.patchValue({ addActivity: true });
    const ctrl = component.resolveForm.get('activityDescription');
    ctrl?.markAsTouched();
    fixture.detectChanges();

    expect(component.isFieldInvalid('activityDescription')).toBe(true);
    expect(component.getFieldError('activityDescription')).toBe('RESOLVE_MODAL.ACTIVITY_REQUIRED');

    // Set too short value to trigger minlength
    ctrl?.setValue('ab');
    ctrl?.markAsTouched();
    ctrl?.updateValueAndValidity();
    expect(component.isFieldInvalid('activityDescription')).toBe(true);
    expect(component.getFieldError('activityDescription')).toBe('RESOLVE_MODAL.ACTIVITY_MIN_LENGTH');

    // Set valid value
    ctrl?.setValue('valid text');
    ctrl?.updateValueAndValidity();
    expect(component.isFieldInvalid('activityDescription')).toBe(false);
    expect(component.getFieldError('activityDescription')).toBe('');
  });

  it('should reset form and visibility on close after confirm/cancel', () => {
    spyOn(component.result, 'emit');

    // Make visible, enable activity, set value
    component.isVisible = true;
    component.resolveForm.patchValue({ addActivity: true, activityDescription: 'some text' });

    // Confirm
    component.onConfirm();
    expect(component.isVisible).toBeFalse();
    expect(component.showActivityField).toBeFalse();
    expect(component.resolveForm.get('addActivity')?.value).toBeFalsy();
    expect(component.resolveForm.get('activityDescription')?.value).toBeNull();

    // Show again and cancel
    component.isVisible = true;
    component.onCancel();
    expect(component.isVisible).toBeFalse();
    expect(component.showActivityField).toBeFalse();
    expect(component.resolveForm.get('addActivity')?.value).toBeFalsy();
  });

});
