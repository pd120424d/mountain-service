import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { BaseTranslatableComponent } from '../../base-translatable.component';
import { UrgencyService } from '../urgency.service';
import { Urgency, UrgencyLevel, Status } from '../urgency.model';

@Component({
  selector: 'app-urgency-list',
  standalone: true,
  imports: [RouterModule, TranslateModule, CommonModule],
  templateUrl: './urgency-list.component.html',
  styleUrls: ['./urgency-list.component.css']
})
export class UrgencyListComponent extends BaseTranslatableComponent implements OnInit {
  urgencies: Urgency[] = [];
  isLoading = true;
  error: string | null = null;

  UrgencyLevel = UrgencyLevel;
  Status = Status;

  constructor(
    private urgencyService: UrgencyService,
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
        this.urgencies = urgencies;
        this.isLoading = false;
      },
      error: (error) => {
        console.error('Error loading urgencies:', error);
        this.error = error.message || 'Failed to load urgencies';
        this.isLoading = false;
      }
    });
  }

  viewUrgency(id: number): void {
    // TODO: Implement view urgency
    console.log('Viewing urgency:', id);
  }

  getStatusClass(status: Status): string {
    switch (status) {
      case Status.OPEN:
        return 'status-open';
      case Status.IN_PROGRESS:
        return 'status-in-progress';
      case Status.RESOLVED:
        return 'status-resolved';
      case Status.CLOSED:
        return 'status-closed';
      default:
        return '';
    }
  }

  getLevelClass(level: UrgencyLevel): string {
    switch (level) {
      case UrgencyLevel.LOW:
        return 'level-low';
      case UrgencyLevel.MEDIUM:
        return 'level-medium';
      case UrgencyLevel.HIGH:
        return 'level-high';
      case UrgencyLevel.CRITICAL:
        return 'level-critical';
      default:
        return '';
    }
  }

  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString();
  }
}
