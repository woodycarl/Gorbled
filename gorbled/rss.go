package gorbled

import (
	"appengine"
	"net/http"
	"text/template"
)

/*
 * RSS handler
 */
func handleRSS(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	articles, err := getArticles(c)
	if err != nil {
		serveError(w, err)
		return
	}

	// New Page
	page := Page{
		"Title":    "RSS",
		"Articles": articles,
		"Config":   config,
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
