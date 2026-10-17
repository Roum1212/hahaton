package service_find_volunteers

import (
	"errors"
	"fmt"

	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	model_email_address "github.com/Roum1212/hahaton/internal/domain/model/email-address"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	model_phone_number "github.com/Roum1212/hahaton/internal/domain/model/phone-number"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

type Request struct {
	emailAddress utils_value.Optional[model_email_address.EmailAddress]
	phoneNumber  utils_value.Optional[model_phone_number.PhoneNumber]
	volunteerID  utils_value.Optional[model_user_id.UserID]

	page model_page.Page
}

func (x Request) GetEmailAddress() utils_value.Optional[model_email_address.EmailAddress] {
	return x.emailAddress
}

func (x Request) GetPage() model_page.Page {
	return x.page
}

func (x Request) GetPhoneNumber() utils_value.Optional[model_phone_number.PhoneNumber] {
	return x.phoneNumber
}

func (x Request) GetVolunteerID() utils_value.Optional[model_user_id.UserID] {
	return x.volunteerID
}

func (x Request) Validate() error {
	if value, ok := x.volunteerID.Get(); ok {
		if err := value.Validate(); err != nil {
			return err
		}
	}

	if value, ok := x.emailAddress.Get(); ok {
		if err := value.Validate(); err != nil {
			return err
		}
	}

	if value, ok := x.phoneNumber.Get(); ok {
		if err := value.Validate(); err != nil {
			return err
		}
	}

	if err := x.page.Validate(); err != nil {
		return fmt.Errorf("page: %w", err)
	}

	criteria := 0
	if _, ok := x.volunteerID.Get(); ok {
		criteria++
	}

	if _, ok := x.emailAddress.Get(); ok {
		criteria++
	}

	if _, ok := x.phoneNumber.Get(); ok {
		criteria++
	}

	if criteria != 1 {
		return errors.New("exactly one volunteer search criterion is required")
	}

	return nil
}

type RequestOption func(x *Request)

func WithEmailAddress(value model_email_address.EmailAddress) RequestOption {
	return func(x *Request) {
		x.emailAddress = utils_value.NewOptional(value)
	}
}

func WithPhoneNumber(value model_phone_number.PhoneNumber) RequestOption {
	return func(x *Request) {
		x.phoneNumber = utils_value.NewOptional(value)
	}
}

func WithVolunteerID(value model_user_id.UserID) RequestOption {
	return func(x *Request) {
		x.volunteerID = utils_value.NewOptional(value)
	}
}

func NewRequest(page model_page.Page, opts ...RequestOption) Request {
	x := Request{
		page: page,
	}
	for _, option := range opts {
		option(&x)
	}

	return x
}

type Response struct {
	volunteers             []*aggregate_volunteer.Volunteer
	hasCoordinatorProfiles map[model_user_id.UserID]bool
	totalCount             int64
}

func (x Response) GetHasCoordinatorProfiles() map[model_user_id.UserID]bool {
	return x.hasCoordinatorProfiles
}

func (x Response) GetTotalCount() int64 {
	return x.totalCount
}

func (x Response) GetVolunteers() []*aggregate_volunteer.Volunteer {
	return x.volunteers
}

func NewResponse(
	volunteers []*aggregate_volunteer.Volunteer,
	hasCoordinatorProfiles map[model_user_id.UserID]bool,
	totalCount int64,
) Response {
	return Response{
		volunteers:             volunteers,
		hasCoordinatorProfiles: hasCoordinatorProfiles,
		totalCount:             totalCount,
	}
}
