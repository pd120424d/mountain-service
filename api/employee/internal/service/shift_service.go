package service

import (
	"fmt"
	"time"

	commonv1 "github.com/pd120424d/mountain-service/api/contracts/common/v1"
	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/employee/internal/repositories"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type shiftService struct {
	log        utils.Logger
	emplRepo   repositories.EmployeeRepository
	shiftsRepo repositories.ShiftRepository
}

func NewShiftService(log utils.Logger, emplRepo repositories.EmployeeRepository, shiftsRepo repositories.ShiftRepository) ShiftService {
	return &shiftService{
		log:        log.WithName("shiftService"),
		emplRepo:   emplRepo,
		shiftsRepo: shiftsRepo,
	}
}

func (s *shiftService) AssignShift(employeeID uint, req employeeV1.AssignShiftRequest) (*employeeV1.AssignShiftResponse, error) {
	s.log.Infof("Assigning shift for employee ID %d", employeeID)

	// Step 1: Validate employee exists
	employee := &model.Employee{}
	err := s.emplRepo.GetEmployeeByID(employeeID, employee)
	if err != nil {
		s.log.Errorf("failed to get employee: %v", err)
		return nil, commonv1.NewAppError("EMPLOYEE_ERRORS.NOT_FOUND", "employee not found", nil)
	}

	// Step 2: Parse and validate shift date
	shiftDate, err := time.Parse("2006-01-02", req.ShiftDate)
	if err != nil {
		s.log.Errorf("failed to parse shift date: %v", err)
		return nil, commonv1.NewAppError("VALIDATION.INVALID_SHIFT_DATE", "invalid shift date format", nil)
	}

	// Step 3: Validate shift date is in the future
	if shiftDate.Before(time.Now().Truncate(24 * time.Hour)) {
		s.log.Errorf("shift date %s is in the past", req.ShiftDate)
		return nil, commonv1.NewAppError("VALIDATION.SHIFT_IN_PAST", "shift date must be in the future", nil)
	}

	// Step 4: Validate shift date is within 3 months
	threeMonthsFromNow := time.Now().AddDate(0, 3, 0)
	if shiftDate.After(threeMonthsFromNow) {
		s.log.Errorf("shift date %s is more than 3 months in the future", req.ShiftDate)
		return nil, commonv1.NewAppError("VALIDATION.SHIFT_TOO_FAR", "shift date cannot be more than 3 months in the future", nil)
	}

	// Step 5: Check consecutive shifts rule (max 2 consecutive shifts, then 1 day rest)
	if err := s.validateConsecutiveShifts(employeeID, shiftDate, req.ShiftType); err != nil {
		s.log.Errorf("consecutive shifts validation failed: %v", err)
		return nil, err
	}

	// Step 6: Get or create shift
	shift, err := s.shiftsRepo.GetOrCreateShift(shiftDate, req.ShiftType)
	if err != nil {
		s.log.Errorf("failed to get or create shift: %v", err)
		return nil, fmt.Errorf("failed to create shift")
	}

	// Step 7: Check if employee is already assigned
	assigned, err := s.shiftsRepo.AssignedToShift(employeeID, shift.ID)
	if err != nil {
		s.log.Errorf("failed to check assignment: %v", err)
		return nil, fmt.Errorf("failed to check assignment")
	}
	if assigned {
		s.log.Errorf("employee with ID %d is already assigned to shift ID %d", employeeID, shift.ID)
		return nil, commonv1.NewAppError("SHIFT_ERRORS.ALREADY_ASSIGNED", "employee is already assigned to this shift", nil)
	}

	// Step 8: Check shift capacity
	currentCount, err := s.shiftsRepo.CountAssignmentsByProfile(shift.ID, employee.ProfileType)
	if err != nil {
		s.log.Errorf("failed to count assignments: %v", err)
		return nil, fmt.Errorf("failed to check shift capacity")
	}

	maxCapacity := s.getMaxCapacityForProfile(employee.ProfileType)
	if currentCount >= int64(maxCapacity) {
		s.log.Errorf("shift capacity full for profile type %s", employee.ProfileType.String())
		return nil, commonv1.NewAppError("SHIFT_ERRORS.CAPACITY_FULL", fmt.Sprintf("shift capacity is full for %s staff", employee.ProfileType.String()), map[string]interface{}{"role": employee.ProfileType.String(), "max": maxCapacity})
	}

	// Step 9: Create assignment
	assignmentID, err := s.shiftsRepo.CreateAssignment(employeeID, shift.ID)
	if err != nil {
		s.log.Errorf("failed to create assignment: %v", err)
		return nil, fmt.Errorf("failed to assign shift")
	}

	s.log.Infof("Successfully assigned shift for employee ID %d", employeeID)

	return &employeeV1.AssignShiftResponse{
		ID:        assignmentID,
		ShiftDate: shift.ShiftDate.Format("2006-01-02"),
		ShiftType: shift.ShiftType,
	}, nil
}

