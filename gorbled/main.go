package gorbled

import (
    "net/http"
)

func init() {
    http.HandleFunc("/", handle)
}

var urls = map[string](func(http.ResponseWriter, *http.Request)){

    "/admin/config": handleConfigEdit,

    "/admin/article-list":   handleArticleList,
    "/admin/article-add":    handleArticleAdd,
    "/admin/article-edit":   handleArticleEdit,
    "/admin/article-delete": handleArticleDelete,

    "/admin/file-list": handleFileList,
    "/admin/file-edit": handleFileEdit,
    "/admin/file-new-url": handleFileNewUrl,
    "/admin/file-upload": handleFileUpload,
    "/admin/file-delete": handleFileDelete,
    "/admin/file-data": handleFileData,

    "/admin/widget-list":   handleWidgetList,
    "/admin/widget-add":    handleWidgetAdd,
    "/admin/widget-edit":   handleWidgetEdit,
    "/admin/widget-delete": handleWidgetDelete,

    "/article": handleArticleView,
    "/file": handleFileGet,
    "/rss.xml": handleRSS,
    "/":    handleIndex,

}

func handle(w http.ResponseWriter, r *http.Request) {
    config = initConfig(r)
    urls[r.RequestURI](w, r)
}