package service

import (
	"fmt"
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/employee/internal/repositories"
	"github.com/pd120424d/mountain-service/api/shared/utils"
)

type EmployeeService interface {
	AssignShift(employeeID uint, req employeeV1.AssignShiftRequest) (*employeeV1.AssignShiftResponse, error)
	GetShifts(employeeID uint) ([]employeeV1.ShiftResponse, error)
	GetShiftsAvailability(days int) (*employeeV1.ShiftAvailabilityResponse, error)
	RemoveShift(employeeID uint, req employeeV1.RemoveShiftRequest) error
	GetOnCallEmployees(currentTime time.Time, shiftBuffer time.Duration) ([]employeeV1.EmployeeResponse, error)
	GetShiftWarnings(employeeID uint) ([]string, error)
}

type employeeService struct {
	log        utils.Logger
	emplRepo   repositories.EmployeeRepository
	shiftsRepo repositories.ShiftRepository
}

func NewEmployeeService(log utils.Logger, emplRepo repositories.EmployeeRepository, shiftsRepo repositories.ShiftRepository) EmployeeService {
	return &employeeService{
		log:        log.WithName("employeeService"),
		emplRepo:   emplRepo,
		shiftsRepo: shiftsRepo,
	}
}

func (s *employeeService) AssignShift(employeeID uint, req employeeV1.AssignShiftRequest) (*employeeV1.AssignShiftResponse, error) {
	s.log.Infof("Assigning shift for employee ID %d", employeeID)

	// Step 1: Validate employee exists
	employee := &model.Employee{}
	err := s.emplRepo.GetEmployeeByID(employeeID, employee)
	if err != nil {
		s.log.Errorf("failed to get employee: %v", err)
		return nil, fmt.Errorf("employee not found")
	}

	// Step 2: Parse and validate shift date
	shiftDate, err := time.Parse("2006-01-02", req.ShiftDate)
	if err != nil {
		s.log.Errorf("failed to parse shift date: %v", err)
		return nil, fmt.Errorf("invalid shift date format")
	}

	// Step 3: Validate shift is not in the past
	if shiftDate.Before(time.Now().Truncate(24 * time.Hour)) {
		s.log.Errorf("cannot assign shift in the past")
		return nil, fmt.Errorf("cannot assign shift in the past")
	}

	// Step 4: Validate shift is not more than 3 months in advance
	threeMonthsFromNow := time.Now().AddDate(0, 3, 0)
	if shiftDate.After(threeMonthsFromNow) {
		s.log.Errorf("cannot assign shift more than 3 months in advance")
		return nil, fmt.Errorf("cannot assign shift more than 3 months in advance")
	}

	// Step 5: Validate consecutive shifts rule (max 2 consecutive shifts)
	err = s.validateConsecutiveShifts(employeeID, shiftDate)
	if err != nil {
		s.log.Errorf("consecutive shifts validation failed: %v", err)
		return nil, err
	}

	// Step 6: Get or create shift
	shift, err := s.shiftsRepo.GetOrCreateShift(shiftDate, req.ShiftType)
	if err != nil {
		s.log.Errorf("failed to get or create shift: %v", err)
		return nil, fmt.Errorf("failed to process shift")
	}

	// Step 7: Check if employee is already assigned
	assigned, err := s.shiftsRepo.AssignedToShift(employeeID, shift.ID)
	if err != nil {
		s.log.Errorf("failed to check assignment: %v", err)
		return nil, fmt.Errorf("failed to check assignment")
	}
	if assigned {
		s.log.Errorf("employee with ID %d is already assigned to shift ID %d", employeeID, shift.ID)
		return nil, fmt.Errorf("employee is already assigned to this shift")
	}

	// Step 8: Check capacity based on profile type
	profileType := employee.ProfileType
	count, err := s.shiftsRepo.CountAssignmentsByProfile(shift.ID, profileType)
	if err != nil {
		s.log.Errorf("failed to count assignments by profile: %v", err)
		return nil, fmt.Errorf("failed to check shift capacity")
	}

	maxCapacity := s.getMaxCapacityForProfile(profileType)
	if count >= int64(maxCapacity) {
		s.log.Errorf("maximum capacity for role %s reached in the selected shift", profileType.String())
		return nil, fmt.Errorf("maximum capacity for this role reached in the selected shift")
	}

	// Step 9: Create assignment
	assignmentID, err := s.shiftsRepo.CreateAssignment(employeeID, shift.ID)
	if err != nil {
		s.log.Errorf("failed to assign shift: %v", err)
		return nil, fmt.Errorf("failed to assign employee")
	}

	s.log.Infof("Successfully assigned shift ID %d to employee ID %d with assignment ID %d", shift.ID, employeeID, assignmentID)

	return &employeeV1.AssignShiftResponse{
		ID:        assignmentID,
		ShiftDate: req.ShiftDate,
		ShiftType: req.ShiftType,
	}, nil
}

