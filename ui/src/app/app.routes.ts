import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { EmployeeListComponent } from './employee/employee-list/employee-list.component';
import { EmployeeFormComponent } from './employee/employee-form/employee-form.component';
import { NotFoundComponent } from './not-found/not-found.component';
import { LoginComponent } from './login/login.component';
import { AuthGuard } from './auth.guard';
import { ShiftManagementComponent } from './shifts/shift.component';
import { ToastSpinnerTestComponent } from './toast-spinner-test.component';

export const routes: Routes = [
  { path: '', component: HomeComponent }, // Default route that loads HomeComponent
  { path: 'home', component: HomeComponent }, // /home route
  { path: 'login', component: LoginComponent }, // used to login existing employees
  { path: 'employees', component: EmployeeListComponent, canActivate: [AuthGuard] },
  { path: 'employees/edit/:id', component: EmployeeFormComponent, canActivate: [AuthGuard] },
  { path: 'employees/new', component: EmployeeFormComponent }, // used to register new employees
  { path: 'shifts', component: ShiftManagementComponent, canActivate: [AuthGuard] }, // Shift management page
  { path: 'test-toast-spinner', component: ToastSpinnerTestComponent },

  { path: '**', component: NotFoundComponent }, // Wildcard route for undefined paths
];
