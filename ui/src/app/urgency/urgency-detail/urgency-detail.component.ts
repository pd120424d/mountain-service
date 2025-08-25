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
import {
  Urgency,
  Activity,
  ActivityCreateRequest,
  UrgencyStatus,
  createUrgencyDisplayName,
  getUrgencyLevelColor,
  getUrgencyStatusColor,
  getActivityIcon,
  getActivityDisplayTime
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

  UrgencyStatus = UrgencyStatus;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private fb: FormBuilder,
    private urgencyService: UrgencyService,
    private activityService: ActivityService,
    private authService: AuthService,
    private toastr: ToastrService,
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
          new Date(b.created_at || '').getTime() - new Date(a.created_at || '').getTime()
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
      const activityRequest: ActivityCreateRequest = {
        urgency_id: this.urgencyId,
        employee_id: parseInt(this.authService.getUserId()),
        description: formValue.description
      };

      this.activityService.createActivity(activityRequest).subscribe({
        next: (activity) => {
          this.toastr.success(this.translate.instant('URGENCY_DETAIL.ACTIVITY_ADDED_SUCCESS'));
          this.activityForm.reset();
          this.loadActivities(); // Reload activities
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

  private markFormGroupTouched(): void {
    Object.keys(this.activityForm.controls).forEach(key => {
      const control = this.activityForm.get(key);
      control?.markAsTouched();
    });
  }
}
