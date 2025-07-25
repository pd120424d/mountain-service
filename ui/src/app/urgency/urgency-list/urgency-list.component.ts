import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { BaseTranslatableComponent } from '../../base-translatable.component';
import { UrgencyService } from '../urgency.service';
import { Urgency } from '../../shared/models';
import {
  GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyLevel as GeneratedUrgencyLevel,
  GithubComPd120424DMountainServiceApiContractsUrgencyV1UrgencyStatus as GeneratedUrgencyStatus
} from '../../shared/models/generated/urgency';

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

  UrgencyLevel = GeneratedUrgencyLevel;
  Status = GeneratedUrgencyStatus;

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
      case GeneratedUrgencyLevel.Low:
        return 'level-low';
      case GeneratedUrgencyLevel.Medium:
        return 'level-medium';
      case GeneratedUrgencyLevel.High:
        return 'level-high';
      case GeneratedUrgencyLevel.Critical:
        return 'level-critical';
      default:
        return '';
    }
  }

  formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString();
  }
}
