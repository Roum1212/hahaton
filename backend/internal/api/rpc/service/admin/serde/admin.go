package rpc_service_admin_serde

import (
	admin_v1 "github.com/Roum1212/hahaton/genproto/admin/v1"
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
)

func VolunteerToPB(volunteer *aggregate_volunteer.Volunteer, hasCoordinator bool) *admin_v1.AdminVolunteer {
	name := volunteer.GetFullName()

	return &admin_v1.AdminVolunteer{
		Id:                    volunteer.GetUserID().String(),
		AvatarUrl:             volunteer.GetAvatarPath().GetURL(),
		FirstName:             name.GetFirstName(),
		LastName:              name.GetLastName(),
		MiddleName:            name.GetMiddleName(),
		EmailAddress:          volunteer.GetEmailAddress().String(),
		PhoneNumber:           volunteer.GetPhoneNumber().String(),
		HasCoordinatorProfile: hasCoordinator,
	}
}
