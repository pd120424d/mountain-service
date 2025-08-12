import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { provideRouter } from '@angular/router';
import { ActivatedRoute } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { ToastrService } from 'ngx-toastr';
import { of, throwError } from 'rxjs';

import { UrgencyDetailComponent } from './urgency-detail.component';
import { UrgencyService } from '../urgency.service';
import { ActivityService } from '../../services/activity.service';
import { AuthService } from '../../services/auth.service';
import { Urgency, Activity, UrgencyLevel, UrgencyStatus, ActivityLevel, ActivityType } from '../../shared/models';

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
    updatedAt: '2024-01-15T10:00:00Z'
  };

  const mockActivities: Activity[] = [
    {
      id: 1,
      type: ActivityType.UrgencyCreated,
      level: ActivityLevel.Info,
      title: 'Emergency Created',
      description: 'Emergency report was created',
      actorName: 'System',
      targetId: 1,
      targetType: 'urgency',
      createdAt: '2024-01-15T10:00:00Z'
    }
  ];

  beforeEach(async () => {
    const urgencyServiceSpy = jasmine.createSpyObj('UrgencyService', ['getUrgencyById']);
    const activityServiceSpy = jasmine.createSpyObj('ActivityService', ['getActivitiesByUrgency', 'createActivity']);
    const authServiceSpy = jasmine.createSpyObj('AuthService', ['getUserId']);
    const toastrServiceSpy = jasmine.createSpyObj('ToastrService', ['success', 'error']);
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
    activityService.getActivitiesByUrgency.and.returnValue(of(mockActivities));

    component.ngOnInit();

    expect(urgencyService.getUrgencyById).toHaveBeenCalledWith(1);
    expect(activityService.getActivitiesByUrgency).toHaveBeenCalledWith(1);
    expect(component.urgency).toEqual(mockUrgency);
    expect(component.activities).toEqual(mockActivities);
  });

  it('should handle error when loading urgency', () => {
    const errorMessage = 'Failed to load urgency';
    urgencyService.getUrgencyById.and.returnValue(throwError(() => new Error(errorMessage)));
    activityService.getActivitiesByUrgency.and.returnValue(of([]));

    component.ngOnInit();

    expect(component.error).toBe(errorMessage);
    expect(component.isLoading).toBeFalse();
  });

  it('should create activity successfully', () => {
    const newActivity: Activity = {
      id: 2,
      type: ActivityType.UrgencyUpdated,
      level: ActivityLevel.Info,
      title: 'Test Activity',
      description: 'Test Description',
      actorId: 1,
      targetId: 1,
      targetType: 'urgency',
      createdAt: '2024-01-15T11:00:00Z'
    };

    // Set up the mocks before calling ngOnInit
    urgencyService.getUrgencyById.and.returnValue(of(mockUrgency));
    activityService.getActivitiesByUrgency.and.returnValue(of(mockActivities));

    component.urgencyId = 1;
    component.ngOnInit(); // Initialize the form and load data

    component.activityForm.patchValue({
      title: 'Test Activity',
      description: 'Test Description',
      level: ActivityLevel.Info
    });

    authService.getUserId.and.returnValue('1');
    activityService.createActivity.and.returnValue(of(newActivity));
    activityService.getActivitiesByUrgency.and.returnValue(of([newActivity, ...mockActivities]));

    component.onSubmitActivity();

    expect(activityService.createActivity).toHaveBeenCalled();
    expect(toastrService.success).toHaveBeenCalled();
    expect(component.isSubmittingActivity).toBeFalse();
  });
});
