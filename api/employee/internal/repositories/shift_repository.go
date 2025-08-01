package repositories

//go:generate mockgen -source=shift_repository.go -destination=shift_repository_gomock.go -package=repositories

import (
	"errors"
	"fmt"
	"time"

	"github.com/pd120424d/mountain-service/api/employee/internal/model"
	"github.com/pd120424d/mountain-service/api/shared/utils"

	"gorm.io/gorm"
)

type ShiftRepository interface {
	GetOrCreateShift(shiftDate time.Time, shiftType int) (*model.Shift, error)
	AssignedToShift(employeeID, shiftID uint) (bool, error)
	CountAssignmentsByProfile(shiftID uint, profileType model.ProfileType) (int64, error)
	CreateAssignment(employeeID, shiftID uint) (uint, error)
	GetShiftsByEmployeeID(employeeID uint, result *[]model.Shift) error
	GetShiftsByEmployeeIDInDateRange(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error
	GetShiftAvailability(start, end time.Time) (*model.ShiftsAvailabilityRange, error)
	GetShiftAvailabilityWithEmployeeStatus(employeeID uint, start, end time.Time) (*model.ShiftsAvailabilityWithEmployeeStatus, error)
	RemoveEmployeeFromShiftByDetails(employeeID uint, shiftDate time.Time, shiftType int) error
	GetOnCallEmployees(currentTime time.Time, shiftBuffer time.Duration) ([]model.Employee, error)
}

type shiftRepository struct {
	log utils.Logger
	db  *gorm.DB
}

func NewShiftRepository(log utils.Logger, db *gorm.DB) ShiftRepository {
	return &shiftRepository{log: log.WithName("shiftRepository"), db: db}
}

func (r *shiftRepository) GetOrCreateShift(shiftDate time.Time, shiftType int) (*model.Shift, error) {
	var shift model.Shift
	err := r.db.FirstOrCreate(&shift, model.Shift{
		ShiftDate: shiftDate,
		ShiftType: shiftType,
	}).Error
	if err != nil {
		return nil, fmt.Errorf("failed to find or create shift: %w", err)
	}
	return &shift, nil
}

func (r *shiftRepository) AssignedToShift(employeeID, shiftID uint) (bool, error) {
	var existing model.EmployeeShift
	err := r.db.Where("employee_id = ? AND shift_id = ?", employeeID, shiftID).First(&existing).Error
	if err == nil {
		return true, nil
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false, nil
	}
	return false, fmt.Errorf("failed to check assignment: %w", err)
}

func (r *shiftRepository) CountAssignmentsByProfile(shiftID uint, profileType model.ProfileType) (int64, error) {
	var count int64
	err := r.db.Table("employee_shifts").
		Joins("JOIN employees ON employee_shifts.employee_id = employees.id").
		Where("employee_shifts.shift_id = ? AND employees.profile_type = ?", shiftID, profileType).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to count assignments: %w", err)
	}
	return count, nil
}

func (r *shiftRepository) CreateAssignment(employeeID, shiftID uint) (uint, error) {
	assignment := model.EmployeeShift{
		EmployeeID: employeeID,
		ShiftID:    shiftID,
	}
	if err := r.db.Create(&assignment).Error; err != nil {
		return 0, fmt.Errorf("failed to create assignment: %w", err)
	}
	return assignment.ID, nil
}

func (r *shiftRepository) GetShiftsByEmployeeID(employeeID uint, result *[]model.Shift) error {
	return r.db.Table("employee_shifts").
		Select("employee_shifts.id, shifts.shift_date, shifts.shift_type").
		Joins("JOIN shifts ON employee_shifts.shift_id = shifts.id").
		Where("employee_shifts.employee_id = ?", employeeID).
		Order("shifts.shift_date ASC, shifts.shift_type ASC").
		Scan(result).Error
}

func (r *shiftRepository) GetShiftsByEmployeeIDInDateRange(employeeID uint, startDate, endDate time.Time, result *[]model.Shift) error {
	return r.db.Table("employee_shifts").
		Select("shifts.id, shifts.shift_date, shifts.shift_type, shifts.created_at").
		Joins("JOIN shifts ON employee_shifts.shift_id = shifts.id").
		Where("employee_shifts.employee_id = ? AND shifts.shift_date >= ? AND shifts.shift_date <= ?", employeeID, startDate, endDate).
		Order("shifts.shift_date ASC, shifts.shift_type ASC").
		Scan(result).Error
}

