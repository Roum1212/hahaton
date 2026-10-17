package model_page

import (
	"errors"
	"fmt"
)

const (
	DefaultLimit = 20
	MaxLimit     = 100
)

type Page struct {
	limit    int32
	offset   int64
	maxLimit int32
}

func (x Page) GetLimit() int32 {
	return x.limit
}

func (x Page) GetOffset() int64 {
	return x.offset
}

func (x Page) Validate() error {
	if x.maxLimit < 1 {
		return errors.New("maximum limit must be positive")
	}

	if x.limit < 1 || x.limit > x.maxLimit {
		return fmt.Errorf("limit must be between 1 and %d", x.maxLimit)
	}

	if x.offset < 0 {
		return errors.New("offset must not be negative")
	}

	return nil
}

func NewPage(limit int32, offset int64) (Page, error) {
	return NewPageWithLimits(limit, offset, DefaultLimit, MaxLimit)
}

func NewPageWithLimits(limit int32, offset int64, defaultLimit, maxLimit int32) (Page, error) {
	if limit == 0 {
		limit = defaultLimit
	}

	x := Page{
		limit:    limit,
		offset:   offset,
		maxLimit: maxLimit,
	}
	if err := x.Validate(); err != nil {
		return Page{}, err
	}

	return x, nil
}
