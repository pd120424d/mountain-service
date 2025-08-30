import { Component, OnInit } from "@angular/core";
import { AdministratorRole, Employee, MedicRole, EmployeeRole, ShiftAvailabilityResponse, ShiftResponse } from "../shared/models";
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
import { formatISODateLocal } from '../shared/utils/date-time';

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
  selectedTimeSpan = 7; // Default to 1 week
  selectedEmployeeShifts: Map<string, ShiftResponse[]> = new Map(); // Cache for employee shifts

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
    this.loadShiftAvailability();
    this.loadShiftWarnings();
    if (this.userId) {
      this.loadEmployeeShifts(this.userId);
    }
    if (this.userRole === AdministratorRole) {
      this.loadAllEmployees();
    }
  }

  loadShiftAvailability() {
    this.isLoading = true;
    const obs = (this.userRole === AdministratorRole)
      ? this.shiftService.getAdminShiftAvailability(this.selectedTimeSpan)
      : this.shiftService.getShiftAvailability(this.selectedTimeSpan);
    obs.subscribe({
      next: (data) => {
        this.shiftAvailability = data;
        this.dates = Object.keys(data.days || {})
          .map(d => new Date(d))
          .sort((a, b) => a.getTime() - b.getTime());
        this.isLoading = false;
      },
      error: () => {
        this.toastr.error(this.translate.instant('SHIFT_MANAGEMENT.TOAST_LOAD_ERROR') || 'Failed to load shift data');
        this.isLoading = false;
      }
    });
  }

  changeTimeSpan(days: number) {
    if (this.selectedTimeSpan !== days) {
      this.selectedTimeSpan = days;
      this.loadShiftAvailability();
    }
  }

  loadAllEmployees() {
    this.shiftService.getAllEmployees().subscribe(data => {
      this.employees = data;
    });
  }

  loadEmployeeShifts(employeeId: string) {
    if (!employeeId || this.selectedEmployeeShifts.has(employeeId)) {
      return; // Already loaded or no employee selected
    }

    this.shiftService.getEmployeeShifts(employeeId).subscribe({
      next: (shifts) => {
        this.selectedEmployeeShifts.set(employeeId, shifts);
      },
      error: (error) => {
        console.error('Failed to load employee shifts:', error);
        // Don't show error toast as this is not critical for the UI
      }
    });
  }

  onEmployeeSelectionChange(employeeId: string) {
    if (employeeId) {
      this.loadEmployeeShifts(employeeId);
    }
  }

  loadShiftWarnings() {
    if (this.userId) {
      this.shiftService.getShiftWarnings(this.userId).subscribe({
        next: (data) => {
          this.shiftWarnings = data.warnings || [];
        },
        error: () => {
          // Don't show error toast for warnings as they're not critical
        }
      });
    }
  }

  // Helper method to find the correct day data from the API response.
  // Handles the mismatch between date-only keys and full ISO timestamp keys.
  private getDayData(date: Date) {
    const availableKeys = Object.keys(this.shiftAvailability?.days || {});
    const targetDateStr = formatISODateLocal(date); // e.g., "2025-08-17"

    // Find the key that matches our target date (could be "2025-08-17T00:00:00Z" format)
    const key = availableKeys.find(k => k.startsWith(targetDateStr));
    return key ? this.shiftAvailability?.days?.[key] : null;
  }

  private getAssignmentFor(shiftType: number, date: Date, employeeId?: string): ShiftResponse | undefined {
    const id = (this.canModifyOthers() && employeeId) ? employeeId : this.userId;
    if (!id) { return undefined; }
    const shifts = this.selectedEmployeeShifts.get(id);
    if (!shifts) { return undefined; }
    const dateStr = formatISODateLocal(date);
    return shifts.find(s => s.shiftDate === dateStr && s.shiftType === shiftType);
  }

  getAssignmentTime(shiftType: number, date: Date, employeeId?: string): string | null {
    const assignment: any = this.getAssignmentFor(shiftType, date, employeeId);
    const iso = assignment?.assignedAt || assignment?.createdAt;
    if (!iso) { return null; }
    try {
      const dt = new Date(iso);
      const lang = this.translate.currentLang;
      const locale = (lang === 'sr-cyr') ? sr : (lang === 'sr-lat') ? srLatn : (lang === 'ru') ? ru : enUS;
      return format(dt, 'PPpp', { locale });
    } catch {
      return iso as string;
    }
  }

  canModifyOthers(): boolean {
    return this.userRole === AdministratorRole;
  }

  assignToShift(shiftType: number, date: Date, employeeId?: string) {
    const idToAssign = employeeId ?? this.userId;
    if (!idToAssign) return;

    // Show confirmation dialog
    const shiftLabel = this.getShiftLabel(shiftType);
    const dateStr = this.getTranslatedDate(date);
    const confirmMessage = this.translate.instant('SHIFT_MANAGEMENT.CONFIRM_ASSIGN', {
      shift: shiftLabel,
      date: dateStr
    });

    if (!confirm(confirmMessage)) {
      return;
    }

    console.log('Assigning to shift:', { shiftType, date, employeeId: idToAssign });
    this.isAssigning = true;
    this.spinner.show();
    this.shiftService.assignEmployeeToShift(shiftType, idToAssign, date).subscribe({
      next: (response) => {
        console.log('Assignment successful:', response);
        // Reload data to reflect the assignment
        this.loadShiftAvailability();
        this.loadShiftWarnings();
        this.isAssigning = false;
        this.spinner.hide();
        this.toastr.success(this.translate.instant('SHIFT_MANAGEMENT.TOAST_ASSIGN_SUCCESS'));
      },
      error: (error) => {
        console.error('Assignment failed:', error);
        this.isAssigning = false;
        this.spinner.hide();
        const rawErrorMessage = error?.message || this.translate.instant('SHIFT_MANAGEMENT.TOAST_ASSIGN_ERROR');
        const translatedErrorMessage = this.getTranslatedErrorMessage(rawErrorMessage);
        this.toastr.error(translatedErrorMessage);
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

    // Show confirmation dialog
    const shiftLabel = this.getShiftLabel(shiftType);
    const dateStr = this.getTranslatedDate(date);
    const confirmMessage = this.translate.instant('SHIFT_MANAGEMENT.CONFIRM_REMOVE', {
      shift: shiftLabel,
      date: dateStr
    });

    if (!confirm(confirmMessage)) {
      return;
    }

    console.log('Removing from shift:', { shiftType, employeeId: idToRemove, date });

    this.isRemoving = true;
    this.spinner.show();
    this.shiftService.removeEmployeeFromShiftByDetails(idToRemove, shiftType, date).subscribe({
      next: (response) => {
        console.log('Removal successful:', response);
        // Reload data to reflect the removal
        this.loadShiftAvailability();
        this.loadShiftWarnings();
        this.isRemoving = false;
        this.spinner.hide();
        this.toastr.success(this.translate.instant('SHIFT_MANAGEMENT.TOAST_REMOVE_SUCCESS'));
      },
      error: (error) => {
        console.error('Removal failed:', error);
        this.isRemoving = false;
        this.spinner.hide();
        const errorMessage = error?.message || this.translate.instant('SHIFT_MANAGEMENT.TOAST_REMOVE_ERROR');
        this.toastr.error(errorMessage);
      },
      complete: () => {
        this.isRemoving = false;
        this.spinner.hide();
      }
    });
  }

  getAvailableMedics(shiftType: number, date: Date): number {
    const day = this.getDayData(date);

    if (!day) {
      const targetDateStr = formatISODateLocal(date);
      const availableKeys = Object.keys(this.shiftAvailability?.days || {});
      console.warn(`No day data found for date: ${targetDateStr}. Available keys:`, availableKeys);
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

  getAssignedMedics(shiftType: number, date: Date): number {
    const available = this.getAvailableMedics(shiftType, date);
    const assigned = 2 - (available ?? 0);
    return Math.max(0, Math.min(2, assigned));
  }


  getAvailableTechnicals(shiftType: number, date: Date): number {
    const day = this.getDayData(date);

    if (!day) {
      const targetDateStr = formatISODateLocal(date);
      const availableKeys = Object.keys(this.shiftAvailability?.days || {});
      console.warn(`No day data found for date: ${targetDateStr}. Available keys:`, availableKeys);
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

  getAssignedTechnicals(shiftType: number, date: Date): number {
    const available = this.getAvailableTechnicals(shiftType, date);
    const assigned = 4 - (available ?? 0);
    return Math.max(0, Math.min(4, assigned));
  }

  isAssignedToShift(shiftType: number, date: Date, employeeId?: string): boolean {
    // For admin users, check if the specific employee is assigned
    if (this.canModifyOthers() && employeeId) {
      return this.isEmployeeAssignedToShift(employeeId, shiftType, date);
    }

    // For normal users, check if the current user is assigned
    const day = this.getDayData(date);

    if (!day) {
      return false;
    }

    switch (shiftType) {
      case 1: return day.firstShift?.isAssignedToEmployee || false;
      case 2: return day.secondShift?.isAssignedToEmployee || false;
      case 3: return day.thirdShift?.isAssignedToEmployee || false;
      default: return false;
    }
  }

  private isEmployeeAssignedToShift(employeeId: string, shiftType: number, date: Date): boolean {
    const employeeShifts = this.selectedEmployeeShifts.get(employeeId);
    if (!employeeShifts) {
      return false; // No shifts loaded for this employee yet
    }

    const dateStr = formatISODateLocal(date);
    return employeeShifts.some(shift =>
      shift.shiftDate === dateStr && shift.shiftType === shiftType
    );
  }

  isShiftFullyBooked(shiftType: number, date: Date): boolean {
    const day = this.getDayData(date);

    if (!day) {
      return false;
    }

    switch (shiftType) {
      case 1: return day.firstShift?.isFullyBooked || false;
      case 2: return day.secondShift?.isFullyBooked || false;
      case 3: return day.thirdShift?.isFullyBooked || false;
      default: return false;
    }
  }

  isShiftLowCapacity(shiftType: number, date: Date): boolean {
    const medicSlots = this.getAvailableMedics(shiftType, date);
    const technicalSlots = this.getAvailableTechnicals(shiftType, date);

    // Consider low capacity if less than or equal to 1 slot available for either role
    return medicSlots <= 1 || technicalSlots <= 1;
  }

  getTranslatedWarning(warning: string): string {
    // Check if the warning is in the new format with translation key
    if (warning.includes('|')) {
      const parts = warning.split('|');
      if (parts.length === 4) {
        const [translationKey, shiftsCount, daysCount, requiredDays] = parts;
        return this.translate.instant(translationKey, {
          shiftsCount: shiftsCount,
          daysCount: daysCount,
          requiredDays: requiredDays
        });
      }
    }

    // Fallback to original warning text for backward compatibility
    return warning;
  }

  getTranslatedErrorMessage(errorMessage: string): string {
    // Extract the actual error message from the service error format
    // Format: "Server error: 400 - Bad Request - SHIFT_ERRORS.CONSECUTIVE_SHIFTS_LIMIT|7"
    let actualError = errorMessage;

    // Try to extract the backend error message
    const serverErrorMatch = errorMessage.match(/Server error: \d+ - [^-]+ - (.+)$/);
    if (serverErrorMatch) {
      actualError = serverErrorMatch[1];
    }

    // Check if the error is in the new format with translation key
    if (actualError.includes('|')) {
      const parts = actualError.split('|');
      if (parts.length === 2) {
        const [translationKey, consecutiveCount] = parts;
        if (translationKey === 'SHIFT_ERRORS.CONSECUTIVE_SHIFTS_LIMIT') {
          return this.translate.instant(translationKey, {
            consecutiveCount: consecutiveCount
          });
        }
      }
    }

    // Fallback to original error message for backward compatibility
    return errorMessage;
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

  canAssignToShift(shiftType: number, date: Date, employeeId?: string): boolean {
    // Can't assign if currently processing
    if (this.isAssigning || this.isRemoving) {
      return false;
    }

    // For admin view, need employee selected
    if (this.canModifyOthers() && !employeeId) {
      return false;
    }

    // Can't assign if already assigned to this shift
    if (this.isAssignedToShift(shiftType, date, employeeId)) {
      return false;
    }

    // Can't assign if shift is fully booked
    if (this.isShiftFullyBooked(shiftType, date)) {
      return false;
    }

    return true;
  }

  canRemoveFromShift(shiftType: number, date: Date, employeeId?: string): boolean {
    // Can't remove if currently processing
    if (this.isAssigning || this.isRemoving) {
      return false;
    }

    // For admin view, need employee selected
    if (this.canModifyOthers() && !employeeId) {
      return false;
    }

    // Can only remove if the employee (or current user) is assigned to this shift
    if (!this.isAssignedToShift(shiftType, date, employeeId)) {
      return false;
    }

    return true;
  }

  getAssignButtonTooltip(shiftType: number, date: Date, employeeId?: string): string {
    if (this.isAssigning || this.isRemoving) {
      return this.translate.instant('SHIFT_MANAGEMENT.TOOLTIP_PROCESSING');
    }

    if (this.canModifyOthers() && !employeeId) {
      return this.translate.instant('SHIFT_MANAGEMENT.TOOLTIP_SELECT_EMPLOYEE');
    }

    if (this.isAssignedToShift(shiftType, date, employeeId)) {
      return this.translate.instant('SHIFT_MANAGEMENT.TOOLTIP_ALREADY_ASSIGNED');
    }

    if (this.isShiftFullyBooked(shiftType, date)) {
      return this.translate.instant('SHIFT_MANAGEMENT.TOOLTIP_FULLY_BOOKED');
    }

    return this.translate.instant('SHIFT_MANAGEMENT.TOOLTIP_ASSIGN');
  }

  getRemoveButtonTooltip(shiftType: number, date: Date, employeeId?: string): string {
    if (this.isAssigning || this.isRemoving) {
      return this.translate.instant('SHIFT_MANAGEMENT.TOOLTIP_PROCESSING');
    }

    if (this.canModifyOthers() && !employeeId) {
      return this.translate.instant('SHIFT_MANAGEMENT.TOOLTIP_SELECT_EMPLOYEE');
    }

    if (!this.isAssignedToShift(shiftType, date, employeeId)) {
      return this.translate.instant('SHIFT_MANAGEMENT.TOOLTIP_NOT_ASSIGNED');
    }

    return this.translate.instant('SHIFT_MANAGEMENT.TOOLTIP_REMOVE');
  }

  shouldShowAssignedIndicator(shiftType: number, day: Date, selectedEmployeeId?: string): boolean {
    // For admin users with an employee selected, check if that employee is assigned
    if (this.canModifyOthers() && selectedEmployeeId) {
      return this.isAssignedToShift(shiftType, day, selectedEmployeeId);
    }

    // For normal users or admin without employee selected, check if current user is assigned
    return this.isAssignedToShift(shiftType, day);
  }

  shouldApplyAssignedClass(shiftType: number, day: Date, selectedEmployeeId?: string): boolean {
    // For admin users with an employee selected, apply assigned class if that employee is assigned
    if (this.canModifyOthers() && selectedEmployeeId) {
      return this.isAssignedToShift(shiftType, day, selectedEmployeeId);
    }

    // For normal users or admin without employee selected, apply assigned class if current user is assigned
    return this.isAssignedToShift(shiftType, day);
  }

  goBack(): void {
    this.router.navigate(['/']);
  }
}
