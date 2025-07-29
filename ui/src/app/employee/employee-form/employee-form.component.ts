// src/app/employee/employee-form/employee-form.component.ts
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
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

@Component({
  selector: 'app-employee-form',
  templateUrl: './employee-form.component.html',
  styleUrls: ['./employee-form.component.css'],
  standalone: true,
  imports: [RouterModule, ReactiveFormsModule, CommonModule, TranslateModule, NgxSpinnerModule],
})
export class EmployeeFormComponent extends BaseTranslatableComponent implements OnInit {
  employeeForm: FormGroup;
  employeeId?: number;
  isEditMode = false;
  isSubmitting = false;

  constructor(
    private fb: FormBuilder,
    private employeeService: EmployeeService,
    private router: Router,
    translate: TranslateService,
    private toastr: ToastrService,
    private spinner: NgxSpinnerService
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
      profileType: ['', Validators.required],
      profilePicture: [null],
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

  ngOnInit(): void {}

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
  }

  onSubmit(): void {
    if (this.employeeForm.valid && !this.isSubmitting) {
      this.isSubmitting = true;
      this.spinner.show();

      const formValue = this.employeeForm.value;

      if (this.employeeId) {
        // Update existing employee
        const employeeUpdate: EmployeeUpdateRequest = {
          firstName: formValue.firstName,
          lastName: formValue.lastName,
          email: formValue.email,
          phone: formValue.phoneNumber,
          username: formValue.username,
          profileType: formValue.profileType,
          gender: formValue.gender,
          profilePicture: formValue.profilePicture || undefined
        };

        this.employeeService.updateEmployee(this.employeeId, employeeUpdate).subscribe({
          next: () => {
            this.spinner.hide();
            this.toastr.success(this.translate.instant('EMPLOYEE_FORM.UPDATE_SUCCESS'));
            this.router.navigate(['/employees']);
          },
          error: (error) => {
            this.handleError(error, 'UPDATE_ERROR');
          }
        });
      } else {
        // Create new employee
        const employeeCreateRequest: EmployeeCreateRequest = {
          firstName: formValue.firstName,
          lastName: formValue.lastName,
          email: formValue.email,
          phone: formValue.phoneNumber,
          username: formValue.username,
          password: formValue.password,
          profileType: formValue.profileType,
          gender: formValue.gender,
          profilePicture: formValue.profilePicture || undefined
        };

        this.employeeService.addEmployee(employeeCreateRequest).subscribe({
          next: () => {
            this.spinner.hide();
            this.toastr.success(this.translate.instant('EMPLOYEE_FORM.CREATE_SUCCESS'));
            this.router.navigate(['/employees']);
          },
          error: (error) => {
            this.handleError(error, 'CREATE_ERROR');
          }
        });
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
