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

	// New Pagina
	pagina := Pagina{
		"Title":    "RSS",
		"Articles": articles,
		"Config":   config,
	}

	// Render pagina
	tmpl, err := template.New("rss.html").Funcs(funcMap).ParseFiles(
		"gorbled/admin/rss.html",
	)
	if err != nil {
		serveError(w, err)
		return
	}

	tmpl.Execute(w, pagina)
}
