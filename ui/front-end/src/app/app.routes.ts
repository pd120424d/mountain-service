// src/app/app.routes.ts
import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { EmployeeListComponent } from './employee/employee-list/employee-list.component';
import { EmployeeFormComponent } from './employee/employee-form/employee-form.component';

export const routes: Routes = [
  { path: '', component: HomeComponent }, // Default route that loads HomeComponent
  { path: 'home', component: HomeComponent }, // /home route
  { path: 'employees', component: EmployeeListComponent }, // /employees route
  { path: 'employees/new', component: EmployeeFormComponent }, // /employees/new route
  { path: 'employees/edit/:id', component: EmployeeFormComponent }, // /employees/edit/:id route
  { path: '**', redirectTo: '' } // Wildcard route for undefined paths
];
