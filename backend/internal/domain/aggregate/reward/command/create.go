package aggregate_reward_command

import (
	"time"

	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
)

type Create struct {
	title         string
	description   string
	imagePath     model_image_path.ImagePath
	priceCoins    int64
	quantityLimit int32
	now           time.Time
}

func (x Create) GetDescription() string {
	return x.description
}

func (x Create) GetImagePath() model_image_path.ImagePath {
	return x.imagePath
}

func (x Create) GetNow() time.Time {
	return x.now
}

func (x Create) GetPriceCoins() int64 {
	return x.priceCoins
}

func (x Create) GetQuantityLimit() int32 {
	return x.quantityLimit
}

func (x Create) GetTitle() string {
	return x.title
}

func NewCreate(
	title string,
	description string,
	imagePath model_image_path.ImagePath,
	priceCoins int64,
	quantityLimit int32,
	now time.Time,
) Create {
	return Create{
		title:         title,
		description:   description,
		imagePath:     imagePath,
		priceCoins:    priceCoins,
		quantityLimit: quantityLimit,
		now:           now,
	}
}
