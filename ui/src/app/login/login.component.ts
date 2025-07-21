import { Component } from '@angular/core';
import { AuthService } from '../services/auth.service';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { FormsModule } from '@angular/forms';
import { CommonModule } from '@angular/common';
import { catchError } from 'rxjs/operators';
import { of } from 'rxjs';


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
  loginError = false;

  constructor(private route: ActivatedRoute, private authService: AuthService, private router: Router, private translate: TranslateService) {
    this.translate.setDefaultLang('sr-cyr')
  }

  ngOnInit(): void {
    this.route.queryParams.subscribe(params => {
      this.sessionExpired = params['sessionExpired'] === 'true';
    });
  }

  onLogin(): void {
    this.loginError = false;
    this.authService.login(this.credentials).pipe(
      catchError((error) => {
        if (error.status === 401) {
          this.loginError = true;
        }
        return of(null);
      })
    ).subscribe((response) => {
      if (response) {
        this.router.navigate(['/']);
      }
    });
  }

  switchLanguage(language: string): void {
    this.translate.use(language);
  }
}
