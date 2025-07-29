import { Component, OnInit } from "@angular/core";
import { AdministratorRole, Employee, MedicRole, EmployeeRole, ShiftAvailabilityResponse } from "../shared/models";
import { AuthService } from "../services/auth.service";
import { TranslateModule, TranslateService } from "@ngx-translate/core";
import { Router, RouterModule } from "@angular/router";
import { FormsModule } from "@angular/forms";
import { CommonModule } from "@angular/common";
import { ShiftManagementService } from "./shift.service";
import { ToastrService } from "ngx-toastr";
import { NgxSpinnerService } from "ngx-spinner";
import { BaseTranslatableComponent } from "../base-translatable.component";
import { format } from 'date-fns';
import { enUS, sr, srLatn, ru } from 'date-fns/locale';

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
  isLoading = true;
  isAssigning = false;
  isRemoving = false;
  shiftWarnings: string[] = [];

  constructor(private shiftService: ShiftManagementService,
    private auth: AuthService,
    private toastr: ToastrService,
    private spinner: NgxSpinnerService,
    private router: Router,
    translate: TranslateService) {
    super(translate);
  }

  ngOnInit() {
    this.userRole = this.auth.getRole();
    this.userId = this.auth.getUserId();
    this.loadShifts();
    this.loadShiftWarnings();
    if (this.userRole === AdministratorRole) {
      this.loadAllEmployees();
    }
  }

  loadShifts() {
    this.isLoading = true;
    this.shiftService.getShiftAvailability().subscribe({
      next: (data) => {
        console.log('Shift availability data received:', data);
        this.shiftAvailability = data;
        this.dates = Object.keys(data.days || {})
          .map(d => new Date(d))
          .sort((a, b) => a.getTime() - b.getTime());
        console.log('Processed dates:', this.dates);
        console.log('Shift availability object:', this.shiftAvailability);
        this.isLoading = false;
      },
      error: (error) => {
        console.error('Error loading shift availability:', error);
        this.toastr.error(this.translate.instant('SHIFT_MANAGEMENT.TOAST_LOAD_ERROR') || 'Failed to load shift data');
        this.isLoading = false;
      }
    });
  }

  loadAllEmployees() {
    this.shiftService.getAllEmployees().subscribe(data => {
      this.employees = data;
    });
  }

  loadShiftWarnings() {
    if (this.userId) {
      this.shiftService.getShiftWarnings(this.userId).subscribe({
        next: (data) => {
          this.shiftWarnings = data.warnings || [];
          console.log('Shift warnings loaded:', this.shiftWarnings);
        },
        error: (error) => {
          console.error('Error loading shift warnings:', error);
          // Don't show error toast for warnings as they're not critical
        }
      });
    }
  }

  canModifyOthers(): boolean {
    return this.userRole === AdministratorRole;
  }

  assignToShift(shiftType: number, date: Date, employeeId?: string) {
    const idToAssign = employeeId ?? this.userId;
    if (!idToAssign) return;

    console.log('Assigning to shift:', { shiftType, date, employeeId: idToAssign });
    this.isAssigning = true;
    this.spinner.show();
    this.shiftService.assignEmployeeToShift(shiftType, idToAssign, date).subscribe({
      next: (response) => {
        console.log('Assignment successful:', response);
        this.loadShifts();
        this.isAssigning = false;
        this.spinner.hide();
        this.toastr.success(this.translate.instant('SHIFT_MANAGEMENT.TOAST_ASSIGN_SUCCESS'));
      },
      error: (error) => {
        console.error('Assignment failed:', error);
        this.isAssigning = false;
        this.spinner.hide();
        this.toastr.error(this.translate.instant('SHIFT_MANAGEMENT.TOAST_ASSIGN_ERROR'));
      },
      complete: () => {
        this.isAssigning = false;
        this.spinner.hide();
      }
    });
  }

  removeFromShift(shiftType: number, employeeId?: string, date?: Date) {
    const idToRemove = employeeId ?? this.userId;
    if (!idToRemove || !date) {
      console.error('Missing required parameters for shift removal:', { employeeId: idToRemove, date });
      this.toastr.error(this.translate.instant('SHIFT_MANAGEMENT.TOAST_REMOVE_ERROR'));
      return;
    }

    console.log('Removing from shift:', { shiftType, employeeId: idToRemove, date });

    this.isRemoving = true;
    this.spinner.show();
    this.shiftService.removeEmployeeFromShiftByDetails(idToRemove, shiftType, date).subscribe({
      next: (response) => {
        console.log('Removal successful:', response);
        this.loadShifts();
        this.isRemoving = false;
        this.spinner.hide();
        this.toastr.success(this.translate.instant('SHIFT_MANAGEMENT.TOAST_REMOVE_SUCCESS'));
      },
      error: (error) => {
        console.error('Removal failed:', error);
        this.isRemoving = false;
        this.spinner.hide();
        this.toastr.error(this.translate.instant('SHIFT_MANAGEMENT.TOAST_REMOVE_ERROR'));
      },
      complete: () => {
        this.isRemoving = false;
        this.spinner.hide();
      }
    });
  }

  getAvailableMedics(shiftType: number, date: Date): number {
    const key = date.toISOString().split('T')[0];
    const day = this.shiftAvailability?.days?.[key];

    if (!day) {
      console.warn(`No day data found for key: ${key}. Available keys:`, Object.keys(this.shiftAvailability?.days || {}));
      return 0;
    }

    let result = 0;
    switch (shiftType) {
      case 1: result = day.firstShift?.medicSlotsAvailable || 0; break;
      case 2: result = day.secondShift?.medicSlotsAvailable || 0; break;
      case 3: result = day.thirdShift?.medicSlotsAvailable || 0; break;
      default: result = 0;
    }

    return result;
  }

  getAvailableTechnicals(shiftType: number, date: Date): number {
    const key = date.toISOString().split('T')[0];
    const day = this.shiftAvailability?.days?.[key];

    if (!day) {
      console.warn(`No day data found for key: ${key}. Available keys:`, Object.keys(this.shiftAvailability?.days || {}));
      return 0;
    }

    let result = 0;
    switch (shiftType) {
      case 1: result = day.firstShift?.technicalSlotsAvailable || 0; break;
      case 2: result = day.secondShift?.technicalSlotsAvailable || 0; break;
      case 3: result = day.thirdShift?.technicalSlotsAvailable || 0; break;
      default: result = 0;
    }

    return result;
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
    switch (lang) {
      case 'sr-cyr': return format(date, 'EEEE, d. MMMM yyyy', { locale: sr });
      case 'sr-lat': return format(date, 'EEEE, d. MMMM yyyy', { locale: srLatn });
      case 'ru': return format(date, 'EEEE, d MMMM yyyy', { locale: ru });
      case 'en': return format(date, 'EEEE, MMMM d, yyyy', { locale: enUS });
      default: return format(date, 'EEEE, MMMM d, yyyy', { locale: enUS });
    }
  }

  goBack(): void {
    this.router.navigate(['/']);
  }
}
