package service

import (
	"fmt"
	"time"

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
		return nil, fmt.Errorf("employee not found")
	}

	// Step 2: Parse and validate shift date
	shiftDate, err := time.Parse("2006-01-02", req.ShiftDate)
	if err != nil {
		s.log.Errorf("failed to parse shift date: %v", err)
		return nil, fmt.Errorf("invalid shift date format")
	}

	// Step 3: Validate shift date is in the future
	if shiftDate.Before(time.Now().Truncate(24 * time.Hour)) {
		s.log.Errorf("shift date %s is in the past", req.ShiftDate)
		return nil, fmt.Errorf("shift date must be in the future")
	}

	// Step 4: Validate shift date is within 3 months
	threeMonthsFromNow := time.Now().AddDate(0, 3, 0)
	if shiftDate.After(threeMonthsFromNow) {
		s.log.Errorf("shift date %s is more than 3 months in the future", req.ShiftDate)
		return nil, fmt.Errorf("shift date cannot be more than 3 months in the future")
	}

	// Step 5: Check consecutive shifts rule (max 6 consecutive shifts)
	if err := s.validateConsecutiveShifts(employeeID, shiftDate); err != nil {
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

	// Step 8: Check shift capacity
	currentCount, err := s.shiftsRepo.CountAssignmentsByProfile(shift.ID, employee.ProfileType)
	if err != nil {
		s.log.Errorf("failed to count assignments: %v", err)
		return nil, fmt.Errorf("failed to check shift capacity")
	}

	maxCapacity := s.getMaxCapacityForProfile(employee.ProfileType)
	if currentCount >= int64(maxCapacity) {
		s.log.Errorf("shift capacity full for profile type %s", employee.ProfileType.String())
		return nil, fmt.Errorf("shift capacity is full for %s staff", employee.ProfileType.String())
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

func (s *shiftService) GetShiftsAvailability(days int) (*employeeV1.ShiftAvailabilityResponse, error) {
	s.log.Infof("Getting shift availability for %d days", days)

	if days <= 0 || days > 90 {
		s.log.Errorf("invalid days parameter: %d", days)
		return nil, fmt.Errorf("days must be between 1 and 90")
	}

	start := time.Now().Truncate(24 * time.Hour)
	end := start.AddDate(0, 0, days)

	availability, err := s.shiftsRepo.GetShiftAvailability(start, end)
	if err != nil {
		s.log.Errorf("failed to get shift availability: %v", err)
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
		}
		dayAvailability.SecondShift = employeeV1.ShiftAvailability{
			MedicSlotsAvailable:     0,
			TechnicalSlotsAvailable: 0,
		}
		dayAvailability.ThirdShift = employeeV1.ShiftAvailability{
			MedicSlotsAvailable:     0,
			TechnicalSlotsAvailable: 0,
		}

		for shiftIndex, shiftMap := range shifts {
			var shiftAvailability *employeeV1.ShiftAvailability
			switch shiftIndex {
			case 0:
				shiftAvailability = &dayAvailability.FirstShift
			case 1:
				shiftAvailability = &dayAvailability.SecondShift
			case 2:
				shiftAvailability = &dayAvailability.ThirdShift
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

	s.log.Infof("Successfully retrieved shift availability for %d days", days)
	return response, nil
}

func (s *shiftService) RemoveShift(employeeID uint, req employeeV1.RemoveShiftRequest) error {
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
		return nil, fmt.Errorf("employee not found")
	}

	var warnings []string

	// Get next two weeks date range
	now := time.Now()
	twoWeeksFromNow := now.AddDate(0, 0, 14)

	// Check general coverage for next two weeks
	availability, err := s.shiftsRepo.GetShiftAvailability(now, twoWeeksFromNow)
	if err != nil {
		s.log.Errorf("failed to get shift availability: %v", err)
		return nil, fmt.Errorf("failed to check shift coverage")
	}

	// Check if there's insufficient coverage
	hasInsufficientCoverage := false
	for _, shifts := range availability.Days {
		for _, shift := range shifts {
			if shift[model.Medic] < 2 || shift[model.Technical] < 4 {
				hasInsufficientCoverage = true
				break
			}
		}
		if hasInsufficientCoverage {
			break
		}
	}

	// Only check individual quota if there's insufficient coverage
	if hasInsufficientCoverage {
		// Get employee's shifts in the next two weeks
		var employeeShifts []model.Shift
		err = s.shiftsRepo.GetShiftsByEmployeeIDInDateRange(employeeID, now, twoWeeksFromNow, &employeeShifts)
		if err != nil {
			s.log.Errorf("failed to get employee shifts: %v", err)
			return nil, fmt.Errorf("failed to check employee shifts")
		}

		// Check if employee has less than 5 days in next two weeks
		if len(employeeShifts) < 5 {
			warnings = append(warnings, fmt.Sprintf("You have only %d shifts scheduled in the next 2 weeks. Consider scheduling more shifts to meet the 5 days/week quota.", len(employeeShifts)))
		}
	}

	s.log.Infof("Successfully retrieved %d warnings for employee ID %d", len(warnings), employeeID)
	return warnings, nil
}

// Helper methods

func (s *shiftService) validateConsecutiveShifts(employeeID uint, shiftDate time.Time) error {
	// Get shifts for the employee in a range around the requested date
	startDate := shiftDate.AddDate(0, 0, -6) // 6 days before
	endDate := shiftDate.AddDate(0, 0, 6)    // 6 days after

	var shifts []model.Shift
	err := s.shiftsRepo.GetShiftsByEmployeeIDInDateRange(employeeID, startDate, endDate, &shifts)
	if err != nil {
		return fmt.Errorf("failed to get employee shifts: %w", err)
	}

	// Create a map of dates for quick lookup
	shiftDates := make(map[string]bool)
	for _, shift := range shifts {
		dateStr := shift.ShiftDate.Format("2006-01-02")
		shiftDates[dateStr] = true
	}

	// Add the requested shift date
	requestedDateStr := shiftDate.Format("2006-01-02")
	shiftDates[requestedDateStr] = true

	// Check for consecutive shifts
	consecutiveCount := 0
	maxConsecutive := 0

	// Check from 6 days before to 6 days after
	for i := -6; i <= 6; i++ {
		checkDate := shiftDate.AddDate(0, 0, i)
		checkDateStr := checkDate.Format("2006-01-02")

		if shiftDates[checkDateStr] {
			consecutiveCount++
			if consecutiveCount > maxConsecutive {
				maxConsecutive = consecutiveCount
			}
		} else {
			consecutiveCount = 0
		}
	}

	if maxConsecutive > 6 {
		return fmt.Errorf("assigning this shift would result in more than 6 consecutive shifts")
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
