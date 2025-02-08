import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { CommonModule } from '@angular/common'; // Import CommonModule



@Component({
  selector: 'app-root',
  standalone: true,
  imports: [TranslateModule, RouterOutlet, CommonModule],
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent {
  isDropdownOpen = false; // Control dropdown visibility

  constructor(private translate: TranslateService) {
    const savedLanguage = localStorage.getItem('language') || 'sr-cyr';
    this.translate.use(savedLanguage);
  }

  toggleDropdown(): void {
    this.isDropdownOpen = !this.isDropdownOpen; // Toggle dropdown visibility
  }

  switchLanguage(language: string): void {
    this.translate.use(language);
    localStorage.setItem('language', language);
  }
}
