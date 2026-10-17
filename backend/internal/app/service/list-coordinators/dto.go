package service_list_coordinators

import (
	"fmt"

	aggregate_coordinator "github.com/Roum1212/hahaton/internal/domain/aggregate/coordinator"
	aggregate_volunteer "github.com/Roum1212/hahaton/internal/domain/aggregate/volunteer"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID model_user_id.UserID
	page   model_page.Page
}

func (x Request) GetPage() model_page.Page {
	return x.page
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.page.Validate(); err != nil {
		return fmt.Errorf("page: %w", err)
	}

	return nil
}

func NewRequest(userID model_user_id.UserID, page model_page.Page) Request {
	return Request{
		userID: userID,
		page:   page,
	}
}

type Response struct {
	coordinators []*aggregate_coordinator.Coordinator
	volunteers   map[model_user_id.UserID]*aggregate_volunteer.Volunteer
	totalCount   int64
}

func (x Response) GetCoordinators() []*aggregate_coordinator.Coordinator {
	return x.coordinators
}

func (x Response) GetTotalCount() int64 {
	return x.totalCount
}

func (x Response) GetVolunteers() map[model_user_id.UserID]*aggregate_volunteer.Volunteer {
	return x.volunteers
}

func NewResponse(
	coordinators []*aggregate_coordinator.Coordinator,
	volunteers map[model_user_id.UserID]*aggregate_volunteer.Volunteer,
	totalCount int64,
) Response {
	return Response{
		coordinators: coordinators,
		volunteers:   volunteers,
		totalCount:   totalCount,
	}
}
