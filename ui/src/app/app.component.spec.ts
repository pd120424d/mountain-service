import { TestBed, ComponentFixture } from '@angular/core/testing';
import { AppComponent } from './app.component';
import { sharedTestingProviders } from './test-utils/shared-test-imports';
import { TranslateService } from '@ngx-translate/core';
import { AuthService } from './services/auth.service';

describe('AppComponent', () => {
  let component: AppComponent;
  let fixture: ComponentFixture<AppComponent>;
  let translateService: jasmine.SpyObj<TranslateService>;
  let authService: jasmine.SpyObj<AuthService>;

  beforeEach(async () => {
    const translateSpy = jasmine.createSpyObj('TranslateService', ['use', 'setDefaultLang']);
    const authSpy = jasmine.createSpyObj('AuthService', ['getRole', 'getUserId']);

    await TestBed.configureTestingModule({
      imports: [AppComponent],
      providers: [
        ...sharedTestingProviders,
        { provide: TranslateService, useValue: translateSpy },
        { provide: AuthService, useValue: authSpy }
      ],
    }).compileComponents();

    fixture = TestBed.createComponent(AppComponent);
    component = fixture.componentInstance;
    translateService = TestBed.inject(TranslateService) as jasmine.SpyObj<TranslateService>;
    authService = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;
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
});
