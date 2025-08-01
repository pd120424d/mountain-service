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
    translate: TranslateService) {
    super(translate);
  }

  ngOnInit(): void {
    this.loadEmployees();
  }

  loadEmployees(): void {
    this.employeeService.getEmployees().subscribe(data => {
      this.employees = data;
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
    this.employeeService.deleteEmployee(id).subscribe(() => {
      this.loadEmployees();
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
