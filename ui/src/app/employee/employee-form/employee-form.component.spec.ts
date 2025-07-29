import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { EmployeeFormComponent } from './employee-form.component';
import { sharedTestingProviders } from '../../test-utils/shared-test-imports';
import { Employee, EmployeeCreateRequest, MedicRole } from '../../shared/models';
import { of, EMPTY } from 'rxjs';
import { NgxSpinnerService } from 'ngx-spinner';
import { ToastrService } from 'ngx-toastr';

describe('EmployeeFormComponent', () => {
  let component: EmployeeFormComponent;
  let fixture: ComponentFixture<EmployeeFormComponent>;
  let mockSpinnerService: jasmine.SpyObj<NgxSpinnerService>;
  let mockToastrService: jasmine.SpyObj<ToastrService>;

  beforeEach(async () => {
    // Create spy objects for the services
    mockSpinnerService = jasmine.createSpyObj('NgxSpinnerService', ['show', 'hide', 'getSpinner']);
    mockToastrService = jasmine.createSpyObj('ToastrService', ['success', 'error', 'warning']);

    // Configure getSpinner to return an empty observable
    mockSpinnerService.getSpinner.and.returnValue(EMPTY);

    await TestBed.configureTestingModule({
      imports: [EmployeeFormComponent],
      providers: [
        ...sharedTestingProviders,
        { provide: NgxSpinnerService, useValue: mockSpinnerService },
        { provide: ToastrService, useValue: mockToastrService }
      ]
    })
    .compileComponents();

    fixture = TestBed.createComponent(EmployeeFormComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should start in create mode by default', () => {
    expect(component.isEditMode).toBe(false);
  });

  it('should start in edit mode and populate form when employee data is passed via state', async () => {
    const employee: Employee = {
      id: 1,
      firstName: 'John',
      lastName: 'Doe',
      email: 'john.doe@example.com',
      phone: '+1234567890',
      profileType: MedicRole as any,
      username: 'johndoe',
      gender: 'Male'
    };

    const mockRouter = {
      getCurrentNavigation: jasmine.createSpy('getCurrentNavigation').and.returnValue({
        extras: {
          state: { employee }
        }
      }),
      navigate: jasmine.createSpy('navigate')
    };

    TestBed.resetTestingModule();
    await TestBed.configureTestingModule({
      imports: [EmployeeFormComponent],
      providers: [
        ...sharedTestingProviders,
        { provide: Router, useValue: mockRouter }
      ]
    }).compileComponents();

    const newFixture = TestBed.createComponent(EmployeeFormComponent);
    const newComponent = newFixture.componentInstance;
    newFixture.detectChanges();

    expect(newComponent.isEditMode).toBe(true);
    expect(newComponent.employeeId).toBe(1);
    expect(newComponent.employeeForm.get('firstName')?.value).toBe('John');
    expect(newComponent.employeeForm.get('lastName')?.value).toBe('Doe');
    expect(newComponent.employeeForm.get('username')?.value).toBe('johndoe');
  });

  it('should populate form with employee data', () => {
    const employee: Employee = {
      id: 1,
      firstName: 'John',
      lastName: 'Doe',
      email: 'john.doe@example.com',
      phone: '+1234567890',
      profileType: MedicRole as any,
      username: 'johndoe',
      gender: 'Male'
    };

    component.populateForm(employee);

    expect(component.employeeForm.value).toEqual({
      password: '',
      confirmPassword: '',
      username: 'johndoe',
      firstName: 'John',
      lastName: 'Doe',
      gender: 'Male',
      phoneNumber: '+1234567890',
      email: 'john.doe@example.com',
      profileType: MedicRole,
      profilePicture: null
    });
  });

  it('should submit form and navigate back to employee list', () => {
    const mockEmployeeResponse: Employee = {
      id: 1,
      firstName: 'John',
      lastName: 'Doe',
      email: 'john.doe@example.com',
      phone: '+1234567890',
      profileType: MedicRole as any,
      username: 'johndoe',
      gender: 'Male'
    };
    spyOn(component['employeeService'], 'addEmployee').and.returnValue(of(mockEmployeeResponse));
    spyOn(component['router'], 'navigate');

    component.employeeForm.setValue({
      username: 'johndoe',
      password: 'password',
      confirmPassword: 'password',
      firstName: 'John',
      lastName: 'Doe',
      gender: 'Male',
      phoneNumber: '+1234567890',
      email: 'john.doe@example.com',
      profileType: MedicRole as any,
      profilePicture: null
    });

    component.onSubmit();

    const expectedCreateRequest: EmployeeCreateRequest = {
      firstName: 'John',
      lastName: 'Doe',
      email: 'john.doe@example.com',
      phone: '+1234567890',
      username: 'johndoe',
      password: 'password',
      profileType: MedicRole as any,
      gender: 'Male',
      profilePicture: undefined
    };

    expect(component['employeeService'].addEmployee).toHaveBeenCalledWith(expectedCreateRequest);
    expect(mockSpinnerService.show).toHaveBeenCalled();
    expect(mockSpinnerService.hide).toHaveBeenCalled();
    expect(mockToastrService.success).toHaveBeenCalled();
    expect(component['router'].navigate).toHaveBeenCalledWith(['/employees']);
  });

  it('should cancel and navigate back to employee list', () => {
    spyOn(component['router'], 'navigate');
    component.cancel();
    expect(component['router'].navigate).toHaveBeenCalledWith(['/employees']);
  });


});
