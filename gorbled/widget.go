package gorbled

import (
    "net/http"
    "strconv"
    "other_package/blackfriday"

    "appengine"
    "appengine/datastore"
)

type WidgetDB struct {
    ID, Title   string
    Sequence    int
    Content     []byte
}

type WidgetData struct {
    ID, Title, Content string
    Sequence    int
}

/*
 * Get widget data
 */
func getWidgetData(dbQuery *datastore.Query, MDOutput bool, c appengine.Context) (widgetData []WidgetData , err error) {
    var widgetDB []*WidgetDB
    _, err = dbQuery.GetAll(c, &widgetDB)
    if err != nil {
        return
    }

    widgetData = make([]WidgetData, len(widgetDB))
    for i := 0; i < len(widgetDB); i++ {
        widgetData[i].ID       = widgetDB[i].ID
        widgetData[i].Title    = widgetDB[i].Title
        widgetData[i].Sequence = widgetDB[i].Sequence
        if MDOutput {
            widgetData[i].Content  = string(blackfriday.MarkdownCommon(widgetDB[i].Content))
        } else {
            widgetData[i].Content  = string(widgetDB[i].Content)
        }
    }

    return
}

func handleWidget(w http.ResponseWriter, r *http.Request) {
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
            widgetManager(w, r, operation, id)
    }
}

func widgetManager(w http.ResponseWriter, r *http.Request, operation string, id string) {
    switch operation {
        case "add":
            widgetAdd(w, r)
        case "edit":
            widgetEdit(w, r, id)
        case "delete":
            widgetDelete(w, r, id)

        default :
            widgetList(w, r)
    }
}

func widgetList(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    // Get widget data

    // Get page id
    pageId, _ := strconv.Atoi(getUrlQuery(r.URL, "pid"))
    pageSize  := 10

    // Get offset and page numbers
    offset, pageNums := getOffset("Widget", pageId, pageSize, c)

    // New PageSetting
    pageSetting := new(PageSetting)

    // Setting PageSetting
    pageSetting.Title  = "Widget Manager - " + config.Title
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

    // Get widget data
    dbQuery         := datastore.NewQuery("Widget").Order("Sequence").Offset(offset).Limit(pageSize)
    widgetData, err := getWidgetData(dbQuery, false, c)
    if err != nil {
        serveError(c, w, err)
        return
    }

    // New PageData
    pageData := &PageData{ Widget: widgetData }

    // New Page
    page := NewPage(pageSetting, pageData)

    // Render page
    page.Render("widget/manager", w)
}

func widgetAdd(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    if r.Method != "POST" {
        // Show widget add page

        // New PageSetting
        pageSetting := new(PageSetting)
        pageSetting.Title  = "Widget Manager - Add - " + config.Title
        pageSetting.Layout = "column1"

        // New Page
        page := NewPage(pageSetting, nil)

        // Render page
        err := page.Render("widget/add", w)
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

    // Sequence
    var sequence int
    sequence, _ = strconv.Atoi(r.FormValue("sequence"))

    // Create widgetDB
    widgetDB := &WidgetDB{
        ID:       genID(),
        Title:    r.FormValue("title"),
        Sequence: sequence,
        Content:  []byte(r.FormValue("content")),
    }

    // Save to datastore
    if _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Widget", nil), widgetDB); err != nil {
        serveError(c, w, err)
    } else {
        http.Redirect(w, r, "/widget?action=manager", http.StatusFound)
    }
}

func widgetEdit(w http.ResponseWriter, r *http.Request, id string) {
    c := appengine.NewContext(r)

    if r.Method != "POST" {
        // Show widget edit page

        dbQuery := datastore.NewQuery("Widget").Filter("ID =", id)

        // Check error
        if count, _ := dbQuery.Count(c); count < 1 {
            http.Redirect(w, r, "/widget?action=manager", http.StatusFound)
        }

        // Get widget data
        widgetData, err := getWidgetData(dbQuery, false, c)
        if err != nil {
            serveError(c, w, err)
            return
        }

        // New PageSetting
        pageSetting := new(PageSetting)
        pageSetting.Title  = "Widget Manager - Edit - " + config.Title
        pageSetting.Layout = "column1"

        // New PageData
        pageData := &PageData{ Widget: widgetData }

        // New Page
        page := NewPage(pageSetting, pageData)

        // Render page
        err = page.Render("widget/edit", w)
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

    // Delete old widget
    dbQuery := datastore.NewQuery("Widget").Filter("ID =", id)
    var widgetDBTmp []*WidgetDB
    keys, err := dbQuery.GetAll(c, &widgetDBTmp)
    if err != nil {
        serveError(c, w, err)
        return
    }
    datastore.DeleteMulti(c, keys)

    // Add new widget

    // Sequence
    var sequence int
    sequence, _ = strconv.Atoi(r.FormValue("sequence"))

    // Create widgetDB
    widgetDB := &WidgetDB{
        ID:       widgetDBTmp[0].ID,
        Title:    r.FormValue("title"),
        Sequence: sequence,
        Content:  []byte(r.FormValue("content")),
    }

    // Save to datastore
    if _, err := datastore.Put(c, datastore.NewIncompleteKey(c, "Widget", nil), widgetDB); err != nil {
        serveError(c, w, err)
    } else {
        http.Redirect(w, r, "/widget?action=manager", http.StatusFound)
    }
}

func widgetDelete(w http.ResponseWriter, r *http.Request, id string) {
    c := appengine.NewContext(r)

    // Get delete widget
    dbQuery := datastore.NewQuery("Widget").Filter("ID =", id)

    // Check error
    if count, _ := dbQuery.Count(c); count < 1 {
        http.Redirect(w, r, "/", http.StatusFound)
    }

    // Delete widget
    var widgetDB []*WidgetDB
    keys, _ := dbQuery.GetAll(c, &widgetDB)
    datastore.DeleteMulti(c, keys)

    http.Redirect(w, r, "/widget?action=manager", http.StatusFound)
}
