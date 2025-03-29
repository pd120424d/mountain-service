import { Component, HostListener } from '@angular/core';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { AuthService } from './auth.service';
import { CommonModule } from '@angular/common';
import { RouterOutlet } from '@angular/router';

@Component({
  selector: 'app-root',
  standalone: true,
  imports: [CommonModule, RouterOutlet, TranslateModule],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css'],
})
export class AppComponent {
  isDropdownOpen = false; // Controls dropdown visibility

  constructor(public authService: AuthService, private translate: TranslateService) {
    const savedLanguage = localStorage.getItem('language') || 'en';
    this.translate.use(savedLanguage);
  }

  toggleDropdown(): void {
    this.isDropdownOpen = !this.isDropdownOpen; // Toggle dropdown visibility
  }

  switchLanguage(language: string): void {
    this.translate.use(language);
    localStorage.setItem('language', language);
    this.isDropdownOpen = false; // Close dropdown after selection
  }

  @HostListener('document:click', ['$event'])
  closeDropdown(event: Event): void {
    const dropdown = document.querySelector('.language-switcher');
    if (dropdown && !dropdown.contains(event.target as Node)) {
      this.isDropdownOpen = false; // Close if clicked outside
    }
  }
}
