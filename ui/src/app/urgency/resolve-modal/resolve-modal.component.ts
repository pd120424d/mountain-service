import { Component, EventEmitter, Input, Output } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { BaseTranslatableComponent } from '../../base-translatable.component';

export interface ResolveModalResult {
  confirmed: boolean;
  activityDescription?: string;
}

@Component({
  selector: 'app-resolve-modal',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, TranslateModule],
  templateUrl: './resolve-modal.component.html',
  styleUrls: ['./resolve-modal.component.css']
})
export class ResolveModalComponent extends BaseTranslatableComponent {
  @Input() isVisible = false;
  @Input() urgencyId: number | null = null;
  @Output() result = new EventEmitter<ResolveModalResult>();

  resolveForm: FormGroup;
  showActivityField = false;

  constructor(
    private fb: FormBuilder,
    translate: TranslateService
  ) {
    super(translate);
    
    this.resolveForm = this.fb.group({
      addActivity: [false],
      activityDescription: ['']
    });

    // Watch for checkbox changes to show/hide activity field
    this.resolveForm.get('addActivity')?.valueChanges.subscribe(checked => {
      this.showActivityField = checked;
      const activityControl = this.resolveForm.get('activityDescription');
      if (checked) {
        activityControl?.setValidators([Validators.required, Validators.minLength(3)]);
      } else {
        activityControl?.clearValidators();
        activityControl?.setValue('');
      }
      activityControl?.updateValueAndValidity();
    });
  }

  onConfirm(): void {
    if (this.resolveForm.valid) {
      const formValue = this.resolveForm.value;
      this.result.emit({
        confirmed: true,
        activityDescription: formValue.addActivity ? formValue.activityDescription : undefined
      });
      this.close();
    } else {
      this.markFormGroupTouched();
    }
  }

  onCancel(): void {
    this.result.emit({ confirmed: false });
    this.close();
  }

  private close(): void {
    this.isVisible = false;
    this.resolveForm.reset();
    this.showActivityField = false;
  }

  private markFormGroupTouched(): void {
    Object.keys(this.resolveForm.controls).forEach(key => {
      const control = this.resolveForm.get(key);
      control?.markAsTouched();
    });
  }

  isFieldInvalid(fieldName: string): boolean {
    const field = this.resolveForm.get(fieldName);
    return !!(field && field.invalid && (field.dirty || field.touched));
  }

  getFieldError(fieldName: string): string {
    const field = this.resolveForm.get(fieldName);
    if (field && field.errors) {
      if (field.errors['required']) {
        return 'RESOLVE_MODAL.ACTIVITY_REQUIRED';
      }
      if (field.errors['minlength']) {
        return 'RESOLVE_MODAL.ACTIVITY_MIN_LENGTH';
      }
    }
    return '';
  }
}
