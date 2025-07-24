import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router } from '@angular/router';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { ToastrService } from 'ngx-toastr';
import { BaseTranslatableComponent } from '../../base-translatable.component';
import { UrgencyService } from '../urgency.service';
import { UrgencyCreateRequest, UrgencyLevel } from '../../shared/models';

@Component({
  selector: 'app-urgency-form',
  standalone: true,
  imports: [RouterModule, TranslateModule, CommonModule, ReactiveFormsModule],
  templateUrl: './urgency-form.component.html',
  styleUrls: ['./urgency-form.component.css']
})
export class UrgencyFormComponent extends BaseTranslatableComponent implements OnInit {
  urgencyForm!: FormGroup;
  urgencyLevels = Object.values(UrgencyLevel);
  isSubmitting = false;

  constructor(
    private fb: FormBuilder,
    private urgencyService: UrgencyService,
    private router: Router,
    private toastr: ToastrService,
    translate: TranslateService
  ) {
    super(translate);
  }

  ngOnInit(): void {
    this.initializeForm();
  }

  private initializeForm(): void {
    this.urgencyForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(2)]],
      email: ['', [Validators.required, Validators.email]],
      contactPhone: ['', [Validators.required, Validators.minLength(8)]],
      location: ['', [Validators.required, Validators.minLength(3)]],
      description: ['', [Validators.required, Validators.minLength(10)]],
      level: [UrgencyLevel.MEDIUM, [Validators.required]]
    });
  }

  onSubmit(): void {
    if (this.urgencyForm.valid && !this.isSubmitting) {
      this.isSubmitting = true;
      const urgencyRequest: UrgencyCreateRequest = this.urgencyForm.value;

      this.urgencyService.addUrgency(urgencyRequest).subscribe({
        next: (response) => {
          console.log('Urgency created successfully:', response);
          this.toastr.success(this.translate.instant('URGENCY_FORM.SUCCESS_MESSAGE'));
          this.router.navigate(['/']);
        },
        error: (error) => {
          console.error('Error creating urgency:', error);
          this.isSubmitting = false;
        }
      });
    } else {
      this.markFormGroupTouched();
    }
  }

  cancel(): void {
    this.router.navigate(['/']);
  }

  private markFormGroupTouched(): void {
    Object.keys(this.urgencyForm.controls).forEach(key => {
      const control = this.urgencyForm.get(key);
      control?.markAsTouched();
    });
  }

  // Helper methods for template
  isFieldInvalid(fieldName: string): boolean {
    const field = this.urgencyForm.get(fieldName);
    return !!(field && field.invalid && (field.dirty || field.touched));
  }

  getFieldError(fieldName: string): string {
    const field = this.urgencyForm.get(fieldName);
    if (field?.errors) {
      if (field.errors['required']) return `${fieldName.toUpperCase()}_FORM.REQUIRED`;
      if (field.errors['email']) return `${fieldName.toUpperCase()}_FORM.INVALID_EMAIL`;
      if (field.errors['minlength']) return `${fieldName.toUpperCase()}_FORM.MIN_LENGTH`;
    }
    return '';
  }
}