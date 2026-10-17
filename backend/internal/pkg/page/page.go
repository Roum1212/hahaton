package page

import (
	model_page "github.com/Roum1212/hahaton/internal/domain/model/page"
)

func Slice[T any](values []T, paging model_page.Page) ([]T, int64) {
	total := int64(len(values))
	start := min(paging.GetOffset(), total)
	end := min(start+int64(paging.GetLimit()), total)

	return values[start:end], total
}