func (s *employeeService) GetShifts(employeeID uint) ([]employeeV1.ShiftResponse, error) {
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

func (s *employeeService) GetShiftsAvailability(days int) (*employeeV1.ShiftAvailabilityResponse, error) {
	s.log.Infof("Getting shifts availability for the next %d days", days)

	if days <= 0 {
		return nil, fmt.Errorf("invalid days parameter")
	}

	start := time.Now().Truncate(24 * time.Hour)
	end := start.AddDate(0, 0, days)

	availability, err := s.shiftsRepo.GetShiftAvailability(start, end)
	if err != nil {
		s.log.Errorf("failed to get shifts availability: %v", err)
		return nil, fmt.Errorf("failed to retrieve shift availability")
	}

	s.log.Infof("Successfully retrieved shifts availability for %d days", days)

	response := &employeeV1.ShiftAvailabilityResponse{
		Days: make(map[time.Time]employeeV1.ShiftAvailabilityPerDay),
	}

	for date, shifts := range availability.Days {
		dayAvailability := employeeV1.ShiftAvailabilityPerDay{}

		// TODO: rename to FirstShift
		dayAvailability.Shift1 = employeeV1.ShiftAvailability{
			MedicSlotsAvailable:     0,
			TechnicalSlotsAvailable: 0,
		}
		// TODO: rename to SecondShift
		dayAvailability.Shift2 = employeeV1.ShiftAvailability{
			MedicSlotsAvailable:     0,
			TechnicalSlotsAvailable: 0,
		}
		// TODO: rename to ThirdShifts
		dayAvailability.Shift3 = employeeV1.ShiftAvailability{
			MedicSlotsAvailable:     0,
			TechnicalSlotsAvailable: 0,
		}

		for shiftIndex, shiftMap := range shifts {
			var shiftAvailability *employeeV1.ShiftAvailability
			switch shiftIndex {
			case 0:
				shiftAvailability = &dayAvailability.Shift1
			case 1:
				shiftAvailability = &dayAvailability.Shift2
			case 2:
				shiftAvailability = &dayAvailability.Shift3
			default:
				s.log.Warnf("more than 3 shifts found for date %v, skipping the entry", date)
				continue
			}

			for profileType, availableSlots := range shiftMap {
				switch profileType {
				case model.Medic:
					shiftAvailability.MedicSlotsAvailable = max(0, availableSlots)
				case model.Technical:
					shiftAvailability.TechnicalSlotsAvailable = max(0, availableSlots)
				}
			}
		}

		response.Days[date] = dayAvailability
	}

	return response, nil
}

func (s *employeeService) RemoveShift(employeeID uint, req employeeV1.RemoveShiftRequest) error {
	s.log.Infof("Removing shift for employee ID %d", employeeID)

	shiftDate, err := time.Parse("2006-01-02", req.ShiftDate)
	if err != nil {
		s.log.Errorf("failed to parse shift date: %v", err)
		return fmt.Errorf("invalid shift date format")
	}

	err = s.shiftsRepo.RemoveEmployeeFromShiftByDetails(employeeID, shiftDate, req.ShiftType)
	if err != nil {
		s.log.Errorf("failed to remove shift: %v", err)
		return fmt.Errorf("failed to remove shift")
	}

	s.log.Infof("Successfully removed shift for employee ID %d", employeeID)
	return nil
}

func (s *employeeService) GetOnCallEmployees(currentTime time.Time, shiftBuffer time.Duration) ([]employeeV1.EmployeeResponse, error) {
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

func (s *employeeService) GetShiftWarnings(employeeID uint) ([]string, error) {
	s.log.Infof("Getting shift warnings for employee ID %d", employeeID)

	employee := &model.Employee{}
	err := s.emplRepo.GetEmployeeByID(employeeID, employee)
	if err != nil {
		s.log.Errorf("failed to get employee: %v", err)
		return nil, fmt.Errorf("employee not found")
	}

	warnings := []string{}

	now := time.Now()
	twoWeeksFromNow := now.AddDate(0, 0, 14)

	// First, check if the next two weeks have adequate coverage in general
	uncoveredShifts, err := s.findUncoveredShifts(employee.ProfileType, now, twoWeeksFromNow)
	if err != nil {
		s.log.Errorf("failed to find uncovered shifts: %v", err)
		return nil, fmt.Errorf("failed to check shift coverage")
	}

	// If there are no uncovered shifts, no warnings needed
	if len(uncoveredShifts) == 0 {
		s.log.Infof("All shifts in next 2 weeks have adequate coverage, no warnings needed for employee %d", employeeID)
		return warnings, nil
	}

	// There are uncovered shifts, now check if THIS employee should be warned
	// Get employee's shifts in the next 2 weeks
	var employeeShifts []model.Shift
	err = s.shiftsRepo.GetShiftsByEmployeeIDInDateRange(employeeID, now, twoWeeksFromNow, &employeeShifts)
	if err != nil {
		s.log.Errorf("failed to get employee shifts: %v", err)
		return nil, fmt.Errorf("failed to get shifts in date range")
	}

	// Check if employee has met weekly quota (5 days per week)
	weeklyQuotaMet := s.checkWeeklyQuota(employeeShifts, now)
	if weeklyQuotaMet {
		s.log.Infof("Employee %d has met weekly quota, no warnings needed even though there are uncovered shifts", employeeID)
		return warnings, nil
	}

	// Employee hasn't met quota AND there are uncovered shifts - warn them
	warning := fmt.Sprintf("There are %d shifts in the next 2 weeks that need %s coverage. Consider scheduling shifts to help your team.",
		len(uncoveredShifts), employee.ProfileType.String())
	warnings = append(warnings, warning)

	// Check for periods without any shifts for this employee
	if len(employeeShifts) == 0 {
		warnings = append(warnings, "You have no shifts scheduled for the next 2 weeks. Consider scheduling shifts to meet your weekly quota.")
	}

	s.log.Infof("Generated %d warnings for employee %d", len(warnings), employeeID)
	return warnings, nil
}

func (s *employeeService) getMaxCapacityForProfile(profileType model.ProfileType) int {
	switch profileType {
	case model.Medic:
		return 2
	case model.Technical:
		return 4
	default:
		return 0
	}
}

func (s *employeeService) validateConsecutiveShifts(employeeID uint, newShiftDate time.Time) error {
	s.log.Infof("Validating consecutive shifts for employee ID %d", employeeID)

	// Get shifts in a range around the new shift date to check for consecutive patterns
	// We need to check 3 days before and after to ensure we catch all consecutive patterns
	startDate := newShiftDate.AddDate(0, 0, -3)
	endDate := newShiftDate.AddDate(0, 0, 3)

	var existingShifts []model.Shift
	err := s.shiftsRepo.GetShiftsByEmployeeIDInDateRange(employeeID, startDate, endDate, &existingShifts)
	if err != nil {
		s.log.Errorf("failed to get existing shifts: %v", err)
		return fmt.Errorf("failed to validate consecutive shifts")
	}

	// Create a map of dates to track which days have shifts
	shiftDates := make(map[string]bool)
	for _, shift := range existingShifts {
		dateKey := shift.ShiftDate.Format("2006-01-02")
		shiftDates[dateKey] = true
	}

	// Add the new shift date to the map
	newShiftDateKey := newShiftDate.Format("2006-01-02")
	shiftDates[newShiftDateKey] = true

	// Check for consecutive shifts starting from the new shift date
	consecutiveCount := s.countConsecutiveShifts(shiftDates, newShiftDate)

	if consecutiveCount > 2 {
		s.log.Errorf("employee would have %d consecutive shifts, maximum allowed is 2", consecutiveCount)
		return fmt.Errorf("cannot assign shift: would result in more than 2 consecutive shifts")
	}

	s.log.Infof("Consecutive shifts validation passed: %d consecutive shifts", consecutiveCount)
	return nil
}

// countConsecutiveShifts counts the maximum consecutive shifts that would include the given date
func (s *employeeService) countConsecutiveShifts(shiftDates map[string]bool, centerDate time.Time) int {
	// Find the start of the consecutive sequence that includes centerDate
	startDate := centerDate
	for {
		prevDate := startDate.AddDate(0, 0, -1)
		prevDateKey := prevDate.Format("2006-01-02")
		if !shiftDates[prevDateKey] {
			break
		}
		startDate = prevDate
	}

	consecutiveCount := 0
	currentDate := startDate
	for {
		currentDateKey := currentDate.Format("2006-01-02")
		if !shiftDates[currentDateKey] {
			break
		}
		consecutiveCount++
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	return consecutiveCount
}

// checkWeeklyQuota checks if the employee has met the quota of 5 days per week
func (s *employeeService) checkWeeklyQuota(shifts []model.Shift, startDate time.Time) bool {
	// Count shifts per week
	firstWeekCount := 0
	secondWeekCount := 0

	endOfFirstWeek := startDate.AddDate(0, 0, 7)

	for _, shift := range shifts {
		if shift.ShiftDate.Before(endOfFirstWeek) {
			firstWeekCount++
		} else {
			secondWeekCount++
		}
	}

	return firstWeekCount >= 5 && secondWeekCount >= 5
}

// findUncoveredShifts finds shifts that don't have coverage for the specified profile type
func (s *employeeService) findUncoveredShifts(profileType model.ProfileType, startDate, endDate time.Time) ([]model.Shift, error) {
	var uncoveredShifts []model.Shift

	availability, err := s.shiftsRepo.GetShiftAvailability(startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to get shift availability: %w", err)
	}

	for date, shifts := range availability.Days {
		for shiftIndex, shiftMap := range shifts {
			if currentCount, exists := shiftMap[profileType]; exists && currentCount == 0 {
				if currentCount == 0 {
					// This shift has no coverage for this profile type
					shiftType := shiftIndex + 1 // shift types are 1-indexed
					shift := model.Shift{
						ShiftDate: date,
						ShiftType: shiftType,
					}
					uncoveredShifts = append(uncoveredShifts, shift)
				}
			}
		}
	}

	return uncoveredShifts, nil
}
