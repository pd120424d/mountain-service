import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, Router } from '@angular/router';
import { DateTimeModule } from '../../shared/utils/date-time.module';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { BaseTranslatableComponent } from '../../base-translatable.component';
import { UrgencyService } from '../urgency.service';
import { AuthService } from '../../services/auth.service';
import { ActivityService } from '../../services/activity.service';
import { finalize } from 'rxjs/operators';

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
  page = 1;
  pageSize = 20;
  total = 0;
  totalPages = 0;
  activeTab: 'mine' | 'all' = 'mine';

  countsLoading = false;
  countsError = false;
  countsByUrgencyId: Record<number, number> = {};

  UrgencyLevel = UrgencyLevel;
  Status = UrgencyStatus;

  get unassignedCount(): number {
    return (this.urgencies || []).filter(u => !hasAcceptedAssignment(u as any)).length;
  }

  constructor(
    private urgencyService: UrgencyService,
    private router: Router,
    public authService: AuthService,
    private activityService: ActivityService,
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

    const myUrgencies = this.activeTab === 'mine';
    this.urgencyService.getUrgenciesPaginated({ page: this.page, pageSize: this.pageSize, myUrgencies }).subscribe({
      next: (resp) => {
        const items = resp?.urgencies || [];
        this.urgencies = items;
        this.total = resp?.total ?? items.length;
        this.page = resp?.page ?? this.page;
        this.pageSize = resp?.pageSize ?? this.pageSize;
        this.totalPages = resp?.totalPages ?? Math.ceil(this.total / this.pageSize);
        this.isLoading = false;
        this.fetchActivityCounts();
      },
      error: (error) => {
        this.error = error.message || 'Failed to load urgencies';
        this.isLoading = false;
        this.countsByUrgencyId = {};
      }
    });
  }

  private fetchActivityCounts(): void {
    const ids = (this.urgencies || []).map(u => u.id!).filter((v): v is number => typeof v === 'number');
    if (!ids.length) { this.countsByUrgencyId = {}; this.countsLoading = false; this.countsError = false; return; }
    this.countsLoading = true;
    this.countsError = false;
    const obs = (this.activityService as any)?.getCountsByUrgencyIds?.(ids);
    if (!obs || typeof obs.pipe !== 'function') { this.countsLoading = false; this.countsByUrgencyId = {}; this.countsError = true; return; }
    obs
      .pipe(finalize(() => { this.countsLoading = false; }))
      .subscribe({
        next: (counts: Record<string, number>) => {
          const mapped: Record<number, number> = {};
          Object.entries(counts || {}).forEach(([k, v]) => { mapped[parseInt(k, 10)] = Number(v) || 0; });
          this.countsByUrgencyId = mapped;
          this.countsError = false;
        },
        error: () => { this.countsByUrgencyId = {}; this.countsError = true; }
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

  // Simple pager handlers
  onPrev(): void {
    if (this.page > 1) {
      this.page--;
      this.loadUrgencies();
    }
  }
  onNext(): void {
    if (this.page < this.totalPages) {
      this.page++;
      this.loadUrgencies();
    }
  }

  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString();
  }
}
