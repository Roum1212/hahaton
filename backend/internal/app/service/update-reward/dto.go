package service_update_reward

import (
	"errors"
	"fmt"
	"strings"

	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	utils_value "github.com/Roum1212/hahaton/internal/pkg/utils/value"
)

type Request struct {
	userID   model_user_id.UserID
	rewardID model_reward_id.RewardID

	description utils_value.Optional[string]
	imagePath   utils_value.Optional[model_image_path.ImagePath]
	title       utils_value.Optional[string]

	priceCoins    utils_value.Optional[int64]
	quantityLimit utils_value.Optional[int32]
}

func (x Request) GetDescription() utils_value.Optional[string] {
	return x.description
}

func (x Request) GetImagePath() utils_value.Optional[model_image_path.ImagePath] {
	return x.imagePath
}

func (x Request) GetPriceCoins() utils_value.Optional[int64] {
	return x.priceCoins
}

func (x Request) GetQuantityLimit() utils_value.Optional[int32] {
	return x.quantityLimit
}

func (x Request) GetRewardID() model_reward_id.RewardID {
	return x.rewardID
}

func (x Request) GetTitle() utils_value.Optional[string] {
	return x.title
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if err := x.rewardID.Validate(); err != nil {
		return fmt.Errorf("reward id: %w", err)
	}

	if value, ok := x.title.Get(); ok && len(strings.TrimSpace(value)) == 0 {
		return errors.New("reward title is required")
	}

	if value, ok := x.priceCoins.Get(); ok && value <= 0 {
		return errors.New("price must be positive")
	}

	if value, ok := x.quantityLimit.Get(); ok && value < 0 {
		return errors.New("quantity must not be negative")
	}

	if value, ok := x.imagePath.Get(); ok {
		if err := value.Validate(); err != nil {
			return fmt.Errorf("image path: %w", err)
		}
	}

	return nil
}

type RequestOption func(x *Request)

func WithDescription(value string) RequestOption {
	return func(x *Request) {
		x.description = utils_value.NewOptional(value)
	}
}

func WithImagePath(value model_image_path.ImagePath) RequestOption {
	return func(x *Request) {
		x.imagePath = utils_value.NewOptional(value)
	}
}

func WithPriceCoins(value int64) RequestOption {
	return func(x *Request) {
		x.priceCoins = utils_value.NewOptional(value)
	}
}

func WithQuantityLimit(value int32) RequestOption {
	return func(x *Request) {
		x.quantityLimit = utils_value.NewOptional(value)
	}
}

func WithTitle(value string) RequestOption {
	return func(x *Request) {
		x.title = utils_value.NewOptional(value)
	}
}

func NewRequest(userID model_user_id.UserID, rewardID model_reward_id.RewardID, opts ...RequestOption) Request {
	x := Request{
		userID:   userID,
		rewardID: rewardID,
	}
	for _, option := range opts {
		option(&x)
	}

	return x
}
