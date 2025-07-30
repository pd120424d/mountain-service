import { ComponentFixture, TestBed, fakeAsync, tick } from '@angular/core/testing';
import { ToastSpinnerTestComponent } from './toast-spinner-test.component';
import { sharedTestingProviders } from './test-utils/shared-test-imports';
import { ToastrService } from 'ngx-toastr';
import { NgxSpinnerService } from 'ngx-spinner';
import { environment } from '../environments/environment';
import { EMPTY } from 'rxjs';

describe('ToastSpinnerTestComponent', () => {
  let component: ToastSpinnerTestComponent;
  let fixture: ComponentFixture<ToastSpinnerTestComponent>;
  let toastrService: jasmine.SpyObj<ToastrService>;
  let spinnerService: jasmine.SpyObj<NgxSpinnerService>;

  beforeEach(async () => {
    const toastrSpy = jasmine.createSpyObj('ToastrService', ['success', 'info', 'warning', 'error']);
    const spinnerSpy = jasmine.createSpyObj('NgxSpinnerService', ['show', 'hide', 'getSpinner']);
    spinnerSpy.getSpinner.and.returnValue(EMPTY);

    await TestBed.configureTestingModule({
      imports: [ToastSpinnerTestComponent],
      providers: [
        ...sharedTestingProviders,
        { provide: ToastrService, useValue: toastrSpy },
        { provide: NgxSpinnerService, useValue: spinnerSpy }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(ToastSpinnerTestComponent);
    component = fixture.componentInstance;
    toastrService = TestBed.inject(ToastrService) as jasmine.SpyObj<ToastrService>;
    spinnerService = TestBed.inject(NgxSpinnerService) as jasmine.SpyObj<NgxSpinnerService>;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should set isStaging from environment', () => {
    expect(component.isStaging).toBe(environment.staging);
  });

  describe('testToastr', () => {
    it('should call all toastr methods', () => {
      spyOn(console, 'log');
      
      component.testToastr();

      expect(console.log).toHaveBeenCalledWith('Testing Toastr...');
      expect(toastrService.success).toHaveBeenCalledWith('Success message!', 'Success');
      expect(toastrService.info).toHaveBeenCalledWith('Info message!', 'Info');
      expect(toastrService.warning).toHaveBeenCalledWith('Warning message!', 'Warning');
      expect(toastrService.error).toHaveBeenCalledWith('Error message!', 'Error');
    });
  });

  describe('testSpinner', () => {
    it('should show and hide spinner with success message', fakeAsync(() => {
      spyOn(console, 'log');
      
      component.testSpinner();

      expect(console.log).toHaveBeenCalledWith('Testing Spinner...');
      expect(spinnerService.show).toHaveBeenCalledWith(undefined, {
        type: 'ball-scale-multiple',
        size: 'large',
        bdColor: 'rgba(0, 0, 0, 0.8)',
        color: '#fff'
      });

      tick(2000);

      expect(spinnerService.hide).toHaveBeenCalled();
      expect(toastrService.success).toHaveBeenCalledWith('Spinner test completed!');
    }));
  });

  describe('testBoth', () => {
    it('should test both spinner and toastr', fakeAsync(() => {
      spyOn(console, 'log');
      
      component.testBoth();

      expect(console.log).toHaveBeenCalledWith('Testing both...');
      expect(spinnerService.show).toHaveBeenCalledWith(undefined, {
        type: 'ball-scale-multiple',
        size: 'large',
        bdColor: 'rgba(0, 0, 0, 0.8)',
        color: '#fff'
      });
      expect(toastrService.info).toHaveBeenCalledWith('Starting combined test...');

      tick(3000);

      expect(spinnerService.hide).toHaveBeenCalled();
      expect(toastrService.success).toHaveBeenCalledWith('Combined test completed!');
    }));
  });
});
