package gorbled

import (
    "net/http"
    "time"
    "strconv"
    "appengine"
    "appengine/datastore"
)

func getArticles(c appengine.Context) (articles []Entry, err error) {

    dbQuery := datastore.NewQuery("Entry").
        Filter("Type =", "article")
    _, err = dbQuery.GetAll(c, &articles)

    return
}

func getArticlesAndNav(paginaId, paginaSize int, c appengine.Context) (articles []Entry, nav PaginaNav, err error) {

    // Get offset and pagina nav
    dbQuery := datastore.NewQuery("Entry").Filter("Type =", "article")
    count, _ := dbQuery.Count(c)
    offset, nav := getPaginaNav(count, paginaId, paginaSize, c)

    // Get article data
    dbQuery = dbQuery.Order("-Date").Offset(offset).Limit(paginaSize)
    _, err = dbQuery.GetAll(c, &articles)

    return
}

/*
 * Article handler
 */
func handleArticleList(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    // Get pagina id, paginaSize
    paginaId, _ := strconv.Atoi(getUrlVar(r, "pid"))
    paginaSize  := config.AdminArticles

    articles, nav, err := getArticlesAndNav(paginaId, paginaSize, c)
    if err != nil {
        serveError(w, err)
        return
    }

    pagina := Pagina {
        "Title":    "Article Manager",
        "Articles": articles,
        "Nav":      nav,
        "Config":   config,
    }

    pagina.Render("admin/articles", w)
}

func handleArticleAdd(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    if r.Method != "POST" {
        pagina := Pagina {
            "Title":  "Add Article",
            "Config": config,
            "New":    true,
        }

        pagina.Render("admin/article", w)

        return
    }

    // Process post data

    if err := r.ParseForm(); err != nil {
        serveError(w, err)
        return
    }

    // Create article
    article := &Entry{
        Date:    time.Now(),
        Title:   r.FormValue("title"),
        Content: []byte(r.FormValue("content")),
        Type:   "article",
    }

    if r.FormValue("customid") != "" && !checkIdIsExists("Entry", r.FormValue("customid"), c) {
        article.ID = r.FormValue("customid")

        // Save to datastore
        if err := article.put(c); err != nil {
            serveError(w, err)
            return
        }
        
    } else if err := article.save(c); err != nil {
        serveError(w, err)
        return
    }

    http.Redirect(w, r, "/admin/article", http.StatusFound)
}

func handleArticleEdit(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    // Get article data
    article, key, err := getEntry(getUrlVar(r, "id"), c)
    if err != nil {
        serveError(w, err)
        return
    }

    // Check article is exists
    if article.ID == "" {
        serve404(w)
        return
    }

    if r.Method != "POST" {
        // New Pagina
        pagina := Pagina {
            "Title":   "Edit Article",
            "Article": article,
            "Config":  config,
            "New":     false,
        }

        // Show article edit pagina
        pagina.Render("admin/article", w)

        return
    }

    // Process post data

    if err := r.ParseForm(); err != nil {
        serveError(w, err)
        return
    }

    // Update article
    if r.FormValue("customid") != "" {
        article.ID = r.FormValue("customid")
    }
    article.Title   = r.FormValue("title")
    article.Content = []byte(r.FormValue("content"))

    if err := article.update(key, c); err != nil {
        serveError(w, err)
        return
    }

    http.Redirect(w, r, "/admin/article", http.StatusFound)
}

func handleRedirectArticleList(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/admin/article", http.StatusFound)
}
