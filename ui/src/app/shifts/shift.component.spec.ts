import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ShiftManagementComponent } from './shift.component';
import { sharedTestingProviders } from '../test-utils/shared-test-imports';
import { DatePipe } from '@angular/common';
import { of, throwError } from 'rxjs';

describe('ShiftManagementComponent', () => {
  let component: ShiftManagementComponent;
  let fixture: ComponentFixture<ShiftManagementComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ShiftManagementComponent],
      providers: [...sharedTestingProviders, DatePipe]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ShiftManagementComponent);
    component = fixture.componentInstance;

    // Mock services before detectChanges to prevent real HTTP calls
    spyOn(component['shiftService'], 'getShiftAvailability').and.returnValue(of({ days: {} }));
    spyOn(component['shiftService'], 'getAdminShiftAvailability').and.returnValue(of({ days: {} }));
    spyOn(component['shiftService'], 'getShiftWarnings').and.returnValue(of({ warnings: [] }));
    spyOn(component['auth'], 'getRole').and.returnValue('Technical');
    spyOn(component['auth'], 'getUserId').and.returnValue('1');

    // Mock the global confirm function to prevent hanging dialogs
    spyOn(window, 'confirm').and.returnValue(true);

    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  describe('Admin shift assignment logic', () => {
    beforeEach(() => {
      // Set up admin user
      (component['auth'].getRole as jasmine.Spy).and.returnValue('Administrator');
      component.userRole = 'Administrator';

      // Mock employee shifts data
      const mockEmployeeShifts = [
        { id: 1, shiftDate: '2024-01-15', shiftType: 1, createdAt: '2024-01-15T00:00:00Z' },
        { id: 2, shiftDate: '2024-01-15', shiftType: 2, createdAt: '2024-01-15T00:00:00Z' }
      ];
      component.selectedEmployeeShifts.set('123', mockEmployeeShifts);
    });

    it('should correctly identify when an employee is assigned to a shift', () => {
      const testDate = new Date('2024-01-15');

      // Employee 123 is assigned to shift type 1 on 2024-01-15
      expect(component.isAssignedToShift(1, testDate, '123')).toBe(true);

      // Employee 123 is assigned to shift type 2 on 2024-01-15
      expect(component.isAssignedToShift(2, testDate, '123')).toBe(true);

      // Employee 123 is NOT assigned to shift type 3 on 2024-01-15
      expect(component.isAssignedToShift(3, testDate, '123')).toBe(false);

      // Employee 456 (not in cache) is NOT assigned to any shift
      expect(component.isAssignedToShift(1, testDate, '456')).toBe(false);
    });

    it('should disable assign button when employee is already assigned', () => {
      const testDate = new Date('2024-01-15');

      // Should not be able to assign employee 123 to shift 1 (already assigned)
      expect(component.canAssignToShift(1, testDate, '123')).toBe(false);

      // Should be able to assign employee 123 to shift 3 (not assigned)
      expect(component.canAssignToShift(3, testDate, '123')).toBe(true);
    });

    it('should enable remove button when employee is assigned', () => {
      const testDate = new Date('2024-01-15');

      // Should be able to remove employee 123 from shift 1 (assigned)
      expect(component.canRemoveFromShift(1, testDate, '123')).toBe(true);

      // Should NOT be able to remove employee 123 from shift 3 (not assigned)
      expect(component.canRemoveFromShift(3, testDate, '123')).toBe(false);
    });
  });

  it('should fetch role, userid and load shifts on init when user is non administrator', () => {
    // Services are already mocked in beforeEach, just verify the calls and state
    expect(component['shiftService'].getShiftAvailability).toHaveBeenCalled();
    expect(component.userRole).toBe('Technical');
    expect(component.userId).toBe('1');
  });

  it('should fetch role, userid and load shifts and employees on init when user is administrator', () => {
    // Reset the existing spies and set up for administrator test
    (component['shiftService'].getShiftAvailability as jasmine.Spy).calls.reset();
    (component['shiftService'].getAdminShiftAvailability as jasmine.Spy).calls.reset();
    (component['shiftService'].getShiftWarnings as jasmine.Spy).calls.reset();
    (component['auth'].getRole as jasmine.Spy).and.returnValue('Administrator');
    spyOn(component['shiftService'], 'getAllEmployees').and.returnValue(of([]));

    // Call ngOnInit to test administrator flow
    component.ngOnInit();

    expect(component['shiftService'].getAdminShiftAvailability).toHaveBeenCalled();
    expect(component['shiftService'].getAllEmployees).toHaveBeenCalled();
    expect(component.userRole).toBe('Administrator');
    expect(component.userId).toBe('1');
  });

  it('should load shifts', () => {
    // Reset the spy and set new return value
    (component['shiftService'].getShiftAvailability as jasmine.Spy).and.returnValue(of({ days: {} }));
    component.loadShiftAvailability();
    expect(component['shiftService'].getShiftAvailability).toHaveBeenCalled();
  });

  it('should load all employees', () => {
    spyOn(component['shiftService'], 'getAllEmployees').and.returnValue(of([]));
    component.loadAllEmployees();
    expect(component['shiftService'].getAllEmployees).toHaveBeenCalled();
  });

  it('should check if user can modify others', () => {
    component.userRole = 'Administrator';
    expect(component.canModifyOthers()).toBe(true);
    component.userRole = 'Technical';
    expect(component.canModifyOthers()).toBe(false);
  });

  it('should assign to shift', () => {
    spyOn(component['shiftService'], 'assignEmployeeToShift').and.returnValue(of({
      id: 123,
      shiftDate: '2024-01-15',
      shiftType: 1
    }));
    // Reset the existing spy calls instead of creating new spies
    (component['toastr'].success as jasmine.Spy).calls.reset();
    (component['spinner'].show as jasmine.Spy).calls.reset();
    (component['spinner'].hide as jasmine.Spy).calls.reset();

    component.assignToShift(1, new Date(), '1');
    expect(component['shiftService'].assignEmployeeToShift).toHaveBeenCalled();
    expect(component['toastr'].success).toHaveBeenCalled();
    expect(component['spinner'].show).toHaveBeenCalled();
    expect(component['spinner'].hide).toHaveBeenCalled();
  });

  it('should remove from shift', () => {
    spyOn(component['shiftService'], 'removeEmployeeFromShiftByDetails').and.returnValue(of({}));
    // Reset the existing spy calls instead of creating new spies
    (component['toastr'].success as jasmine.Spy).calls.reset();
    (component['spinner'].show as jasmine.Spy).calls.reset();
    (component['spinner'].hide as jasmine.Spy).calls.reset();

    component.removeFromShift(1, '1', new Date());
    expect(component['shiftService'].removeEmployeeFromShiftByDetails).toHaveBeenCalled();
    expect(component['toastr'].success).toHaveBeenCalled();
    expect(component['spinner'].show).toHaveBeenCalled();
    expect(component['spinner'].hide).toHaveBeenCalled();
  });

  it('should get available medics', () => {
    component.shiftAvailability = {
      days: {
        '2024-01-15T00:00:00Z': {
          firstShift: { medicSlotsAvailable: 2, technicalSlotsAvailable: 1 },
          secondShift: { medicSlotsAvailable: 1, technicalSlotsAvailable: 2 },
          thirdShift: { medicSlotsAvailable: 0, technicalSlotsAvailable: 1 }
        }
      }
    };
    expect(component.getAvailableMedics(1, new Date('2024-01-15'))).toBe(2);
    expect(component.getAvailableMedics(2, new Date('2024-01-15'))).toBe(1);
    expect(component.getAvailableMedics(3, new Date('2024-01-15'))).toBe(0);
  });

  it('should get available technicals', () => {
    component.shiftAvailability = {
      days: {
        '2024-01-15T00:00:00Z': {
          firstShift: { medicSlotsAvailable: 2, technicalSlotsAvailable: 1 },
          secondShift: { medicSlotsAvailable: 1, technicalSlotsAvailable: 2 },
          thirdShift: { medicSlotsAvailable: 0, technicalSlotsAvailable: 1 }
        }
      }
    };
    expect(component.getAvailableTechnicals(1, new Date('2024-01-15'))).toBe(1);
    expect(component.getAvailableTechnicals(2, new Date('2024-01-15'))).toBe(2);
    expect(component.getAvailableTechnicals(3, new Date('2024-01-15'))).toBe(1);
  });

  it('should get shift label', () => {
    expect(component.getShiftLabel(1)).toBe('06:00 - 14:00');
    expect(component.getShiftLabel(2)).toBe('14:00 - 22:00');
    expect(component.getShiftLabel(3)).toBe('22:00 - 06:00');
    expect(component.getShiftLabel(4)).toBe('');
  });

  it('should get translated date', () => {
    const date = new Date('2024-01-15');
    expect(component.getTranslatedDate(date)).toBe('Monday, January 15, 2024');
  });

  it('should load shift warnings', () => {
    // Reset the spy and set new return value
    (component['shiftService'].getShiftWarnings as jasmine.Spy).and.returnValue(of({ warnings: ['Test warning'] }));
    component.userId = '1';
    component.loadShiftWarnings();
    expect(component['shiftService'].getShiftWarnings).toHaveBeenCalledWith('1');
    expect(component.shiftWarnings).toEqual(['Test warning']);
  });

  it('should handle empty shift warnings', () => {
    // Reset the spy and set new return value
    (component['shiftService'].getShiftWarnings as jasmine.Spy).and.returnValue(of({ warnings: [] }));
    component.userId = '1';
    component.loadShiftWarnings();
    expect(component.shiftWarnings).toEqual([]);
  });

  it('should handle shift warnings error silently', () => {
    // Reset the spy and set new return value
    (component['shiftService'].getShiftWarnings as jasmine.Spy).and.returnValue(throwError(() => new Error('Test error')));
    component.userId = '1';
    component.loadShiftWarnings();
    // Test passes if no error is thrown and warnings remain empty
    expect(component.shiftWarnings).toEqual([]);
  });

  it('should not load shift warnings when userId is empty', () => {
    // Reset call count
    (component['shiftService'].getShiftWarnings as jasmine.Spy).calls.reset();
    component.userId = '';
    component.loadShiftWarnings();
    expect(component['shiftService'].getShiftWarnings).not.toHaveBeenCalled();
  });

  it('should translate warning messages with new format', () => {
    spyOn(component['translate'], 'instant').and.returnValue('You have only 4 shifts scheduled in the next 14 days. Consider scheduling more shifts to meet the 5 days/week quota.');

    const warning = 'SHIFT_WARNINGS.INSUFFICIENT_SHIFTS|4|14|5';
    const result = component.getTranslatedWarning(warning);

    expect(component['translate'].instant).toHaveBeenCalledWith('SHIFT_WARNINGS.INSUFFICIENT_SHIFTS', {
      shiftsCount: '4',
      daysCount: '14',
      requiredDays: '5'
    });
    expect(result).toBe('You have only 4 shifts scheduled in the next 14 days. Consider scheduling more shifts to meet the 5 days/week quota.');
  });

  it('should return original warning for backward compatibility', () => {
    const warning = 'You have only 4 shifts scheduled in the next 2 weeks. Consider scheduling more shifts to meet the 5 days/week quota.';
    const result = component.getTranslatedWarning(warning);

    expect(result).toBe(warning);
  });

  it('should translate error messages with new format', () => {
    spyOn(component['translate'], 'instant').and.returnValue('Assigning this shift would result in 7 consecutive shifts, which exceeds the maximum limit of 6 consecutive shifts.');

    const errorMessage = 'Server error: 400 - Bad Request - SHIFT_ERRORS.CONSECUTIVE_SHIFTS_LIMIT|7';
    const result = component.getTranslatedErrorMessage(errorMessage);

    expect(component['translate'].instant).toHaveBeenCalledWith('SHIFT_ERRORS.CONSECUTIVE_SHIFTS_LIMIT', {
      consecutiveCount: '7'
    });
    expect(result).toBe('Assigning this shift would result in 7 consecutive shifts, which exceeds the maximum limit of 6 consecutive shifts.');
  });

  it('should return original error message for backward compatibility', () => {
    const errorMessage = 'Server error: 400 - Bad Request - Some other error';
    const result = component.getTranslatedErrorMessage(errorMessage);

    expect(result).toBe(errorMessage);
  });

  it('should handle simple error messages without server error format', () => {
    const errorMessage = 'SHIFT_ERRORS.CONSECUTIVE_SHIFTS_LIMIT|8';
    spyOn(component['translate'], 'instant').and.returnValue('Assigning this shift would result in 8 consecutive shifts, which exceeds the maximum limit of 6 consecutive shifts.');

    const result = component.getTranslatedErrorMessage(errorMessage);

    expect(component['translate'].instant).toHaveBeenCalledWith('SHIFT_ERRORS.CONSECUTIVE_SHIFTS_LIMIT', {
      consecutiveCount: '8'
    });
    expect(result).toBe('Assigning this shift would result in 8 consecutive shifts, which exceeds the maximum limit of 6 consecutive shifts.');
  });

  it('should correctly identify low capacity shifts with custom data', () => {
    // Mock shift availability data
    component.shiftAvailability = {
      days: {
        '2024-01-15T00:00:00Z': {
          firstShift: {
            medicSlotsAvailable: 1,
            technicalSlotsAvailable: 2
          },
          secondShift: {
            medicSlotsAvailable: 2,
            technicalSlotsAvailable: 1
          },
          thirdShift: {
            medicSlotsAvailable: 2,
            technicalSlotsAvailable: 4
          }
        }
      }
    };

    const testDate = new Date('2024-01-15');

    // First shift should be low capacity (medics <= 1)
    expect(component.isShiftLowCapacity(1, testDate)).toBe(true);

    // Second shift should be low capacity (technicals <= 1)
    expect(component.isShiftLowCapacity(2, testDate)).toBe(true);

    // Third shift should NOT be low capacity (medics > 1 AND technicals > 1)
    expect(component.isShiftLowCapacity(3, testDate)).toBe(false);
  });

  it('should handle missing day data in getAvailableMedics', () => {
    component.shiftAvailability = { days: {} };
    spyOn(console, 'warn');
    const result = component.getAvailableMedics(1, new Date('2024-01-15'));
    expect(result).toBe(0);
    expect(console.warn).toHaveBeenCalled();
  });

  it('should handle missing day data in getAvailableTechnicals', () => {
    component.shiftAvailability = { days: {} };
    spyOn(console, 'warn');
    const result = component.getAvailableTechnicals(1, new Date('2024-01-15'));
    expect(result).toBe(0);
    expect(console.warn).toHaveBeenCalled();
  });

  it('should handle invalid shift type in getAvailableMedics', () => {
    component.shiftAvailability = {
      days: {
        '2024-01-15T00:00:00Z': {
          firstShift: { medicSlotsAvailable: 2, technicalSlotsAvailable: 1 },
          secondShift: { medicSlotsAvailable: 1, technicalSlotsAvailable: 2 },
          thirdShift: { medicSlotsAvailable: 0, technicalSlotsAvailable: 1 }
        }
      }
    };
    expect(component.getAvailableMedics(4, new Date('2024-01-15'))).toBe(0);
  });

  it('should handle invalid shift type in getAvailableTechnicals', () => {
    component.shiftAvailability = {
      days: {
        '2024-01-15T00:00:00Z': {
          firstShift: { medicSlotsAvailable: 2, technicalSlotsAvailable: 1 },
          secondShift: { medicSlotsAvailable: 1, technicalSlotsAvailable: 2 },
          thirdShift: { medicSlotsAvailable: 0, technicalSlotsAvailable: 1 }
        }
      }
    };
    expect(component.getAvailableTechnicals(4, new Date('2024-01-15'))).toBe(0);
  });

  it('should go back using router', () => {
    spyOn(component['router'], 'navigate');
    component.goBack();
    expect(component['router'].navigate).toHaveBeenCalledWith(['/']);
  });

  describe('New functionality tests', () => {
    beforeEach(() => {
      // Override the mock to provide specific test data
      // Use the date-only format that the component expects
      component.shiftAvailability = {
        days: {
          '2025-08-01': {
            firstShift: {
              medicSlotsAvailable: 1,
              technicalSlotsAvailable: 3,
              isAssignedToEmployee: true,
              isFullyBooked: false
            },
            secondShift: {
              medicSlotsAvailable: 0,
              technicalSlotsAvailable: 0,
              isAssignedToEmployee: false,
              isFullyBooked: true
            },
            thirdShift: {
              medicSlotsAvailable: 2,
              technicalSlotsAvailable: 4,
              isAssignedToEmployee: false,
              isFullyBooked: false
            }
          }
        }
      };
    });

    it('should correctly identify assigned shifts', () => {
      const testDate = new Date('2025-08-01T00:00:00.000Z');

      expect(component.isAssignedToShift(1, testDate)).toBe(true);
      expect(component.isAssignedToShift(2, testDate)).toBe(false);
      expect(component.isAssignedToShift(3, testDate)).toBe(false);
    });

    it('should correctly identify fully booked shifts', () => {
      const testDate = new Date('2025-08-01T00:00:00.000Z');

      expect(component.isShiftFullyBooked(1, testDate)).toBe(false);
      expect(component.isShiftFullyBooked(2, testDate)).toBe(true);
      expect(component.isShiftFullyBooked(3, testDate)).toBe(false);
    });

    it('should correctly identify low capacity shifts', () => {
      const testDate = new Date('2025-08-01T00:00:00.000Z');

      expect(component.isShiftLowCapacity(1, testDate)).toBe(true); // 1 medic, 3 technical (medics <= 1)
      expect(component.isShiftLowCapacity(2, testDate)).toBe(true); // 0 medic, 0 technical (both <= 1)
      expect(component.isShiftLowCapacity(3, testDate)).toBe(false); // 2 medic, 4 technical (both > 1)
    });

    it('should change time span and reload shifts', () => {
      spyOn(component, 'loadShiftAvailability');
      component.selectedTimeSpan = 7;

      component.changeTimeSpan(14);

      expect(component.selectedTimeSpan).toBe(14);
      expect(component.loadShiftAvailability).toHaveBeenCalled();
    });

    it('should not reload shifts if time span is the same', () => {
      spyOn(component, 'loadShiftAvailability');
      component.selectedTimeSpan = 7;

      component.changeTimeSpan(7);

      expect(component.selectedTimeSpan).toBe(7);
      expect(component.loadShiftAvailability).not.toHaveBeenCalled();
    });

  describe('Tooltip and assignment time helpers', () => {
    beforeEach(() => {
      // set locale/lang if needed implicitly via translate; we only check translation keys
      (component as any).translate.instant = (key: string) => key;
    });

    it('getAssignButtonTooltip covers processing/select/assigned/fully booked/default', () => {
      const day = new Date('2025-08-01');
      // Default available data
      component.shiftAvailability = {
        days: {
          '2025-08-01': {
            firstShift: { medicSlotsAvailable: 1, technicalSlotsAvailable: 1, isAssignedToEmployee: false, isFullyBooked: false },
            secondShift: { medicSlotsAvailable: 0, technicalSlotsAvailable: 0, isAssignedToEmployee: false, isFullyBooked: true },
            thirdShift: { medicSlotsAvailable: 1, technicalSlotsAvailable: 1, isAssignedToEmployee: true, isFullyBooked: false }
          }
        }
      } as any;

      // processing state
      component.isAssigning = true;
      expect(component.getAssignButtonTooltip(1, day)).toBe('SHIFT_MANAGEMENT.TOOLTIP_PROCESSING');
      component.isAssigning = false;

      // admin needs selected employee
      (component['auth'].getRole as jasmine.Spy).and.returnValue('Administrator');
      component.userRole = 'Administrator';
      expect(component.getAssignButtonTooltip(1, day)).toBe('SHIFT_MANAGEMENT.TOOLTIP_SELECT_EMPLOYEE');

      // switch to normal user for following checks
      (component['auth'].getRole as jasmine.Spy).and.returnValue('Medic');
      component.userRole = 'Medic';

      // already assigned (current user)
      expect(component.getAssignButtonTooltip(3, day)).toBe('SHIFT_MANAGEMENT.TOOLTIP_ALREADY_ASSIGNED');

      // fully booked
      expect(component.getAssignButtonTooltip(2, day)).toBe('SHIFT_MANAGEMENT.TOOLTIP_FULLY_BOOKED');

      // default assign
      expect(component.getAssignButtonTooltip(1, day)).toBe('SHIFT_MANAGEMENT.TOOLTIP_ASSIGN');
    });

    it('getRemoveButtonTooltip covers processing/select/not assigned/default', () => {
      const day = new Date('2025-08-01');
      component.shiftAvailability = {
        days: {
          '2025-08-01': {
            firstShift: { medicSlotsAvailable: 1, technicalSlotsAvailable: 1, isAssignedToEmployee: true, isFullyBooked: false },
            secondShift: { medicSlotsAvailable: 1, technicalSlotsAvailable: 1, isAssignedToEmployee: false, isFullyBooked: false }
          }
        }
      } as any;

      // processing
      component.isRemoving = true;
      expect(component.getRemoveButtonTooltip(1, day)).toBe('SHIFT_MANAGEMENT.TOOLTIP_PROCESSING');
      component.isRemoving = false;

      // admin needs employee
      (component['auth'].getRole as jasmine.Spy).and.returnValue('Administrator');
      component.userRole = 'Administrator';
      expect(component.getRemoveButtonTooltip(1, day)).toBe('SHIFT_MANAGEMENT.TOOLTIP_SELECT_EMPLOYEE');

      // switch to normal user for following checks
      (component['auth'].getRole as jasmine.Spy).and.returnValue('Medic');
      component.userRole = 'Medic';

      // not assigned
      expect(component.getRemoveButtonTooltip(2, day)).toBe('SHIFT_MANAGEMENT.TOOLTIP_NOT_ASSIGNED');

      // default remove
      expect(component.getRemoveButtonTooltip(1, day)).toBe('SHIFT_MANAGEMENT.TOOLTIP_REMOVE');
    });

    it('getAssignmentTime formats date and falls back safely', () => {
      // prepare assignment cache for current user
      component.userId = 'u1';
      component.selectedEmployeeShifts.set('u1', [
        { id: 1, shiftDate: '2024-01-15', shiftType: 1, assignedAt: '2024-01-15T10:30:00Z' } as any,
        { id: 2, shiftDate: '2024-01-16', shiftType: 2, createdAt: '2024-01-16T11:45:00Z' } as any,
        { id: 3, shiftDate: '2024-01-17', shiftType: 3 } as any
      ]);

      // Valid assignedAt
      const t1 = component.getAssignmentTime(1, new Date('2024-01-15'));
      expect(typeof t1).toBe('string');
      expect(t1).not.toBeNull();

      // Fallback to createdAt
      const t2 = component.getAssignmentTime(2, new Date('2024-01-16'));
      expect(typeof t2).toBe('string');

      // No timestamps
      const t3 = component.getAssignmentTime(3, new Date('2024-01-17'));
      expect(t3).toBeNull();
    });
  });

  });

});
