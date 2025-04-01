// src/app/app.config.ts
import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { routes } from './app.routes'; // Import routes from app.routes.ts

@NgModule({
  imports: [RouterModule.forRoot(routes)], // Configure the router with routes
  exports: [RouterModule] // Export RouterModule to make it available throughout the app
})
export class AppConfigModule { }
