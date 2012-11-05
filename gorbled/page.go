package gorbled

import (
    "net/http"
    "text/template"
    "time"
    "appengine"
    "appengine/datastore"
    "strings"
    "github.com/russross/blackfriday"
)

type Page struct {
    Title       string

    User        User
    Articles    []Article
    Article     Article
    Widgets     []Widget
    Widget      Widget
    Files       []File
    File        File

    Nav         PageNav
    Config      Config

    New         bool    // Show add or edit template
}

type PageId struct {
    Id         int
    Current    bool
}

type PageNav struct {
    ShowPrev      bool
    ShowNext      bool
    NextPageID    int
    PrevPageID    int
    PageIDs       []PageId
}

/*
 * Render FuncMap
 */

func decodeMD(content []byte) string {
    return string(blackfriday.MarkdownCommon(content))
}

func showDate(t time.Time) string {
    newLocation := time.FixedZone("myTimeZone", (int)(config.TimeZone*60*60))
    return t.In(newLocation).Format(config.TimeFormat)
}

var funcMap = template.FuncMap{
    "showDate": showDate,
    "decodeMD": decodeMD,
    "l":        L,
}

/*
 * 计算页面导航序号
 *
 * <<       2       3   ...   x             >>
 * prev     ids[0]  ids[1]    ids[NavLen-1] next
 */
func getPageNav(kind string, pageId int, pageSize int, c appengine.Context) (offset int, pageNav PageNav) {
    NavLen := config.NavLen

    dbQuery  := datastore.NewQuery(kind)
    count, _ := dbQuery.Count(c)
    pageNums := (count / pageSize)
    if count % pageSize != 0 {
        pageNums++
    }
    if pageId <= 0 || pageId > pageNums {
        pageId = 1
    }
    offset = (pageId - 1) * pageSize

    var start, length, nextId, prevId int
    var prev, next bool

    start = ((pageId-1)/NavLen)*NavLen + 1
    if start+NavLen-1<=pageNums {
        length = NavLen
    } else {
        length = pageNums - start +1
    }
    if start - 1 > 0 {
        prev = true
        prevId = start - 1
    } else {
        prev = false
    }
    if start + length <= pageNums {
        next = true
        nextId = start + length
    } else {
        next = false
    }

    var ids = make([]PageId, length)
    for i:=0; i < length;i++{
        ids[i].Id = i+start
        if ids[i].Id == pageId {
            ids[i].Current = true
        } else {
            ids[i].Current = false
        }
    }

    pageNav = PageNav {
        ShowPrev      : prev,
        ShowNext      : next,
        NextPageID    : nextId,
        PrevPageID    : prevId,
        PageIDs       : ids,
    }

    return
}

/*
 * Render page
 *
 * @param page          (string)
 * @param w             (http.ResponseWriter)
 *
 * @return (error)
 */
func (page *Page) Render(pageFilePath string, w http.ResponseWriter) {
    base := "gorbled/templates/" + config.Theme + "/"
    
    if strings.Contains(pageFilePath, "admin") {
        base = "gorbled/admin/"
        pageFilePath = strings.Replace(pageFilePath, "admin/", "", -1)
    }

    tmpl, err := template.New("main.html").Funcs(funcMap).ParseFiles(
        base + "main.html",
        base + "sidebar.html",
        base + pageFilePath + ".html",
    )

    if err != nil {
        serveError(w, err)
        return
    }

    if err = tmpl.Execute(w, page); err != nil {
        serveError(w, err)
        return
    }

}
