import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router } from '@angular/router';
import { DateTimeModule } from '../../shared/utils/date-time.module';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { BaseTranslatableComponent } from '../../base-translatable.component';
import { UrgencyService } from '../urgency.service';
import { AuthService } from '../../services/auth.service';
import { Urgency, UrgencyLevel, UrgencyStatus, UrgencyStatus as GeneratedUrgencyStatus, UrgencyLevel as GeneratedUrgencyLevel, createUrgencyDisplayName, hasAcceptedAssignment } from '../../shared/models';

@Component({
  selector: 'app-urgency-list',
  standalone: true,
  imports: [RouterModule, TranslateModule, CommonModule, DateTimeModule],
  templateUrl: './urgency-list.component.html',
  styleUrls: ['./urgency-list.component.css']
})
export class UrgencyListComponent extends BaseTranslatableComponent implements OnInit {
  urgencies: Urgency[] = [];
  isLoading = true;
  error: string | null = null;

  UrgencyLevel = UrgencyLevel;
  Status = UrgencyStatus;

  get unassignedCount(): number {
    return (this.urgencies || []).filter(u => !hasAcceptedAssignment(u as any)).length;
  }

  constructor(
    private urgencyService: UrgencyService,
    private router: Router,
    private authService: AuthService,
    translate: TranslateService
  ) {
    super(translate);
  }

  ngOnInit(): void {
    this.loadUrgencies();
  }

  loadUrgencies(): void {
    this.isLoading = true;
    this.error = null;

    this.urgencyService.getUrgencies().subscribe({
      next: (urgencies) => {
        // Sort urgencies by createdAt date in descending order
        this.urgencies = (urgencies || []).sort((a, b) => {
          const dateA = new Date(a.createdAt || '').getTime();
          const dateB = new Date(b.createdAt || '').getTime();
          return dateB - dateA; // Descending order
        });
        this.isLoading = false;
      },
      error: (error) => {
        this.error = error.message || 'Failed to load urgencies';
        this.isLoading = false;
      }
    });
  }

  viewUrgency(id: number): void {
    this.router.navigate(['/urgencies', id]);
  }

  getDisplayName(urgency: Urgency): string {
    return createUrgencyDisplayName(urgency);
  }

  getStatusClass(status: GeneratedUrgencyStatus): string {
    switch (status) {
      case GeneratedUrgencyStatus.Open:
        return 'status-open';
      case GeneratedUrgencyStatus.InProgress:
        return 'status-in-progress';
      case GeneratedUrgencyStatus.Resolved:
        return 'status-resolved';
      case GeneratedUrgencyStatus.Closed:
        return 'status-closed';
      default:
        return '';
    }
  }

  getLevelClass(level: GeneratedUrgencyLevel): string {
    switch (level) {
      case UrgencyLevel.Low:
        return 'level-low';
      case UrgencyLevel.Medium:
        return 'level-medium';
      case UrgencyLevel.High:
        return 'level-high';
      case UrgencyLevel.Critical:
        return 'level-critical';
      default:
        return '';
    }
  }

  getRowClass(urgency: Urgency): string {
    const currentUserId = parseInt(this.authService.getUserId() || '0', 10);
    const isAssigned = hasAcceptedAssignment(urgency);
    const assignedToMe = isAssigned && urgency.assignedEmployeeId === currentUserId;

    if (urgency.status === GeneratedUrgencyStatus.Closed || urgency.status === GeneratedUrgencyStatus.Resolved) {
      return 'urgency-row-closed';
    }

    if (urgency.status === GeneratedUrgencyStatus.Open && !isAssigned) {
      return 'urgency-row-open-unassigned';
    }

    if (urgency.status === GeneratedUrgencyStatus.InProgress) {
      return assignedToMe ? 'urgency-row-assigned-me' : 'urgency-row-assigned-other';
    }

    return '';
  }

  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString();
  }

  hasAcceptedAssignment(urgency: Urgency): boolean {
    return hasAcceptedAssignment(urgency);
  }
}
