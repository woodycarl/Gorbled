package gorbled


import (
    "net/http"

    "appengine"
    "appengine/datastore"
    "text/template"
)

func init() {
    http.HandleFunc("/rss.xml", handleRSS)
}

/*
 * Article handler
 */
func handleRSS(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    initConfig(r)

    var articles []Article
    _, err := datastore.NewQuery("Article").Order("-Date").GetAll(c, &articles)
    if err != nil {
        serveError(c, w, err)
        return
    }

    // New Page
    page := Page {
        Title :     "RSS",
        Articles :  articles,
        Config :    config,
    }

    // Render page
    tmpl, err := template.New("rss.html").Funcs(funcMap).ParseFiles(
            "gorbled/admin/rss.html",
        )
    if err != nil {
        return
    }

    tmpl.Execute(w, page)
}