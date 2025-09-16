import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { catchError } from 'rxjs/operators';
import { of } from 'rxjs';

@Component({
  selector: 'app-admin',
  standalone: true,
  imports: [CommonModule, FormsModule, TranslateModule],
  templateUrl: './admin.component.html',
  styleUrls: ['./admin.component.css']
})
export class AdminComponent {
  // Reset state
  isResetting = false;
  resetSuccess = false;
  resetError = false;
  errorMessage = '';

  // Restart state
  services = [
    { id: 'employee-service', labelKey: 'ADMIN_PANEL.SERVICES.EMPLOYEE' },
    { id: 'urgency-service',  labelKey: 'ADMIN_PANEL.SERVICES.URGENCY' },
    { id: 'activity-service', labelKey: 'ADMIN_PANEL.SERVICES.ACTIVITY' },
    { id: 'version-service',  labelKey: 'ADMIN_PANEL.SERVICES.VERSION' },
    { id: 'docs-aggregator',  labelKey: 'ADMIN_PANEL.SERVICES.DOCS_AGGREGATOR' },
    { id: 'docs-ui',          labelKey: 'ADMIN_PANEL.SERVICES.DOCS_UI' },
  ];
  selectedDeployment = '';
  showRestartModal = false;
  isRestarting = false;
  restartSuccess = false;
  restartError = false;
  restartMessage = '';

  constructor(
    private authService: AuthService,
    private router: Router,
    private translate: TranslateService
  ) {
    this.translate.setDefaultLang('sr-cyr');
    if (!this.authService.isAdmin()) {
      this.router.navigate(['/']);
    }
  }

  onResetAllData(): void {
    if (!confirm(this.translate.instant('RESET_WARNING'))) {
      return;
    }

    this.isResetting = true;
    this.resetSuccess = false;
    this.resetError = false;
    this.errorMessage = '';

    this.authService.resetAllData().pipe(
      catchError((error) => {
        console.error('Reset failed:', error);
        this.resetError = true;
        this.errorMessage = error.error?.error || this.translate.instant('RESET_ERROR');
        this.isResetting = false;
        return of(null);
      })
    ).subscribe((response) => {
      this.isResetting = false;
      if (response) {
        this.resetSuccess = true;
        console.log('Reset successful:', response.message);
      }
    });
  }

  openRestartModal(): void {
    if (!this.selectedDeployment) return;
    this.restartError = false;
    this.restartSuccess = false;
    this.restartMessage = '';
    this.showRestartModal = true;
  }

  cancelRestart(): void {
    this.showRestartModal = false;
  }

  confirmRestart(): void {
    if (!this.selectedDeployment) return;
    this.isRestarting = true;
    this.restartError = false;
    this.restartSuccess = false;
    this.restartMessage = '';

    this.authService.restartService(this.selectedDeployment).pipe(
      catchError((error) => {
        console.error('Restart failed:', error);
        this.restartError = true;
        this.restartMessage = error.error?.error || this.translate.instant('ADMIN_PANEL.RESTART.ERROR');
        this.isRestarting = false;
        return of(null);
      })
    ).subscribe((response) => {
      this.isRestarting = false;
      this.showRestartModal = false;
      if (response) {
        this.restartSuccess = true;
        this.restartMessage = response.message || this.translate.instant('ADMIN_PANEL.RESTART.SUCCESS');
      }
    });
  }

  goBack(): void {
    this.router.navigate(['/']);
  }
}
