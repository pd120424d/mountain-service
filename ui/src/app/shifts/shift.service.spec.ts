import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';
import { ShiftManagementService } from './shift.service';
import {
  ShiftAvailabilityResponse,
  AssignShiftRequest,
  AssignShiftResponse,
  RemoveShiftRequest,
  Employee,
  EmployeeResponseProfileTypeEnum
} from '../shared/models';

describe('ShiftManagementService', () => {
  let service: ShiftManagementService;
  let httpMock: HttpTestingController;

  // Test data
  const mockShiftAvailability: ShiftAvailabilityResponse = {
    days: {
      '2024-01-15': {
        firstShift: { medicSlotsAvailable: 2, technicalSlotsAvailable: 1 },
        secondShift: { medicSlotsAvailable: 1, technicalSlotsAvailable: 2 },
        thirdShift: { medicSlotsAvailable: 0, technicalSlotsAvailable: 1 }
      },
      '2024-01-16': {
        firstShift: { medicSlotsAvailable: 3, technicalSlotsAvailable: 2 },
        secondShift: { medicSlotsAvailable: 2, technicalSlotsAvailable: 1 },
        thirdShift: { medicSlotsAvailable: 1, technicalSlotsAvailable: 0 }
      }
    }
  };

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

  const mockAssignResponse: AssignShiftResponse = {
    id: 123,
    shiftDate: '2024-01-15',
    shiftType: 1
  };

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        ShiftManagementService,
        provideHttpClient(),
        provideHttpClientTesting()
      ]
    });
    service = TestBed.inject(ShiftManagementService);
    httpMock = TestBed.inject(HttpTestingController);
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('getShiftAvailability', () => {
    it('should fetch shift availability with default 7 days', () => {
      service.getShiftAvailability().subscribe(response => {
        expect(response).toEqual(mockShiftAvailability);
      });

      const req = httpMock.expectOne('/api/v1/shifts/availability?days=7');
      expect(req.request.method).toBe('GET');
      req.flush(mockShiftAvailability);
    });

    it('should fetch shift availability with custom days parameter', () => {
      const customDays = 14;
      
      service.getShiftAvailability(customDays).subscribe(response => {
        expect(response).toEqual(mockShiftAvailability);
      });

      const req = httpMock.expectOne(`/api/v1/shifts/availability?days=${customDays}`);
      expect(req.request.method).toBe('GET');
      req.flush(mockShiftAvailability);
    });

    it('should handle HTTP errors gracefully', () => {
      const errorMessage = 'Server error';
      
      service.getShiftAvailability().subscribe({
        next: () => fail('Expected error'),
        error: (error) => {
          expect(error.message).toContain('Server error: 500');
        }
      });

      const req = httpMock.expectOne('/api/v1/shifts/availability?days=7');
      req.flush(errorMessage, { status: 500, statusText: 'Internal Server Error' });
    });
  });

  describe('getAllEmployees', () => {
    it('should fetch all employees', () => {
      service.getAllEmployees().subscribe(employees => {
        expect(employees).toEqual(mockEmployees);
        expect(employees.length).toBe(2);
      });

      const req = httpMock.expectOne('/api/v1/employees');
      expect(req.request.method).toBe('GET');
      req.flush(mockEmployees);
    });

    it('should handle empty employee list', () => {
      service.getAllEmployees().subscribe(employees => {
        expect(employees).toEqual([]);
      });

      const req = httpMock.expectOne('/api/v1/employees');
      req.flush([]);
    });

    it('should handle HTTP errors when fetching employees', () => {
      service.getAllEmployees().subscribe({
        next: () => fail('Expected error'),
        error: (error) => {
          expect(error.message).toContain('Server error: 404');
        }
      });

      const req = httpMock.expectOne('/api/v1/employees');
      req.flush('Not found', { status: 404, statusText: 'Not Found' });
    });
  });

  describe('assignEmployeeToShift', () => {
    it('should assign employee to shift successfully', () => {
      const shiftType = 1;
      const employeeId = 'emp123';
      const date = new Date('2024-01-15');

      service.assignEmployeeToShift(shiftType, employeeId, date).subscribe(response => {
        expect(response).toEqual(mockAssignResponse);
      });

      const req = httpMock.expectOne(`/api/v1/employees/${employeeId}/shifts`);
      expect(req.request.method).toBe('POST');
      expect(req.request.body).toEqual({
        shiftDate: '2024-01-15',
        shiftType: 1
      } as AssignShiftRequest);
      req.flush(mockAssignResponse);
    });

    it('should handle assignment errors', () => {
      const shiftType = 1;
      const employeeId = 'emp123';
      const date = new Date('2024-01-15');

      service.assignEmployeeToShift(shiftType, employeeId, date).subscribe({
        next: () => fail('Expected error'),
        error: (error) => {
          expect(error.message).toContain('Server error: 409');
        }
      });

      const req = httpMock.expectOne(`/api/v1/employees/${employeeId}/shifts`);
      req.flush('Conflict - Employee already assigned', { status: 409, statusText: 'Conflict' });
    });
  });

  describe('removeEmployeeFromShiftByDetails', () => {
    it('should remove employee from shift by details', () => {
      const employeeId = 'emp123';
      const shiftType = 2;
      const date = new Date('2024-01-16');

      service.removeEmployeeFromShiftByDetails(employeeId, shiftType, date).subscribe(response => {
        expect(response).toBeTruthy();
      });

      const req = httpMock.expectOne(`/api/v1/employees/${employeeId}/shifts`);
      expect(req.request.method).toBe('DELETE');
      expect(req.request.body).toEqual({
        shiftDate: '2024-01-16',
        shiftType: 2
      } as RemoveShiftRequest);
      req.flush({});
    });
  });

  describe('error handling', () => {
    it('should handle client-side errors', () => {
      const clientError = new ErrorEvent('Network error', {
        message: 'Connection failed'
      });

      service.getShiftAvailability().subscribe({
        next: () => fail('Expected error'),
        error: (error) => {
          expect(error.message).toContain('Client error: Connection failed');
        }
      });

      const req = httpMock.expectOne('/api/v1/shifts/availability?days=7');
      req.error(clientError);
    });

    it('should handle server errors with custom error message', () => {
      const errorResponse = { message: 'Custom server error' };

      service.getShiftAvailability().subscribe({
        next: () => fail('Expected error'),
        error: (error) => {
          expect(error.message).toContain('Custom server error');
        }
      });

      const req = httpMock.expectOne('/api/v1/shifts/availability?days=7');
      req.flush(errorResponse, { status: 500, statusText: 'Internal Server Error' });
    });

    it('should handle server errors with string error message', () => {
      service.getShiftAvailability().subscribe({
        next: () => fail('Expected error'),
        error: (error) => {
          expect(error.message).toContain('String error message');
        }
      });

      const req = httpMock.expectOne('/api/v1/shifts/availability?days=7');
      req.flush('String error message', { status: 400, statusText: 'Bad Request' });
    });
  });
});
