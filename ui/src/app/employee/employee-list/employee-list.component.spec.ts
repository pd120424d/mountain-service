import { ComponentFixture, TestBed } from '@angular/core/testing';
import { Router } from '@angular/router';
import { MatDialog } from '@angular/material/dialog';
import { TranslateService } from '@ngx-translate/core';
import { of, throwError } from 'rxjs';

import { EmployeeListComponent } from './employee-list.component';
import { EmployeeService } from '../employee.service';
import { Employee } from '../../shared/models';
import { ConfirmDialogComponent } from '../confirm-dialog/confirm-dialog.component';
import { sharedTestingProviders } from '../../test-utils/shared-test-imports';

describe('EmployeeListComponent', () => {
  let component: EmployeeListComponent;
  let fixture: ComponentFixture<EmployeeListComponent>;
  let mockEmployeeService: jasmine.SpyObj<EmployeeService>;
  let mockRouter: jasmine.SpyObj<Router>;
  let mockDialog: jasmine.SpyObj<MatDialog>;
  let mockTranslateService: jasmine.SpyObj<TranslateService>;

  const mockEmployees: Employee[] = [
    {
      id: 1,
      username: 'john.doe',
      firstName: 'John',
      lastName: 'Doe',
      email: 'john.doe@example.com',
      phone: '+1234567890',
      profileType: 'Medic',
      gender: 'Male',
      profilePicture: 'https://example.com/profile1.jpg'
    },
    {
      id: 2,
      username: 'jane.smith',
      firstName: 'Jane',
      lastName: 'Smith',
      email: 'jane.smith@example.com',
      phone: '+1234567891',
      profileType: 'Technical',
      gender: 'Female',
      profilePicture: 'https://example.com/profile2.jpg'
    }
  ];

  beforeEach(async () => {
    const employeeServiceSpy = jasmine.createSpyObj('EmployeeService', ['getEmployees', 'deleteEmployee']);
    const routerSpy = jasmine.createSpyObj('Router', ['navigate']);
    const dialogSpy = jasmine.createSpyObj('MatDialog', ['open']);
    const translateServiceSpy = jasmine.createSpyObj('TranslateService', ['get', 'use', 'setDefaultLang']);

    await TestBed.configureTestingModule({
      imports: [EmployeeListComponent],
      providers: [
        ...sharedTestingProviders,
        { provide: EmployeeService, useValue: employeeServiceSpy },
        { provide: Router, useValue: routerSpy },
        { provide: MatDialog, useValue: dialogSpy },
        { provide: TranslateService, useValue: translateServiceSpy }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(EmployeeListComponent);
    component = fixture.componentInstance;
    mockEmployeeService = TestBed.inject(EmployeeService) as jasmine.SpyObj<EmployeeService>;
    mockRouter = TestBed.inject(Router) as jasmine.SpyObj<Router>;
    mockDialog = TestBed.inject(MatDialog) as jasmine.SpyObj<MatDialog>;
    mockTranslateService = TestBed.inject(TranslateService) as jasmine.SpyObj<TranslateService>;

    // Setup default translate service behavior
    mockTranslateService.get.and.returnValue(of('translated text'));
    mockTranslateService.use.and.returnValue(of({}));
    mockTranslateService.setDefaultLang.and.stub();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('ngOnInit', () => {
    it('should call loadEmployees on initialization', () => {
      spyOn(component, 'loadEmployees');
      
      component.ngOnInit();
      
      expect(component.loadEmployees).toHaveBeenCalled();
    });
  });

  describe('loadEmployees', () => {
    it('it succeeds when employees are loaded successfully', () => {
      mockEmployeeService.getEmployees.and.returnValue(of(mockEmployees));
      
      component.loadEmployees();
      
      expect(mockEmployeeService.getEmployees).toHaveBeenCalled();
      expect(component.employees).toEqual(mockEmployees);
    });

    it('it handles empty employee list', () => {
      mockEmployeeService.getEmployees.and.returnValue(of([]));
      
      component.loadEmployees();
      
      expect(mockEmployeeService.getEmployees).toHaveBeenCalled();
      expect(component.employees).toEqual([]);
    });

    it('it handles service error gracefully', () => {
      mockEmployeeService.getEmployees.and.returnValue(throwError(() => new Error('Service error')));
      spyOn(console, 'error'); // Spy on console.error to prevent actual error logging

      component.loadEmployees();

      expect(mockEmployeeService.getEmployees).toHaveBeenCalled();
      // Component should handle error gracefully without crashing
    });
  });

  describe('openDeleteDialog', () => {
    it('it succeeds when dialog is opened and user confirms deletion', () => {
      const dialogRefSpy = jasmine.createSpyObj('MatDialogRef', ['afterClosed']);
      dialogRefSpy.afterClosed.and.returnValue(of(true));
      mockDialog.open.and.returnValue(dialogRefSpy);
      spyOn(component, 'deleteEmployee');
      
      component.openDeleteDialog(1);
      
      expect(mockDialog.open).toHaveBeenCalledWith(ConfirmDialogComponent);
      expect(component.deleteEmployee).toHaveBeenCalledWith(1);
    });

    it('it succeeds when dialog is opened and user cancels deletion', () => {
      const dialogRefSpy = jasmine.createSpyObj('MatDialogRef', ['afterClosed']);
      dialogRefSpy.afterClosed.and.returnValue(of(false));
      mockDialog.open.and.returnValue(dialogRefSpy);
      spyOn(component, 'deleteEmployee');
      
      component.openDeleteDialog(1);
      
      expect(mockDialog.open).toHaveBeenCalledWith(ConfirmDialogComponent);
      expect(component.deleteEmployee).not.toHaveBeenCalled();
    });

    it('it succeeds when dialog is closed without result', () => {
      const dialogRefSpy = jasmine.createSpyObj('MatDialogRef', ['afterClosed']);
      dialogRefSpy.afterClosed.and.returnValue(of(null));
      mockDialog.open.and.returnValue(dialogRefSpy);
      spyOn(component, 'deleteEmployee');
      
      component.openDeleteDialog(1);
      
      expect(mockDialog.open).toHaveBeenCalledWith(ConfirmDialogComponent);
      expect(component.deleteEmployee).not.toHaveBeenCalled();
    });
  });

  describe('deleteEmployee', () => {
    it('it succeeds when employee is deleted successfully', () => {
      mockEmployeeService.deleteEmployee.and.returnValue(of(undefined));
      spyOn(component, 'loadEmployees');

      component.deleteEmployee(1);

      expect(mockEmployeeService.deleteEmployee).toHaveBeenCalledWith(1);
      expect(component.loadEmployees).toHaveBeenCalled();
    });

    it('it handles deletion error gracefully', () => {
      mockEmployeeService.deleteEmployee.and.returnValue(throwError(() => new Error('Delete error')));
      spyOn(component, 'loadEmployees');
      spyOn(console, 'error'); // Spy on console.error to prevent actual error logging

      component.deleteEmployee(1);

      expect(mockEmployeeService.deleteEmployee).toHaveBeenCalledWith(1);
      expect(component.loadEmployees).not.toHaveBeenCalled();
    });
  });

  describe('editEmployee', () => {
    it('it succeeds when navigating to edit employee', () => {
      const employee = mockEmployees[0];
      
      component.editEmployee(employee);
      
      expect(mockRouter.navigate).toHaveBeenCalledWith(['/employees/edit', employee.id], {
        state: { employee }
      });
    });

    it('it logs employee data when editing', () => {
      spyOn(console, 'log');
      const employee = mockEmployees[0];
      
      component.editEmployee(employee);
      
      expect(console.log).toHaveBeenCalledWith('Navigating to edit with employee:', employee);
    });
  });

  describe('goBackToHome', () => {
    it('it succeeds when navigating back to home', () => {
      component.goBackToHome();
      
      expect(mockRouter.navigate).toHaveBeenCalledWith(['/']);
    });
  });

  describe('component properties', () => {
    it('should initialize with default values', () => {
      expect(component.showModal).toBe(false);
      expect(component.employeeToDelete).toBe(null);
      expect(component.employees).toEqual([]);
    });

    it('should set employeeToDelete property', () => {
      component.employeeToDelete = 5;
      expect(component.employeeToDelete).toBe(5);
    });

    it('should set showModal property', () => {
      component.showModal = true;
      expect(component.showModal).toBe(true);
    });
  });
});
