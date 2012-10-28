package gorbled

import (
    "net/http"
    "strconv"

    "appengine"
    "appengine/datastore"
)

func init() {
    http.HandleFunc("/admin/widget-list",   handleWidgetList)
    http.HandleFunc("/admin/widget-add",    handleWidgetAdd)
    http.HandleFunc("/admin/widget-edit",   handleWidgetEdit)
    http.HandleFunc("/admin/widget-delete", handleWidgetDelete)
}

/*
 * Widget data struct
 */
type Widget struct {
    ID          string
    Title       string
    Sequence    int
    Content     []byte
}

func (p *Widget) save(c appengine.Context) (err error) {
    _, err = datastore.Put(c, datastore.NewIncompleteKey(c, "Widget", nil), p)
    return
}

func (p *Widget) update(key *datastore.Key, c appengine.Context) (err error) {
    _, err = datastore.Put(c, key, p)
    return
}

/*
 * Get widget data
 */
func getWidgets(c appengine.Context) (widgets []Widget, err error) {
    dbQuery := datastore.NewQuery("Widget").Order("Sequence")
    _, err = dbQuery.GetAll(c, &widgets)
    return
}

func getWidget(id string,
    c appengine.Context) (widget Widget, key *datastore.Key, err error) {

    dbQuery := datastore.NewQuery("Widget").Filter("ID =", id)
    var widgets []Widget
    keys, err := dbQuery.GetAll(c, &widgets)
    if len(widgets) > 0 {
        widget = widgets[0]
        key = keys[0]
    }

    return
}

func getWidgetsPerPage(offset, pageSize int,
        c appengine.Context) (widgets []Widget, err error) {

    dbQuery := datastore.NewQuery("Widget").
        Order("-Sequence").
        Offset(offset).
        Limit(pageSize)

    _, err = dbQuery.GetAll(c, &widgets)

    return
}

/*
 * Widget handler
 */
func handleWidgetList(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    

    // Get page id, pageSize
    pageId, _ := strconv.Atoi(getUrlQuery(r.URL, "pid"))
    pageSize  := config.AdminWidgets

    // Get offset and page nav
    offset, nav := getPageNav("Widget", pageId, pageSize, c)

    // Get widget data
    widgets, err := getWidgetsPerPage(offset, pageSize, c)
    if err != nil {
        serveError(c, w, err)
        return
    }

    // New Page
    page := Page {
        Title :     "Widget Manager",
        Widgets :   widgets,
        Nav :       nav,
        Config :    config,
    }

    // Render page
    page.Render("admin/widgets", w)
}

func handleWidgetAdd(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    

    if r.Method != "POST" {
        // Show widget add page

        // New Page
        page := Page {
            Title:     "Add Widget",
            Config:    config,
            New:       true,
        }

        // Render page
        err := page.Render("admin/widget", w)
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

    // Create widget
    sequence, _ := strconv.Atoi(r.FormValue("sequence"))
    widget := &Widget{
        ID:         getID("Widget", r.FormValue("customid"), c),
        Title:      r.FormValue("title"),
        Content:    []byte(r.FormValue("content")),
        Sequence:   sequence,
    }

    // Save to datastore
    if err := widget.save(c); err != nil {
        serveError(c, w, err)
    }

    http.Redirect(w, r, "/admin/widget-list", http.StatusFound)
}

func handleWidgetEdit(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    

    // Get widget id
    id := getUrlQuery(r.URL, "id")

    // Get widget data
    widget, key, err := getWidget(id, c)
    if err != nil {
        serveError(c, w, err)
        return
    }

    // Check widget is exists
    if widget.ID == "" {
        serve404(w)
        return
    }

    if r.Method != "POST" {
        // Show widget edit page

        dbQuery := datastore.NewQuery("Widget").Filter("ID =", id)

        // Check error
        if count, _ := dbQuery.Count(c); count < 1 {
            http.Redirect(w, r, "/admin/widget-list", http.StatusFound)
        }

        // New Page
        page := Page {
            Title :     "Edit Widget",
            Config :    config,
            Widget:     widget,
            New:        false,
        }

        // Render page
        err = page.Render("admin/widget", w)
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

    // Add new widget

    // Sequence
    sequence, _ := strconv.Atoi(r.FormValue("sequence"))

    widget.Sequence = sequence
    widget.Title = r.FormValue("title")
    widget.Content = []byte(r.FormValue("content"))

    // Save to datastore
    if err := widget.update(key, c); err != nil {
      serveError(c, w, err)
    }
    http.Redirect(w, r, "/admin/widget-list", http.StatusFound)
}

func handleWidgetDelete(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    

    // Get widget id
    id := getUrlQuery(r.URL, "id")

    // Get widget data
    _, key, err := getWidget(id, c)
    if err != nil {
      serveError(c, w, err)
      return
    }

    datastore.Delete(c, key)

    http.Redirect(w, r, "/admin/widget-list", http.StatusFound)
}
