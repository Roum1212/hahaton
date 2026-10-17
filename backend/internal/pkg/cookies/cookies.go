package cookies

import (
	"net/http"
	"strings"
)

type Service interface {
	DeleteCookie(r *http.Request, w http.ResponseWriter)
	GetCookie(r *http.Request) (*http.Cookie, error)
	SetCookie(h http.Header, v string)
}

type service struct {
	name string
	opts []Option
}

func (x service) DeleteCookie(r *http.Request, w http.ResponseWriter) {
	cookies := r.CookiesNamed(x.name)
	for _, cookie := range cookies {
		cookie.MaxAge = -1
		cookie.Value = ""

		http.SetCookie(w, cookie)
	}
}

func (x service) GetCookie(r *http.Request) (*http.Cookie, error) {
	switch cookies := r.CookiesNamed(x.name); len(cookies) {
	case 0:
		return nil, http.ErrNoCookie
	case 1:
		return cookies[0], nil
	default:
		var v strings.Builder
		for _, cookie := range cookies {
			v.WriteString(cookie.Value)
		}

		cookie := cookies[0]
		cookie.Partitioned = false
		cookie.Value = v.String()

		return cookie, nil
	}
}

func (x service) SetCookie(h http.Header, v string) {
	cookie := NewCookie(x.name, x.opts...)
	cookie.Value = v

	h.Add("Set-Cookie", cookie.String())
}

func NewService(name string, opts ...Option) Service {
	return service{
		name: name,
		opts: opts,
	}
}
