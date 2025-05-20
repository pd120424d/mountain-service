import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ShiftManagementComponent } from './shift.component';
import { sharedTestingProviders } from '../test-utils/shared-test-imports';

describe('ShiftManagementComponent', () => {
  let component: ShiftManagementComponent;
  let fixture: ComponentFixture<ShiftManagementComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ShiftManagementComponent],
      providers: [...sharedTestingProviders]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ShiftManagementComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });

  // it('should load shifts on init', () => {
  //   spyOn(component.shiftService, 'getShiftAvailability').and.returnValue(of({}));
  //   component.ngOnInit();
  //   expect(component.shiftService.getShiftAvailability).toHaveBeenCalled();
  // });

  // it('should load all employees on init if user is admin', () => {
  //   spyOn(component.shiftService, 'getAllEmployees').and.returnValue(of([]));
  //   component.userRole = 'Administrator';
  //   component.ngOnInit();
  //   expect(component.shiftService.getAllEmployees).toHaveBeenCalled();
  // });

  // it('should not load all employees on init if user is not admin', () => {
  //   spyOn(component.shiftService, 'getAllEmployees').and.returnValue(of([]));
  //   component.userRole = 'Medic';
  //   component.ngOnInit();
  //   expect(component.shiftService.getAllEmployees).not.toHaveBeenCalled();
  // });

  // it('should assign to shift', () => {
  //   spyOn(component.shiftService, 'assignEmployeeToShift').and.returnValue(of({}));
  //   spyOn(component.toastr, 'success');
  //   spyOn(component.spinner, 'show');
  //   spyOn(component.spinner, 'hide');
  //   component.assignToShift(1);
  //   expect(component.shiftService.assignEmployeeToShift).toHaveBeenCalledWith(1, component.userId);
  //   expect(component.toastr.success).toHaveBeenCalledWith('Successfully assigned!');
  //   expect(component.spinner.show).toHaveBeenCalled();
  //   expect(component.spinner.hide).toHaveBeenCalled();
  // });

  // it('should remove from shift', () => {
  //   spyOn(component.shiftService, 'removeEmployeeFromShift').and.returnValue(of({}));
  //   spyOn(component.toastr, 'success');
  //   spyOn(component.spinner, 'show');
  //   spyOn(component.spinner, 'hide');
  //   component.removeFromShift(1);
  //   expect(component.shiftService.removeEmployeeFromShift).toHaveBeenCalledWith(1, component.userId);
  //   expect(component.toastr.success).toHaveBeenCalledWith('Successfully removed!');
  //   expect(component.spinner.show).toHaveBeenCalled();
  //   expect(component.spinner.hide).toHaveBeenCalled();
  // });

});
