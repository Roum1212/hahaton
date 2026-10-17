package cookies

import (
	"net/http"
)

func NewCookie(name string, opts ...Option) *http.Cookie {
	x := http.Cookie{
		Name: name,
		Path: "/",
	}
	for i := range opts {
		opts[i](&x)
	}

	return &x
}
