import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { provideRouter } from '@angular/router';
import { ActivatedRoute } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { ToastrService } from 'ngx-toastr';
import { of, throwError } from 'rxjs';

import { UrgencyDetailComponent } from './urgency-detail.component';
import { UrgencyService } from '../urgency.service';
import { ActivityService } from '../../services/activity.service';
import { AuthService } from '../../services/auth.service';
import { Urgency, Activity, UrgencyLevel, UrgencyStatus } from '../../shared/models';

describe('UrgencyDetailComponent', () => {
  let component: UrgencyDetailComponent;
  let fixture: ComponentFixture<UrgencyDetailComponent>;
  let urgencyService: jasmine.SpyObj<UrgencyService>;
  let activityService: jasmine.SpyObj<ActivityService>;
  let authService: jasmine.SpyObj<AuthService>;
  let toastrService: jasmine.SpyObj<ToastrService>;
  let activatedRoute: jasmine.SpyObj<ActivatedRoute>;

  const mockUrgency: Urgency = {
    id: 1,
    firstName: 'John',
    lastName: 'Doe',
    email: 'john.doe@example.com',
    contactPhone: '1234567890',
    location: 'Test Location',
    description: 'Test Description',
    level: UrgencyLevel.High,
    status: UrgencyStatus.Open,
    createdAt: '2024-01-15T10:00:00Z',
    updatedAt: '2024-01-15T10:00:00Z',
    assignedEmployeeId: undefined
  } as Urgency;

  const mockActivities: Activity[] = [
    {
      id: 1,
      description: 'Emergency report was created',
      employeeId: 1,
      urgencyId: 1,
      createdAt: '2024-01-15T10:00:00Z',
      updatedAt: '2024-01-15T10:00:00Z'
    }
  ];

  beforeEach(async () => {
    // Mock IntersectionObserver used by the component
    (window as any).IntersectionObserver = class {
      constructor(_: any) {}
      observe() {}
      unobserve() {}
      disconnect() {}
    };

    const urgencyServiceSpy = jasmine.createSpyObj('UrgencyService', ['getUrgencyById', 'assignUrgency', 'unassignUrgency', 'closeUrgency']);
    const activityServiceSpy = jasmine.createSpyObj('ActivityService', ['getActivitiesByUrgency', 'getActivitiesWithPagination', 'getActivitiesCursor', 'createActivity', 'pollForActivityInUrgency']);
    const authServiceSpy = jasmine.createSpyObj('AuthService', ['getUserId', 'isAdmin']);
    const toastrServiceSpy = jasmine.createSpyObj('ToastrService', ['success', 'error', 'warning']);
    const activatedRouteSpy = jasmine.createSpyObj('ActivatedRoute', [], {
      params: of({ id: '1' })
    });

    await TestBed.configureTestingModule({
      imports: [UrgencyDetailComponent, TranslateModule.forRoot()],
      providers: [
        provideHttpClient(),
        provideHttpClientTesting(),
        provideRouter([]),
        { provide: UrgencyService, useValue: urgencyServiceSpy },
        { provide: ActivityService, useValue: activityServiceSpy },
        { provide: AuthService, useValue: authServiceSpy },
        { provide: ToastrService, useValue: toastrServiceSpy },
        { provide: ActivatedRoute, useValue: activatedRouteSpy }
      ]
    }).compileComponents();

    fixture = TestBed.createComponent(UrgencyDetailComponent);
    component = fixture.componentInstance;
    urgencyService = TestBed.inject(UrgencyService) as jasmine.SpyObj<UrgencyService>;
    activityService = TestBed.inject(ActivityService) as jasmine.SpyObj<ActivityService>;
    authService = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;
    toastrService = TestBed.inject(ToastrService) as jasmine.SpyObj<ToastrService>;
    activatedRoute = TestBed.inject(ActivatedRoute) as jasmine.SpyObj<ActivatedRoute>;
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should load urgency and activities on init', () => {
    urgencyService.getUrgencyById.and.returnValue(of(mockUrgency));
    activityService.getActivitiesCursor.and.returnValue(of({ activities: mockActivities, nextPageToken: undefined }));

    component.ngOnInit();

    expect(urgencyService.getUrgencyById).toHaveBeenCalledWith(1);
    expect(activityService.getActivitiesCursor).toHaveBeenCalledWith({ urgencyId: 1, pageSize: 10 });
    expect(component.urgency).toEqual(mockUrgency);
    expect(component.activities).toEqual(mockActivities);
  });

  it('should handle error when loading urgency', () => {
    const errorMessage = 'Failed to load urgency';
    urgencyService.getUrgencyById.and.returnValue(throwError(() => new Error(errorMessage)));
    activityService.getActivitiesCursor.and.returnValue(of({ activities: [], nextPageToken: undefined }));

    component.ngOnInit();

    expect(component.error).toBe(errorMessage);
    expect(component.isLoading).toBeFalse();
  });

  it('should create activity successfully', () => {
    const newActivity: Activity = {
      id: 2,
      description: 'Test Description',
      employeeId: 1,
      urgencyId: 1,
      createdAt: '2024-01-15T11:00:00Z',
      updatedAt: '2024-01-15T11:00:00Z'
    };

    // Set up the mocks before calling ngOnInit
    const assignedInProgress = { ...mockUrgency, status: UrgencyStatus.InProgress, assignedEmployeeId: 1 } as Urgency;
    urgencyService.getUrgencyById.and.returnValue(of(assignedInProgress));
    activityService.getActivitiesCursor.and.returnValue(of({ activities: mockActivities, nextPageToken: undefined }));

    component.urgencyId = 1;
    component.ngOnInit(); // Initialize the form and load data

    component.activityForm.patchValue({
      description: 'Test Description'
    });

    authService.getUserId.and.returnValue('1');
    activityService.createActivity.and.returnValue(of(newActivity));
    activityService.pollForActivityInUrgency.and.returnValue(of(void 0));
    activityService.getActivitiesCursor.and.returnValue(of({ activities: [newActivity, ...mockActivities], nextPageToken: undefined }));

    component.onSubmitActivity();

    expect(activityService.createActivity).toHaveBeenCalled();
    expect(activityService.pollForActivityInUrgency).toHaveBeenCalledWith(1, 2, jasmine.any(Object));
    expect(toastrService.success).toHaveBeenCalled();
    expect(component.isSubmittingActivity).toBeFalse();
    expect(component.isSyncingActivity).toBeFalse();
  });


  it('should show syncing and refresh list after poll resolves', () => {
    const newActivity: Activity = {
      id: 99,
      description: 'Immediate note',
      employeeId: 1,
      urgencyId: 1,
      createdAt: '2024-01-15T12:00:00Z',
      updatedAt: '2024-01-15T12:00:00Z'
    };

    const assignedInProgress = { ...mockUrgency, status: UrgencyStatus.InProgress, assignedEmployeeId: 1 } as Urgency;
    urgencyService.getUrgencyById.and.returnValue(of(assignedInProgress));
    activityService.getActivitiesCursor.and.returnValue(of({ activities: mockActivities, nextPageToken: undefined }));

    component.urgencyId = 1;
    component.ngOnInit();
    component.activityForm.patchValue({ description: 'Immediate note' });
    authService.getUserId.and.returnValue('1');
    activityService.createActivity.and.returnValue(of(newActivity));
    activityService.pollForActivityInUrgency.and.returnValue(of(void 0));

    spyOn(component, 'loadActivities');

    component.onSubmitActivity();

    expect(activityService.createActivity).toHaveBeenCalled();
    expect(activityService.pollForActivityInUrgency).toHaveBeenCalledWith(1, 99, jasmine.any(Object));
    expect(component.loadActivities).toHaveBeenCalled();
    expect(component.isSyncingActivity).toBeFalse();
  });

  it('canAddActivity returns true only when assigned and InProgress', () => {
    component.urgency = {
      id: 1,
      description: 'd',
      level: UrgencyLevel.High,
      status: UrgencyStatus.InProgress,
      createdAt: new Date().toISOString(),
      assignedEmployeeId: 1
    } as any;

    expect(component.canAddActivity()).toBeTrue();

    component.urgency = { ...(component.urgency as any), status: UrgencyStatus.Open } as any;
    expect(component.canAddActivity()).toBeFalse();

    component.urgency = { ...(component.urgency as any), assignedEmployeeId: undefined, status: UrgencyStatus.InProgress } as any;
    expect(component.canAddActivity()).toBeFalse();
  });

  it('onPrevActivities and onNextActivities enforce bounds and reload', () => {
    spyOn(component, 'loadActivities');

    component.activitiesPage = 2;
    component.totalActivitiesPages = 3;

    component.onPrevActivities();
    expect(component.activitiesPage).toBe(1);
    expect(component.loadActivities).toHaveBeenCalled();

    (component.loadActivities as jasmine.Spy).calls.reset();
    component.onPrevActivities(); // at lower bound
    expect(component.activitiesPage).toBe(1);
    expect(component.loadActivities).not.toHaveBeenCalled();

    component.activitiesPage = 2;
    component.onNextActivities();
    expect(component.activitiesPage).toBe(3);
    expect(component.loadActivities).toHaveBeenCalled();

    (component.loadActivities as jasmine.Spy).calls.reset();
    component.onNextActivities(); // at upper bound
    expect(component.activitiesPage).toBe(3);
    expect(component.loadActivities).not.toHaveBeenCalled();
  });


});
