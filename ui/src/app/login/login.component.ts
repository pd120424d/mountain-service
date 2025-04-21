import { Component } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';


@Component({
  selector: 'app-login',
  standalone: true,
  imports: [RouterModule, FormsModule, TranslateModule, CommonModule],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
})
export class LoginComponent {
  credentials = { username: '', password: '' };
  sessionExpired = false;

  constructor(private route: ActivatedRoute, private authService: AuthService, private router: Router, private translate: TranslateService) {
    this.translate.setDefaultLang('sr-cyr')
  }

  ngOnInit(): void {
    this.route.queryParams.subscribe(params => {
      this.sessionExpired = params['sessionExpired'] === 'true';
    });
  }

  onLogin(): void {
    this.authService.login(this.credentials).subscribe(() => {
      this.router.navigate(['/']);
    });
  }

  switchLanguage(language: string): void {
    this.translate.use(language);
  }
}
