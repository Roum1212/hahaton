package aggregate_reward_command

import (
	"time"

	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
)

type Update struct {
	title         string
	description   string
	imagePath     model_image_path.ImagePath
	priceCoins    int64
	quantityLimit int32
	now           time.Time
}

func (x Update) GetDescription() string {
	return x.description
}

func (x Update) GetImagePath() model_image_path.ImagePath {
	return x.imagePath
}

func (x Update) GetNow() time.Time {
	return x.now
}

func (x Update) GetPriceCoins() int64 {
	return x.priceCoins
}

func (x Update) GetQuantityLimit() int32 {
	return x.quantityLimit
}

func (x Update) GetTitle() string {
	return x.title
}

func NewUpdate(
	title string,
	description string,
	imagePath model_image_path.ImagePath,
	priceCoins int64,
	quantityLimit int32,
	now time.Time,
) Update {
	return Update{
		title:         title,
		description:   description,
		imagePath:     imagePath,
		priceCoins:    priceCoins,
		quantityLimit: quantityLimit,
		now:           now,
	}
}
