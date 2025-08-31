// src/app/employee/employee-form/employee-form.component.ts
import { Component, OnInit } from '@angular/core';
import { AbstractControl, FormBuilder, FormGroup, ValidationErrors, ValidatorFn, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { EmployeeService } from '../employee.service';
import { Employee, EmployeeCreateRequest, EmployeeUpdateRequest } from '../../shared/models';
import { ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { BaseTranslatableComponent } from '../../base-translatable.component';
import { ToastrService } from 'ngx-toastr';
import { NgxSpinnerService } from 'ngx-spinner';
import { NgxSpinnerModule } from 'ngx-spinner';
import { environment } from '../../../environments/environment';
import { ImageUploadComponent, ImageUploadEvent } from '../../shared/components/image-upload/image-upload.component';
import { ImageUploadService } from '../../services/image-upload.service';

@Component({
  selector: 'app-employee-form',
  templateUrl: './employee-form.component.html',
  styleUrls: ['./employee-form.component.css'],
  standalone: true,
  imports: [RouterModule, ReactiveFormsModule, CommonModule, TranslateModule, NgxSpinnerModule, ImageUploadComponent],
})
export class EmployeeFormComponent extends BaseTranslatableComponent implements OnInit {
  employeeForm: FormGroup;
  employeeId?: number;
  isEditMode = false;
  isSubmitting = false;
  selectedImageFile: File | null = null;
  currentProfilePictureUrl: string | undefined = undefined;

  // UI state for password visibility
  showPassword = false;
  showConfirmPassword = false;

  constructor(
    private fb: FormBuilder,
    private employeeService: EmployeeService,
    private router: Router,
    translate: TranslateService,
    private toastr: ToastrService,
    private spinner: NgxSpinnerService,
    private imageUploadService: ImageUploadService
  ) {
    super(translate);

    this.employeeForm = this.fb.group({
      username: ['', Validators.required],
      password: ['', Validators.required], // Only relevant in create mode
      confirmPassword: ['', Validators.required], // Only relevant in create mode
      firstName: ['', Validators.required],
      lastName: ['', Validators.required],
      gender: ['', Validators.required],
      phoneNumber: ['', Validators.required],
      email: ['', [Validators.required, Validators.email]],
      profileType: ['', Validators.required]
    });

    const navigation = this.router.getCurrentNavigation();
    const state = navigation?.extras.state as { employee: Employee };

    if (state && state.employee) {
      this.isEditMode = true;
      this.employeeId = state.employee.id;
      this.populateForm(state.employee);
    } else {
      this.isEditMode = false; // We're in create mode
    }
  }

  ngOnInit(): void {
    // Apply backend-consistent password validation only in create mode
    if (!this.isEditMode) {
      const pwdCtrl = this.employeeForm.get('password');
      const confirmCtrl = this.employeeForm.get('confirmPassword');
      pwdCtrl?.setValidators([Validators.required, this.passwordPolicyValidator()]);
      confirmCtrl?.setValidators([Validators.required]);
      // Re-run validation after setting validators
      pwdCtrl?.updateValueAndValidity();
      confirmCtrl?.updateValueAndValidity();
    }
  }

  // Validator mirroring api/shared/utils/crypt.go ValidatePassword
  private passwordPolicyValidator(): ValidatorFn {
    return (control: AbstractControl): ValidationErrors | null => {
      const value: string = control.value || '';
      if (value.length === 0) return null; // let required handle empties

      // Length 6..20
      if (value.length < 6 || value.length > 20) {
        return { passwordPolicy: { code: 'ErrPasswordLength' } };
      }
      // Must start with a letter
      if (!/^[A-Za-z].*$/.test(value)) {
        return { passwordPolicy: { code: 'ErrPasswordStartWithLetter' } };
      }
      // At least 1 uppercase
      if (!/[A-Z]/.test(value)) {
        return { passwordPolicy: { code: 'ErrPasswordUppercase' } };
      }
      // At least 3 lowercase letters
      const lowerCount = (value.match(/[a-z]/g) || []).length;
      if (lowerCount < 3) {
        return { passwordPolicy: { code: 'ErrPasswordLowercase' } };
      }
      // At least 1 digit
      if (!/[0-9]/.test(value)) {
        return { passwordPolicy: { code: 'ErrPasswordDigit' } };
      }
      // At least 1 special (punctuation or symbol)
      if (!/[\p{P}\p{S}]/u.test(value)) {
        return { passwordPolicy: { code: 'ErrPasswordSpecial' } };
      }

      return null;
    };
  }

  // Helper to get a translated error message for password
  getPasswordErrorKey(): string {
    const ctrl = this.employeeForm.get('password');
    if (!ctrl || !ctrl.errors) return '';

    if (ctrl.errors['required']) return 'EMPLOYEE_FORM.PASSWORD_REQUIRED';
    if (ctrl.errors['passwordPolicy']) {
      const code = ctrl.errors['passwordPolicy'].code as string;
      switch (code) {
        case 'ErrPasswordLength': return 'EMPLOYEE_FORM.PASSWORD_ERR_LENGTH';
        case 'ErrPasswordUppercase': return 'EMPLOYEE_FORM.PASSWORD_ERR_UPPER';
        case 'ErrPasswordLowercase': return 'EMPLOYEE_FORM.PASSWORD_ERR_LOWER';
        case 'ErrPasswordDigit': return 'EMPLOYEE_FORM.PASSWORD_ERR_DIGIT';
        case 'ErrPasswordSpecial': return 'EMPLOYEE_FORM.PASSWORD_ERR_SPECIAL';
        case 'ErrPasswordStartWithLetter': return 'EMPLOYEE_FORM.PASSWORD_ERR_START_LETTER';
      }
    }
    return 'EMPLOYEE_FORM.PASSWORD_INVALID';
  }

  // Confirm password simple equality check for UX
  get confirmPasswordMismatch(): boolean {
    const pwd = this.employeeForm.get('password')?.value;
    const c = this.employeeForm.get('confirmPassword')?.value;
    return !!(pwd && c && pwd !== c);
  }

  populateForm(employee: Employee): void {
    this.employeeForm.patchValue({
      username: employee.username,
      firstName: employee.firstName,
      lastName: employee.lastName,
      gender: employee.gender,
      phoneNumber: employee.phone,
      email: employee.email,
      profileType: employee.profileType,
    });

    if (this.isEditMode) {
      this.employeeForm.get('email')?.disable();
    }

    // Set current profile picture URL for display
    this.currentProfilePictureUrl = employee.profilePicture || undefined;
  }

  onSubmit(): void {
    if (this.employeeForm.valid && !this.isSubmitting) {
      this.isSubmitting = true;
      this.spinner.show();

      const formValue = this.employeeForm.value;

      if (this.employeeId) {
        // Update existing employee
        this.updateEmployee(formValue);
      } else {
        // Create new employee
        this.createEmployee(formValue);
      }
    } else if (!this.employeeForm.valid) {
      // Show validation error
      this.toastr.warning(this.translate.instant('EMPLOYEE_FORM.VALIDATION_ERROR'));
      this.markFormGroupTouched();
    }
  }

  cancel(): void {
    this.router.navigate(['/employees']);
  }

  onImageSelected(event: ImageUploadEvent): void {
    // In edit mode only; do nothing during registration
    if (!this.isEditMode) return;
    if (event.isValid && event.file) {
      this.selectedImageFile = event.file;
    } else {
      this.selectedImageFile = null;
      this.toastr.error(event.error || 'Invalid image file');
    }
  }

  onImageRemoved(): void {
    if (!this.isEditMode) return;
    this.selectedImageFile = null;
    this.currentProfilePictureUrl = undefined;
  }

  /**
   * Handle errors in a user-friendly way
   */
  private handleError(error: any, defaultErrorKey: string): void {
    this.spinner.hide();
    this.isSubmitting = false;

    let errorMessage = this.translate.instant(`EMPLOYEE_FORM.${defaultErrorKey}`);

    // Categorize errors for better user experience
    if (!navigator.onLine) {
      // Network is offline
      errorMessage = this.translate.instant('EMPLOYEE_FORM.NETWORK_ERROR');
    } else if (error?.status === 0) {
      // Network error (server unreachable)
      errorMessage = this.translate.instant('EMPLOYEE_FORM.NETWORK_ERROR');
    } else if (error?.status >= 400 && error?.status < 500) {
      // Client error (validation, authentication, etc.)
      if (error?.error?.message) {
        // Use server-provided message if available and user-friendly
        errorMessage = error.error.message;
      } else if (error?.status === 409) {
        errorMessage = 'An employee with this username or email already exists.';
      } else if (error?.status === 400) {
        errorMessage = 'Please check your input and try again.';
      }
    } else if (error?.status >= 500) {
      // Server error
      errorMessage = this.translate.instant('EMPLOYEE_FORM.SERVER_ERROR');
    }

    // Show user-friendly error message
    this.toastr.error(errorMessage);

    // Log technical details for debugging (only in development)
    if (!environment.production) {
      console.error('Employee form error:', error);
    }
  }

  private createEmployee(formValue: any): void {
    // Registration has no profile picture handling; just create the employee
    this.createEmployeeJSON(formValue);
  }

  private updateEmployee(formValue: any): void {
    // For updates, we'll use JSON and handle image separately if needed
    const employeeUpdate: EmployeeUpdateRequest = {
      firstName: formValue.firstName,
      lastName: formValue.lastName,
      email: formValue.email,
      phone: formValue.phoneNumber,
      username: formValue.username,
      profileType: formValue.profileType,
      gender: formValue.gender,
      profilePicture: this.currentProfilePictureUrl || undefined
    };

    this.employeeService.updateEmployee(this.employeeId!, employeeUpdate).subscribe({
      next: (response) => {
        // If there's a new image, upload it after update
        if (this.selectedImageFile && response.id) {
          this.uploadImageAfterUpdate(response.id);
        } else {
          this.spinner.hide();
          this.toastr.success(this.translate.instant('EMPLOYEE_FORM.UPDATE_SUCCESS'));
          this.router.navigate(['/employees']);
        }
      },
      error: (error) => {
        this.handleError(error, 'UPDATE_ERROR');
      }
    });
  }

  private createEmployeeJSON(formValue: any): void {
    const employeeCreateRequest: EmployeeCreateRequest = {
      firstName: formValue.firstName,
      lastName: formValue.lastName,
      email: formValue.email,
      phone: formValue.phoneNumber,
      username: formValue.username,
      password: formValue.password,
      profileType: formValue.profileType,
      gender: formValue.gender
    };

    this.employeeService.addEmployee(employeeCreateRequest).subscribe({
      next: (response) => {
        this.spinner.hide();
        this.toastr.success(this.translate.instant('EMPLOYEE_FORM.CREATE_SUCCESS'));
        this.router.navigate(['/employees']);
      },
      error: (error) => {
        this.handleError(error, 'CREATE_ERROR');
      }
    });
  }


  private uploadImageAfterUpdate(employeeId: number): void {
    if (!this.selectedImageFile) {
      this.spinner.hide();
      this.toastr.success(this.translate.instant('EMPLOYEE_FORM.UPDATE_SUCCESS'));
      this.router.navigate(['/employees']);
      return;
    }

    this.imageUploadService.uploadProfilePicture(employeeId, this.selectedImageFile).subscribe({
      next: (progress) => {
        if (progress.status === 'completed') {
          this.spinner.hide();
          this.toastr.success(this.translate.instant('EMPLOYEE_FORM.UPDATE_SUCCESS'));
          this.router.navigate(['/employees']);
        }
      },
      error: (error) => {
        this.spinner.hide();
        this.toastr.warning('Employee updated but image upload failed: ' + error.message);
        this.router.navigate(['/employees']);
      }
    });
  }

  /**
   * Mark all form fields as touched to show validation errors
   */
  private markFormGroupTouched(): void {
    Object.keys(this.employeeForm.controls).forEach(key => {
      const control = this.employeeForm.get(key);
      control?.markAsTouched();
    });
  }
}
