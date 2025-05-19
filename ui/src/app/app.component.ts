import { Component, HostListener } from '@angular/core';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { AuthService } from './services/auth.service';
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
export class AppComponent {
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
    private translate: TranslateService
  ) {
    const savedLanguage = localStorage.getItem('language') || 'en';
    this.translate.use(savedLanguage);
    this.setLanguageLabel(savedLanguage);
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
