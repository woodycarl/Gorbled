package gorbled

import (
    "net/http"
    "appengine"
    "appengine/datastore"
    "text/template"
)

/*
 * RSS handler
 */
func handleRSS(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    initSystem(r)

    var articles []Article
    _, err := datastore.NewQuery("Article").Order("-Date").GetAll(c, &articles)
    if err != nil {
        serveError(w, err)
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
        serveError(w, err)
        return
    }

    tmpl.Execute(w, page)
}