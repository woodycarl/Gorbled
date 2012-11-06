package gorbled

import (
    "net/http"
    "github.com/gorilla/mux"

)

func init() {
    r := mux.NewRouter()

    // widget.go
    r.HandleFunc("/admin/widget/", handleRedirectWidgetList)
    r.HandleFunc("/admin/widget", requireConfig(handleWidgetList))
    r.HandleFunc("/admin/widget/{pid:[0-9]+}", requireConfig(handleWidgetList))
    r.HandleFunc("/admin/widget/add", requireConfig(handleWidgetAdd))
    r.HandleFunc("/admin/widget/edit/{id}", requireConfig(handleWidgetEdit))
    r.HandleFunc("/admin/widget/delete/{id}", handleWidgetDelete)

    // user.go
    r.HandleFunc("/login", handleUserLogin)
    r.HandleFunc("/logout", handleUserLogout)

    // rss.go
    r.HandleFunc("/feed", requireConfig(handleRSS))

    // file.go
    r.HandleFunc("/admin/file/", handleRedirectFileList)
    r.HandleFunc("/admin/file", requireConfig(handleFileList))
    r.HandleFunc("/admin/file/{pid:[0-9]+}", requireConfig(handleFileList))
    r.HandleFunc("/admin/file/edit/{id}", handleFileEdit)
    r.HandleFunc("/admin/file/new-url/{num}", handleFileNewUrl)
    r.HandleFunc("/admin/file/upload", handleFileUpload)
    r.HandleFunc("/admin/file/delete/{id}", handleFileDelete)
    r.HandleFunc("/admin/file/data/{pid:[0-9]+}", handleFileData)

    r.HandleFunc("/file/{key}", handleFileGet)

    // article.go
    r.HandleFunc("/admin/article/", handleRedirectArticleList)
    r.HandleFunc("/admin/article", requireConfig(handleArticleList))
    r.HandleFunc("/admin/article/{pid:[0-9]+}", requireConfig(handleArticleList))
    r.HandleFunc("/admin/article/add", requireConfig(handleArticleAdd))
    r.HandleFunc("/admin/article/edit/{id}", requireConfig(handleArticleEdit))
    r.HandleFunc("/admin/article/delete/{id}", handleArticleDelete)

    r.HandleFunc("/decodeContent", handleDecodeContent)
    r.HandleFunc("/article/{id}", requireConfig(handleArticleView))

    // config.go
    r.HandleFunc("/admin/config", requireConfig(handleConfigEdit))

    // lang.go
    r.HandleFunc("/admin/init/lang", requireConfig(handleInitLang))

    // index.go
    r.HandleFunc("/", requireConfig(handleIndex))
    r.HandleFunc("/{pid:[0-9]+}", requireConfig(handleIndex))

    r.HandleFunc("/test/{id}", handle)
    r.HandleFunc("/test", handle)

    http.Handle("/", r)
}

func handle(w http.ResponseWriter, r *http.Request) {

}
