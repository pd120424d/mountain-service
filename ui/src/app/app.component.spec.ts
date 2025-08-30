import { TestBed, ComponentFixture } from '@angular/core/testing';
import { AppComponent } from './app.component';
import { sharedTestingProviders } from './test-utils/shared-test-imports';
import { TranslateService } from '@ngx-translate/core';
import { AuthService } from './services/auth.service';
import { EmployeeService } from './employee/employee.service';
import { UrgencyService } from './urgency/urgency.service';
import { AppInitializationService } from './services/app-initialization.service';
import { of, throwError } from 'rxjs';
import { Employee } from './shared/models';

describe('AppComponent', () => {
  let component: AppComponent;
  let fixture: ComponentFixture<AppComponent>;
  let translateService: jasmine.SpyObj<TranslateService>;
  let authService: jasmine.SpyObj<AuthService>;
  let employeeService: jasmine.SpyObj<EmployeeService>;
  let appInitService: jasmine.SpyObj<AppInitializationService>;
  let urgencyService: jasmine.SpyObj<UrgencyService>;

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
    const translateSpy = jasmine.createSpyObj('TranslateService', ['use', 'setDefaultLang']);
    const authSpy = jasmine.createSpyObj('AuthService', ['getRole', 'getUserId', 'isAuthenticated']);
    const employeeSpy = jasmine.createSpyObj('EmployeeService', ['getEmployeeById']);
    const appInitSpy = jasmine.createSpyObj('AppInitializationService', ['initialize', 'cleanup']);
    const urgencySpy = jasmine.createSpyObj('UrgencyService', ['getUrgencies']);

    await TestBed.configureTestingModule({
      imports: [AppComponent],
      providers: [
        ...sharedTestingProviders,
        { provide: TranslateService, useValue: translateSpy },
        { provide: AuthService, useValue: authSpy },
        { provide: EmployeeService, useValue: employeeSpy },
        { provide: AppInitializationService, useValue: appInitSpy },
        { provide: UrgencyService, useValue: urgencySpy }
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(AppComponent);
    component = fixture.componentInstance;
    translateService = TestBed.inject(TranslateService) as jasmine.SpyObj<TranslateService>;
    authService = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;
    employeeService = TestBed.inject(EmployeeService) as jasmine.SpyObj<EmployeeService>;
    appInitService = TestBed.inject(AppInitializationService) as jasmine.SpyObj<AppInitializationService>;
    urgencyService = TestBed.inject(UrgencyService) as jasmine.SpyObj<UrgencyService>;

    // Setup default return values
    appInitService.initialize.and.returnValue(Promise.resolve());
    appInitService.cleanup.and.returnValue();
  });

  it('should create the app', () => {
    expect(component).toBeTruthy();
  });

  it('should initialize with saved language from localStorage', () => {
    spyOn(localStorage, 'getItem').and.returnValue('sr-lat');

    // Create a new component instance to test constructor
    const newFixture = TestBed.createComponent(AppComponent);
    const newComponent = newFixture.componentInstance;

    expect(translateService.use).toHaveBeenCalledWith('sr-lat');
    expect(newComponent.currentLanguageLabel).toBe('SR');
  });

  it('should initialize with default language when localStorage is empty', () => {
    spyOn(localStorage, 'getItem').and.returnValue(null);

    // Create a new component instance to test constructor
    const newFixture = TestBed.createComponent(AppComponent);
    const newComponent = newFixture.componentInstance;

    expect(translateService.use).toHaveBeenCalledWith('en');
    expect(newComponent.currentLanguageLabel).toBe('EN');
  });

  it('should toggle dropdown', () => {
    component.isDropdownOpen = false;
    component.toggleDropdown();
    expect(component.isDropdownOpen).toBe(true);

    component.toggleDropdown();
    expect(component.isDropdownOpen).toBe(false);
  });

  it('should switch language and close dropdown', () => {
    const setItemSpy = spyOn(localStorage, 'setItem');
    const setLanguageLabelSpy = spyOn(component, 'setLanguageLabel');
    component.isDropdownOpen = true;

    component.switchLanguage('sr-lat');

    expect(translateService.use).toHaveBeenCalledWith('sr-lat');
    expect(setItemSpy).toHaveBeenCalledWith('language', 'sr-lat');
    expect(setLanguageLabelSpy).toHaveBeenCalledWith('sr-lat');
    expect(component.isDropdownOpen).toBe(false);
  });

  it('should set language label correctly', () => {
    component.setLanguageLabel('en');
    expect(component.currentLanguageLabel).toBe('EN');

    component.setLanguageLabel('sr-lat');
    expect(component.currentLanguageLabel).toBe('SR');

    component.setLanguageLabel('sr-cyr');
    expect(component.currentLanguageLabel).toBe('СР');

    component.setLanguageLabel('ru');
    expect(component.currentLanguageLabel).toBe('RU');

    component.setLanguageLabel('unknown');
    expect(component.currentLanguageLabel).toBe('UNKNOWN');
  });

  it('should close dropdown when clicking outside', () => {
    // Mock DOM elements
    const mockDropdown = document.createElement('div');
    mockDropdown.className = 'language-switcher';
    const mockTarget = document.createElement('div');

    spyOn(document, 'querySelector').and.returnValue(mockDropdown);
    spyOn(mockDropdown, 'contains').and.returnValue(false);

    component.isDropdownOpen = true;
    const event = new Event('click');
    Object.defineProperty(event, 'target', { value: mockTarget });

    component.closeDropdown(event);

    expect(component.isDropdownOpen).toBe(false);
  });

  it('should not close dropdown when clicking inside', () => {
    // Mock DOM elements
    const mockDropdown = document.createElement('div');
    mockDropdown.className = 'language-switcher';
    const mockTarget = document.createElement('div');

    spyOn(document, 'querySelector').and.returnValue(mockDropdown);
    spyOn(mockDropdown, 'contains').and.returnValue(true);

    component.isDropdownOpen = true;
    const event = new Event('click');
    Object.defineProperty(event, 'target', { value: mockTarget });

    component.closeDropdown(event);

    expect(component.isDropdownOpen).toBe(true);
  });

  it('should handle missing dropdown element', () => {
    spyOn(document, 'querySelector').and.returnValue(null);

    component.isDropdownOpen = true;
    const event = new Event('click');

    expect(() => component.closeDropdown(event)).not.toThrow();
    expect(component.isDropdownOpen).toBe(true);
  });

  it('should load current user when authenticated', () => {
    authService.isAuthenticated.and.returnValue(true);
    authService.getUserId.and.returnValue('1');
    employeeService.getEmployeeById.and.returnValue(of(mockEmployee));

    component['loadCurrentUser']();

    expect(employeeService.getEmployeeById).toHaveBeenCalledWith(1);
    expect(component.currentUser).toEqual(mockEmployee);
  });

  it('should not load user when not authenticated', () => {
    authService.isAuthenticated.and.returnValue(false);

    component['loadCurrentUser']();

    expect(employeeService.getEmployeeById).not.toHaveBeenCalled();
    expect(component.currentUser).toBeNull();
  });

  it('should handle error when loading current user', () => {
    authService.isAuthenticated.and.returnValue(true);
    authService.getUserId.and.returnValue('1');
    employeeService.getEmployeeById.and.returnValue(throwError(() => new Error('Load error')));

    component['loadCurrentUser']();

    expect(component.currentUser).toBeNull();
  });

  it('should return user display name', () => {
    component.currentUser = mockEmployee;

    const displayName = component.getUserDisplayName();

    expect(displayName).toBe('John Doe');
  });

  it('should return empty string when no current user', () => {
    component.currentUser = null;

    const displayName = component.getUserDisplayName();

    expect(displayName).toBe('');
  });

  it('should return user profile picture', () => {
    component.currentUser = mockEmployee;

    const profilePicture = component.getUserProfilePicture();

    expect(profilePicture).toBe('https://example.com/profile.jpg');
  });

  it('should return null when no profile picture', () => {
    component.currentUser = { ...mockEmployee, profilePicture: undefined };

    const profilePicture = component.getUserProfilePicture();

    expect(profilePicture).toBeNull();
  });

  it('should return null when no current user', () => {
    component.currentUser = null;

    const profilePicture = component.getUserProfilePicture();

    expect(profilePicture).toBeNull();
  });

  it('should navigate to profile when authenticated', () => {
    authService.isAuthenticated.and.returnValue(true);
    const navSpy = spyOn((component as any)['router'], 'navigate');

    component.goToProfile();

    expect(navSpy).toHaveBeenCalledWith(['/profile']);
  });

  it('should not navigate to profile when not authenticated', () => {
    authService.isAuthenticated.and.returnValue(false);
    const navSpy = spyOn((component as any)['router'], 'navigate');

    component.goToProfile();

    expect(navSpy).not.toHaveBeenCalled();
  });

  it('should set openUrgenciesCount to 0 when not authenticated', () => {
    authService.isAuthenticated.and.returnValue(false);
    (component as any)['refreshOpenUrgencies']();
    expect(component.openUrgenciesCount).toBe(0);
  });

  it('should compute open urgencies count when authenticated', () => {
    authService.isAuthenticated.and.returnValue(true);
    const urgencies = [
      { id: 1, assignedEmployeeId: undefined }, // unassigned
      { id: 2, assignedEmployeeId: null }, // unassigned
      { id: 3, assignedEmployeeId: 5 } // assigned
    ] as any[];
    urgencyService.getUrgencies.and.returnValue(of(urgencies));

    (component as any)['refreshOpenUrgencies']();

    expect(component.openUrgenciesCount).toBe(2);
  });

  it('should keep previous openUrgenciesCount on error', () => {
    authService.isAuthenticated.and.returnValue(true);
    component.openUrgenciesCount = 5;
    urgencyService.getUrgencies.and.returnValue(throwError(() => new Error('fail')));

    (component as any)['refreshOpenUrgencies']();

    expect(component.openUrgenciesCount).toBe(5);
  });
});
