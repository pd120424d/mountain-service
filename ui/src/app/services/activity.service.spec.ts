import { TestBed } from '@angular/core/testing';
import { HttpTestingController, provideHttpClientTesting } from '@angular/common/http/testing';
import { provideHttpClient } from '@angular/common/http';
import { ActivityService } from './activity.service';
import { Activity } from '../shared/models';
import { environment } from '../../environments/environment';

describe('ActivityService', () => {
  let service: ActivityService;
  let httpMock: HttpTestingController;
  let expectedBaseUrl: string;
  let expectedActivityUrl: string;

  beforeEach(() => {
    TestBed.configureTestingModule({
      providers: [
        ActivityService,
        provideHttpClient(),
        provideHttpClientTesting()
      ]
    });
    service = TestBed.inject(ActivityService);
    httpMock = TestBed.inject(HttpTestingController);

    expectedBaseUrl = environment.useMockApi ? '/api/v1' : environment.apiUrl;
    expectedActivityUrl = expectedBaseUrl + '/activities';
  });

  afterEach(() => {
    httpMock.verify();
  });

  it('should be created', () => {
    expect(service).toBeTruthy();
  });

  describe('getActivities', () => {
    it('should fetch all activities', () => {
      const mockActivities: Activity[] = [
        {
          id: 1,
          description: 'Emergency report was created',
          employeeId: 1,
          urgencyId: 1,
          createdAt: '2024-01-15T10:00:00Z',
          updatedAt: '2024-01-15T10:00:00Z'
        },
        {
          id: 2,
          description: 'Emergency report was updated',
          employeeId: 1,
          urgencyId: 1,
          createdAt: '2024-01-15T11:00:00Z',
          updatedAt: '2024-01-15T11:00:00Z'
        }
      ];

      service.getActivities().subscribe(activities => {
        expect(activities).toEqual(mockActivities);
        expect(activities.length).toBe(2);
      });

      const req = httpMock.expectOne(expectedActivityUrl);
      expect(req.request.method).toBe('GET');
      req.flush({ activities: mockActivities, total: 2, page: 1, pageSize: 20, totalPages: 1 });
    });

    it('should handle error when fetching activities', () => {
      service.getActivities().subscribe({
        next: () => fail('Expected an error'),
        error: (error) => {
          expect(error.message).toBe('Server error');
        }
      });

      const req = httpMock.expectOne(expectedActivityUrl);
      req.flush({ error: 'Server error' }, { status: 500, statusText: 'Server Error' });
    });
  });

  describe('getActivitiesByUrgency', () => {
    it('should fetch activities for specific urgency', () => {
      const urgencyId = 123;
      const mockActivities: Activity[] = [
        {
          id: 1,
          description: 'Emergency report was created',
          employeeId: 1,
          urgencyId: urgencyId,
          createdAt: '2024-01-15T10:00:00Z',
          updatedAt: '2024-01-15T10:00:00Z'
        }
      ];

      service.getActivitiesByUrgency(urgencyId).subscribe(activities => {
        expect(activities).toEqual(mockActivities);
        expect(activities[0].urgencyId).toBe(urgencyId);
      });

      const expectedUrl = `${expectedActivityUrl}?urgencyId=${urgencyId}`;
      const req = httpMock.expectOne(expectedUrl);
      expect(req.request.method).toBe('GET');
      req.flush({ activities: mockActivities, total: 1, page: 1, pageSize: 20, totalPages: 1 });
    });
  });

	  describe('getCountsByUrgencyIds', () => {
	    it('should fetch counts for provided urgency IDs', () => {
	      const ids = [1, 2, 5];
	      const expectedUrl = `${expectedActivityUrl}/counts?urgencyId=1&urgencyId=2&urgencyId=5`;

	      service.getCountsByUrgencyIds(ids).subscribe(counts => {
	        expect(counts).toEqual({ '1': 10, '2': 0, '5': 3 } as any);
	      });

	      const req = httpMock.expectOne(expectedUrl);
	      expect(req.request.method).toBe('GET');
	      req.flush({ counts: { '1': 10, '2': 0, '5': 3 } });
	    });

	    it('should return empty object and not call HTTP when ids are empty', () => {
	      service.getCountsByUrgencyIds([]).subscribe(counts => {
	        expect(counts).toEqual({});
	      });
	      httpMock.expectNone((req) => req.url.includes('/activities/counts'));
	    });
	  });


  describe('getActivitiesWithPagination', () => {
    it('should fetch activities with pagination parameters', () => {
      const params = { urgencyId: 123, employeeId: 1, page: 1, pageSize: 10 };

      const mockResponse = {
        activities: [
          {
            id: 1,
            description: 'Emergency report was created',
            employeeId: 1,
            urgencyId: 123,
            createdAt: '2024-01-15T10:00:00Z',
            updatedAt: '2024-01-15T10:00:00Z'
          }
        ],
        total: 1,
        page: 1,
        pageSize: 10,
        totalPages: 1,
      };

      service.getActivitiesWithPagination(params).subscribe(response => {
        expect(response).toEqual(mockResponse);
        expect(response.activities.length).toBe(1);
        expect(response.total).toBe(1);
        expect(response.page).toBe(1);
      });

      const expectedUrl = `${expectedActivityUrl}?urgencyId=123&employeeId=1&page=1&pageSize=10`;
      const req = httpMock.expectOne(expectedUrl);
      expect(req.request.method).toBe('GET');
      req.flush(mockResponse);
    });

    it('should fetch activities without parameters', () => {
      const mockResponse = { activities: [], total: 0, page: 1, pageSize: 20, totalPages: 0 };

      service.getActivitiesWithPagination({}).subscribe(response => {
        expect(response).toEqual(mockResponse);
      });

      const req = httpMock.expectOne(expectedActivityUrl);
      expect(req.request.method).toBe('GET');
      req.flush(mockResponse);
    });
  });

  describe('createActivity', () => {
    it('should create a new activity', () => {
      const activityRequest = { description: 'Test activity description', employeeId: 1, urgencyId: 123 };

      const mockResponse: Activity = {
        id: 1,
        description: 'Test activity description',
        employeeId: 1,
        urgencyId: 123,
        createdAt: '2024-01-15T10:00:00Z',
        updatedAt: '2024-01-15T10:00:00Z'
      };

      service.createActivity(activityRequest).subscribe(activity => {
        expect(activity).toEqual(mockResponse);
        expect(activity.id).toBe(1);
        expect(activity.description).toBe(activityRequest.description);
      });

      const req = httpMock.expectOne(expectedActivityUrl);
      expect(req.request.method).toBe('POST');
      expect(req.request.body).toEqual(activityRequest);
      req.flush(mockResponse);
    });

    it('should handle error when creating activity', () => {
      const activityRequest = { description: 'Test activity description', employeeId: 1, urgencyId: 123 };

      service.createActivity(activityRequest).subscribe({
        next: () => fail('Expected an error'),
        error: (error) => {
          expect(error.message).toBe('Validation failed');
        }
      });

      const req = httpMock.expectOne(expectedActivityUrl);
      req.flush({ error: 'Validation failed' }, { status: 400, statusText: 'Bad Request' });
    });
  });

  describe('error handling', () => {
    it('should handle 409 conflict error', () => {
      const activityRequest = { description: 'Test activity description', employeeId: 1, urgencyId: 123 };

      service.createActivity(activityRequest).subscribe({
        next: () => fail('Expected an error'),
        error: (error) => {
          expect(error.message).toBe('Resource already exists');
        }
      });

      const req = httpMock.expectOne(expectedActivityUrl);
      req.flush({ error: 'Resource already exists' }, { status: 409, statusText: 'Conflict' });
    });

    it('should handle client-side error', () => {
      service.getActivities().subscribe({
        next: () => fail('Expected an error'),
        error: (error) => {
          expect(error.message).toBe('Something went wrong. Please try again later.');
        }
      });

      const req = httpMock.expectOne(expectedActivityUrl);
      req.error(new ProgressEvent('Network error'), { status: 0, statusText: 'Unknown Error' });
    });

    it('should handle generic server error', () => {
      service.getActivities().subscribe({
        next: () => fail('Expected an error'),
        error: (error) => {
          expect(error.message).toBe('Something went wrong. Please try again later.');
        }
      });

      const req = httpMock.expectOne(expectedActivityUrl);
      req.flush('Server Error', { status: 500, statusText: 'Internal Server Error' });
    });
  });
});
