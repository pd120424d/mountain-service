import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { TranslateService } from '@ngx-translate/core';
import { NgxSpinnerService } from 'ngx-spinner';
import { ToastrService } from 'ngx-toastr';
import { of } from 'rxjs';

import { HomeComponent } from './home.component';
import { sharedTestingProviders } from '../test-utils/shared-test-imports';
import { AuthService } from '../services/auth.service';

describe('HomeComponent', () => {
  let component: HomeComponent;
  let fixture: ComponentFixture<HomeComponent>;
  let translateService: TranslateService;
  let spinnerService: jasmine.SpyObj<NgxSpinnerService>;
  let toastrService: jasmine.SpyObj<ToastrService>;
  let authService: jasmine.SpyObj<AuthService>;

  beforeEach(async () => {
    const spinnerSpy = jasmine.createSpyObj('NgxSpinnerService', ['show', 'hide']);
    const toastrSpy = jasmine.createSpyObj('ToastrService', ['success', 'error', 'info', 'warning']);
    const authSpy = jasmine.createSpyObj('AuthService', ['isAuthenticated', 'isAdmin']);

    await TestBed.configureTestingModule({
      imports: [HomeComponent],
      providers: [
        ...sharedTestingProviders,
        { provide: NgxSpinnerService, useValue: spinnerSpy },
        { provide: ToastrService, useValue: toastrSpy },
        { provide: AuthService, useValue: authSpy }
      ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(HomeComponent);
    component = fixture.componentInstance;
    translateService = TestBed.inject(TranslateService);
    spinnerService = TestBed.inject(NgxSpinnerService) as jasmine.SpyObj<NgxSpinnerService>;
    toastrService = TestBed.inject(ToastrService) as jasmine.SpyObj<ToastrService>;
    authService = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;

    spyOn(translateService, 'use');
    spyOn(translateService, 'setDefaultLang');

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should initialize with correct default values', () => {
    expect(component.images.length).toBe(11);
    expect(component.currentImageIndex).toBe(0);
    expect(component.prevImageIndex).toBe(0);
    expect(component.isStaging).toBeDefined();
  });

  describe('preloadImage', () => {
    it('should preload image successfully', async () => {
      const imageUrl = 'test-image.jpg';

      const mockImage = {
        src: '',
        onload: null as any
      };
      spyOn(window, 'Image').and.returnValue(mockImage as any);

      const promise = component.preloadImage(imageUrl);

      mockImage.onload();

      await expectAsync(promise).toBeResolved();
      expect(mockImage.src).toBe(imageUrl);
    });
  });

  describe('switchLanguage', () => {
    it('should switch language', () => {
      component.switchLanguage('en');
      expect(translateService.use).toHaveBeenCalledWith('en');
    });
  });

  describe('testSpinner', () => {
    it('should show and hide spinner with toastr message', fakeAsync(() => {
      component.testSpinner();

      expect(spinnerService.show).toHaveBeenCalledWith(undefined, {
        type: 'ball-scale-multiple',
        size: 'large',
        bdColor: 'rgba(0, 0, 0, 0.8)',
        color: '#fff'
      });

      tick(2000);

      expect(spinnerService.hide).toHaveBeenCalled();
      expect(toastrService.success).toHaveBeenCalledWith('Spinner done!');
    }));
  });

  describe('testToastr', () => {
    it('should show toastr info message', () => {
      component.testToastr();
      expect(toastrService.info).toHaveBeenCalledWith('Toastr test!');
    });
  });

  describe('Navigation Access Control', () => {
    it('should expose authService for template access', () => {
      expect(component.authService).toBeDefined();
    });

    it('should allow access to authService.isAuthenticated', () => {
      authService.isAuthenticated.and.returnValue(true);
      expect(component.authService.isAuthenticated()).toBe(true);
    });

    it('should allow access to authService.isAdmin', () => {
      authService.isAdmin.and.returnValue(true);
      expect(component.authService.isAdmin()).toBe(true);
    });
  });
});