func (r *shiftRepository) GetShiftAvailability(start, end time.Time) (*model.ShiftsAvailabilityRange, error) {
	result := model.ShiftsAvailabilityRange{
		Days: map[time.Time][]map[model.ProfileType]int{},
	}

	// Initial availability per shift
	for d := start; d.Before(end); d = d.Add(24 * time.Hour) {
		day := d.Truncate(24 * time.Hour)
		result.Days[day] = []map[model.ProfileType]int{
			{model.Medic: 2, model.Technical: 4},
			{model.Medic: 2, model.Technical: 4},
			{model.Medic: 2, model.Technical: 4},
		}
	}

	// Query assigned employees grouped by shift and role
	var counts []struct {
		ShiftDate    time.Time
		ShiftType    int
		EmployeeRole string
		Count        int
	}

	err := r.db.Table("shifts").
		Joins("JOIN employee_shifts ON shifts.id = employee_shifts.shift_id").
		Joins("JOIN employees ON employee_shifts.employee_id = employees.id").
		Select("shifts.shift_date, shifts.shift_type, employees.profile_type AS employee_role, COUNT(*) AS count").
		Where("shift_date >= ? AND shift_date < ?", start, end).
		Group("shifts.shift_date, shifts.shift_type, employees.profile_type").
		Scan(&counts).Error
	if err != nil {
		return nil, err
	}

	// Deduct from initialized capacities
	for _, c := range counts {
		day := c.ShiftDate.Truncate(24 * time.Hour)
		role := model.ProfileTypeFromString(c.EmployeeRole)
		shiftIndex := c.ShiftType - 1

		if shifts, ok := result.Days[day]; ok && shiftIndex >= 0 && shiftIndex < len(shifts) {
			shifts[shiftIndex][role] -= c.Count
		}
	}

	return &result, nil
}

func (r *shiftRepository) GetShiftAvailabilityWithEmployeeStatus(employeeID uint, start, end time.Time) (*model.ShiftsAvailabilityWithEmployeeStatus, error) {
	result := model.ShiftsAvailabilityWithEmployeeStatus{
		Days: map[time.Time][]model.ShiftAvailabilityWithStatus{},
	}

	// Initial availability per shift
	for d := start; d.Before(end); d = d.Add(24 * time.Hour) {
		day := d.Truncate(24 * time.Hour)
		result.Days[day] = []model.ShiftAvailabilityWithStatus{
			{MedicSlotsAvailable: 2, TechnicalSlotsAvailable: 4, IsAssignedToEmployee: false, IsFullyBooked: false},
			{MedicSlotsAvailable: 2, TechnicalSlotsAvailable: 4, IsAssignedToEmployee: false, IsFullyBooked: false},
			{MedicSlotsAvailable: 2, TechnicalSlotsAvailable: 4, IsAssignedToEmployee: false, IsFullyBooked: false},
		}
	}

	// Query assigned employees grouped by shift and role
	var counts []struct {
		ShiftDate    time.Time
		ShiftType    int
		EmployeeRole string
		Count        int
	}

	err := r.db.Table("shifts").
		Joins("JOIN employee_shifts ON shifts.id = employee_shifts.shift_id").
		Joins("JOIN employees ON employee_shifts.employee_id = employees.id").
		Select("shifts.shift_date, shifts.shift_type, employees.profile_type AS employee_role, COUNT(*) AS count").
		Where("shift_date >= ? AND shift_date < ?", start, end).
		Group("shifts.shift_date, shifts.shift_type, employees.profile_type").
		Scan(&counts).Error
	if err != nil {
		return nil, err
	}

	// Apply the counts to reduce availability
	for _, count := range counts {
		day := count.ShiftDate.Truncate(24 * time.Hour)
		if dayShifts, exists := result.Days[day]; exists {
			shiftIndex := count.ShiftType - 1 // Convert to 0-based index
			if shiftIndex >= 0 && shiftIndex < len(dayShifts) {
				switch count.EmployeeRole {
				case "Medic":
					result.Days[day][shiftIndex].MedicSlotsAvailable = max(0, 2-count.Count)
				case "Technical":
					result.Days[day][shiftIndex].TechnicalSlotsAvailable = max(0, 4-count.Count)
				}
			}
		}
	}

	var employeeAssignments []struct {
		ShiftDate time.Time
		ShiftType int
	}

	// Check if employee is assigned to each shift and if shifts are fully booked
	err = r.db.Table("shifts").
		Joins("JOIN employee_shifts ON shifts.id = employee_shifts.shift_id").
		Select("shifts.shift_date, shifts.shift_type").
		Where("employee_shifts.employee_id = ? AND shift_date >= ? AND shift_date < ?", employeeID, start, end).
		Scan(&employeeAssignments).Error
	if err != nil {
		return nil, err
	}

	for _, assignment := range employeeAssignments {
		day := assignment.ShiftDate.Truncate(24 * time.Hour)
		if dayShifts, exists := result.Days[day]; exists {
			shiftIndex := assignment.ShiftType - 1
			if shiftIndex >= 0 && shiftIndex < len(dayShifts) {
				result.Days[day][shiftIndex].IsAssignedToEmployee = true
			}
		}
	}

	// Mark fully booked shifts
	for day, dayShifts := range result.Days {
		for i, shift := range dayShifts {
			result.Days[day][i].IsFullyBooked = shift.MedicSlotsAvailable == 0 && shift.TechnicalSlotsAvailable == 0
		}
	}

	return &result, nil
}

