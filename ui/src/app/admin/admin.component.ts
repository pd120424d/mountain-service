import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AuthService } from '../services/auth.service';
import { Router } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { catchError } from 'rxjs/operators';
import { of } from 'rxjs';

@Component({
  selector: 'app-admin',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  templateUrl: './admin.component.html',
  styleUrls: ['./admin.component.css']
})
export class AdminComponent {
  isResetting = false;
  resetSuccess = false;
  resetError = false;
  errorMessage = '';

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

  goBack(): void {
    this.router.navigate(['/']);
  }
}
