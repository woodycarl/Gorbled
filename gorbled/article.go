package gorbled

import (
    "net/http"
    "time"
    "strconv"
    "html"
    "other_package/blackfriday"

    "appengine"
    "appengine/datastore"
)

type ArticleDB struct {
    ID, Title string
    Date      time.Time
    Content   []byte
}

type ArticleData struct {
    ID, Title, Content string
    Date time.Time
}

/*
 * Get article data
 */
func getArticleData(dbQuery *datastore.Query, MDOutput bool, c appengine.Context) (articleData []ArticleData , err error) {
    var articleDB []*ArticleDB
    _, err = dbQuery.GetAll(c, &articleDB)
    if err != nil {
        return
    }

    articleData = make([]ArticleData, len(articleDB))
    for i := 0; i < len(articleDB); i++ {
        articleData[i].ID    = articleDB[i].ID
        articleData[i].Title = articleDB[i].Title
        if MDOutput {
            articleData[i].Content = string(blackfriday.MarkdownCommon(articleDB[i].Content))
        } else {
            articleData[i].Content = string(articleDB[i].Content)
        }
        articleData[i].Date = articleDB[i].Date
    }

    return
}

func handleArticle(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    // Get action && id
    action := getUrlQuery(r.URL, "action")
    id     := getUrlQuery(r.URL, "id")

    switch action {
        case "manager":
            // Check user permissions
            userInfo := getUserInfo(c)
            if userInfo == nil || !userInfo.IsAdmin {
                serve404(w)
                return
            }

            operation := getUrlQuery(r.URL, "operation")
            articleManager(w, r, operation, id)

        default :
            articleView(w, r, id)
    }
}

func articleView(w http.ResponseWriter, r *http.Request, id string) {
    c := appengine.NewContext(r)

    // Get user info
    userInfo := getUserInfo(c)

    // Create get article data query
    dbQuery := datastore.NewQuery("Article").Filter("ID =", id)

    // Check article is exists
    if count, _ := dbQuery.Count(c); count < 1 {
        serve404(w)
        return
    }

    // Get article data
    articleData, err := getArticleData(dbQuery, true, c)
    if err != nil {
        serveError(c, w, err)
        return
    }

    // Get widget data
    dbQuery          = datastore.NewQuery("Widget").Order("Sequence")
    widgetData, err := getWidgetData(dbQuery, true, c)
    if err != nil {
        serveError(c, w, err)
        return
    }

    // New PageSetting
    pageSetting := new(PageSetting)

    // Setting PageSetting
    for _, v := range articleData {
        // title
        pageSetting.Title = v.Title + " - " + config.Title

        // description
        if len(v.Content) < 100 {
            pageSetting.Description = html.EscapeString(v.Content)
        } else {
            pageSetting.Description = html.EscapeString(v.Content[:100])
        }
    }
    pageSetting.Layout      = "column2"
    pageSetting.ShowSidebar = true

    // New PageData
    pageData := &PageData{ User: userInfo, Article: articleData, Widget: widgetData }

    // New Page
    page := NewPage(pageSetting, pageData)

    // Render page
    page.Render("article/view", w)
}

func articleManager(w http.ResponseWriter, r *http.Request, operation string, id string) {
    switch operation {
        case "add":
            articleAdd(w, r)
        case "edit":
            articleEdit(w, r, id)
        case "delete":
            articleDelete(w, r, id)

        default :
            articleList(w, r)
    }
}

func articleList(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    // Get article data

    // Get page id
    pageId, _ := strconv.Atoi(getUrlQuery(r.URL, "pid"))
    pageSize  := 10

    // Get offset and page numbers
    offset, pageNums := getOffset("Article", pageId, pageSize, c)

    // New PageSetting
    pageSetting := new(PageSetting)

    // Setting PageSetting
    pageSetting.Title  = "Article Manager - " + config.Title
    pageSetting.Layout = "column1"

    // showNext and showPrev button
    if pageId <= 0 || pageId > pageNums {
        pageId = 1
    }
    if pageId < pageNums {
        pageSetting.ShowPrev = true
    }
    if pageId != 1 {
        pageSetting.ShowNext = true
    }
    pageSetting.PrevPageID = pageId + 1
    pageSetting.NextPageID = pageId - 1

    // Get article data
    dbQuery          := datastore.NewQuery("Article").Order("-Date").Offset(offset).Limit(pageSize)
    articleData, err := getArticleData(dbQuery, false, c)
    if err != nil {
        serveError(c, w, err)
        return
    }

    // New PageData
    pageData := &PageData{ Article: articleData }

    // New Page
    page := NewPage(pageSetting, pageData)

    // Render page
    page.Render("article/manager", w)
}

