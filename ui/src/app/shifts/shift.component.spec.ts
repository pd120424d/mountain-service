import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ShiftManagementComponent } from './shift.component';
import { sharedTestingProviders } from '../test-utils/shared-test-imports';
import { DatePipe } from '@angular/common';
import { of } from 'rxjs';

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
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  it('should fetch role, userid and load shifts on init when user is non administrator', () => {
    spyOn(component['shiftService'], 'getShiftAvailability').and.returnValue(of({ days: {} }));
    spyOn(component['auth'], 'getRole').and.returnValue('Technical');
    spyOn(component['auth'], 'getUserId').and.returnValue('1');
    component.ngOnInit();
    expect(component['shiftService'].getShiftAvailability).toHaveBeenCalled();
    expect(component.userRole).toBe('Technical');
    expect(component.userId).toBe('1');
  });

  it('should fetch role, userid and load shifts and employees on init when user is administrator', () => {
    spyOn(component['shiftService'], 'getShiftAvailability').and.returnValue(of({ days: {} }));
    spyOn(component['shiftService'], 'getAllEmployees').and.returnValue(of([]));
    spyOn(component['auth'], 'getRole').and.returnValue('Administrator');
    spyOn(component['auth'], 'getUserId').and.returnValue('1');
    component.ngOnInit();
    expect(component['shiftService'].getShiftAvailability).toHaveBeenCalled();
    expect(component['shiftService'].getAllEmployees).toHaveBeenCalled();
    expect(component.userRole).toBe('Administrator');
    expect(component.userId).toBe('1');
  });

  it('should load shifts', () => {
    spyOn(component['shiftService'], 'getShiftAvailability').and.returnValue(of({ days: {} }));
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
    spyOn(component['toastr'], 'success');
    spyOn(component['spinner'], 'show');
    spyOn(component['spinner'], 'hide');
    component.assignToShift(1, new Date(), '1');
    expect(component['shiftService'].assignEmployeeToShift).toHaveBeenCalled();
    expect(component['toastr'].success).toHaveBeenCalled();
    expect(component['spinner'].show).toHaveBeenCalled();
    expect(component['spinner'].hide).toHaveBeenCalled();
  });

  it('should remove from shift', () => {
    spyOn(component['shiftService'], 'removeEmployeeFromShiftByDetails').and.returnValue(of({}));
    spyOn(component['toastr'], 'success');
    spyOn(component['spinner'], 'show');
    spyOn(component['spinner'], 'hide');
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
          firstShift: { medic: 2, technical: 1 },
          secondShift: { medic: 1, technical: 2 },
          thirdShift: { medic: 0, technical: 1 }
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
          firstShift: { medic: 2, technical: 1 },
          secondShift: { medic: 1, technical: 2 },
          thirdShift: { medic: 0, technical: 1 }
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

});
