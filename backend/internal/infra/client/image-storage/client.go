package client_image_storage

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	model_image_data "github.com/Roum1212/hahaton/internal/domain/model/image-data"
	model_image_path "github.com/Roum1212/hahaton/internal/domain/model/image-path"
)

const (
	directoryMode = 0o750
	fileMode      = 0o600
)

type Client struct {
	directory string
}

func (x Client) Save(ctx context.Context, path model_image_path.ImagePath, data model_image_data.ImageData) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	fileName := filepath.Join(x.directory, filepath.FromSlash(path.String()))
	if err := os.MkdirAll(filepath.Dir(fileName), directoryMode); err != nil {
		return fmt.Errorf("create image directory: %w", err)
	}

	if err := os.WriteFile(fileName, data.GetData(), fileMode); err != nil {
		return fmt.Errorf("save image: %w", err)
	}

	return nil
}

func NewClient(directory string) Client {
	return Client{directory: directory}
}
