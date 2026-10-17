package aggregate_reward

import (
	"errors"
	"time"

	aggregate_reward_command "github.com/Roum1212/hahaton/internal/domain/aggregate/reward/command"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
	model_reward_id "github.com/Roum1212/hahaton/internal/domain/model/reward-id"
	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
)

var (
	ErrInvalidState     = errors.New("invalid reward state")
	ErrOutOfStock       = errors.New("reward is out of stock")
	ErrPriceChanged     = errors.New("reward price has changed")
	ErrReservedQuantity = errors.New("quantity is below reserved and issued units")
)

type Reward struct {
	rewardID         model_reward_id.RewardID
	coordinatorID    model_user_id.UserID
	title            string
	description      string
	imagePath        model_image_path.ImagePath
	priceCoins       int64
	quantityLimit    int32
	reservedQuantity int32
	issuedQuantity   int32
	deleted          bool
	createdAt        time.Time
	updatedAt        time.Time
}

func (x *Reward) GetAvailableQuantity() int32 {
	return x.quantityLimit - x.reservedQuantity - x.issuedQuantity
}

func (x *Reward) GetCoordinatorID() model_user_id.UserID {
	return x.coordinatorID
}

func (x *Reward) GetCreatedAt() time.Time {
	return x.createdAt
}

func (x *Reward) GetDeleted() bool {
	return x.deleted
}

func (x *Reward) GetDescription() string {
	return x.description
}

func (x *Reward) GetImagePath() model_image_path.ImagePath {
	return x.imagePath
}

func (x *Reward) GetIssuedQuantity() int32 {
	return x.issuedQuantity
}

func (x *Reward) GetPriceCoins() int64 {
	return x.priceCoins
}

func (x *Reward) GetQuantityLimit() int32 {
	return x.quantityLimit
}

func (x *Reward) GetReservedQuantity() int32 {
	return x.reservedQuantity
}

func (x *Reward) GetRewardID() model_reward_id.RewardID {
	return x.rewardID
}

func (x *Reward) GetTitle() string {
	return x.title
}

func (x *Reward) GetUpdatedAt() time.Time {
	return x.updatedAt
}

func (x *Reward) Create(command aggregate_reward_command.Create) error {
	if !x.createdAt.IsZero() {
		return ErrInvalidState
	}

	x.title = command.GetTitle()
	x.description = command.GetDescription()
	x.imagePath = command.GetImagePath()
	x.priceCoins = command.GetPriceCoins()
	x.quantityLimit = command.GetQuantityLimit()
	x.createdAt = command.GetNow()
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Reward) Delete(command aggregate_reward_command.Delete) error {
	if x.reservedQuantity != 0 {
		return ErrReservedQuantity
	}

	x.deleted = true
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Reward) Issue(command aggregate_reward_command.Issue) error {
	if x.deleted || x.reservedQuantity < 1 {
		return ErrInvalidState
	}

	x.reservedQuantity--
	x.issuedQuantity++
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Reward) Release(command aggregate_reward_command.Release) error {
	if x.reservedQuantity < 1 {
		return ErrInvalidState
	}

	x.reservedQuantity--
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Reward) Reserve(command aggregate_reward_command.Reserve) error {
	if x.deleted {
		return ErrInvalidState
	}

	if command.GetExpectedPrice() != x.priceCoins {
		return ErrPriceChanged
	}

	if x.GetAvailableQuantity() < 1 {
		return ErrOutOfStock
	}

	x.reservedQuantity++
	x.updatedAt = command.GetNow()

	return nil
}

func (x *Reward) Update(command aggregate_reward_command.Update) error {
	if x.deleted {
		return ErrInvalidState
	}

	if command.GetQuantityLimit() < x.reservedQuantity+x.issuedQuantity {
		return ErrReservedQuantity
	}

	x.title = command.GetTitle()
	x.description = command.GetDescription()
	x.imagePath = command.GetImagePath()
	x.priceCoins = command.GetPriceCoins()
	x.quantityLimit = command.GetQuantityLimit()
	x.updatedAt = command.GetNow()

	return nil
}

func NewReward(rewardID model_reward_id.RewardID, coordinatorID model_user_id.UserID) *Reward {
	return &Reward{
		rewardID:      rewardID,
		coordinatorID: coordinatorID,
	}
}

func NewRewardFromFields(
	rewardID model_reward_id.RewardID,
	coordinatorID model_user_id.UserID,
	title string,
	description string,
	imagePath model_image_path.ImagePath,
	priceCoins int64,
	quantityLimit int32,
	reservedQuantity int32,
	issuedQuantity int32,
	deleted bool,
	createdAt time.Time,
	updatedAt time.Time,
) *Reward {
	return &Reward{
		rewardID:         rewardID,
		coordinatorID:    coordinatorID,
		title:            title,
		description:      description,
		imagePath:        imagePath,
		priceCoins:       priceCoins,
		quantityLimit:    quantityLimit,
		reservedQuantity: reservedQuantity,
		issuedQuantity:   issuedQuantity,
		deleted:          deleted,
		createdAt:        createdAt,
		updatedAt:        updatedAt,
	}
}
