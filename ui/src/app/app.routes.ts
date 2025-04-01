import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { EmployeeListComponent } from './employee/employee-list/employee-list.component';
import { EmployeeFormComponent } from './employee/employee-form/employee-form.component';
import { NotFoundComponent } from './not-found/not-found.component';
import { LoginComponent } from './login/login.component';
import { AuthGuard } from './auth.guard';

export const routes: Routes = [
  { path: '', component: HomeComponent }, // Default route that loads HomeComponent
  { path: 'home', component: HomeComponent }, // /home route
  { path: 'employees', component: EmployeeListComponent, canActivate: [AuthGuard] }, // /employees route
  { path: 'employees/new', component: EmployeeFormComponent, canActivate: [AuthGuard] }, // /employees/new route
  { path: 'employees/edit/:id', component: EmployeeFormComponent, canActivate: [AuthGuard] }, // /employees/edit/:id route
  { path: '**', component: NotFoundComponent }, // Wildcard route for undefined paths
  { path: 'login', component: LoginComponent },
];
