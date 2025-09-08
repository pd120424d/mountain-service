import { Component, HostListener, OnInit, OnDestroy } from '@angular/core';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { AuthService } from './services/auth.service';
import { EmployeeService } from './employee/employee.service';
import { AppInitializationService } from './services/app-initialization.service';
import { CommonModule } from '@angular/common';
import { RouterModule, RouterOutlet, Router } from '@angular/router';
import { VersionBannerComponent } from './version-banner/version-banner.component';
import { NgxSpinnerModule } from 'ngx-spinner';
import { Employee, hasAcceptedAssignment } from './shared/models';
import { UrgencyService } from './urgency/urgency.service';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [RouterModule, CommonModule, RouterOutlet, TranslateModule, VersionBannerComponent, NgxSpinnerModule],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
})
export class AppComponent implements OnInit, OnDestroy {
  isDropdownOpen = false;
  currentLanguageLabel = 'EN';
  currentUser: Employee | null = null;

  openUrgenciesCount: number = 0;

  private languageMap: { [key: string]: string } = {
    'en': 'EN',
    'sr-lat': 'SR',
    'sr-cyr': 'СР',
    'ru': 'RU'
  };

  constructor(
    public authService: AuthService,
    private employeeService: EmployeeService,
    private urgencyService: UrgencyService,
    private translate: TranslateService,
    private appInitService: AppInitializationService,
    private router: Router
  ) {
    const savedLanguage = localStorage.getItem('language') || 'en';
    this.translate.use(savedLanguage);
    this.setLanguageLabel(savedLanguage);
  }

  ngOnInit(): void {
    // Initialize the application
    this.appInitService.initialize().then(() => {
      console.log('Application initialized successfully');
      this.loadCurrentUser();
      this.refreshOpenUrgencies();
    }).catch((error) => {
      console.error('Failed to initialize application:', error);
    });

    // Update header when auth state changes (e.g., after login/logout)
    this.authService.authChanged$.subscribe(() => {
      this.loadCurrentUser();
      this.refreshOpenUrgencies();
    });
  }

  ngOnDestroy(): void {
    // Cleanup resources
    this.appInitService.cleanup();
  }

  toggleDropdown(): void {
    this.isDropdownOpen = !this.isDropdownOpen;
  }

  switchLanguage(language: string): void {
    this.translate.use(language);
    localStorage.setItem('language', language);
    this.setLanguageLabel(language);
    this.isDropdownOpen = false;
  }

  setLanguageLabel(language: string): void {
    this.currentLanguageLabel = this.languageMap[language] || language.toUpperCase();
  }

  @HostListener('document:click', ['$event'])
  closeDropdown(event: Event): void {
    const dropdown = document.querySelector('.language-switcher');
    if (dropdown && !dropdown.contains(event.target as Node)) {
      this.isDropdownOpen = false;
    }
  }

  private loadCurrentUser(): void {
    if (this.authService.isAuthenticated()) {
      const userId = this.authService.getUserId();
      if (userId) {
        this.employeeService.getEmployeeById(parseInt(userId)).subscribe({
          next: (user) => {
            this.currentUser = user;
          },
          error: (error) => {
            console.error('Error loading current user:', error);
            this.currentUser = null;
          }
        });
      }
    } else {
      this.currentUser = null;
    }
  }

  getUserDisplayName(): string {
    if (this.currentUser) {
      return `${this.currentUser.firstName} ${this.currentUser.lastName}`;
    }
    return '';
  }

  goToProfile(): void {
    if (this.authService.isAuthenticated()) {
      this.router.navigate(['/profile']);
    }
  }

  private refreshOpenUrgencies(): void {
    if (!this.authService.isAuthenticated()) {
      this.openUrgenciesCount = 0;
      return;
    }
    this.urgencyService.getUrgenciesPaginated({ page: 1, pageSize: 1000, myUrgencies: false }).subscribe({
      next: (resp) => {
        const items = resp?.urgencies || [];
        this.openUrgenciesCount = items.filter(u => !hasAcceptedAssignment(u as any)).length;
      },
      error: () => {
        this.openUrgenciesCount = this.openUrgenciesCount || 0;
      }
    });
  }

  goToUrgencies(): void {
    this.router.navigate(['/urgencies']);
  }

  getUserProfilePicture(): string | null {
    return this.currentUser?.profilePicture || null;
  }
}
