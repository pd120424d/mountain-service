import { ComponentFixture, TestBed } from '@angular/core/testing';
import { provideHttpClient } from '@angular/common/http';
import { provideHttpClientTesting } from '@angular/common/http/testing';
import { provideRouter } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { of, throwError } from 'rxjs';

import { UrgencyListComponent } from './urgency-list.component';
import { UrgencyService } from '../urgency.service';
import { Urgency, UrgencyLevel, UrgencyStatus } from '../../shared/models';

describe('UrgencyListComponent', () => {
  let component: UrgencyListComponent;
  let fixture: ComponentFixture<UrgencyListComponent>;
  let urgencyService: jasmine.SpyObj<UrgencyService>;

  const mockUrgencies: Urgency[] = [
    {
      id: 1,
      name: 'Test Urgency 1',
      email: 'test1@example.com',
      contactPhone: '1234567890',
      location: 'Test Location 1',
      description: 'Test Description 1',
      level: UrgencyLevel.High,
      status: UrgencyStatus.Open,
      createdAt: '2024-01-15T10:00:00Z',
      updatedAt: '2024-01-15T10:00:00Z'
    },
    {
      id: 2,
      name: 'Test Urgency 2',
      email: 'test2@example.com',
      contactPhone: '0987654321',
      location: 'Test Location 2',
      description: 'Test Description 2',
      level: UrgencyLevel.Critical,
      status: UrgencyStatus.InProgress,
      createdAt: '2024-01-16T11:00:00Z',
      updatedAt: '2024-01-16T11:00:00Z'
    }
  ];

  beforeEach(async () => {
    const urgencyServiceSpy = jasmine.createSpyObj('UrgencyService', ['getUrgencies']);

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
        { provide: UrgencyService, useValue: urgencyServiceSpy }
      ]
    })
      .compileComponents();

    fixture = TestBed.createComponent(UrgencyListComponent);
    component = fixture.componentInstance;
    urgencyService = TestBed.inject(UrgencyService) as jasmine.SpyObj<UrgencyService>;
  });

  it('should create', () => {
    urgencyService.getUrgencies.and.returnValue(of([]));
    expect(component).toBeTruthy();
  });

  it('should load urgencies on init', () => {
    urgencyService.getUrgencies.and.returnValue(of(mockUrgencies));

    component.ngOnInit();

    expect(urgencyService.getUrgencies).toHaveBeenCalled();
    expect(component.urgencies).toEqual(mockUrgencies);
    expect(component.isLoading).toBeFalse();
  });

  it('should handle error when loading urgencies', () => {
    const errorMessage = 'Failed to load urgencies';
    urgencyService.getUrgencies.and.returnValue(throwError(() => new Error(errorMessage)));

    component.ngOnInit();

    expect(component.error).toBe(errorMessage);
    expect(component.isLoading).toBeFalse();
  });

  it('should handle view urgency', () => {
    spyOn(console, 'log');

    component.viewUrgency(1);

    expect(console.log).toHaveBeenCalledWith('Viewing urgency:', 1);
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
});
