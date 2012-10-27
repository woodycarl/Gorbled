package gorbled

import (
    "net/http"
    "text/template"
    "time"
    "appengine"
    "appengine/datastore"
    "strings"

    "other_package/blackfriday"
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
    ShowIDs       bool
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
    return t.In(newLocation).Format("3:04pm, Mon 2 Jan")
}

var funcMap = template.FuncMap{
    "showDate": showDate,
    "decodeMD": decodeMD,
}

/*
 * 计算页面导航序号
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

  pageHalfLen := (NavLen / 2)
  if NavLen / 2 != 0 {
    pageHalfLen++
  }

  var start, length, nextId, prevId int
  var prev, next, showIDs bool

  switch {
    case pageNums <= NavLen :
      start = 1
      prev = false
      length = pageNums
      next = false
    case pageId <= pageHalfLen :
      start = 1
      prev = false
      length = NavLen
      next = true
      nextId = NavLen + 1
    case pageId + pageHalfLen >= pageNums :
      start = pageNums - NavLen
      prev = true
      prevId = start - 1
      length = NavLen
      next = false
    default :
      start = pageId - pageHalfLen + 1
      prev = true
      prevId = start - 1
      length = NavLen
      next = true
      nextId = start + length
  }

  var ids = make([]PageId, length)
  for i:=0;i<length;i++{
    ids[i].Id=i+start
    if ids[i].Id == pageId {
      ids[i].Current = true
    } else {
      ids[i].Current = false
    }
  }
  if length>1 {
    showIDs = true
  } else {
    showIDs = false
  }

  pageNav = PageNav {
    ShowPrev      : prev,
    ShowNext      : next,
    ShowIDs       : showIDs,
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
func (page *Page) Render(pageFilePath string, w http.ResponseWriter) (err error) {
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
        return
    }

    tmpl.Execute(w, page)
    return
}
