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
	if req.Gender != "" {
		existing.Gender = req.Gender
	}
	if req.Phone != "" {
		existing.Phone = req.Phone
	}
	if req.ProfilePicture != "" {
		existing.ProfilePicture = req.ProfilePicture
	}
	if req.ProfileType != "" {
		newProfileType := ProfileTypeFromString(req.ProfileType)
		if newProfileType != "" {
			existing.ProfileType = newProfileType
		}
	}
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
			FirstShift: employeeV1.ShiftAvailability{
				MedicSlotsAvailable:     max(0, day[0][Medic]),
				TechnicalSlotsAvailable: max(0, day[0][Technical]),
			},
			SecondShift: employeeV1.ShiftAvailability{
				MedicSlotsAvailable:     max(0, day[1][Medic]),
				TechnicalSlotsAvailable: max(0, day[1][Technical]),
			},
			ThirdShift: employeeV1.ShiftAvailability{
				MedicSlotsAvailable:     max(0, day[2][Medic]),
				TechnicalSlotsAvailable: max(0, day[2][Technical]),
			},
		}
	}

	return response
}
