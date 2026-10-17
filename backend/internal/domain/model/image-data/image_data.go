package model_image_data

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"slices"

	_ "image/jpeg"
	_ "image/png"
)

const (
	MaxBytes     = 5 << 20
	MaxDimension = 4096

	formatJPEG = "jpeg"
	formatPNG  = "png"
)

type ImageData struct {
	data   []byte
	format string
}

func (x ImageData) GetData() []byte {
	return slices.Clone(x.data)
}

func (x ImageData) GetExtension() string {
	if x.format == formatJPEG {
		return ".jpg"
	}

	return ".png"
}

func (x ImageData) Validate() error {
	if len(x.data) == 0 || len(x.data) > MaxBytes {
		return errors.New("image size must be between 1 byte and 5 mib")
	}

	config, format, err := image.DecodeConfig(bytes.NewReader(x.data))
	if err != nil {
		return fmt.Errorf("decode image config: %w", err)
	}

	if (format != formatJPEG && format != formatPNG) || format != x.format {
		return errors.New("only jpeg and png images are supported")
	}

	if config.Width > MaxDimension || config.Height > MaxDimension {
		return errors.New("image dimensions must not exceed 4096 pixels")
	}

	if _, _, err = image.Decode(bytes.NewReader(x.data)); err != nil {
		return fmt.Errorf("decode image: %w", err)
	}

	return nil
}

func NewImageData(data []byte) (ImageData, error) {
	if len(data) == 0 || len(data) > MaxBytes {
		return ImageData{}, errors.New("image size must be between 1 byte and 5 mib")
	}

	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return ImageData{}, fmt.Errorf("decode image config: %w", err)
	}

	x := ImageData{
		data:   slices.Clone(data),
		format: format,
	}
	if err = x.Validate(); err != nil {
		return ImageData{}, fmt.Errorf("invalid image data: %w", err)
	}

	return x, nil
}
