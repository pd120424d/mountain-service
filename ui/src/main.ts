import { bootstrapApplication } from '@angular/platform-browser';
import { provideRouter } from '@angular/router';
import { provideHttpClient, withInterceptors } from '@angular/common/http';
import { AppComponent } from './app/app.component';
import { routes } from './app/app.routes';
import { provideAnimations } from '@angular/platform-browser/animations';
import { provideTranslate } from './shared/translate-provider';
import { authInterceptor } from './app/auth.interceptor';
import { ToastrModule } from 'ngx-toastr';
import { NgxSpinnerModule } from 'ngx-spinner';
import { importProvidersFrom } from '@angular/core';

import { requestIdInterceptor } from './app/request-id.interceptor';

bootstrapApplication(AppComponent, {
  providers: [
    provideRouter(routes),
    provideHttpClient(withInterceptors([requestIdInterceptor, authInterceptor])),
    provideAnimations(),
    provideTranslate(),
    importProvidersFrom(
      ToastrModule.forRoot({
        positionClass: 'toast-bottom-right',
        timeOut: 3000,
        closeButton: true,
        preventDuplicates: true
      }),
      NgxSpinnerModule.forRoot({ type: 'ball-scale-multiple' })
    )]
}).catch(err => console.error(err));
