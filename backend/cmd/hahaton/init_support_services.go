package hahaton

import (
	"net/http"
	"time"

	"github.com/Roum1212/hahaton/internal/config"
	"github.com/Roum1212/hahaton/internal/pkg/cookies"
	session_cookie "github.com/Roum1212/hahaton/internal/pkg/session/cookie"
	"github.com/Roum1212/hahaton/internal/pkg/updates"
)

type supportServiceSet struct {
	updates        *updates.Hub
	cookies        cookies.Service
	expiredCookies cookies.Service
}

func initSupportServices(cfg config.Config) supportServiceSet {
	cookieService := cookies.NewService(
		session_cookie.Name,
		cookies.WithHTTPOnly(),
		cookies.WithSecure(),
		cookies.WithPath("/"),
		cookies.WithSameSite(http.SameSiteNoneMode),
		cookies.WithMaxAge(cfg.SessionToken.TTL),
	)
	expiredCookies := cookies.NewService(
		session_cookie.Name,
		cookies.WithHTTPOnly(),
		cookies.WithSecure(),
		cookies.WithPath("/"),
		cookies.WithSameSite(http.SameSiteNoneMode),
		cookies.WithMaxAge(-time.Second),
	)

	return supportServiceSet{
		updates:        updates.NewHub(),
		cookies:        cookieService,
		expiredCookies: expiredCookies,
	}
}
