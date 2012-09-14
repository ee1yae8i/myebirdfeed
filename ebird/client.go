package ebird

import (
	"net/http"
	"net/url"
	"appengine"
	"appengine/urlfetch"
)

type simpleCookieJar struct {
	cookies map[string][]*http.Cookie
}

func (jar *simpleCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.cookies[u.Host] = cookies
}

func (jar *simpleCookieJar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies[u.Host]
}

func newCookieJar() http.CookieJar {
	return &simpleCookieJar{make(map[string][]*http.Cookie)}
}

func NewClient(context appengine.Context) *http.Client {
	client := urlfetch.Client(context)
	client.Jar = newCookieJar()
	return client
}