func articleAdd(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    if r.Method != "POST" {
        // Show article add page

        // New PageSetting
        pageSetting := new(PageSetting)
        pageSetting.Title  = "Article Manager - Add - " + config.Title
        pageSetting.Layout = "column1"

        // New Page
        page := NewPage(pageSetting, nil)

        // Render page
        err := page.Render("article/add", w)
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

    // Custom ID
    var customID string
    if r.FormValue("customid") != "" {
        // Check custom id is exists
        if checkIdIsExists("Article", r.FormValue("customid"), c) {
            customID = genID()
        } else {
            customID = r.FormValue("customid")
        }
    } else {
        customID = genID()
    }

    // Create articleDB
    articleDB := &ArticleDB{
        ID:      customID,
        Date:    time.Now(),
        Title:   r.FormValue("title"),
        Content: []byte(r.FormValue("content")),
    }

    // Save to datastore
    if _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Article", nil), articleDB); err != nil {
        serveError(c, w, err)
    } else {
        http.Redirect(w, r, "/article?action=manager", http.StatusFound)
    }
}

func articleEdit(w http.ResponseWriter, r *http.Request, id string) {
    c := appengine.NewContext(r)

    if r.Method != "POST" {
        // Show article edit page

        dbQuery := datastore.NewQuery("Article").Filter("ID =", id)

        // Check error
        if count, _ := dbQuery.Count(c); count < 1 {
            http.Redirect(w, r, "/article?action=manager", http.StatusFound)
        }

        // Get article data
        articleData, err := getArticleData(dbQuery, false, c)
        if err != nil {
            serveError(c, w, err)
            return
        }

        // New PageSetting
        pageSetting := new(PageSetting)
        pageSetting.Title  = "Article Manager - Edit - " + config.Title
        pageSetting.Layout = "column1"

        // New PageData
        pageData := &PageData{ Article: articleData }

        // New Page
        page := NewPage(pageSetting, pageData)

        // Render page
        err = page.Render("article/edit", w)
        if err != nil {
            serveError(c, w, err)
        }

        return
    }

    // Process post data

    // Parse form data
    if err := r.ParseForm(); err != nil {
        serveError(c, w, err)
        return
    }

    // Delete old article
    dbQuery := datastore.NewQuery("Article").Filter("ID =", id)
    var articleDBTmp []*ArticleDB
    keys, err := dbQuery.GetAll(c, &articleDBTmp)
    if err != nil {
        serveError(c, w, err)
        return
    }
    datastore.DeleteMulti(c, keys)

    // Add new article

    // Custom ID
    var customID string
    if r.FormValue("customid") != "" {
        // Check custom id is exists
        if checkIdIsExists("Article", r.FormValue("customid"), c) {
            customID = articleDBTmp[0].ID
        } else {
            customID = r.FormValue("customid")
        }
    } else {
        customID = articleDBTmp[0].ID
    }

    // Create articleDB
    articleDB := &ArticleDB{
        ID:      customID,
        Date:    time.Now(),
        Title:   r.FormValue("title"),
        Content: []byte(r.FormValue("content")),
    }

    // Save to datastore
    if _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Article", nil), articleDB); err != nil {
        serveError(c, w, err)
    } else {
        http.Redirect(w, r, "/article?action=manager", http.StatusFound)
    }
}

func articleDelete(w http.ResponseWriter, r *http.Request, id string) {
    c := appengine.NewContext(r)

    // Get delete article
    dbQuery := datastore.NewQuery("Article").Filter("ID =", id)

    // Check error
    if count, _ := dbQuery.Count(c); count < 1 {
        http.Redirect(w, r, "/", http.StatusFound)
    }

    // Delete article
    var articleDB []*ArticleDB
    keys, _ := dbQuery.GetAll(c, &articleDB)
    datastore.DeleteMulti(c, keys)

    http.Redirect(w, r, "/article?action=manager", http.StatusFound)
}
