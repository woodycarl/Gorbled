package gorbled

import (
    "net/http"
    "strconv"

    "appengine"
    "appengine/datastore"
)

func handleIndex(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    // Get user info
    userInfo := getUserInfo(c)

    // Get post id and page id
    pageId, _ := strconv.Atoi(getUrlQuery(r.URL, "pid"))
    pageSize  := config.PageSize

    // Get offset and page numbers
    offset, pageNums := getOffset("Article", pageId, pageSize, c)

    // New PageSetting
    pageSetting := new(PageSetting)

    // Setting PageSetting
    pageSetting.Title       = config.Title
    pageSetting.Description = config.Description
    pageSetting.Layout      = "column2"
    pageSetting.ShowSidebar = true

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

    // New PageData
    pageData := &PageData{ User: userInfo, Article: articleData, Widget: widgetData }

    // New Page
    page := NewPage(pageSetting, pageData)

    // Render page
    page.Render("index", w)
}
