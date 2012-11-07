package gorbled

import (
	"appengine"
	"net/http"
	"strconv"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Get post id and pagina id
	paginaId, _ := strconv.Atoi(getUrlVar(r, "pid"))

	articles, nav, err := getArticlesAndNav(paginaId, c)
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
	pagina := Pagina{
		"User":     getUserInfo(c),
		"Articles": articles,
		"Widgets":  widgets,
		"Nav":      nav,
		"Config":   config,
	}

	// Render pagina
	pagina.Render("index", w)
}
