// model/mapper.go
package model

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

func MapShiftsAvailabilityToResponse(availability *ShiftsAvailability) *ShiftAvailabilityResponse {
	return &ShiftAvailabilityResponse{
		FirstShift: ShiftAvailabilityDto{
			Medic:     availability.Availability[1][Medic],
			Technical: availability.Availability[1][Technical],
		},
		SecondShift: ShiftAvailabilityDto{
			Medic:     availability.Availability[2][Medic],
			Technical: availability.Availability[2][Technical],
		},
		ThirdShift: ShiftAvailabilityDto{
			Medic:     availability.Availability[3][Medic],
			Technical: availability.Availability[3][Technical],
		},
	}
}
