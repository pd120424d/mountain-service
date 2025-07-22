// src/app/employee/employee-form/employee-form.component.ts
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { EmployeeService } from '../employee.service';
import { Employee, EmployeeCreateRequest } from '../employee.model';
import { ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { BaseTranslatableComponent } from '../../base-translatable.component';

@Component({
  selector: 'app-employee-form',
  templateUrl: './employee-form.component.html',
  styleUrls: ['./employee-form.component.css'],
  standalone: true,
  imports: [RouterModule, ReactiveFormsModule, CommonModule, TranslateModule],
})
export class EmployeeFormComponent extends BaseTranslatableComponent implements OnInit {
  employeeForm: FormGroup;
  employeeId?: number;
  isEditMode = false;

  constructor(
    private fb: FormBuilder,
    private employeeService: EmployeeService,
    private router: Router,
    translate: TranslateService,
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
      console.log("Employee passed as state:", state.employee);

      this.isEditMode = true;
      this.employeeId = state.employee.id;
      this.populateForm(state.employee);
    } else {
      this.isEditMode = false; // We're in create mode
      console.log("No employee data passed via state.");
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
    if (this.employeeForm.valid) {
      const formValue = this.employeeForm.value;

      if (this.employeeId) {
        const employee: Employee = {
          id: this.employeeId,
          firstName: formValue.firstName,
          lastName: formValue.lastName,
          email: formValue.email,
          phone: formValue.phoneNumber,
          username: formValue.username,
          profileType: formValue.profileType,
          gender: formValue.gender,
          profilePicture: formValue.profilePicture || undefined
        };

        this.employeeService.updateEmployee(this.employeeId, employee).subscribe({
          next: () => {
            this.router.navigate(['/employees']);
          },
          error: (error) => {
            console.error('Error updating employee:', error);
            alert(`Error updating employee: ${error.message}`);
          }
        });
      } else {
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
            this.router.navigate(['/employees']);
          },
          error: (error) => {
            console.error('Error creating employee:', error);
            alert(`Error creating employee: ${error.message}`);
          }
        });
      }
    }
  }

  cancel(): void {
    this.router.navigate(['/employees']); // Navigate back to the employee list or another appropriate route
  }
}
