package gorbled

import (
    "net/http"
    "strconv"

    "appengine"
    "appengine/datastore"

)

/*
 * Get widget data
 */
func getWidgets(c appengine.Context) (widgets []Entry, err error) {

    dbQuery := datastore.NewQuery("Entry").
        Filter("Type =", "widget").
        Order("Sequence")
    _, err = dbQuery.GetAll(c, &widgets)

    return
}

func getWidgetsAndNav(paginaId, paginaSize int, c appengine.Context) (widgets []Entry, nav PaginaNav, err error) {

    // Get offset and pagina nav
    dbQuery := datastore.NewQuery("Entry").Filter("Type =", "widget")
    count, _ := dbQuery.Count(c)
    offset, nav := getPaginaNav(count, paginaId, paginaSize, c)

    // Get widget data
    dbQuery = dbQuery.Order("Sequence").Offset(offset).Limit(paginaSize)
    _, err = dbQuery.GetAll(c, &widgets)

    return
}

/*
 * Widget handler
 */
func handleWidgetList(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    // Get pagina id, paginaSize
    paginaId, _ := strconv.Atoi(getUrlVar(r, "pid"))
    paginaSize  := config.AdminWidgets

    widgets, nav, err := getWidgetsAndNav(paginaId, paginaSize, c)
    if err != nil {
        serveError(w, err)
        return
    }

    // New Pagina
    pagina := Pagina {
        "Title" :     "Widget Manager",
        "Widgets" :   widgets,
        "Nav" :       nav,
        "Config" :    config,
    }

    // Render pagina
    pagina.Render("admin/widgets", w)
}

func handleWidgetAdd(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    if r.Method != "POST" {
        // New Pagina
        pagina := Pagina {
            "Title":     "Add Widget",
            "Config":    config,
            "New":       true,
        }

        // Show widget add pagina
        pagina.Render("admin/widget", w)

        return
    }

    // Process post data

    // Parse form data
    if err := r.ParseForm(); err != nil {
        serveError(w, err)
        return
    }

    // Create widget
    sequence, _ := strconv.Atoi(r.FormValue("sequence"))
    widget := &Entry{
        ID:         getID("Widget", r.FormValue("customid"), c),
        Title:      r.FormValue("title"),
        Content:    []byte(r.FormValue("content")),
        Sequence:   sequence,
        Type:       "widget",
    }

    // Save to datastore
    if err := widget.save(c); err != nil {
        serveError(w, err)
        return
    }

    http.Redirect(w, r, "/admin/widget", http.StatusFound)
}

func handleWidgetEdit(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    // Get widget id
    id := getUrlVar(r, "id")

    // Get widget data
    widget, key, err := getEntry(id, c)
    if err != nil {
        serveError(w, err)
        return
    }

    // Check widget is exists
    if widget.ID == "" {
        serve404(w)
        return
    }

    if r.Method != "POST" {
        // New Pagina
        pagina := Pagina {
            "Title" :     "Edit Widget",
            "Config" :    config,
            "Widget":     widget,
            "New":        false,
        }

        // Show widget edit pagina
        pagina.Render("admin/widget", w)

        return
    }

    // Process post data

    // Parse form data
    if err := r.ParseForm(); err != nil {
        serveError(w, err)
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
        serveError(w, err)
        return
    }
    http.Redirect(w, r, "/admin/widget", http.StatusFound)
}

func handleRedirectWidgetList(w http.ResponseWriter, r *http.Request) {
    http.Redirect(w, r, "/admin/widget", http.StatusFound)
}