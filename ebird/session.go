package ebird

import (
	"net/http"
	"net/url"
	"appengine"
)

type Session interface {
	FetchObservationSummaries() []ObservationSummary
	Logout()
}

type session struct {
	Client *http.Client
}

func NewSession(c appengine.Context, username, password string) Session {
	client := newClient(c)
	login(client, username, password)
	return &session{client}
}

func login(client *http.Client, username, password string) {
	form := url.Values{}
	form.Set("j_username", username)
	form.Set("j_password", password)
	form.Set("cmd", "login")
	resp, err := client.PostForm("https://ebird.org/ebird/j_acegi_security_check", form)
	if (err != nil) {
		panic(err.Error())
	}
}

func (s *session) FetchObservationSummaries() []ObservationSummary {
	_, err = client.Get("http://ebird.org/ebird/eBirdReports?cmd=subReport")
	if (err != nil) {
		panic(err.Error())
	}

//	data, err := ioutil.ReadAll(resp.Body)

	return []ObservationSummary{}
}

func (s *session) Logout() {
	s.client.Get("http://ebird.org/ebird/j_acegi_logout")
}
