import { Component, OnInit } from "@angular/core";
import { AdministratorRole, Employee, MedicRole } from "../employee/employee.model";
import { AuthService } from "../services/auth.service";
import { TranslateModule, TranslateService } from "@ngx-translate/core";
import { RouterModule } from "@angular/router";
import { FormsModule } from "@angular/forms";
import { CommonModule } from "@angular/common";
import { EmployeeRole } from "../employee/employee.model";
import { ShiftManagementService } from "./shift.service";
import { ToastrService } from "ngx-toastr";
import { NgxSpinnerService } from "ngx-spinner";
import { BaseTranslatableComponent } from "../base-translatable.component";
import { ShiftAvailabilityResponse } from "./shift.model";
import { format } from 'date-fns';
import { enUS, sr, srLatn } from 'date-fns/locale';

@Component({
  selector: 'shift',
  standalone: true,
  imports: [RouterModule, FormsModule, TranslateModule, CommonModule],
  templateUrl: './shift.component.html',
  styleUrls: ['./shift.component.css'],
})
export class ShiftManagementComponent extends BaseTranslatableComponent implements OnInit {
  shiftAvailability: ShiftAvailabilityResponse = { days: {} };
  userRole: EmployeeRole = MedicRole;
  userId = '';
  employees: Employee[] = [];
  dates: Date[] = [];

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
      this.shiftAvailability = data;
      this.dates = Object.keys(data.days)
      .map(d => new Date(d))
      .sort((a, b) => a.getTime() - b.getTime());
      console.log(this.dates);
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

  assignToShift(shiftType: number, date: Date, employeeId?: string) {
    const idToAssign = employeeId ?? this.userId;
    if (!idToAssign) return;

    this.spinner.show();
    this.shiftService.assignEmployeeToShift(shiftType, idToAssign, date).subscribe({
      next: () => {
        this.loadShifts();
        this.toastr.success(this.translate.instant('SHIFT_MANAGEMENT.TOAST_ASSIGN_SUCCESS'));
      },
      error: () => this.toastr.error(this.translate.instant('SHIFT_MANAGEMENT.TOAST_ASSIGN_ERROR')),
      complete: () => this.spinner.hide()
    });
  }

  removeFromShift(shiftId: number, employeeId?: string) {
    const idToRemove = employeeId ?? this.userId;
    if (!idToRemove) return;

    this.spinner.show();
    this.shiftService.removeEmployeeFromShift(idToRemove, shiftId).subscribe({
      next: () => {
        this.loadShifts();
        this.toastr.success(this.translate.instant('SHIFT_MANAGEMENT.TOAST_REMOVE_SUCCESS'));
      },
      error: () => this.toastr.error(this.translate.instant('SHIFT_MANAGEMENT.TOAST_REMOVE_ERROR')),
      complete: () => this.spinner.hide()
    });
  }
  getAvailableMedics(shiftType: number, date: Date): number {
    const key = date.toISOString().split('T')[0];
    const day = this.shiftAvailability?.days?.[key];
    if (!day) return 0;
  
    switch (shiftType) {
      case 1: return day.firstShift.medic;
      case 2: return day.secondShift.medic;
      case 3: return day.thirdShift.medic;
      default: return 0;
    }
  }
  
  getAvailableTechnicals(shiftType: number, date: Date): number {
    const key = date.toISOString().split('T')[0];
    const day = this.shiftAvailability?.days?.[key];
    if (!day) return 0;
  
    switch (shiftType) {
      case 1: return day.firstShift.technical;
      case 2: return day.secondShift.technical;
      case 3: return day.thirdShift.technical;
      default: return 0;
    }
  }
  

  getShiftLabel(type: number): string {
    switch (type) {
      case 1: return '06:00 - 14:00';
      case 2: return '14:00 - 22:00';
      case 3: return '22:00 - 06:00';
      default: return '';
    }
  }

  getTranslatedDate(date: Date): string {
  const lang = this.translate.currentLang;
  const locale = lang === 'sr' ? sr : lang === 'sr-Latn' ? srLatn : enUS;
  return format(date, 'EEEE, MMMM d, yyyy', { locale });
}
}
