import { Routes } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { EmployeeListComponent } from './employee/employee-list/employee-list.component';
import { EmployeeFormComponent } from './employee/employee-form/employee-form.component';
import { NotFoundComponent } from './not-found/not-found.component';
import { LoginComponent } from './login/login.component';
import { AuthGuard } from './auth.guard';
import { AdminGuard } from './admin.guard';
import { NonAdminGuard } from './non-admin.guard';
import { ShiftManagementComponent } from './shifts/shift.component';
import { ToastSpinnerTestComponent } from './toast-spinner-test.component';
import { AdminComponent } from './admin/admin.component';
import { UrgencyFormComponent } from './urgency/urgency-form/urgency-form.component';
import { UrgencyListComponent } from './urgency/urgency-list/urgency-list.component';
import { UrgencyDetailComponent } from './urgency/urgency-detail/urgency-detail.component';
import { ProfileComponent } from './profile/profile.component';

export const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: 'home', component: HomeComponent },
  { path: 'login', component: LoginComponent },
  { path: 'employees', component: EmployeeListComponent, canActivate: [AdminGuard] },
  { path: 'employees/edit/:id', component: EmployeeFormComponent, canActivate: [AdminGuard] },
  { path: 'employees/new', component: EmployeeFormComponent },
  { path: 'shifts', component: ShiftManagementComponent, canActivate: [AuthGuard] },
  { path: 'urgencies', component: UrgencyListComponent, canActivate: [AuthGuard] },
  { path: 'urgencies/new', component: UrgencyFormComponent },
  { path: 'urgencies/:id', component: UrgencyDetailComponent, canActivate: [AuthGuard] },
  { path: 'profile', component: ProfileComponent, canActivate: [NonAdminGuard] },
  { path: 'admin', component: AdminComponent, canActivate: [AdminGuard] },

  { path: 'test-toast-spinner', component: ToastSpinnerTestComponent },

  { path: '**', component: NotFoundComponent },
];
