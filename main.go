package shrtn

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"appengine/user"

	"fmt"
	"net/http"
	"strings"
	"time"
)

type Shortcut struct {
	Owner    string
	FullUrl  string
	ShortUrl string
	Create   time.Time
	LastEdit time.Time
}

func init() {
	http.HandleFunc("/_create", Create)
	http.HandleFunc("/", Redirect)
}

const guestbookForm = `
<html>
  <body>
    <form action="/_create" method="post">
      <div><label for="shorturl">Short URL</label> <input id="shorturl" name="shorturl" type="text"/></div>
      <div><label for="fullurl">Full URL</label> <input id="fullurl" name="fullurl" type="text"/></div>
      <div><input type="submit" value="Save"></div>
    </form>
  </body>
</html>`

const ShortcutType = "Shortcut"

func Create(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	if r.Method == "GET" {
		fmt.Fprint(w, guestbookForm)
		return
	}
	if r.Method == "POST" {
		s := Shortcut{}
		if u := user.Current(c); u != nil {
			s.Owner = u.ID
		}
		s.Create = time.Now()
		s.LastEdit = time.Now()
		s.FullUrl = r.FormValue("fullurl")
		s.ShortUrl = r.FormValue("shorturl")
		if len(s.ShortUrl) == 0 {
			http.Error(w, "Short url cannot be empty.", http.StatusBadRequest)
			return
		}
		if len(s.FullUrl) == 0 {
			http.Error(w, "Full url cannot be empty.", http.StatusBadRequest)
		}
		s.ShortUrl = strings.Trim(s.ShortUrl, "/")
		s.ShortUrl = "/" + s.ShortUrl

		key := datastore.NewKey(c, ShortcutType, s.ShortUrl, 0, nil)
		_, err := datastore.Put(c, key, &s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// Change the Value of the item
		item := &memcache.Item{
			Key:   s.ShortUrl,
			Value: []byte(s.FullUrl),
		}
		// Set the item, unconditionally
		if err := memcache.Set(c, item); err != nil {
			c.Errorf("error setting item: %v", err)
		}
		fmt.Fprintf(w, "<html>Dodano: <a href=\"/%s\">%s%s</a>", s.ShortUrl, r.URL.Host, s.ShortUrl)
	}
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	var dest Shortcut
	key := r.URL.Path
	if item, err := memcache.Get(c, key); err == nil {
    http.Redirect(w, r, string(item.Value), http.StatusSeeOther)
    return
  } else if err != memcache.ErrCacheMiss {
		c.Errorf("error getting item: %v", err)
	}

	dkey := datastore.NewKey(c, ShortcutType, key, 0, nil)
	if err := datastore.Get(c, dkey, &dest); err == datastore.ErrNoSuchEntity {
		c.Debugf("not found %q: %s", key, err)
		http.Error(w, "nothing to redirect", http.StatusNotFound)
		return
	} else if err != nil {
		c.Errorf("error: %s", err)
		http.Error(w, "ups...", http.StatusInternalServerError)
		return
	}
	item := memcache.Item{
		Key:   dest.ShortUrl,
		Value: []byte(dest.FullUrl),
	}
	if err := memcache.Set(c, &item); err != nil {
		c.Errorf("error setting item: %v", err)
	}
	http.Redirect(w, r, dest.FullUrl, http.StatusSeeOther)
}
