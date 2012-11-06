package gorbled

import (
    "net/http"
    "text/template"
    "time"
    "appengine"
    "strings"
    "github.com/russross/blackfriday"
)

type Pagina map[string]interface{}

type PageId struct {
    Id         int
    Current    bool
}

type PaginaNav struct {
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
func getPaginaNav(count, paginaId, paginaSize int, c appengine.Context) (offset int, paginaNav PaginaNav) {
    NavLen := config.NavLen

    paginaNums := (count / paginaSize)
    if count % paginaSize != 0 {
        paginaNums++
    }
    if paginaId <= 0 || paginaId > paginaNums {
        paginaId = 1
    }
    offset = (paginaId - 1) * paginaSize

    var start, length, nextId, prevId int
    var prev, next bool

    start = ((paginaId-1)/NavLen)*NavLen + 1
    if start+NavLen-1<=paginaNums {
        length = NavLen
    } else {
        length = paginaNums - start +1
    }
    if start - 1 > 0 {
        prev = true
        prevId = start - 1
    } else {
        prev = false
    }
    if start + length <= paginaNums {
        next = true
        nextId = start + length
    } else {
        next = false
    }

    var ids = make([]PageId, length)
    for i:=0; i < length;i++{
        ids[i].Id = i+start
        if ids[i].Id == paginaId {
            ids[i].Current = true
        } else {
            ids[i].Current = false
        }
    }

    paginaNav = PaginaNav {
        ShowPrev      : prev,
        ShowNext      : next,
        NextPageID    : nextId,
        PrevPageID    : prevId,
        PageIDs       : ids,
    }

    return
}

/*
 * Render pagina
 *
 * @param pagina          (string)
 * @param w             (http.ResponseWriter)
 *
 * @return (error)
 */
func (pagina *Pagina) Render(paginaFilePath string, w http.ResponseWriter) {
    base := "gorbled/templates/" + config.Theme + "/"
    
    if strings.Contains(paginaFilePath, "admin") {
        base = "gorbled/admin/"
        paginaFilePath = strings.Replace(paginaFilePath, "admin/", "", -1)
    }

    tmpl, err := template.New("main.html").Funcs(funcMap).ParseFiles(
        base + "main.html",
        base + "sidebar.html",
        base + paginaFilePath + ".html",
    )

    if err != nil {
        serveError(w, err)
        return
    }

    if err = tmpl.Execute(w, pagina); err != nil {
        serveError(w, err)
        return
    }
}

func requireConfig(handler http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        initSystem(r)
        handler(w, r)
    }
}