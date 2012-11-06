package gorbled

import (
    "net/http"
    "time"
    "fmt"
    "appengine"
    "appengine/datastore"
)

/*
 * Entry data struct
 */
type Entry struct {
    ID         string
    Title      string
    Date       time.Time
    Content    []byte
    
    Type       string    //article:文章；page:页面；widget: 小工具
    SubPage      []string
    PageClass   int
    Sequence   int

    AllowComment bool
    Slug        string
}

func (entry *Entry) save(c appengine.Context) (err error) {
    config.EntryID = config.EntryID + 1
    config.update(configKey, c)
    entry.ID = fmt.Sprint(config.EntryID)
    entry.put(c)
    return
}

func (entry *Entry) put(c appengine.Context) (err error) {
    _, err = datastore.Put(c, datastore.NewIncompleteKey(c, "Entry", nil), entry)
    return
}

func (entry *Entry) update(key *datastore.Key, c appengine.Context) (err error) {
    _, err = datastore.Put(c, key, entry)
    return
}

func getEntry(id string, c appengine.Context) (entry Entry, key *datastore.Key, err error) {

    dbQuery := datastore.NewQuery("Entry").Filter("ID =", id)
    var entries []Entry
    keys, err := dbQuery.GetAll(c, &entries)
    if len(entries) > 0 {
        entry = entries[0]
        key = keys[0]
    }

    return
}


/*
 * Entry handler
 */

/*
 * Decode markdown code
 *
 * @return (string) 
 */
func handleDecodeContent(w http.ResponseWriter, r *http.Request) {
    content := []byte(r.FormValue("content"))
    fmt.Fprint(w, decodeMD(content))
}

func handleEntryView(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    // Get entry data
    entry, _, err := getEntry(getUrlVar(r, "id"), c)
    if err != nil {
        serveError(w, err)
        return
    }

    // Get widget data
    widgets, err := getWidgets(c)
    if err != nil {
        serveError(w, err)
        return
    }

    // New Pagina
    pagina := Pagina {
        "User":    getUserInfo(c),
        "Entry": entry,
        "Widgets": widgets,
        "Config":  config,
    }

    // Render pagina
    pagina.Render("article", w)
}

func handleEntryDelete(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    // Get entry key
    _, key, err := getEntry(getUrlVar(r, "id"), c)
    if err != nil {
        serveError(w, err)
        return
    }

    // Delete entry
    if err = datastore.Delete(c, key); err != nil {
        serveError(w, err)
        return
    }

    http.Redirect(w, r, r.Referer(), http.StatusFound)
}