func (r *shiftRepository) RemoveEmployeeFromShiftByDetails(employeeID uint, shiftDate time.Time, shiftType int) error {
	var shift model.Shift
	err := r.db.Where("shift_date = ? AND shift_type = ?", shiftDate, shiftType).First(&shift).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("shift not found for date %s and type %d", shiftDate.Format(time.DateOnly), shiftType)
		}
		return fmt.Errorf("failed to find shift: %w", err)
	}

	var assignment model.EmployeeShift
	err = r.db.Where("employee_id = ? AND shift_id = ?", employeeID, shift.ID).First(&assignment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("employee is not assigned to this shift")
		}
		return fmt.Errorf("failed to find assignment: %w", err)
	}

	if err := r.db.Delete(&assignment).Error; err != nil {
		return fmt.Errorf("failed to remove employee from shift: %w", err)
	}

	return nil
}

// GetOnCallEmployees returns all emloyees who are assigned to the current shift with one exception:
// If the current shift is ending soon (within the shiftBuffer), we also include employees assigned to the next shift
// If the shiftBuffer is 0, we only include employees assigned to the current shift
func (r *shiftRepository) GetOnCallEmployees(currentTime time.Time, shiftBuffer time.Duration) ([]model.Employee, error) {
	r.log.Infof("Getting on-call employees at %v with buffer %v", currentTime, shiftBuffer)

	currentShiftType := r.getShiftTypeForTime(currentTime)
	currentDate := currentTime.Truncate(24 * time.Hour)

	var employees []model.Employee
	var shiftDates []time.Time
	var shiftTypes []int

	shiftDates = append(shiftDates, currentDate)
	shiftTypes = append(shiftTypes, currentShiftType)

	// if buffer is defined we may need to include the following shift as well
	if shiftBuffer > 0 {
		timeUntilShiftEnd := r.getTimeUntilShiftEnd(currentTime, currentShiftType)
		if timeUntilShiftEnd <= shiftBuffer {
			r.log.Infof("Including next shift due to buffer: timeUntilShiftEnd=%v, buffer=%v", timeUntilShiftEnd, shiftBuffer)

			nextShiftType, nextShiftDate := r.getNextShift(currentShiftType, currentDate)
			shiftDates = append(shiftDates, nextShiftDate)
			shiftTypes = append(shiftTypes, nextShiftType)
		}
	}

	query := r.db.Distinct().
		Select("employees.*").
		Table("employees").
		Joins("JOIN employee_shifts ON employees.id = employee_shifts.employee_id").
		Joins("JOIN shifts ON employee_shifts.shift_id = shifts.id").
		Where("(shifts.shift_date = ? AND shifts.shift_type = ?)", shiftDates[0], shiftTypes[0])

	for i := 1; i < len(shiftDates); i++ {
		query = query.Or("(shifts.shift_date = ? AND shifts.shift_type = ?)", shiftDates[i], shiftTypes[i])
	}

	if err := query.Find(&employees).Error; err != nil {
		r.log.Errorf("Failed to get on-call employees: %v", err)
		return nil, fmt.Errorf("failed to get on-call employees: %w", err)
	}

	r.log.Infof("Successfully retrieved on-call employees: count=%d", len(employees))
	return employees, nil
}

func (r *shiftRepository) getShiftTypeForTime(currentTime time.Time) int {
	hour := currentTime.Hour()

	// Shift types: 1: 6am-2pm, 2: 2pm-10pm, 3: 10pm-6am
	if hour >= 6 && hour < 14 {
		return 1
	} else if hour >= 14 && hour < 22 {
		return 2
	}
	return 3
}

func (r *shiftRepository) getTimeUntilShiftEnd(currentTime time.Time, shiftType int) time.Duration {
	hour := currentTime.Hour()
	minute := currentTime.Minute()
	second := currentTime.Second()

	currentMinutes := hour*60 + minute
	currentSeconds := currentMinutes*60 + second

	var shiftEndSeconds int
	switch shiftType {
	case 1: // 6am-2pm, ends at 14:00
		shiftEndSeconds = 14 * 3600
	case 2: // 2pm-10pm, ends at 22:00
		shiftEndSeconds = 22 * 3600
	case 3: // 10pm-6am, ends at 6:00 next day
		if hour >= 22 {
			// Same day, ends at 6am next day
			shiftEndSeconds = 24*3600 + 6*3600
		} else {
			// Next day, ends at 6am
			shiftEndSeconds = 6 * 3600
		}
	default:
		return 0
	}

	remainingSeconds := shiftEndSeconds - currentSeconds
	if remainingSeconds < 0 {
		remainingSeconds += 24 * 3600 // Add 24 hours for overnight shifts
	}

	return time.Duration(remainingSeconds) * time.Second
}

func (r *shiftRepository) getNextShift(currentShiftType int, currentDate time.Time) (int, time.Time) {
	nextShiftType := currentShiftType + 1
	shiftDate := currentDate

	// if we are in the last shift, the next one is the first one and is on the next day
	if nextShiftType > 3 {
		return 1, currentDate.Add(24 * time.Hour)
	}

	return nextShiftType, shiftDate
}
