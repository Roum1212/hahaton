package service_suggest_addresses

import (
	"errors"
	"fmt"
	"slices"
	"strings"
	"unicode/utf8"

	model_address_suggestion "github.com/Roum1212/hahaton/internal/domain/model/address-suggestion"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID    model_user_id.UserID
	query     string
	fromBound string
	toBound   string
	count     int32
}

func (x Request) GetCount() int32 {
	return x.count
}

func (x Request) GetFromBound() string {
	return x.fromBound
}

func (x Request) GetQuery() string {
	return x.query
}

func (x Request) GetToBound() string {
	return x.toBound
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if len(strings.TrimSpace(x.query)) == 0 || utf8.RuneCountInString(strings.TrimSpace(x.query)) > 300 {
		return errors.New("query must contain 1 to 300 characters")
	}

	bounds := []string{
		"country",
		"region",
		"area",
		"city",
		"settlement",
		"street",
		"house",
	}

	from, to := slices.Index(bounds, x.fromBound), slices.Index(bounds, x.toBound)
	if from < 0 || to < 0 || from > to {
		return errors.New("invalid address bound range")
	}

	if x.count < 0 || x.count > 20 {
		return errors.New("count must be between 1 and 20 or zero")
	}

	return nil
}

func NewRequest(userID model_user_id.UserID, query, fromBound, toBound string, count int32) Request {
	return Request{
		userID:    userID,
		query:     query,
		fromBound: fromBound,
		toBound:   toBound,
		count:     count,
	}
}

type Response struct {
	suggestions []model_address_suggestion.AddressSuggestion
}

func (x Response) GetSuggestions() []model_address_suggestion.AddressSuggestion {
	return x.suggestions
}

func NewResponse(suggestions []model_address_suggestion.AddressSuggestion) Response {
	return Response{
		suggestions: suggestions,
	}
}
