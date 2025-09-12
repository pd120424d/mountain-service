import { Component, OnInit, AfterViewInit, OnDestroy, ViewChild, ElementRef, HostListener } from '@angular/core';
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
import { DateTimeModule } from '../../shared/utils/date-time.module';
import { ResolveModalComponent, ResolveModalResult } from '../resolve-modal/resolve-modal.component';

@Component({
  selector: 'app-urgency-detail',
  standalone: true,
  imports: [RouterModule, TranslateModule, CommonModule, ReactiveFormsModule, DateTimeModule, ResolveModalComponent],
  templateUrl: './urgency-detail.component.html',
  styleUrls: ['./urgency-detail.component.css']
})
export class UrgencyDetailComponent extends BaseTranslatableComponent implements OnInit, AfterViewInit, OnDestroy {
  @ViewChild('activityTextarea') activityTextarea!: ElementRef<HTMLTextAreaElement>;
  @ViewChild('loadMoreAnchor') loadMoreAnchor!: ElementRef<HTMLDivElement>;

  urgency: Urgency | null = null;
  activities: Activity[] = [];
  activityForm!: FormGroup;
  isLoading = true;
  isLoadingActivities = false;
  isLoadingMore = false;
  isSubmittingActivity = false;
  isSyncingActivity = false;
  nextPageToken: string | null = null;
  private intersectionObserver?: IntersectionObserver;
  // Pagination (legacy, retained for compatibility)
  activitiesPage = 1;
  activitiesPageSize = 10;
  totalActivities = 0;
  totalActivitiesPages = 0;

  error: string | null = null;
  urgencyId: number | null = null;

  employeeNames: Record<number, string> = {};
  private fetchingEmployeeIds = new Set<number>();

  UrgencyStatus = UrgencyStatus;
  isAssigning = false;
  isUnassigning = false;
  isResolving = false;
  showResolveModal = false;
  showScrollTop = false;

  @HostListener('window:scroll', [])
  onWindowScroll(): void {
    const offset = window.pageYOffset || document.documentElement.scrollTop || document.body.scrollTop || 0;
    this.showScrollTop = offset > 200;
  }

  scrollToTop(): void {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }


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
    this.nextPageToken = null;
    this.activities = [];

