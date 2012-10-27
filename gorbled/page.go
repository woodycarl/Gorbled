package gorbled

import (
    "net/http"
    "text/template"
)

type Page struct {
    PageSetting *PageSetting
    PageData    *PageData
}

type PageSetting struct {
    Title, Description, Layout string
    NextPageID, PrevPageID int
    ShowNext, ShowPrev, ShowSidebar bool
}

type PageData struct {
    User    *UserData
    Article []ArticleData
    Widget  []WidgetData
    File    FileData
}

const (
    // Folder setting
    STATIC_FOLDER = "gorbled/static/html/"
    LAYOUT_FOLDER = "gorbled/static/html/layouts/"

    // Default layout setting
    DEFAULT_LAYOUT = "column2"
)

/*
 * New Page
 *
 * @param layout      (string)
 * @param showSidebar (bool)
 * @param pageData    (*PageData)
 *
 * @return (*Page)
 */
func NewPage(pageSetting *PageSetting, pageData *PageData) *Page {
    if pageSetting.Layout == "" {
        pageSetting.Layout = DEFAULT_LAYOUT
    }

    if pageSetting.Title == "" {
        pageSetting.Title = config.Title
    }

    if pageSetting.Description == "" {
        pageSetting.Description = config.Description
    }

    return &Page{ PageSetting: pageSetting, PageData: pageData }
}

/*
 * Render page
 *
 * @param pageFilePath (string)
 * @param w            (http.ResponseWriter)
 *
 * @return (error)
 */
func (page *Page) Render(pageFilePath string, w http.ResponseWriter) (err error) {
    columnFilePath  := page.PageSetting.Layout + ".html"
    mainFilePath    := "main.html"
    contentFilePath := pageFilePath + ".html"
    sidebarFilePath := "sidebar.html"

    var tmpl *template.Template

    switch page.PageSetting.ShowSidebar {
        case true:
            tmpl, err = template.ParseFiles(
                            LAYOUT_FOLDER + mainFilePath,
                            LAYOUT_FOLDER + columnFilePath,
                            LAYOUT_FOLDER + sidebarFilePath,
                            STATIC_FOLDER + contentFilePath)
        case false:
            tmpl, err = template.ParseFiles(
                            LAYOUT_FOLDER + mainFilePath,
                            LAYOUT_FOLDER + columnFilePath,
                            STATIC_FOLDER + contentFilePath)

    }

    if err != nil {
        return
    }

    tmpl.Execute(w, page)
    return
}
