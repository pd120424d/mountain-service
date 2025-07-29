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

  it('should fetch role, userid and load shifts on init when user is non administrator', () => {
    // Services are already mocked in beforeEach, just verify the calls and state
    expect(component['shiftService'].getShiftAvailability).toHaveBeenCalled();
    expect(component.userRole).toBe('Technical');
    expect(component.userId).toBe('1');
  });

  it('should fetch role, userid and load shifts and employees on init when user is administrator', () => {
    // Reset the existing spies and set up for administrator test
    (component['shiftService'].getShiftAvailability as jasmine.Spy).calls.reset();
    (component['shiftService'].getShiftWarnings as jasmine.Spy).calls.reset();
    (component['auth'].getRole as jasmine.Spy).and.returnValue('Administrator');
    spyOn(component['shiftService'], 'getAllEmployees').and.returnValue(of([]));

    // Call ngOnInit to test administrator flow
    component.ngOnInit();

    expect(component['shiftService'].getShiftAvailability).toHaveBeenCalled();
    expect(component['shiftService'].getAllEmployees).toHaveBeenCalled();
    expect(component.userRole).toBe('Administrator');
    expect(component.userId).toBe('1');
  });

  it('should load shifts', () => {
    // Reset the spy and set new return value
    (component['shiftService'].getShiftAvailability as jasmine.Spy).and.returnValue(of({ days: {} }));
    component.loadShifts();
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
        '2024-01-15': {
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
        '2024-01-15': {
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
        '2024-01-15': {
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
        '2024-01-15': {
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

      expect(component.isShiftLowCapacity(1, testDate)).toBe(true); // 1 medic, 3 technical
      expect(component.isShiftLowCapacity(2, testDate)).toBe(true); // 0 medic, 0 technical
      expect(component.isShiftLowCapacity(3, testDate)).toBe(false); // 2 medic, 4 technical
    });

    it('should change time span and reload shifts', () => {
      spyOn(component, 'loadShifts');
      component.selectedTimeSpan = 7;

      component.changeTimeSpan(14);

      expect(component.selectedTimeSpan).toBe(14);
      expect(component.loadShifts).toHaveBeenCalled();
    });

    it('should not reload shifts if time span is the same', () => {
      spyOn(component, 'loadShifts');
      component.selectedTimeSpan = 7;

      component.changeTimeSpan(7);

      expect(component.selectedTimeSpan).toBe(7);
      expect(component.loadShifts).not.toHaveBeenCalled();
    });
  });

});
