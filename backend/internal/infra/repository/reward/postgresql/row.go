package repository_reward_postgresql

import (
	"fmt"
	"time"

	aggregate_reward "github.com/Roum1212/hahaton/internal/domain/aggregate/reward"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

type row struct {
	rewardID         string
	coordinatorID    string
	title            string
	description      string
	imagePath        string
	priceCoins       int64
	quantityLimit    int32
	reservedQuantity int32
	issuedQuantity   int32
	deleted          bool
	createdAt        time.Time
	updatedAt        time.Time
}

func (x *row) scan(source interface{ Scan(dest ...any) error }) error {
	return source.Scan(
		&x.rewardID,
		&x.coordinatorID,
		&x.title,
		&x.description,
		&x.imagePath,
		&x.priceCoins,
		&x.quantityLimit,
		&x.reservedQuantity,
		&x.issuedQuantity,
		&x.deleted,
		&x.createdAt,
		&x.updatedAt,
	)
}

func (x *row) toAggregate() (*aggregate_reward.Reward, error) {
	rewardID, err := model_reward_id.NewRewardID(x.rewardID)
	if err != nil {
		return nil, fmt.Errorf("restore reward id: %w", err)
	}

	coordinatorID, err := model_user_id.NewUserID(x.coordinatorID)
	if err != nil {
		return nil, fmt.Errorf("restore coordinator id: %w", err)
	}

	var imagePath model_image_path.ImagePath
	if len(x.imagePath) != 0 {
		value, err := model_image_path.NewImagePath(x.imagePath)
		if err != nil {
			return nil, fmt.Errorf("restore value: %w", err)
		}

		imagePath = value
	}

	return aggregate_reward.NewRewardFromFields(
		rewardID,
		coordinatorID,
		x.title,
		x.description,
		imagePath,
		x.priceCoins,
		x.quantityLimit,
		x.reservedQuantity,
		x.issuedQuantity,
		x.deleted,
		x.createdAt,
		x.updatedAt,
	), nil
}
