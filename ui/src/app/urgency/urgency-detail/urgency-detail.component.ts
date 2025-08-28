import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, ActivatedRoute, Router } from '@angular/router';
import { ReactiveFormsModule, FormBuilder, FormGroup, Validators } from '@angular/forms';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { ToastrService } from 'ngx-toastr';
import { BaseTranslatableComponent } from '../../base-translatable.component';
import { UrgencyService } from '../urgency.service';
import { ActivityService } from '../../services/activity.service';
import { AuthService } from '../../services/auth.service';
import { EmployeeService } from '../../employee/employee.service';
import {
  Urgency,
  Activity,
  UrgencyStatus,
  createUrgencyDisplayName,
  getUrgencyLevelColor,
  getUrgencyStatusColor,
  getActivityIcon,
  getActivityDisplayTime,
  hasAcceptedAssignment
} from '../../shared/models';

@Component({
  selector: 'app-urgency-detail',
  standalone: true,
  imports: [RouterModule, TranslateModule, CommonModule, ReactiveFormsModule],
  templateUrl: './urgency-detail.component.html',
  styleUrls: ['./urgency-detail.component.css']
})
export class UrgencyDetailComponent extends BaseTranslatableComponent implements OnInit {
  urgency: Urgency | null = null;
  activities: Activity[] = [];
  activityForm!: FormGroup;
  isLoading = true;
  isLoadingActivities = false;
  isSubmittingActivity = false;
  error: string | null = null;
  urgencyId: number | null = null;

  employeeNames: Record<number, string> = {};
  private fetchingEmployeeIds = new Set<number>();