func (s *shiftService) GetShifts(employeeID uint) ([]employeeV1.ShiftResponse, error) {
	s.log.Infof("Getting shifts for employee ID %d", employeeID)

	var shifts []model.Shift
	err := s.shiftsRepo.GetShiftsByEmployeeID(employeeID, &shifts)
	if err != nil {
		s.log.Errorf("failed to get shifts for employee ID %d: %v", employeeID, err)
		return nil, fmt.Errorf("failed to retrieve shifts")
	}

	s.log.Infof("Successfully retrieved %d shifts for employee ID %d", len(shifts), employeeID)

	response := make([]employeeV1.ShiftResponse, 0, len(shifts))
	for _, shift := range shifts {
		response = append(response, employeeV1.ShiftResponse{
			ID:        shift.ID,
			ShiftDate: shift.ShiftDate,
			ShiftType: shift.ShiftType,
			CreatedAt: shift.CreatedAt,
		})
	}

	return response, nil
}

func (s *shiftService) GetShiftsAvailability(employeeID uint, days int) (*employeeV1.ShiftAvailabilityResponse, error) {
	s.log.Infof("Getting shift availability for employee %d for %d days", employeeID, days)

	if days <= 0 || days > 90 {
		s.log.Errorf("invalid days parameter: %d", days)
		return nil, fmt.Errorf("days must be between 1 and 90")
	}

	start := time.Now().Truncate(24 * time.Hour)
	end := start.AddDate(0, 0, days)

	availability, err := s.shiftsRepo.GetShiftAvailabilityWithEmployeeStatus(employeeID, start, end)
	if err != nil {
		s.log.Errorf("failed to get shift availability with employee status: %v", err)
		return nil, fmt.Errorf("failed to retrieve shift availability")
	}

	response := &employeeV1.ShiftAvailabilityResponse{
		Days: make(map[time.Time]employeeV1.ShiftAvailabilityPerDay),
	}

	for date, shifts := range availability.Days {
		dayAvailability := employeeV1.ShiftAvailabilityPerDay{}

		// Initialize all shifts with zero availability
		dayAvailability.FirstShift = employeeV1.ShiftAvailability{
			MedicSlotsAvailable:     0,
			TechnicalSlotsAvailable: 0,
			IsAssignedToEmployee:    false,
			IsFullyBooked:           false,
		}
		dayAvailability.SecondShift = employeeV1.ShiftAvailability{
			MedicSlotsAvailable:     0,
			TechnicalSlotsAvailable: 0,
			IsAssignedToEmployee:    false,
			IsFullyBooked:           false,
		}
		dayAvailability.ThirdShift = employeeV1.ShiftAvailability{
			MedicSlotsAvailable:     0,
			TechnicalSlotsAvailable: 0,
			IsAssignedToEmployee:    false,
			IsFullyBooked:           false,
		}

		if len(shifts) >= 3 {
			dayAvailability.FirstShift = employeeV1.ShiftAvailability{
				MedicSlotsAvailable:     max(0, shifts[0].MedicSlotsAvailable),
				TechnicalSlotsAvailable: max(0, shifts[0].TechnicalSlotsAvailable),
				IsAssignedToEmployee:    shifts[0].IsAssignedToEmployee,
				IsFullyBooked:           shifts[0].IsFullyBooked,
			}
			dayAvailability.SecondShift = employeeV1.ShiftAvailability{
				MedicSlotsAvailable:     max(0, shifts[1].MedicSlotsAvailable),
				TechnicalSlotsAvailable: max(0, shifts[1].TechnicalSlotsAvailable),
				IsAssignedToEmployee:    shifts[1].IsAssignedToEmployee,
				IsFullyBooked:           shifts[1].IsFullyBooked,
			}
			dayAvailability.ThirdShift = employeeV1.ShiftAvailability{
				MedicSlotsAvailable:     max(0, shifts[2].MedicSlotsAvailable),
				TechnicalSlotsAvailable: max(0, shifts[2].TechnicalSlotsAvailable),
				IsAssignedToEmployee:    shifts[2].IsAssignedToEmployee,
				IsFullyBooked:           shifts[2].IsFullyBooked,
			}
		}

		response.Days[date] = dayAvailability
	}

	s.log.Infof("Successfully retrieved shift availability for employee %d for %d days", employeeID, days)
	return response, nil
}

