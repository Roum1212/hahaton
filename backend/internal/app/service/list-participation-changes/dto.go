package service_list_participation_changes

import (
	"fmt"

	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	model_participation_change "github.com/Roum1212/hahaton/internal/domain/model/participation-change"
	model_participation_id "github.com/Roum1212/hahaton/internal/domain/model/participation-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID          model_user_id.UserID
	participationID model_participation_id.ParticipationID
	page            model_page.Page
}

func (x Request) GetPage() model_page.Page {
	return x.page
}

func (x Request) GetParticipationID() model_participation_id.ParticipationID {
	return x.participationID
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.participationID.Validate(); err != nil {
		return fmt.Errorf("participation id: %w", err)
	}

	if err := x.page.Validate(); err != nil {
		return fmt.Errorf("page: %w", err)
	}

	return nil
}

func NewRequest(
	userID model_user_id.UserID,
	participationID model_participation_id.ParticipationID,
	page model_page.Page,
) Request {
	return Request{
		userID:          userID,
		participationID: participationID,
		page:            page,
	}
}

type Response struct {
	changes    []model_participation_change.ParticipationChange
	totalCount int64
}

func (x Response) GetChanges() []model_participation_change.ParticipationChange {
	return x.changes
}

func (x Response) GetTotalCount() int64 {
	return x.totalCount
}

func NewResponse(changes []model_participation_change.ParticipationChange, totalCount int64) Response {
	return Response{
		changes:    changes,
		totalCount: totalCount,
	}
}