  UrgencyStatus = UrgencyStatus;
  isAssigning = false;
  isUnassigning = false;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private fb: FormBuilder,
    private urgencyService: UrgencyService,
    private activityService: ActivityService,
    private authService: AuthService,
    private toastr: ToastrService,
    private employeeService: EmployeeService,
    translate: TranslateService
  ) {
    super(translate);
  }

  ngOnInit(): void {
    this.initializeActivityForm();
    this.route.params.subscribe(params => {
      this.urgencyId = +params['id'];
      if (this.urgencyId) {
        this.loadUrgency();
        this.loadActivities();
      }
    });
  }

  initializeActivityForm(): void {
    this.activityForm = this.fb.group({
      description: ['', [Validators.required, Validators.minLength(10)]]
    });
  }

  loadUrgency(): void {
    if (!this.urgencyId) return;

    this.isLoading = true;
    this.error = null;

    this.urgencyService.getUrgencyById(this.urgencyId).subscribe({
      next: (urgency) => {
        this.urgency = urgency;
        this.isLoading = false;
        const assignedId = (urgency as any)?.assignedEmployeeId as number | undefined;
        if (assignedId && !this.employeeNames[assignedId]) {
          this.fetchAssignedEmployeeName(assignedId);
        }
      },
      error: (error) => {
        this.error = error.message || 'Failed to load urgency details';
        this.isLoading = false;
      }
    });
  }

  loadActivities(): void {
    if (!this.urgencyId) return;

    this.isLoadingActivities = true;

    this.activityService.getActivitiesByUrgency(this.urgencyId).subscribe({
      next: (activities) => {
        this.activities = activities.sort((a, b) =>
          new Date(b.createdAt || '').getTime() - new Date(a.createdAt || '').getTime()
        );
        this.isLoadingActivities = false;
      },
      error: (error) => {
        console.error('Failed to load activities:', error);
        this.isLoadingActivities = false;
      }
    });
  }

  onSubmitActivity(): void {
    if (this.activityForm.valid && !this.isSubmittingActivity && this.urgencyId) {
      this.isSubmittingActivity = true;

      const formValue = this.activityForm.value;
      const activityRequest = {
        urgencyId: this.urgencyId,
        employeeId: parseInt(this.authService.getUserId()),
        description: formValue.description
      };

      this.activityService.createActivity(activityRequest).subscribe({
        next: () => {
          this.toastr.success(this.translate.instant('URGENCY_DETAIL.ACTIVITY_ADDED_SUCCESS'));
          this.activityForm.reset();
          this.loadActivities();
          this.isSubmittingActivity = false;
        },
        error: (error) => {
          console.error('Error creating activity:', error);
          this.toastr.error(this.translate.instant('URGENCY_DETAIL.ACTIVITY_ADDED_ERROR'));
          this.isSubmittingActivity = false;
        }
      });
    } else {
      this.markFormGroupTouched();
    }
  }

  goBack(): void {
    this.router.navigate(['/urgencies']);
  }

  canAssign(): boolean {
    if (!this.urgency) return false;
    const assigned = hasAcceptedAssignment(this.urgency);
    return !assigned && (this.urgency.status === UrgencyStatus.Open || this.urgency.status === UrgencyStatus.InProgress);
  }

  canUnassign(): boolean {
    if (!this.urgency) return false;
    const assigned = hasAcceptedAssignment(this.urgency);
    return assigned && this.authService.isAdmin();
  }

  onAssignToMe(): void {
    if (!this.urgencyId || this.isAssigning) return;
    this.isAssigning = true;
    const employeeId = parseInt(this.authService.getUserId());
    this.urgencyService.assignUrgency(this.urgencyId, employeeId).subscribe({
      next: () => {
        this.toastr.success(this.translate.instant('URGENCY_DETAIL.ASSIGN_SUCCESS'));
        this.isAssigning = false;
        this.loadUrgency();
      },
      error: (err) => {
        this.toastr.error(err?.message || this.translate.instant('URGENCY_DETAIL.ASSIGN_ERROR'));
        this.isAssigning = false;
      }
    });
  }

  onUnassign(): void {
    if (!this.urgencyId || this.isUnassigning) return;
    this.isUnassigning = true;
    this.urgencyService.unassignUrgency(this.urgencyId).subscribe({
      next: () => {
        this.toastr.success(this.translate.instant('URGENCY_DETAIL.UNASSIGN_SUCCESS'));
        this.isUnassigning = false;
        this.loadUrgency();
      },
      error: (err) => {
        this.toastr.error(err?.message || this.translate.instant('URGENCY_DETAIL.UNASSIGN_ERROR'));
        this.isUnassigning = false;
      }
    });
  }

  getDisplayName(urgency: Urgency): string {
    return createUrgencyDisplayName(urgency);
  }

  getUrgencyLevelColor(level: string): string {
    return getUrgencyLevelColor(level as any);
  }

  getUrgencyStatusColor(status: string): string {
    return getUrgencyStatusColor(status as any);
  }

  getActivityIcon(activity: Activity): string {
    return getActivityIcon(activity);
  }

  getActivityDisplayTime(activity: Activity): string {
    return getActivityDisplayTime(activity);
  }

  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString();
  }

  isFieldInvalid(fieldName: string): boolean {
    const field = this.activityForm.get(fieldName);
    return !!(field && field.invalid && (field.dirty || field.touched));
  }

  getFieldError(fieldName: string): string {
    const field = this.activityForm.get(fieldName);
    if (field && field.errors) {
      if (field.errors['required']) {
        return `${fieldName.toUpperCase()}_FORM.REQUIRED`;
      }
      if (field.errors['minlength']) {
        return `${fieldName.toUpperCase()}_FORM.MIN_LENGTH`;
      }
    }
    return '';
  }

  getEmployeeName(employeeId: number): string {
    const cached = this.employeeNames[employeeId];
    if (cached) return cached;

    if (!this.fetchingEmployeeIds.has(employeeId)) {
      this.fetchingEmployeeIds.add(employeeId);
      this.employeeService.getEmployeeById(employeeId).subscribe({
        next: (emp) => {
          this.employeeNames[employeeId] = `${emp.firstName} ${emp.lastName}`;
          this.fetchingEmployeeIds.delete(employeeId);
        },
        error: () => {
          this.employeeNames[employeeId] = this.translate.instant('URGENCY_DETAIL.EMPLOYEE_PLACEHOLDER', { id: employeeId });
          this.fetchingEmployeeIds.delete(employeeId);
        }
      });
    }
    return this.translate.instant('URGENCY_DETAIL.EMPLOYEE_PLACEHOLDER', { id: employeeId });
  }

  private fetchAssignedEmployeeName(employeeId: number): void {
    if (this.fetchingEmployeeIds.has(employeeId)) return;
    this.fetchingEmployeeIds.add(employeeId);
    this.employeeService.getEmployeeById(employeeId).subscribe({
      next: (emp) => {
        this.employeeNames[employeeId] = `${emp.firstName} ${emp.lastName}`;
        this.fetchingEmployeeIds.delete(employeeId);
      },
      error: () => {
        this.employeeNames[employeeId] = this.translate.instant('URGENCY_DETAIL.EMPLOYEE_PLACEHOLDER', { id: employeeId });
        this.fetchingEmployeeIds.delete(employeeId);
      }
    });
  }

  private markFormGroupTouched(): void {
    Object.keys(this.activityForm.controls).forEach(key => {
      const control = this.activityForm.get(key);
      control?.markAsTouched();
    });
  }
}