func (s *shiftService) RemoveShift(employeeID uint, req employeeV1.RemoveShiftRequest) error {
	s.log.Infof("Removing shift for employee ID %d", employeeID)

	shiftDate, err := time.Parse("2006-01-02", req.ShiftDate)
	if err != nil {
		s.log.Errorf("failed to parse shift date: %v", err)
		return commonv1.NewAppError("VALIDATION.INVALID_SHIFT_DATE", "invalid shift date format", nil)
	}

	err = s.shiftsRepo.RemoveEmployeeFromShiftByDetails(employeeID, shiftDate, req.ShiftType)
	if err != nil {
		s.log.Errorf("failed to remove shift: %v", err)
		return fmt.Errorf("failed to remove shift")
	}

	s.log.Infof("Successfully removed shift for employee ID %d", employeeID)
	return nil
}

func (s *shiftService) GetOnCallEmployees(currentTime time.Time, shiftBuffer time.Duration) ([]employeeV1.EmployeeResponse, error) {
	s.log.Infof("Getting on-call employees")

	employees, err := s.shiftsRepo.GetOnCallEmployees(currentTime, shiftBuffer)
	if err != nil {
		s.log.Errorf("Failed to get on-call employees: %v", err)
		return nil, fmt.Errorf("failed to retrieve on-call employees")
	}

	var employeeResponses []employeeV1.EmployeeResponse
	for _, emp := range employees {
		employeeResponses = append(employeeResponses, emp.UpdateResponseFromEmployee())
	}

	s.log.Infof("Successfully retrieved %d on-call employees", len(employeeResponses))
	return employeeResponses, nil
}

func (s *shiftService) GetShiftWarnings(employeeID uint) ([]string, error) {
	s.log.Infof("Getting shift warnings for employee ID %d", employeeID)

	// Check if employee exists
	employee := &model.Employee{}
	err := s.emplRepo.GetEmployeeByID(employeeID, employee)
	if err != nil {
		s.log.Errorf("failed to get employee: %v", err)
		return nil, commonv1.NewAppError("EMPLOYEE_ERRORS.NOT_FOUND", "employee not found", nil)
	}

	var warnings []string

	// Get next two weeks date range (start of today to start of day +14)
	start := time.Now().Truncate(24 * time.Hour)
	end := start.AddDate(0, 0, 14)

	// Check coverage for the employee's role in the next two weeks
	availability, err := s.shiftsRepo.GetShiftAvailability(start, end)
	if err != nil {
		s.log.Errorf("failed to get shift availability: %v", err)
		return nil, fmt.Errorf("failed to check shift coverage")
	}

	// Only show warnings if there's at least one shift with zero staff for the employee's role
	zeroRoleShiftExists := false
	for _, shifts := range availability.Days {
		for _, shift := range shifts {
			switch employee.ProfileType {
			case model.Medic:
				if shift[model.Medic] == 2 { // 2 available => 0 assigned
					zeroRoleShiftExists = true
				}
			case model.Technical:
				if shift[model.Technical] == 4 { // 4 available => 0 assigned
					zeroRoleShiftExists = true
				}
			}
			if zeroRoleShiftExists {
				break
			}
		}
		if zeroRoleShiftExists {
			break
		}
	}

	if !zeroRoleShiftExists {
		// Adequate baseline coverage for this role; no warnings
		s.log.Infof("No zero-coverage shifts for role %s in next two weeks; no warnings", employee.ProfileType.String())
		return warnings, nil
	}

	// Compute distinct days scheduled by the employee in the next two weeks (total over 14 days)
	var employeeShifts []model.Shift
	err = s.shiftsRepo.GetShiftsByEmployeeIDInDateRange(employeeID, start, end, &employeeShifts)
	if err != nil {
		s.log.Errorf("failed to get employee shifts: %v", err)
		return nil, fmt.Errorf("failed to check employee shifts")
	}

	distinctDays := map[string]struct{}{}
	for _, sft := range employeeShifts {
		day := sft.ShiftDate.Truncate(24 * time.Hour)
		key := day.Format("2006-01-02")
		distinctDays[key] = struct{}{}
	}

	totalDistinctDays := len(distinctDays)

	// Warn only if fewer than 10 distinct days are scheduled in the next 14 days
	if totalDistinctDays < 10 {
		warnings = append(warnings, fmt.Sprintf("%s|%d|14|5", model.WarningInsufficientShifts, totalDistinctDays))
	}

	s.log.Infof("Successfully retrieved %d warnings for employee ID %d", len(warnings), employeeID)
	return warnings, nil
}

func (s *shiftService) GetAdminShiftsAvailability(days int) (*employeeV1.ShiftAvailabilityResponse, error) {
	s.log.Infof("Getting admin shifts availability for %d days", days)

	// TODO: we need to return shift availability for all employees here, if it is possible somehow
	return &employeeV1.ShiftAvailabilityResponse{
		Days: make(map[time.Time]employeeV1.ShiftAvailabilityPerDay),
	}, nil
}

// Helper methods

