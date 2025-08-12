import { ComponentFixture, TestBed } from '@angular/core/testing';
import { ReactiveFormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { NgxSpinnerService } from 'ngx-spinner';
import { ToastrService } from 'ngx-toastr';
import { of, throwError } from 'rxjs';

import { ProfileComponent } from './profile.component';
import { AuthService } from '../services/auth.service';
import { EmployeeService } from '../employee/employee.service';
import { ImageUploadService } from '../services/image-upload.service';
import { Employee } from '../shared/models';

describe('ProfileComponent', () => {
  let component: ProfileComponent;
  let fixture: ComponentFixture<ProfileComponent>;
  let authService: jasmine.SpyObj<AuthService>;
  let employeeService: jasmine.SpyObj<EmployeeService>;
  let imageUploadService: jasmine.SpyObj<ImageUploadService>;
  let router: jasmine.SpyObj<Router>;
  let spinner: jasmine.SpyObj<NgxSpinnerService>;
  let toastr: jasmine.SpyObj<ToastrService>;
  let translateService: jasmine.SpyObj<TranslateService>;

  const mockEmployee: Employee = {
    id: 1,
    firstName: 'John',
    lastName: 'Doe',
    email: 'john.doe@example.com',
    phone: '+1234567890',
    username: 'johndoe',
    profileType: 'Medic',
    gender: 'M',
    profilePicture: 'https://example.com/profile.jpg'
  };

  beforeEach(async () => {
    const authSpy = jasmine.createSpyObj('AuthService', ['getUserId']);
    const employeeSpy = jasmine.createSpyObj('EmployeeService', ['getEmployeeById', 'updateEmployee']);
    const imageSpy = jasmine.createSpyObj('ImageUploadService', ['uploadProfilePicture']);
    const routerSpy = jasmine.createSpyObj('Router', ['navigate']);
    const spinnerSpy = jasmine.createSpyObj('NgxSpinnerService', ['show', 'hide']);
    const toastrSpy = jasmine.createSpyObj('ToastrService', ['success', 'error', 'warning']);
    const translateSpy = jasmine.createSpyObj('TranslateService', ['setDefaultLang', 'instant']);

    await TestBed.configureTestingModule({
      imports: [ProfileComponent, ReactiveFormsModule, TranslateModule.forRoot()],
      providers: [
        { provide: AuthService, useValue: authSpy },
        { provide: EmployeeService, useValue: employeeSpy },
        { provide: ImageUploadService, useValue: imageSpy },
        { provide: Router, useValue: routerSpy },
        { provide: NgxSpinnerService, useValue: spinnerSpy },
        { provide: ToastrService, useValue: toastrSpy },
        { provide: TranslateService, useValue: translateSpy }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(ProfileComponent);
    component = fixture.componentInstance;
    authService = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;
    employeeService = TestBed.inject(EmployeeService) as jasmine.SpyObj<EmployeeService>;
    imageUploadService = TestBed.inject(ImageUploadService) as jasmine.SpyObj<ImageUploadService>;
    router = TestBed.inject(Router) as jasmine.SpyObj<Router>;
    spinner = TestBed.inject(NgxSpinnerService) as jasmine.SpyObj<NgxSpinnerService>;
    toastr = TestBed.inject(ToastrService) as jasmine.SpyObj<ToastrService>;
    translateService = TestBed.inject(TranslateService) as jasmine.SpyObj<TranslateService>;

    translateService.instant.and.returnValue('Translated text');
  });

  it('should create', () => {
    authService.getUserId.and.returnValue('1');
    employeeService.getEmployeeById.and.returnValue(of(mockEmployee));
    
    expect(component).toBeTruthy();
  });

  it('should redirect to login if no user ID', () => {
    authService.getUserId.and.returnValue('');
    
    component.ngOnInit();
    
    expect(router.navigate).toHaveBeenCalledWith(['/login']);
  });

  it('should load current user on init', () => {
    authService.getUserId.and.returnValue('1');
    employeeService.getEmployeeById.and.returnValue(of(mockEmployee));
    
    component.ngOnInit();
    
    expect(spinner.show).toHaveBeenCalled();
    expect(employeeService.getEmployeeById).toHaveBeenCalledWith(1);
    expect(component.currentUser).toEqual(mockEmployee);
    expect(spinner.hide).toHaveBeenCalled();
  });

  it('should handle error when loading user', () => {
    authService.getUserId.and.returnValue('1');
    employeeService.getEmployeeById.and.returnValue(throwError(() => new Error('Load error')));
    
    component.ngOnInit();
    
    expect(spinner.hide).toHaveBeenCalled();
    expect(toastr.error).toHaveBeenCalled();
  });

  it('should handle image selection', () => {
    const mockFile = new File([''], 'test.jpg', { type: 'image/jpeg' });
    const imageEvent = {
      file: mockFile,
      preview: 'data:image/jpeg;base64,test',
      isValid: true
    };
    
    component.onImageSelected(imageEvent);
    
    expect(component.selectedImageFile).toBe(mockFile);
    expect(component.profileForm.get('profilePicture')?.value).toBe('data:image/jpeg;base64,test');
  });

  it('should handle invalid image selection', () => {
    const mockFile = new File([''], 'invalid.txt', { type: 'text/plain' });
    const imageEvent = {
      file: mockFile,
      preview: '',
      isValid: false,
      error: 'Invalid file'
    };

    component.onImageSelected(imageEvent);

    expect(component.selectedImageFile).toBeNull();
    expect(toastr.error).toHaveBeenCalledWith('Invalid file');
  });

  it('should handle image removal', () => {
    component.selectedImageFile = new File([''], 'test.jpg');
    component.currentProfilePictureUrl = 'test-url';
    
    component.onImageRemoved();
    
    expect(component.selectedImageFile).toBeNull();
    expect(component.currentProfilePictureUrl).toBeUndefined();
    expect(component.profileForm.get('profilePicture')?.value).toBeNull();
  });

  it('should not submit if form is invalid', () => {
    component.profileForm.patchValue({
      firstName: '', // Invalid - required
      lastName: 'Doe',
      email: 'john@example.com',
      phone: '+1234567890',
      username: 'johndoe'
    });
    
    component.onSubmit();
    
    expect(employeeService.updateEmployee).not.toHaveBeenCalled();
  });

  it('should submit form successfully without image', () => {
    component.currentUser = mockEmployee;
    component.profileForm.patchValue({
      firstName: 'John',
      lastName: 'Doe',
      email: 'john@example.com',
      phone: '+1234567890',
      username: 'johndoe',
      gender: 'Male',
      profileType: 'Employee'
    });

    employeeService.updateEmployee.and.returnValue(of(mockEmployee));

    component.onSubmit();

    expect(spinner.show).toHaveBeenCalled();
    expect(employeeService.updateEmployee).toHaveBeenCalled();
    expect(spinner.hide).toHaveBeenCalled();
    expect(toastr.success).toHaveBeenCalled();
  });

  it('should submit form successfully with image', () => {
    component.currentUser = mockEmployee;
    component.selectedImageFile = new File([''], 'test.jpg');
    component.profileForm.patchValue({
      firstName: 'John',
      lastName: 'Doe',
      email: 'john@example.com',
      phone: '+1234567890',
      username: 'johndoe',
      gender: 'Male',
      profileType: 'Employee'
    });

    employeeService.updateEmployee.and.returnValue(of(mockEmployee));
    imageUploadService.uploadProfilePicture.and.returnValue(of({
      progress: 100,
      status: 'completed' as const
    }));

    component.onSubmit();

    expect(employeeService.updateEmployee).toHaveBeenCalled();
    expect(imageUploadService.uploadProfilePicture).toHaveBeenCalled();
  });

  it('should handle update error', () => {
    component.currentUser = mockEmployee;
    component.profileForm.patchValue({
      firstName: 'John',
      lastName: 'Doe',
      email: 'john@example.com',
      phone: '+1234567890',
      username: 'johndoe',
      gender: 'Male',
      profileType: 'Employee'
    });

    employeeService.updateEmployee.and.returnValue(throwError(() => new Error('Update error')));

    component.onSubmit();

    expect(spinner.hide).toHaveBeenCalled();
    expect(toastr.error).toHaveBeenCalled();
  });

  it('should go back to home', () => {
    component.goBack();

    expect(router.navigate).toHaveBeenCalledWith(['/']);
  });
});
