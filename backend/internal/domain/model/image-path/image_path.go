package model_image_path

import (
	"errors"
	"fmt"
	"path"
	"strings"

	model_user_id "github.com/Roum1212/hahaton/internal/domain/model/user-id"
	types_uuid_v7 "github.com/Roum1212/hahaton/internal/pkg/types/uuid-v7"
)

type ImagePath struct {
	userID    model_user_id.UserID
	imageID   types_uuid_v7.UUIDv7
	extension string
}

func (x ImagePath) GetURL() string {
	if x.IsZero() {
		return ""
	}

	return "/storage/" + x.String()
}

func (x ImagePath) GetUserID() model_user_id.UserID {
	return x.userID
}

func (x ImagePath) IsZero() bool {
	return x.userID.IsZero() && x.imageID.IsZero() && len(x.extension) == 0
}

func (x ImagePath) String() string {
	if x.IsZero() {
		return ""
	}

	return x.userID.String() + "/" + x.imageID.String() + x.extension
}

func (x ImagePath) Validate() error {
	if x.IsZero() {
		return nil
	}

	if err := x.userID.Validate(); err != nil {
		return fmt.Errorf("image owner id: %w", err)
	}

	if err := x.imageID.Validate(); err != nil {
		return fmt.Errorf("image id: %w", err)
	}

	if x.extension != ".jpg" && x.extension != ".png" {
		return errors.New("invalid image extension")
	}

	return nil
}

func NewImagePath(value string) (ImagePath, error) {
	if len(value) == 0 {
		return NewImagePathFromFields(model_user_id.UserID{}, types_uuid_v7.UUIDv7{}, "")
	}

	owner, fileName, found := strings.Cut(value, "/")
	if !found {
		return ImagePath{}, errors.New("image path must contain owner id and file name")
	}

	userID, err := model_user_id.NewUserID(owner)
	if err != nil {
		return ImagePath{}, fmt.Errorf("parse image owner id: %w", err)
	}

	extension := path.Ext(fileName)

	imageID, err := types_uuid_v7.NewUUIDv7(strings.TrimSuffix(fileName, extension))
	if err != nil {
		return ImagePath{}, fmt.Errorf("parse image id: %w", err)
	}

	return NewImagePathFromFields(userID, imageID, extension)
}

func NewImagePathFromFields(
	userID model_user_id.UserID,
	imageID types_uuid_v7.UUIDv7,
	extension string,
) (ImagePath, error) {
	x := ImagePath{
		userID:    userID,
		imageID:   imageID,
		extension: extension,
	}
	if err := x.Validate(); err != nil {
		return ImagePath{}, fmt.Errorf("invalid image path: %w", err)
	}

	return x, nil
}
