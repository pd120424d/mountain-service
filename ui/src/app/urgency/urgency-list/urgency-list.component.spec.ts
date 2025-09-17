import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { provideRouter } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { of, throwError, Subject } from 'rxjs';

import { UrgencyListComponent } from './urgency-list.component';
import { UrgencyService } from '../urgency.service';
import { Urgency, UrgencyLevel, UrgencyStatus } from '../../shared/models';
import { AuthService } from '../../services/auth.service';
import { ActivityService } from '../../services/activity.service';

describe('UrgencyListComponent', () => {
  let component: UrgencyListComponent;
  let fixture: ComponentFixture<UrgencyListComponent>;
  let urgencyService: jasmine.SpyObj<UrgencyService>;
  let activityService: jasmine.SpyObj<ActivityService>;

  const mockUrgencies: Urgency[] = [
    {
      id: 1,
      firstName: 'John',
      lastName: 'Doe',
      email: 'test1@example.com',
      contactPhone: '1234567890',
      location: 'Test Location 1',
      description: 'Test Description 1',
      level: UrgencyLevel.High,
      status: UrgencyStatus.Open,
      createdAt: '2024-01-15T10:00:00Z',
      updatedAt: '2024-01-15T10:00:00Z',
      assignedEmployeeId: undefined
    },
    {
      id: 2,
      firstName: 'Jane',
      lastName: 'Smith',
      email: 'test2@example.com',
      contactPhone: '0987654321',
      location: 'Test Location 2',
      description: 'Test Description 2',
      level: UrgencyLevel.Critical,
      status: UrgencyStatus.InProgress,
      createdAt: '2024-01-16T11:00:00Z',
      updatedAt: '2024-01-16T11:00:00Z',
      assignedEmployeeId: 456
    }
  ];

  beforeEach(async () => {
    const urgencyServiceSpy = jasmine.createSpyObj('UrgencyService', ['getUrgencies', 'getUrgenciesPaginated']);
    const activityServiceSpy = jasmine.createSpyObj('ActivityService', ['getCountsByUrgencyIds']);

    const authServiceSpy = jasmine.createSpyObj('AuthService', ['getUserId', 'isAuthenticated']);

    await TestBed.configureTestingModule({


      imports: [
        UrgencyListComponent,
        TranslateModule.forRoot()
      ],
      providers: [
        provideHttpClient(),
        provideHttpClientTesting(),
        provideRouter([]),
        TranslateService,
        { provide: UrgencyService, useValue: urgencyServiceSpy },
        { provide: ActivityService, useValue: activityServiceSpy },
        { provide: AuthService, useValue: authServiceSpy }
      ]
    })
      .compileComponents();

    fixture = TestBed.createComponent(UrgencyListComponent);
    component = fixture.componentInstance;
    const auth = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;
    auth.isAuthenticated.and.returnValue(true);
    auth.getUserId.and.returnValue('0');

    urgencyService = TestBed.inject(UrgencyService) as jasmine.SpyObj<UrgencyService>;
    activityService = TestBed.inject(ActivityService) as jasmine.SpyObj<ActivityService>;
  });

  it('should create', () => {
    urgencyService.getUrgenciesPaginated.and.returnValue(of({ urgencies: [], total: 0, page: 1, pageSize: 20, totalPages: 0 }));
    activityService.getCountsByUrgencyIds.and.returnValue(of({}));
    expect(component).toBeTruthy();
  });

  it('should load urgencies on init and fetch activity counts', () => {
    urgencyService.getUrgenciesPaginated.and.returnValue(of({ urgencies: mockUrgencies, total: mockUrgencies.length, page: 1, pageSize: 20, totalPages: 1 }));
    activityService.getCountsByUrgencyIds.and.returnValue(of({ '1': 5, '2': 3 }));

    component.ngOnInit();

    expect(urgencyService.getUrgenciesPaginated).toHaveBeenCalledWith({ page: 1, pageSize: 20, myUrgencies: true });
    expect(activityService.getCountsByUrgencyIds).toHaveBeenCalledWith([1, 2]);
    expect(component.urgencies).toEqual(mockUrgencies);
    expect(component.countsByUrgencyId[1]).toBe(5);
    expect(component.countsByUrgencyId[2]).toBe(3);
    expect(component.isLoading).toBeFalse();
  });

  it('should show spinner while counts are loading and hide after they arrive', () => {
    // Arrange urgencies load
    urgencyService.getUrgenciesPaginated.and.returnValue(of({ urgencies: mockUrgencies, total: mockUrgencies.length, page: 1, pageSize: 20, totalPages: 1 }));
    const counts$ = new Subject<Record<string, number>>();
    activityService.getCountsByUrgencyIds.and.returnValue(counts$.asObservable());

    // Act
    component.ngOnInit();
    fixture.detectChanges();

    // While counts stream is pending, spinner should be visible in the table
    expect(component['countsLoading']).toBeTrue();
    const spinnerEl = fixture.nativeElement.querySelector('td .spinner-inline');
    expect(spinnerEl).toBeTruthy();

    // Now emit counts and complete
    counts$.next({ '1': 7, '2': 0 });
    counts$.complete();
    fixture.detectChanges();

    // Spinner hidden, counts mapped
    expect(component['countsLoading']).toBeFalse();
    expect(component.countsByUrgencyId[1]).toBe(7);
    expect(component.countsByUrgencyId[2]).toBe(0);
  });

  it('should show --- when counts request fails', () => {
    urgencyService.getUrgenciesPaginated.and.returnValue(of({ urgencies: mockUrgencies, total: mockUrgencies.length, page: 1, pageSize: 20, totalPages: 1 }));
    activityService.getCountsByUrgencyIds.and.returnValue(throwError(() => new Error('boom')));

    component.ngOnInit();
    fixture.detectChanges();

    expect((component as any).countsError).toBeTrue();
    const lastCell = fixture.nativeElement.querySelector('tbody tr td:last-child');
    expect(lastCell.textContent.trim()).toBe('---');
  });


  it('should handle error when loading urgencies', () => {
    const errorMessage = 'Failed to load urgencies';
    urgencyService.getUrgenciesPaginated.and.returnValue(throwError(() => new Error(errorMessage)));

    component.ngOnInit();

    expect(component.error).toBe(errorMessage);
    expect(component.isLoading).toBeFalse();
  });

  it('should handle view urgency', () => {
    const routerSpy = spyOn(component['router'], 'navigate');

    component.viewUrgency(1);

    expect(routerSpy).toHaveBeenCalledWith(['/urgencies', 1]);
  });

  it('should return correct status class', () => {
    expect(component.getStatusClass(UrgencyStatus.Open)).toBe('status-open');
    expect(component.getStatusClass(UrgencyStatus.InProgress)).toBe('status-in-progress');
    expect(component.getStatusClass(UrgencyStatus.Resolved)).toBe('status-resolved');
    expect(component.getStatusClass(UrgencyStatus.Closed)).toBe('status-closed');
  });

  it('should return correct level class', () => {
    expect(component.getLevelClass(UrgencyLevel.Low)).toBe('level-low');
    expect(component.getLevelClass(UrgencyLevel.Medium)).toBe('level-medium');
    expect(component.getLevelClass(UrgencyLevel.High)).toBe('level-high');
    expect(component.getLevelClass(UrgencyLevel.Critical)).toBe('level-critical');
  });

  it('should compute unassignedCount correctly', () => {
    urgencyService.getUrgenciesPaginated.and.returnValue(of({ urgencies: mockUrgencies, total: mockUrgencies.length, page: 1, pageSize: 20, totalPages: 1 }));
    activityService.getCountsByUrgencyIds.and.returnValue(of({ '1': 5, '2': 3 }));
    component.ngOnInit();
    expect(component.unassignedCount).toBe(1); // one without assignedEmployeeId
  });

  it('should return correct row class based on status and assignment', () => {
    const auth = TestBed.inject(AuthService) as jasmine.SpyObj<AuthService>;
    auth.getUserId.and.returnValue('456');

    const openUnassigned = { ...mockUrgencies[0], status: UrgencyStatus.Open };
    const inProgressMine = { ...mockUrgencies[1], status: UrgencyStatus.InProgress };
    const inProgressOther = { ...mockUrgencies[1], status: UrgencyStatus.InProgress, assignedEmployeeId: 999 };
    const closed = { ...mockUrgencies[1], status: UrgencyStatus.Closed };

    expect(component.getRowClass(openUnassigned as any)).toBe('urgency-row-open-unassigned');
    expect(component.getRowClass(inProgressMine as any)).toBe('urgency-row-assigned-me');
    expect(component.getRowClass(inProgressOther as any)).toBe('urgency-row-assigned-other');
    expect(component.getRowClass(closed as any)).toBe('urgency-row-closed');
  });

  it('should paginate with onPrev and onNext correctly', () => {
    // Return response without page so component retains its local page value after load
    urgencyService.getUrgenciesPaginated.and.callFake(() => of({ urgencies: mockUrgencies, total: 50, page: component.page, pageSize: 20, totalPages: 3 }));
    component.page = 2;
    component.totalPages = 3;

    component.onPrev();
    expect(component.page).toBe(1);

    component.onNext();
    expect(component.page).toBe(2);

    // at bounds
    component.page = 1;
    component.onPrev();
    expect(component.page).toBe(1);

    component.page = 3;
    component.totalPages = 3;
    component.onNext();
    expect(component.page).toBe(3);
  });

});
