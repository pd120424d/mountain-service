import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';
import { HttpErrorResponse } from '@angular/common/http';
import { EmployeeService } from './employee.service';
import { Employee, EmployeeCreateRequest, EmployeeUpdateRequest, EmployeeResponseProfileTypeEnum, EmployeeCreateRequestProfileTypeEnum, EmployeeUpdateRequestProfileTypeEnum } from '../shared/models';
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
        phone: '+1234567890',
        profileType: EmployeeResponseProfileTypeEnum.Medic,
        username: 'johndoe',
        gender: 'Male'
      },
      {
        id: 2,
        firstName: 'Marko',
        lastName: 'Markovic',
        email: 'marko.markovic@example.com',
        phone: '+1234567891',
        profileType: EmployeeResponseProfileTypeEnum.Technical,
        username: 'markomarkovic',
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
      phone: '+1234567890',
      profileType: EmployeeResponseProfileTypeEnum.Medic,
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
    const mockEmployeeCreateRequest: EmployeeCreateRequest = {
      firstName: 'John',
      lastName: 'Doe',
      email: 'john.doe@example.com',
      phone: '+1234567890',
      profileType: EmployeeCreateRequestProfileTypeEnum.Medic,
      username: 'johndoe',
      password: 'password123',
      gender: 'Male'
    };

    const mockEmployeeResponse: Employee = {
      id: 1,
      firstName: 'John',
      lastName: 'Doe',
      email: 'john.doe@example.com',
      phone: '+1234567890',
      profileType: EmployeeResponseProfileTypeEnum.Medic,
      username: 'johndoe',
      gender: 'Male'
    };

    service.addEmployee(mockEmployeeCreateRequest).subscribe(employee => {
      expect(employee).toEqual(mockEmployeeResponse);
    });

    const req = httpMock.expectOne(expectedEmployeeUrl);
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual(mockEmployeeCreateRequest);
    req.flush(mockEmployeeResponse);
  });

  it('should update employee', () => {
    const mockUpdateRequest: EmployeeUpdateRequest = {
      firstName: 'John',
      lastName: 'Doe',
      email: 'john.doe@example.com',
      phone: '+1234567890',
      profileType: EmployeeUpdateRequestProfileTypeEnum.Medic,
      username: 'johndoe',
      gender: 'Male'
    };

    const mockEmployee: Employee = {
      id: 1,
      firstName: 'John',
      lastName: 'Doe',
      email: 'john.doe@example.com',
      phone: '+1234567890',
      profileType: EmployeeResponseProfileTypeEnum.Medic,
      username: 'johndoe',
      gender: 'Male'
    };

    service.updateEmployee(1, mockUpdateRequest).subscribe(employee => {
      expect(employee).toEqual(mockEmployee);
    });

    const req = httpMock.expectOne(`${expectedEmployeeUrl}/1`);
    expect(req.request.method).toBe('PUT');
    expect(req.request.body).toEqual(mockUpdateRequest);
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
        expect(error.message).toBe('Something went wrong. Please try again later.');
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
        expect(error.message).toBe('Something went wrong. Please try again later.');
      }
    });

    expect(loggingServiceSpy.info).toHaveBeenCalledWith('Fetching employee with ID: 999');
    const req = httpMock.expectOne(`${expectedEmployeeUrl}/999`);
    req.flush(errorMessage, { status: 404, statusText: 'Not Found' });
  });

  it('should handle 409 conflict error', () => {
    const mockEmployeeCreateRequest: EmployeeCreateRequest = {
      firstName: 'John',
      lastName: 'Doe',
      email: 'john.doe@example.com',
      phone: '+1234567890',
      profileType: EmployeeCreateRequestProfileTypeEnum.Medic,
      username: 'johndoe',
      password: 'password123',
      gender: 'Male'
    };

    const errorResponse = new HttpErrorResponse({
      error: { error: 'Employee already exists' },
      status: 409,
      statusText: 'Conflict'
    });

    service.addEmployee(mockEmployeeCreateRequest).subscribe({
      next: () => fail('should have failed with conflict error'),
      error: (error) => {
        expect(error.message).toBe('Employee already exists');
      }
    });

    const req = httpMock.expectOne(expectedEmployeeUrl);
    req.flush({ error: 'Employee already exists' }, errorResponse);
  });

  it('should handle 400 bad request error', () => {
    const mockEmployeeCreateRequest: EmployeeCreateRequest = {
      firstName: 'John',
      lastName: 'Doe',
      email: 'john.doe@example.com',
      phone: '+1234567890',
      profileType: EmployeeCreateRequestProfileTypeEnum.Medic,
      username: 'johndoe',
      password: 'password123',
      gender: 'Male'
    };

    const errorResponse = new HttpErrorResponse({
      error: { error: 'Invalid data provided' },
      status: 400,
      statusText: 'Bad Request'
    });

    service.addEmployee(mockEmployeeCreateRequest).subscribe({
      next: () => fail('should have failed with bad request error'),
      error: (error) => {
        expect(error.message).toBe('Invalid data provided');
      }
    });

    const req = httpMock.expectOne(expectedEmployeeUrl);
    req.flush({ error: 'Invalid data provided' }, errorResponse);
  });

  it('should handle client-side error', () => {
    const errorEvent = new ErrorEvent('Network error', {
      message: 'Connection failed'
    });

    service.getEmployees().subscribe({
      next: () => fail('should have failed with client error'),
      error: (error) => {
        expect(error.message).toBe('Client error: Connection failed');
      }
    });

    const req = httpMock.expectOne(expectedEmployeeUrl);
    req.error(errorEvent);
  });
});
