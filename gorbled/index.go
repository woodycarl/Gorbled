package gorbled

import (
    "net/http"
    "strconv"
    //"fmt"
    "appengine"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    //initSystem(r)

    // Get user info
    user := getUserInfo(c)

    // Get post id and page id
    pageId, _ := strconv.Atoi(getUrlVar(r, "pid"))
    pageSize  := config.Articles

    // Get offset and page nav
    offset, nav := getPageNav("Article", pageId, pageSize, c)

    // Get article data
    articles, err := getArticlesPerPage(offset, pageSize, c)
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
    page := Page {
        "User":       user,
        "Articles" :  articles,
        "Widgets" :   widgets,
        "Nav" :       nav,
        "Config" :    config,
    }

    // Render page
    page.Render("index", w)
}
