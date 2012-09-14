package app

import (
	"appengine"
	"appengine/datastore"
	"ebird"
	"errors"
//	"exp/html"
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
}

type User struct {
	Login string `datastore:",noindex"`
	Password string `datastore:",noindex"`
}

func userKey(context appengine.Context, login string) *datastore.Key {
	return datastore.NewKey(context, "User", login, 0, nil)
}

func serveFeed(w http.ResponseWriter, r *http.Request) {
	login, err := parseUser(r.URL.Path)
	if (err != nil) {
		http.NotFound(w, r)
		return
	}

	context := appengine.NewContext(r)

	user := new(User)
	err = datastore.Get(context, userKey(context, login), user)
	if err != nil {
		log(w, err)
		return
	}

	log(w, user.Login)

	client := ebird.NewClient(context)

	form := url.Values{}
	form.Set("j_username", user.Login)
	form.Set("j_password", user.Password)
	form.Set("cmd", "login")
	resp, err := client.PostForm("https://ebird.org/ebird/j_acegi_security_check", form)
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err = client.Get("http://ebird.org/ebird/eBirdReports?cmd=subReport")
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := ioutil.ReadAll(resp.Body)
	log(w, string(data))
//	node, err := html.Parse(resp.Body)
//	log(w, node)

	_, err = client.Get("http://ebird.org/ebird/j_acegi_logout")
	if (err != nil) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func parseUser(path string) (user string, err error) {
	user = path[len(feedPrefix):]
	if (!regexp.MustCompile("^[a-zA-Z0-9.-_]+$").MatchString(user)) {
		err = errors.New("Invalid user path")
	}
	return
}

func log(w io.Writer, value interface{}) {
	fmt.Fprintf(w, "%+v\n", value)
}
