package cookies

import (
	"net/http"
	"time"
)

type Option func(x *http.Cookie)

func WithDomain(domain string) Option {
	return func(x *http.Cookie) {
		x.Domain = domain
	}
}

func WithExpires(t time.Time) Option {
	return func(x *http.Cookie) {
		x.Expires = t
		x.MaxAge = int(time.Until(t).Truncate(time.Second).Seconds())
	}
}

func WithHTTPOnly() Option {
	return func(x *http.Cookie) {
		x.HttpOnly = true
	}
}

func WithMaxAge(d time.Duration) Option {
	return func(x *http.Cookie) {
		d = d.Truncate(time.Second)

		x.Expires = time.Now().Add(d)
		x.MaxAge = int(d.Seconds())
	}
}

func WithPath(path string) Option {
	return func(x *http.Cookie) {
		x.Path = path
	}
}

func WithSameSite(sameSite http.SameSite) Option {
	return func(x *http.Cookie) {
		x.SameSite = sameSite
	}
}

func WithSecure() Option {
	return func(x *http.Cookie) {
		x.Secure = true
	}
}
