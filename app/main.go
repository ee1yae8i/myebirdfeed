package app

import (
	"appengine"
	"appengine/datastore"
	"ebird"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"io"
	"io/ioutil"
	"regexp"
)

const (
	feedPrefix = "/user/"
)

func init() {
	http.HandleFunc(feedPrefix, serveFeed)
	http.HandleFunc("/debug", debug)
}

type User struct {
	Login string `datastore:",noindex"`
	ObfuscatedPassword []byte `datastore:",noindex"`
}

func userKey(context appengine.Context, login string) *datastore.Key {
	return datastore.NewKey(context, "User", login, 0, nil)
}

func serveFeed(w http.ResponseWriter, r *http.Request) {
	// Lookup user.
	login, err := parseUser(r.URL.Path)
	if (err != nil) {
		http.NotFound(w, r)
		return
	}

	context := appengine.NewContext(r)
	user := new(User)
	err = datastore.Get(context, userKey(context, login), user)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// TODO: Implement some policy for deciding when to refresh list of observations.
	fetchObservations(context, user.Login, password)

	// TODO: transform checklists into RSS/Atom elements.
}

func parseUser(path string) (user string, err error) {
	user = path[len(feedPrefix):]
	if (!regexp.MustCompile("^[a-zA-Z0-9.-_]+$").MatchString(user)) {
		err = errors.New("Invalid user path")
	}
	return
}

func fetchObservations() {
	session := ebird.NewSession(context, username, password)
	defer session.Logout()

	session.FetchObservationSummaries()
	// TODO: If checklists missing, set a timer and begin fetching them and saving to datastore.	
}

func log(w io.Writer, value interface{}) {
	fmt.Fprintf(w, "%+v\n", value)
}
