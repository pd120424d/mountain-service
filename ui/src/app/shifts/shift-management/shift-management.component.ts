import { Component, OnInit } from "@angular/core";
import { AdministratorRole, Employee, MedicRole } from "../../employee/employee.model";
import { AuthService } from "../../services/auth.service";
import { TranslateModule, TranslateService } from "@ngx-translate/core";
import { RouterModule } from "@angular/router";
import { FormsModule } from "@angular/forms";
import { CommonModule } from "@angular/common";
import { EmployeeRole } from "../../employee/employee.model";
import { ShiftManagementService } from "../shift-management.service";
import { ToastrService } from "ngx-toastr";
import { NgxSpinnerService } from "ngx-spinner";
import { BaseTranslatableComponent } from "../../base-translatable.component";

@Component({
  selector: 'shift-management',
  standalone: true,
  imports: [RouterModule, FormsModule, TranslateModule, CommonModule],
  templateUrl: './shift-management.component.html',
  styleUrls: ['./shift-management.component.css'],
})
export class ShiftManagementComponent extends BaseTranslatableComponent implements OnInit {
  shifts: { [key: string]: { [profile: string]: number } } = {};
  userRole: EmployeeRole = MedicRole;
  userId = '';
  employees: Employee[] = [];

  constructor(private shiftService: ShiftManagementService,
    private auth: AuthService,
    private toastr: ToastrService,
    private spinner: NgxSpinnerService,
    translate: TranslateService) {
    super(translate);
  }

  ngOnInit() {
    this.userRole = this.auth.getRole();
    this.userId = this.auth.getUserId();
    this.loadShifts();
    if (this.userRole === AdministratorRole) {
      this.loadAllEmployees();
    }
  }

  loadShifts() {
    this.shiftService.getShiftAvailability().subscribe(data => {
      this.shifts = data;
    });
  }

  loadAllEmployees() {
    this.shiftService.getAllEmployees().subscribe(data => {
      this.employees = data;
    });
  }

  canModifyOthers(): boolean {
    return this.userRole === AdministratorRole;
  }

  assignToShift(shiftType: number, employeeId?: string) {
    const idToAssign = employeeId ?? this.userId;
    if (!idToAssign) return;

    this.spinner.show();
    this.shiftService.assignEmployeeToShift(shiftType, idToAssign).subscribe({
      next: () => {
        this.loadShifts();
        this.toastr.success(this.translate.instant('SHIFT_MANAGEMENT.TOAST_ASSIGN_SUCCESS'));
      },
      error: () => this.toastr.error(this.translate.instant('SHIFT_MANAGEMENT.TOAST_ASSIGN_ERROR')),
      complete: () => this.spinner.hide()
    });
  }

  removeFromShift(shiftType: number, employeeId?: string) {
    const idToRemove = employeeId ?? this.userId;
    if (!idToRemove) return;

    this.spinner.show();
    this.shiftService.removeEmployeeFromShift(shiftType, idToRemove).subscribe({
      next: () => {
        this.loadShifts();
        this.toastr.success(this.translate.instant('SHIFT_MANAGEMENT.TOAST_REMOVE_SUCCESS'));
      },
      error: () => this.toastr.error(this.translate.instant('SHIFT_MANAGEMENT.TOAST_REMOVE_ERROR')),
      complete: () => this.spinner.hide()
    });
  }
}
