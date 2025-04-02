import { Component } from '@angular/core';
import { AuthService } from '../auth.service';
import { Router, RouterModule } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { FormsModule } from '@angular/forms';


@Component({
  selector: 'app-login',
  standalone: true,
  imports: [RouterModule, FormsModule, TranslateModule],
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css'],
})
export class LoginComponent {
  credentials = { username: '', password: '' };

  constructor(private authService: AuthService, private router: Router, private translate: TranslateService) {
    this.translate.setDefaultLang('sr-cyr')
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