    console.debug('[Activities] Fetch first page', { urgencyId: this.urgencyId, pageSize: this.activitiesPageSize });
    this.activityService.getActivitiesCursor({
      urgencyId: this.urgencyId,
      pageSize: this.activitiesPageSize
    }).subscribe({
      next: (resp) => {
        this.activities = (resp.activities || []);
        this.nextPageToken = resp.nextPageToken || null;
        this.isLoadingActivities = false;
        console.debug('[Activities] First page loaded', { count: this.activities.length, nextPageToken: this.nextPageToken });
        // Defer to allow Angular to render the #loadMoreAnchor (*ngIf depends on nextPageToken)
        setTimeout(() => this.setupIntersectionObserver());
      },
      error: (error) => {
        console.error('[Activities] Failed to load first page', error);
        this.isLoadingActivities = false;
      }
    });
  }

  ngAfterViewInit(): void {
    this.setupIntersectionObserver();
  }

  ngOnDestroy(): void {
    if (this.intersectionObserver) {
      this.intersectionObserver.disconnect();
    }
  }

  private setupIntersectionObserver(): void {
    if (!this.loadMoreAnchor) {
      console.warn('[Activities] loadMoreAnchor not available yet; will retry on next change detection');
      return;
    }
    if (this.intersectionObserver) {
      this.intersectionObserver.disconnect();
    }
    console.debug('[Activities] Setting up IntersectionObserver', { hasToken: !!this.nextPageToken });
    this.intersectionObserver = new IntersectionObserver(entries => {
      const entry = entries[0];
      if (entry) {
        console.debug('[Activities] Intersection observed', { isIntersecting: entry.isIntersecting, hasToken: !!this.nextPageToken, isLoadingMore: this.isLoadingMore });
      }
      if (entry && entry.isIntersecting && this.nextPageToken && !this.isLoadingMore) {
        this.loadMoreActivities();
      }
    }, {
      root: null,
      rootMargin: '200px 0px',
      threshold: 0
    });
    this.intersectionObserver.observe(this.loadMoreAnchor.nativeElement);
  }

  private loadMoreActivities(): void {
    if (!this.urgencyId || !this.nextPageToken) return;
    this.isLoadingMore = true;
    const outgoing = this.nextPageToken;
    console.debug('[Activities] Loading more', { urgencyId: this.urgencyId, pageSize: this.activitiesPageSize, pageToken: outgoing });
    this.activityService.getActivitiesCursor({
      urgencyId: this.urgencyId,
      pageSize: this.activitiesPageSize,
      pageToken: outgoing || undefined
    }).subscribe({
      next: (resp) => {
        const more = resp.activities || [];
        this.activities = this.activities.concat(more);
        this.nextPageToken = resp.nextPageToken || null;
        this.isLoadingMore = false;
        console.debug('[Activities] Loaded more', { appended: more.length, nextPageToken: this.nextPageToken });
      },
      error: (err) => {
        console.error('[Activities] Failed to load more', err);
        this.isLoadingMore = false;
      }
    });
  }


  onSubmitActivity(): void {
    if (!this.canAddActivity()) {
      this.toastr.warning('Activities can only be added when the urgency is assigned and in progress.');
      return;
    }

    if (this.activityForm.valid && !this.isSubmittingActivity && this.urgencyId) {
      this.isSubmittingActivity = true;

      const formValue = this.activityForm.value;
      const activityRequest = {
        urgencyId: this.urgencyId,
        employeeId: parseInt(this.authService.getUserId()),
        description: formValue.description
      };

      this.activityService.createActivity(activityRequest).subscribe({
        next: (created) => {
          this.toastr.success(this.translate.instant('URGENCY_DETAIL.ACTIVITY_ADDED_SUCCESS'));
          this.activityForm.reset();
          this.isSubmittingActivity = false;
          // Focus the textarea for easier consecutive activity additions
          setTimeout(() => {
            this.activityTextarea?.nativeElement?.focus();
          }, 100);

          // Poll in background until the created id is visible in the read model, then refresh
          if (created?.id && this.urgencyId) {
            this.isSyncingActivity = true;
            this.activityService
              .pollForActivityInUrgency(this.urgencyId, created.id, { timeoutMs: 10000, intervalMs: 300 })
              .subscribe({
                next: () => { this.isSyncingActivity = false; this.loadActivities(); },
                error: () => { this.isSyncingActivity = false; this.loadActivities(); },
              });
          } else {
            this.loadActivities();
          }
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

  canAddActivity(): boolean {
    if (!this.urgency) return false;
    const assigned = hasAcceptedAssignment(this.urgency);
    return assigned && this.urgency.status === UrgencyStatus.InProgress;
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

  canResolve(): boolean {
    if (!this.urgency) return false;
    const currentUserId = parseInt(this.authService.getUserId());
    const isAssignedToCurrentUser = this.urgency.assignedEmployeeId === currentUserId;
    const isAdmin = this.authService.isAdmin();
    const isNotResolved = this.urgency.status !== UrgencyStatus.Resolved && this.urgency.status !== UrgencyStatus.Closed;

    return isNotResolved && (isAssignedToCurrentUser || isAdmin);
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

  onResolve(): void {
    this.showResolveModal = true;
  }

  onResolveModalResult(result: ResolveModalResult): void {
    if (!result.confirmed || !this.urgencyId || this.isResolving) return;

    this.isResolving = true;

    // If user wants to add an activity, create it first
    if (result.activityDescription) {
      const activityRequest = {
        urgencyId: this.urgencyId,
        employeeId: parseInt(this.authService.getUserId()),
        description: result.activityDescription
      };

      this.activityService.createActivity(activityRequest).subscribe({
        next: () => {
          // Activity created successfully, now close the urgency
          this.closeUrgency();
        },
        error: (error) => {
          console.error('Error creating activity:', error);
          this.toastr.error(this.translate.instant('URGENCY_DETAIL.ACTIVITY_ADDED_ERROR'));
          this.isResolving = false;
        }
      });
    } else {
      // No activity to add, directly close the urgency
      this.closeUrgency();
    }
  }

  private closeUrgency(): void {
    if (!this.urgencyId) return;

    this.urgencyService.closeUrgency(this.urgencyId).subscribe({
      next: () => {
        this.toastr.success(this.translate.instant('URGENCY_DETAIL.RESOLVE_SUCCESS'));
        this.isResolving = false;
        this.loadUrgency();
        this.loadActivities(); // Refresh activities if we added one
      },
      error: (err) => {
        this.toastr.error(err?.message || this.translate.instant('URGENCY_DETAIL.RESOLVE_ERROR'));
        this.isResolving = false;
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


  // Activity pager handlers
  onPrevActivities(): void {
    if (this.activitiesPage > 1) {
      this.activitiesPage--;
      this.loadActivities();
    }
  }
  onNextActivities(): void {
    if (this.activitiesPage < this.totalActivitiesPages) {
      this.activitiesPage++;
      this.loadActivities();
    }
  }

  private markFormGroupTouched(): void {
    Object.keys(this.activityForm.controls).forEach(key => {
      const control = this.activityForm.get(key);
      control?.markAsTouched();
    });
  }
}
