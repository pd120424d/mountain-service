import { bootstrapApplication } from '@angular/platform-browser';
import { provideRouter, withDebugTracing } from '@angular/router';
import { provideHttpClient } from '@angular/common/http';
import { AppComponent } from './app/app.component';
import { routes } from './app/app.routes';
import { provideAnimationsAsync } from '@angular/platform-browser/animations/async';
import { provideTranslate } from './shared/translate-provider';


bootstrapApplication(AppComponent, {
  providers: [
    provideRouter(routes),
    provideHttpClient(),
    provideAnimationsAsync(),
    provideTranslate()]
}).catch(err => console.error(err));
