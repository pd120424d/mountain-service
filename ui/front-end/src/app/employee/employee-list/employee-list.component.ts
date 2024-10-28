// src/app/employee/employee-list/employee-list.component.ts
import { Component, OnInit } from '@angular/core';
import { EmployeeService } from '../employee.service';
import { Employee } from '../employee.model';
import { ReactiveFormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common'
import { RouterModule } from '@angular/router';
import { Router } from '@angular/router';
import { MatDialog, MatDialogModule } from '@angular/material/dialog';
import { ConfirmDialogComponent } from '../confirm-dialog/confirm-dialog.component';



@Component({
  selector: 'app-employee-list',
  templateUrl: './employee-list.component.html',
  styleUrls: ['./employee-list.component.css'],
  standalone: true,
  imports: [RouterModule, ReactiveFormsModule, CommonModule], 
})
export class EmployeeListComponent implements OnInit {
  showModal = false;
  employeeToDelete: number | null = null;
  employees: Employee[] = [];

  constructor(private employeeService: EmployeeService, private router: Router, private dialog: MatDialog) { }

  ngOnInit(): void {
    this.loadEmployees();
    console.log('EmployeeListComponent initialized');
  }

  loadEmployees(): void {
    this.employeeService.getEmployees().subscribe(data => {
      this.employees = data;
    });
  }

  // Open the confirmation dialog
  openDeleteDialog(employeeId: number): void {
    const dialogRef = this.dialog.open(ConfirmDialogComponent);

    // Subscribe to the dialog result
    dialogRef.afterClosed().subscribe((result: any) => {
      if (result) {
        this.deleteEmployee(employeeId); // Proceed with deletion if confirmed
      }
    });
  }

  deleteEmployee(id: number): void {
    this.employeeService.deleteEmployee(id).subscribe(() => {
      this.loadEmployees(); // Reload the list after deletion
    });
  }

  editEmployee(employee: Employee): void {
    console.log('Navigating to edit with employee:', employee); // Debug

    this.router.navigate(['/employees/edit', employee.id], {
      state: { employee }, // Pass the employee data via router state
    });
  }
}
