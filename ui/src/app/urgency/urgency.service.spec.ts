
import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';
import { UrgencyService } from './urgency.service';
import { Urgency, UrgencyCreateRequest, UrgencyUpdateRequest, GeneratedUrgencyLevel, GeneratedUrgencyStatus } from '../shared/models';
import { environment } from '../../environments/environment';

describe('UrgencyService', () => {
  let service: UrgencyService;
  let httpMock: HttpTestingController;
  let expectedBaseUrl: string;
  let expectedUrgencyUrl: string;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        UrgencyService,
        provideHttpClient(),
        provideHttpClientTesting()
      ]
    });
    service = TestBed.inject(UrgencyService);
    httpMock = TestBed.inject(HttpTestingController);

    expectedBaseUrl = environment.useMockApi ? '/api/v1' : environment.apiUrl;
    expectedUrgencyUrl = expectedBaseUrl + '/urgencies';
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  it('should fetch urgencies', () => {
    const mockUrgencies: Urgency[] = [
      {
        id: 1,
        firstName: 'John',
        lastName: 'Doe',
        email: 'urgency1@example.com',
        contactPhone: '1234567890',
        location: 'Location 1',
        description: 'Description 1',
        level: GeneratedUrgencyLevel.Medium,
        status: GeneratedUrgencyStatus.Open,
        createdAt: '2024-01-15',
        updatedAt: '2024-01-15'
      },
      {
        id: 2,
        firstName: 'Jane',
        lastName: 'Smith',
        email: 'urgency2@example.com',
        contactPhone: '0987654321',
        location: 'Location 2',
        description: 'Description 2',
        level: GeneratedUrgencyLevel.High,
        status: GeneratedUrgencyStatus.InProgress,
        createdAt: '2024-01-16',
        updatedAt: '2024-01-16'
      }
    ];

    service.getUrgencies().subscribe(urgencies => {
      expect(urgencies).toEqual(mockUrgencies);
    });

    const req = httpMock.expectOne(expectedUrgencyUrl);
    expect(req.request.method).toBe('GET');
    req.flush(mockUrgencies);
  });

  it('should fetch urgency by ID', () => {
    const mockUrgency: Urgency = {
      id: 1,
      firstName: 'John',
      lastName: 'Doe',
      email: 'urgency1@example.com',
      contactPhone: '1234567890',
      location: 'Location 1',
      description: 'Description 1',
      level: GeneratedUrgencyLevel.Medium,
      status: GeneratedUrgencyStatus.Open,
      createdAt: '2024-01-15',
      updatedAt: '2024-01-15'
    };

    service.getUrgencyById(1).subscribe(urgency => {
      expect(urgency).toEqual(mockUrgency);
    });

    const req = httpMock.expectOne(`${expectedUrgencyUrl}/1`);
    expect(req.request.method).toBe('GET');
    req.flush(mockUrgency);
  });

  it('should add urgency', () => {
    const mockUrgencyCreateRequest: UrgencyCreateRequest = {
      firstName: 'John',
      lastName: 'Doe',
      email: 'urgency1@example.com',
      contactPhone: '1234567890',
      location: 'Location 1',
      description: 'Description 1',
      level: GeneratedUrgencyLevel.Medium
    };

    const mockUrgencyResponse: Urgency = {
      id: 1,
      firstName: 'John',
      lastName: 'Doe',
      email: 'urgency1@example.com',
      contactPhone: '1234567890',
      location: 'Location 1',
      description: 'Description 1',
      level: GeneratedUrgencyLevel.Medium,
      status: GeneratedUrgencyStatus.Open,
      createdAt: '2024-01-15',
      updatedAt: '2024-01-15'
    };

    service.addUrgency(mockUrgencyCreateRequest).subscribe(urgency => {
      expect(urgency).toEqual(mockUrgencyResponse);
    });

    const req = httpMock.expectOne(expectedUrgencyUrl);
    expect(req.request.method).toBe('POST');
    expect(req.request.body).toEqual(mockUrgencyCreateRequest);
    req.flush(mockUrgencyResponse);
  });

  it('should update urgency', () => {
    const mockUrgencyUpdateRequest: UrgencyUpdateRequest = {
      firstName: 'John',
      lastName: 'Doe',
      email: 'urgency1@example.com',
      contactPhone: '1234567890',
      location: 'Location 1',
      description: 'Description 1',
      level: GeneratedUrgencyLevel.Medium,
      status: GeneratedUrgencyStatus.InProgress
    };

    const mockUrgencyResponse: Urgency = {
      id: 1,
      firstName: 'John',
      lastName: 'Doe',
      email: 'urgency1@example.com',
      contactPhone: '1234567890',
      location: 'Location 1',
      description: 'Description 1',
      level: GeneratedUrgencyLevel.Medium,
      status: GeneratedUrgencyStatus.InProgress,
      createdAt: '2024-01-15',
      updatedAt: '2024-01-15'
    };

    service.updateUrgency(1, mockUrgencyUpdateRequest).subscribe(urgency => {
      expect(urgency).toEqual(mockUrgencyResponse);
    });

    const req = httpMock.expectOne(`${expectedUrgencyUrl}/1`);
    expect(req.request.method).toBe('PUT');
    expect(req.request.body).toEqual(mockUrgencyUpdateRequest);
    req.flush(mockUrgencyResponse);
  });

  it('should delete urgency', () => {
    service.deleteUrgency(1).subscribe(() => {
      expect(true).toBe(true);
    });

    const req = httpMock.expectOne(`${expectedUrgencyUrl}/1`);
    expect(req.request.method).toBe('DELETE');
    req.flush({});
  });

  it('should handle error when fetching urgencies fails', () => {
    const errorMessage = 'Server error';

    service.getUrgencies().subscribe({
      next: () => fail('Expected an error'),
      error: (error) => {
        expect(error.message).toBe('Something went wrong. Please try again later.');
      }
    });

    const req = httpMock.expectOne(expectedUrgencyUrl);
    req.flush(errorMessage, { status: 500, statusText: 'Server Error' });


  });

  it('it returns an error when backend provides typed error code', () => {
    service.addUrgency({
      firstName: 'John', lastName: 'Doe', email: 'u@example.com', contactPhone: '1', location: 'L', description: 'D', level: GeneratedUrgencyLevel.Medium
    }).subscribe({
      next: () => fail('Expected error'),
      error: (error) => {
        expect(error.message).toBe('VALIDATION_ERROR.REQUIRED_FIELD');
      }
    });

    const req = httpMock.expectOne(expectedUrgencyUrl);
    req.flush({ error: 'VALIDATION_ERROR.REQUIRED_FIELD' }, { status: 400, statusText: 'Bad Request' });
  });

});
















