import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';
import { EmployeeService } from './employee.service';
import { Employee, MedicRole, TechnicalRole } from './employee.model';
import { environment } from '../../environments/environment';
import { LoggingService } from '../services/logging.service';

describe('EmployeeService', () => {
  let service: EmployeeService;
  let httpMock: HttpTestingController;
  let loggingServiceSpy: jasmine.SpyObj<LoggingService>;
  let expectedBaseUrl: string;
  let expectedEmployeeUrl: string;

  beforeEach(() => {
    const spy = jasmine.createSpyObj('LoggingService', ['info', 'warn', 'error', 'debug']);

    TestBed.configureTestingModule({
      providers: [
        EmployeeService,
        provideHttpClient(),
        provideHttpClientTesting(),
        { provide: LoggingService, useValue: spy }
      ]
    });
    service = TestBed.inject(EmployeeService);
    httpMock = TestBed.inject(HttpTestingController);
    loggingServiceSpy = TestBed.inject(LoggingService) as jasmine.SpyObj<LoggingService>;

    expectedBaseUrl = environment.useMockApi ? '/api/v1' : environment.apiUrl;
    expectedEmployeeUrl = expectedBaseUrl + '/employees';
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created and log initialization', () => {
    expect(service).toBeTruthy();
    expect(loggingServiceSpy.info).toHaveBeenCalledWith(`Starting employee service with url: ${expectedEmployeeUrl}`);
    expect(loggingServiceSpy.info).toHaveBeenCalledWith(`Starting employee service with base apiUrl: ${expectedBaseUrl}`);
  });

  it('should fetch employees', () => {
    const mockEmployees: Employee[] = [
      {
        id: 1,
        firstName: 'John',
        lastName: 'Doe',
        email: 'john.doe@example.com',
        phoneNumber: '+1234567890',
        profileType: MedicRole,
        username: 'johndoe',
        gender: 'Male'
      },
      {
        id: 2,
        firstName: 'Jane',
        lastName: 'Smith',
        email: 'jane.smith@example.com',
        phoneNumber: '+1234567891',
        profileType: TechnicalRole,
        username: 'janesmith',
        gender: 'Female'
      }
    ];

    service.getEmployees().subscribe(employees => {
      expect(employees).toEqual(mockEmployees);
    });

    const req = httpMock.expectOne(expectedEmployeeUrl);
    expect(req.request.method).toBe('GET');
    req.flush(mockEmployees);
  });

  it('should fetch employee by ID', () => {
    const mockEmployee: Employee = {
      id: 1,
      firstName: 'John',
      lastName: 'Doe',
      email: 'john.doe@example.com',
      phoneNumber: '+1234567890',
      profileType: MedicRole,
      username: 'johndoe',
      gender: 'Male'
    };

    service.getEmployeeById(1).subscribe(employee => {
      expect(employee).toEqual(mockEmployee);
    });

    expect(loggingServiceSpy.info).toHaveBeenCalledWith('Fetching employee with ID: 1');
    const req = httpMock.expectOne(`${expectedEmployeeUrl}/1`);
    expect(req.request.method).toBe('GET');
    req.flush(mockEmployee);
  });

  it('should add employee', () => {
    const mockEmployee: Employee = {
      id: 1,
      firstName: 'John',
      lastName: 'Doe',
      email: 'john.doe@example.com',
      phoneNumber: '+1234567890',
      profileType: MedicRole,
      username: 'johndoe',
      gender: 'Male'
    };

    service.addEmployee(mockEmployee).subscribe(employee => {
      expect(employee).toEqual(mockEmployee);
    });

    const req = httpMock.expectOne(expectedEmployeeUrl);
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual(mockEmployee);
    req.flush(mockEmployee);
  });

  it('should update employee', () => {
    const mockEmployee: Employee = {
      id: 1,
      firstName: 'John',
      lastName: 'Doe',
      email: 'john.doe@example.com',
      phoneNumber: '+1234567890',
      profileType: MedicRole,
      username: 'johndoe',
      gender: 'Male'
    };

    service.updateEmployee(1, mockEmployee).subscribe(employee => {
      expect(employee).toEqual(mockEmployee);
    });

    const req = httpMock.expectOne(`${expectedEmployeeUrl}/1`);
    expect(req.request.method).toBe('PUT');
    expect(req.request.body).toEqual(mockEmployee);
    req.flush(mockEmployee);
  });

  it('should delete employee', () => {
    service.deleteEmployee(1).subscribe(() => {
      expect(true).toBe(true);
    });

    const req = httpMock.expectOne(`${expectedEmployeeUrl}/1`);
    expect(req.request.method).toBe('DELETE');
    req.flush({});
  });

  it('should handle error when fetching employees fails', () => {
    const errorMessage = 'Server error';

    service.getEmployees().subscribe({
      next: () => fail('Expected an error'),
      error: (error) => {
        expect(error.message).toBe('Something went wrong; please try again later.');
      }
    });

    const req = httpMock.expectOne(expectedEmployeeUrl);
    req.flush(errorMessage, { status: 500, statusText: 'Server Error' });
  });

  it('should handle error when fetching employee by ID fails', () => {
    const errorMessage = 'Employee not found';

    service.getEmployeeById(999).subscribe({
      next: () => fail('Expected an error'),
      error: (error) => {
        expect(error.message).toBe('Something went wrong; please try again later.');
      }
    });

    expect(loggingServiceSpy.info).toHaveBeenCalledWith('Fetching employee with ID: 999');
    const req = httpMock.expectOne(`${expectedEmployeeUrl}/999`);
    req.flush(errorMessage, { status: 404, statusText: 'Not Found' });
  });
});