func (s *shiftService) validateConsecutiveShifts(employeeID uint, shiftDate time.Time, shiftType int) error {
	// Build a small window around the candidate date to check adjacency and calendar-day rest rules
	startDate := shiftDate.AddDate(0, 0, -3)
	endDate := shiftDate.AddDate(0, 0, 3)

	var shifts []model.Shift
	if err := s.shiftsRepo.GetShiftsByEmployeeIDInDateRange(employeeID, startDate, endDate, &shifts); err != nil {
		return fmt.Errorf("failed to get employee shifts: %w", err)
	}

	// Map assignments per slot and per day
	dayKey := func(d time.Time) string { return d.Truncate(24 * time.Hour).Format("2006-01-02") }
	slotKey := func(d time.Time, t int) string { return fmt.Sprintf("%s|%d", dayKey(d), t) }
	assigned := make(map[string]bool)
	assignedDayAny := make(map[string]bool)
	for _, sh := range shifts {
		assigned[slotKey(sh.ShiftDate, sh.ShiftType)] = true
		assignedDayAny[dayKey(sh.ShiftDate)] = true
	}

	// Helpers for adjacent slots
	prevSlot := func(d time.Time, t int) (time.Time, int) {
		switch t {
		case 1:
			return d.AddDate(0, 0, -1), 3
		case 2:
			return d, 1
		default: // 3
			return d, 2
		}
	}
	nextSlot := func(d time.Time, t int) (time.Time, int) {
		switch t {
		case 1:
			return d, 2
		case 2:
			return d, 3
		default: // 3
			return d.AddDate(0, 0, 1), 1
		}
	}
	endDayOfShift := func(d time.Time, t int) time.Time {
		// 1 and 2 end same day, 3 ends next day
		switch t {
		case 3:
			return d.AddDate(0, 0, 1).Truncate(24 * time.Hour)
		default:
			return d.Truncate(24 * time.Hour)
		}
	}
	dayHasAny := func(d time.Time) bool { return assignedDayAny[dayKey(d)] }

	// 1) Block if any existing double's rest day equals the candidate day
	candDay := dayKey(shiftDate)
	for _, sh := range shifts {
		pd, pt := prevSlot(sh.ShiftDate, sh.ShiftType)
		if assigned[slotKey(pd, pt)] {
			rest := endDayOfShift(sh.ShiftDate, sh.ShiftType).AddDate(0, 0, 1)
			if dayKey(rest) == candDay {
				return commonv1.NewAppError(
					model.ErrorConsecutiveShiftsLimit,
					fmt.Sprintf("%s|%d", model.ErrorConsecutiveShiftsLimit, 3),
					map[string]interface{}{"limit": 3},
				)
			}
		}
	}

	// 2) No triples (max 2 consecutive including candidate)
	leftRun := 0
	if pd, pt := prevSlot(shiftDate, shiftType); assigned[slotKey(pd, pt)] {
		leftRun = 1
		if p2d, p2t := prevSlot(pd, pt); assigned[slotKey(p2d, p2t)] {
			leftRun = 2
		}
	}
	rightRun := 0
	if nd, nt := nextSlot(shiftDate, shiftType); assigned[slotKey(nd, nt)] {
		rightRun = 1
		if n2d, n2t := nextSlot(nd, nt); assigned[slotKey(n2d, n2t)] {
			rightRun = 2
		}
	}
	if leftRun+1+rightRun > 2 {
		return commonv1.NewAppError(
			model.ErrorConsecutiveShiftsLimit,
			fmt.Sprintf("%s|%d", model.ErrorConsecutiveShiftsLimit, leftRun+1+rightRun),
			map[string]interface{}{"limit": leftRun + 1 + rightRun},
		)
	}

	// 3) If candidate forms a new double, enforce calendar-day rest for that double
	if leftRun == 1 {
		rest := endDayOfShift(shiftDate, shiftType).AddDate(0, 0, 1)
		if dayHasAny(rest) {
			return commonv1.NewAppError(
				model.ErrorConsecutiveShiftsLimit,
				fmt.Sprintf("%s|%d", model.ErrorConsecutiveShiftsLimit, 3),
				map[string]interface{}{"limit": 3},
			)
		}
	}
	if rightRun == 1 {
		nd, nt := nextSlot(shiftDate, shiftType)
		rest := endDayOfShift(nd, nt).AddDate(0, 0, 1)
		if dayHasAny(rest) {
			return commonv1.NewAppError(
				model.ErrorConsecutiveShiftsLimit,
				fmt.Sprintf("%s|%d", model.ErrorConsecutiveShiftsLimit, 3),
				map[string]interface{}{"limit": 3},
			)
		}
	}

	return nil
}

func (s *shiftService) getMaxCapacityForProfile(profileType model.ProfileType) int {
	switch profileType {
	case model.Medic:
		return 2
	case model.Technical:
		return 4
	default:
		return 0
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
