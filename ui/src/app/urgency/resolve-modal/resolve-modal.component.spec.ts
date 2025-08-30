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
});
