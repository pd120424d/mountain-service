// src/app/employee/employee-list/employee-list.component.ts
import { Component, OnInit } from '@angular/core';
import { EmployeeService } from '../employee.service';
import { Employee } from '../../shared/models';
import { ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common'
import { RouterModule } from '@angular/router';
import { Router } from '@angular/router';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { ConfirmDialogComponent } from '../confirm-dialog/confirm-dialog.component';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { BaseTranslatableComponent } from '../../base-translatable.component';
import { AuthService } from '../../services/auth.service';


@Component({
  selector: 'app-employee-list',
  templateUrl: './employee-list.component.html',
  styleUrls: ['./employee-list.component.css'],
  standalone: true,
  imports: [RouterModule, ReactiveFormsModule, CommonModule, TranslateModule],
})
export class EmployeeListComponent extends BaseTranslatableComponent implements OnInit {
  showModal = false;
  employeeToDelete: number | null = null;
  employees: Employee[] = [];

  constructor(private employeeService: EmployeeService,
    private router: Router,
    private dialog: MatDialog,
    public authService: AuthService,
    translate: TranslateService) {
    super(translate);

    // Redirect non-admin users to home page
    if (!this.authService.isAdmin()) {
      this.router.navigate(['/']);
    }
  }

  ngOnInit(): void {
    this.loadEmployees();
  }

  loadEmployees(): void {
    this.employeeService.getEmployees().subscribe({
      next: (data) => {
        this.employees = data;
      },
      error: (error) => {
        console.error('Error loading employees:', error);
        this.employees = [];
      }
    });
  }

  openDeleteDialog(employeeId: number): void {
    const dialogRef = this.dialog.open(ConfirmDialogComponent);

    dialogRef.afterClosed().subscribe((result: any) => {
      if (result) {
        this.deleteEmployee(employeeId);
      }
    });
  }

  deleteEmployee(id: number): void {
    this.employeeService.deleteEmployee(id).subscribe({
      next: () => {
        this.loadEmployees();
      },
      error: (error) => {
        console.error('Error deleting employee:', error);
      }
    });
  }

  editEmployee(employee: Employee): void {
    console.log('Navigating to edit with employee:', employee);

    this.router.navigate(['/employees/edit', employee.id], {
      state: { employee },
    });
  }

  goBackToHome(): void {
    this.router.navigate(['/']);
  }
}
