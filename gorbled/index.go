package gorbled

import (
	"appengine"
	"net/http"
	"strconv"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Get post id and page id
	pageId, _ := strconv.Atoi(getUrlVar(r, "pid"))
	pageSize := config.Articles

	articles, nav, err := getArticlesAndNav(pageId, pageSize, c)
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

	// New Page
	page := Page{
		"User":     getUserInfo(c),
		"Articles": articles,
		"Widgets":  widgets,
		"Nav":      nav,
		"Config":   config,
	}

	// Render page
	page.Render("index", w)
}
