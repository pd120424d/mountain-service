// model/mapper.go
package model

import "time"

// MapUpdateRequestToEmployee updates allowed fields from the request into the employee struct.
func MapUpdateRequestToEmployee(req *EmployeeUpdateRequest, existing *Employee) {
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

func MapShiftsAvailabilityToResponse(availability *ShiftsAvailabilityRange) *ShiftAvailabilityResponse {
	response := &ShiftAvailabilityResponse{
		Days: make(map[time.Time]ShiftAvailabilityPerDay),
	}

	for date, day := range availability.Days {
		if len(day) < 3 {
			continue // skip incomplete days, ideally should never be the case
		}

		response.Days[date] = ShiftAvailabilityPerDay{
			FirstShift: ShiftAvailabilityDto{
				Medic:     day[0][Medic],
				Technical: day[0][Technical],
			},
			SecondShift: ShiftAvailabilityDto{
				Medic:     day[1][Medic],
				Technical: day[1][Technical],
			},
			ThirdShift: ShiftAvailabilityDto{
				Medic:     day[2][Medic],
				Technical: day[2][Technical],
			},
		}
	}

	return response
}
