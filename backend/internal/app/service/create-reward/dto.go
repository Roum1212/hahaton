package service_create_reward

import (
	"errors"
	"fmt"
	"strings"

	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type Request struct {
	userID        model_user_id.UserID
	title         string
	description   string
	imagePath     model_image_path.ImagePath
	priceCoins    int64
	quantityLimit int32
}

func (x Request) GetDescription() string {
	return x.description
}

func (x Request) GetImagePath() model_image_path.ImagePath {
	return x.imagePath
}

func (x Request) GetPriceCoins() int64 {
	return x.priceCoins
}

func (x Request) GetQuantityLimit() int32 {
	return x.quantityLimit
}

func (x Request) GetTitle() string {
	return x.title
}

func (x Request) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x Request) Validate() error {
	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("user id: %w", err)
	}

	if len(strings.TrimSpace(x.title)) == 0 {
		return errors.New("reward title is required")
	}

	if x.priceCoins <= 0 || x.quantityLimit <= 0 {
		return errors.New("reward price and quantity must be positive")
	}

	if err := x.imagePath.Validate(); err != nil {
		return fmt.Errorf("image path: %w", err)
	}

	return nil
}

func NewRequest(
	userID model_user_id.UserID,
	title string,
	description string,
	imagePath model_image_path.ImagePath,
	priceCoins int64,
	quantityLimit int32,
) Request {
	return Request{
		userID:        userID,
		title:         title,
		description:   description,
		imagePath:     imagePath,
		priceCoins:    priceCoins,
		quantityLimit: quantityLimit,
	}
}

type Response struct {
	rewardID model_reward_id.RewardID
}

func (x Response) GetRewardID() model_reward_id.RewardID {
	return x.rewardID
}

func NewResponse(rewardID model_reward_id.RewardID) Response {
	return Response{
		rewardID: rewardID,
	}
}
