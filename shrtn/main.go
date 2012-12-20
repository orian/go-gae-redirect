package shrtn

import (
  "appengine"
  "fmt"
  "appengine/datastore"
    "appengine/user"
//    "html/template"
    "net/http"
    "time"
)

type Shortcut struct {
  Owner string
  FullUrl string
  ShortUrl string
  Create time.Time
  LastEdit time.Time
}

func init() {
  http.HandleFunc("/", Redirect)
  http.HandleFunc("/_create", Create)
}

const guestbookForm = `
<html>
  <body>
    <form action="/_create" method="post">
      <div><input name="fullurl" type="text"/></div>
      <div><input name="shorturl" type="text"/></div>
      <div><input type="submit" value="Add"></div>
    </form>
  </body>
</html>`

func Create(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  c.Infof("Requested URL: %v", r.URL)
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
    key := datastore.NewIncompleteKey(c, "Shortcut", nil)
    _, err := datastore.Put(c, key, &s)
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
  }
  fmt.Fprint(w, "<html>Dodano: <a href=\"\">http://smth/</a>")
}

func Redirect(w http.ResponseWriter, r *http.Request) {
  c := appengine.NewContext(r)
  c.Infof("Requested URL: %v", r.URL)
  fmt.Fprintf(w, "<h1>Hello, world</h1>")
}
