package gorbled

import (
    "net/http"
    "time"
    "strconv"

    "appengine"
    "appengine/datastore"
)

func init() {
    http.HandleFunc("/admin/article-list",   handleArticleList)
    http.HandleFunc("/admin/article-add",    handleArticleAdd)
    http.HandleFunc("/admin/article-edit",   handleArticleEdit)
    http.HandleFunc("/admin/article-delete", handleArticleDelete)

    http.HandleFunc("/article", handleArticleView)
}

/*
 * Article data struct
 */
type Article struct {
    ID         string
    Title      string
    Date       time.Time
    Content    []byte
}

func (article *Article) save(c appengine.Context) (err error) {
    _, err = datastore.Put(c, datastore.NewIncompleteKey(c, "Article", nil), article)
    return
}

func (article *Article) update(key *datastore.Key, c appengine.Context) (err error) {
    _, err = datastore.Put(c, key, article)
    return
}

func getArticle(id string, c appengine.Context) (article Article, key *datastore.Key, err error) {
    dbQuery := datastore.NewQuery("Article").Filter("ID =", id)
    var articles []Article
    keys, err := dbQuery.GetAll(c, &articles)
    if len(articles) > 0 {
        article = articles[0]
        key = keys[0]
    }

    return
}

func getArticlesPerPage(offset, pageSize int, c appengine.Context) (articles []Article, err error) {

    dbQuery := datastore.NewQuery("Article").
                                    Order("-Date").
                                    Offset(offset).
                                    Limit(pageSize)
    _, err = dbQuery.GetAll(c, &articles)

    return
}

/*
 * Article handler
 */
func handleArticleList(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    

    // Get page id, pageSize
    pageId, _ := strconv.Atoi(getUrlQuery(r.URL, "pid"))
    pageSize  := config.AdminArticles

    // Get offset and page nav
    offset, nav := getPageNav("Article", pageId, pageSize, c)

    // Get article data
    articles, err := getArticlesPerPage(offset, pageSize, c)
    if err != nil {
        serveError(c, w, err)
        return
    }

    // New Page
    page := Page {
        Title:    "Article Manager",
        Articles: articles,
        Nav:      nav,
        Config:   config,
    }

    // Render page
    page.Render("admin/articles", w)
}

func handleArticleAdd(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    

    if r.Method != "POST" {
        // Show article add page

        // New Page
        page := Page {
            Title:  "Add Article",
            Config: config,
            New:    true,
        }

        // Render page
        err := page.Render("admin/article", w)
        if err != nil {
            serveError(c, w, err)
            return
        }

        return
    }

    // Process post data

    // Parse form data
    if err := r.ParseForm(); err != nil {
        serveError(c, w, err)
        return
    }

    // Create article
    article := &Article{
        Date:    time.Now(),
        Title:   r.FormValue("title"),
        Content: []byte(r.FormValue("content")),
    }

    // Get ID
    if r.FormValue("customid") != "" && !checkIdIsExists("Article", r.FormValue("customid"), c) {
        // Check id is exists
        article.ID = r.FormValue("customid")
    } else {
        article.ID = genID()
    }

    // Save to datastore
    if err := article.save(c); err != nil {
        serveError(c, w, err)
    }

    http.Redirect(w, r, "/admin/article-list", http.StatusFound)
}

func handleArticleEdit(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    

    // Get article id
    id := getUrlQuery(r.URL, "id")

    // Get article data
    article, key, err := getArticle(id, c)
    if err != nil {
        serveError(c, w, err)
        return
    }

    // Check article is exists
    if article.ID == "" {
        serve404(w)
        return
    }

    if r.Method != "POST" {
        // Show article edit page

        // New Page
        page := Page {
            Title:   "Edit Article",
            Article: article,
            Config:  config,
            New:     false,
        }

        // Render page
        err = page.Render("admin/article", w)
        if err != nil {
            serveError(c, w, err)
        }
        return
    }

    // Process post data

    if err := r.ParseForm(); err != nil {
        serveError(c, w, err)
        return
    }

    // Update article
    if r.FormValue("customid") != "" {
        article.ID = r.FormValue("customid")
    }
    article.Title   = r.FormValue("title")
    article.Content = []byte(r.FormValue("content"))

    if err := article.update(key, c); err != nil {
        serveError(c, w, err)
    }

    http.Redirect(w, r, "/admin/article-list", http.StatusFound)
}

func handleArticleDelete(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    

    // Get article id
    id := getUrlQuery(r.URL, "id")

    // Get article data
    _, key, err := getArticle(id, c)
    if err != nil {
      serveError(c, w, err)
      return
    }

    datastore.Delete(c, key)

    http.Redirect(w, r, "/admin/article-list", http.StatusFound)
}

func handleArticleView(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    

    // Get article id
    id := getUrlQuery(r.URL, "id")

    // Get user info
    user := getUserInfo(c)

    // Get article data
    article, _, err := getArticle(id, c)
    if err != nil {
        serveError(c, w, err)
        return
    }
    // Check article is exists
    if article.ID == "" {
        serve404(w)
        return
    }

    // Get widget data
    widgets, err := getWidgets(c)
    if err != nil {
        serveError(c, w, err)
        return
    }

    // New Page
    page := Page {
        User:    user,
        Article: article,
        Widgets: widgets,
        Config:  config,
    }

    // Render page
    page.Render("article", w)
}
