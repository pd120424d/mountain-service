import { Component, HostListener, OnInit, OnDestroy } from '@angular/core';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { AuthService } from './services/auth.service';
import { AppInitializationService } from './services/app-initialization.service';
import { CommonModule } from '@angular/common';
import { RouterModule, RouterOutlet } from '@angular/router';
import { VersionBannerComponent } from './version-banner/version-banner.component';
import { NgxSpinnerModule } from 'ngx-spinner';

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

  private languageMap: { [key: string]: string } = {
    'en': 'EN',
    'sr-lat': 'SR',
    'sr-cyr': 'СР',
    'ru': 'RU'
  };

  constructor(
    public authService: AuthService,
    private translate: TranslateService,
    private appInitService: AppInitializationService
  ) {
    const savedLanguage = localStorage.getItem('language') || 'en';
    this.translate.use(savedLanguage);
    this.setLanguageLabel(savedLanguage);
  }

  ngOnInit(): void {
    // Initialize the application
    this.appInitService.initialize().then(() => {
      console.log('Application initialized successfully');
    }).catch((error) => {
      console.error('Failed to initialize application:', error);
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
}
