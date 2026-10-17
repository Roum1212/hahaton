package service_list_cooperation_invitations

import (
	"fmt"

	aggregate_cooperation_invitation "github.com/Roum1212/hahaton/internal/domain/aggregate/cooperation-invitation"
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID         model_user_id.UserID
	includeExpired bool
	page           model_page.Page
}

func (x Request) GetIncludeExpired() bool {
	return x.includeExpired
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

func NewRequest(userID model_user_id.UserID, includeExpired bool, page model_page.Page) Request {
	return Request{
		userID:         userID,
		includeExpired: includeExpired,
		page:           page,
	}
}

type Response struct {
	invitations []*aggregate_cooperation_invitation.CooperationInvitation
	totalCount  int64
}

func (x Response) GetInvitations() []*aggregate_cooperation_invitation.CooperationInvitation {
	return x.invitations
}

func (x Response) GetTotalCount() int64 {
	return x.totalCount
}

func NewResponse(invitations []*aggregate_cooperation_invitation.CooperationInvitation, totalCount int64) Response {
	return Response{
		invitations: invitations,
		totalCount:  totalCount,
	}
}
