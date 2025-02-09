// src/app/employee/employee-form/employee-form.component.ts
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, Validators } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { EmployeeService } from '../employee.service';
import { Employee } from '../employee.model';
import { ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { TranslateModule, TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-employee-form',
  templateUrl: './employee-form.component.html',
  styleUrls: ['./employee-form.component.css'],
  standalone: true,
  imports: [RouterModule, ReactiveFormsModule, CommonModule, TranslateModule],
})
export class EmployeeFormComponent implements OnInit {
  employeeForm: FormGroup;
  employeeId?: number;
  isEditMode = false; // Flag to check if we're editing

  constructor(
    private fb: FormBuilder,
    private employeeService: EmployeeService,
    private router: Router,
    private translate: TranslateService,
  ) {
    this.employeeForm = this.fb.group({
      username: ['', Validators.required],
      password: ['', Validators.required], // Only relevant in create mode
      confirmPassword: ['', Validators.required], // Only relevant in create mode
      firstName: ['', Validators.required],
      lastName: ['', Validators.required],
      gender: ['', Validators.required],
      phone: ['', Validators.required],
      email: ['', [Validators.required, Validators.email]],
      profileType: ['', Validators.required],
      profilePicture: [null],
    });

    // Check if employee data is passed via router state
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

    this.translate.setDefaultLang('sr-cyr');
  }

  ngOnInit(): void {

  }

  // Populate the form with the employee data
  populateForm(employee: Employee): void {
    this.employeeForm.patchValue({
      username: employee.username,
      firstName: employee.firstName,
      lastName: employee.lastName,
      gender: employee.gender,
      phone: employee.phoneNumber,
      email: employee.email,
      profileType: employee.profileType,
    });
  }

  onSubmit(): void {
    if (this.employeeForm.valid) {
      const employee: Employee = this.employeeForm.value;
      if (this.employeeId) {
        this.employeeService.updateEmployee(this.employeeId, employee).subscribe(() => {
          this.router.navigate(['/employees']);
        });
      } else {
        this.employeeService.addEmployee(employee).subscribe(() => {
          this.router.navigate(['/employees']);
        });
      }
    }
  }

  // Method to handle cancel button click
  cancel(): void {
    this.router.navigate(['/employees']); // Navigate back to the employee list or another appropriate route
  }

  switchLanguage(language: string): void {
    this.translate.use(language);
  }
}
