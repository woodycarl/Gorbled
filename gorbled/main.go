package gorbled

import (
    "net/http"
)

var urls = map[string](func(http.ResponseWriter, *http.Request)){
    // config.go
    "/admin/config": handleConfigEdit,

    // article.go
    "/admin/article-list":   handleArticleList,
    "/admin/article-add":    handleArticleAdd,
    "/admin/article-edit":   handleArticleEdit,
    "/admin/article-delete": handleArticleDelete,
    "/decodeContent": handleDecodeContent,

    "/article": handleArticleView,

    // file.go
    "/admin/file-list": handleFileList,
    "/admin/file-edit": handleFileEdit,
    "/admin/file-new-url": handleFileNewUrl,
    "/admin/file-upload": handleFileUpload,
    "/admin/file-delete": handleFileDelete,
    "/admin/file-data": handleFileData,

    "/file": handleFileGet,

    // widget.go
    "/admin/widget-list":   handleWidgetList,
    "/admin/widget-add":    handleWidgetAdd,
    "/admin/widget-edit":   handleWidgetEdit,
    "/admin/widget-delete": handleWidgetDelete,

    // user.go
    "/login": handleUserLogin,
    "/logout": handleUserLogout,

    // rss.go
    "/rss.xml": handleRSS,

    // index.go
    "/":    handleIndex,

    // lang.go
    "/admin/init-lang": handleInitLang,

}

func init() {
    http.HandleFunc("/", handle)
}

func handle(w http.ResponseWriter, r *http.Request) {
    config = initConfig(r)

    lang = initLang(r, config.Language)

    urls[r.URL.Path](w, r)
}
