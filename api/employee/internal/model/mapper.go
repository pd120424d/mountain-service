// model/mapper.go
package model

import (
	"time"

	employeeV1 "github.com/pd120424d/mountain-service/api/contracts/employee/v1"
)

func MapUpdateRequestToEmployee(req *employeeV1.EmployeeUpdateRequest, existing *Employee) {
	if req.FirstName != "" {
		existing.FirstName = req.FirstName
	}
	if req.LastName != "" {
		existing.LastName = req.LastName
	}
	if req.Email != "" {
		existing.Email = req.Email
	}
	existing.Gender = req.Gender
	existing.Phone = req.Phone
	existing.ProfilePicture = req.ProfilePicture
	existing.ProfileType = ProfileTypeFromString(req.ProfileType)
}

func MapShiftsAvailabilityToResponse(availability *ShiftsAvailabilityRange) *employeeV1.ShiftAvailabilityResponse {
	response := &employeeV1.ShiftAvailabilityResponse{
		Days: make(map[time.Time]employeeV1.ShiftAvailabilityPerDay),
	}

	for date, day := range availability.Days {
		if len(day) < 3 {
			continue
		}

		response.Days[date] = employeeV1.ShiftAvailabilityPerDay{
			Shift1: employeeV1.ShiftAvailability{
				Available: day[0][Medic] > 0 || day[0][Technical] > 0,
				Employees: []string{},
			},
			Shift2: employeeV1.ShiftAvailability{
				Available: day[1][Medic] > 0 || day[1][Technical] > 0,
				Employees: []string{},
			},
			Shift3: employeeV1.ShiftAvailability{
				Available: day[2][Medic] > 0 || day[2][Technical] > 0,
				Employees: []string{},
			},
		}
	}

	return response
}
