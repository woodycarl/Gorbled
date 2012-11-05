package gorbled

import (
    "net/http"
    "github.com/gorilla/mux"

)

func init() {
    r := mux.NewRouter()

    // widget.go
    r.HandleFunc("/admin/widget/", handleRedirectWidgetList)
    r.HandleFunc("/admin/widget", handleWidgetList)
    r.HandleFunc("/admin/widget/{pid:[0-9]+}", handleWidgetList)
    r.HandleFunc("/admin/widget/add", handleWidgetAdd)
    r.HandleFunc("/admin/widget/edit/{id}", handleWidgetEdit)
    r.HandleFunc("/admin/widget/delete/{id}", handleWidgetDelete)

    // user.go
    r.HandleFunc("/login", handleUserLogin)
    r.HandleFunc("/logout", handleUserLogout)

    // rss.go
    r.HandleFunc("/feed", handleRSS)

    // file.go
    r.HandleFunc("/admin/file/", handleRedirectFileList)
    r.HandleFunc("/admin/file", handleFileList)
    r.HandleFunc("/admin/file/{pid:[0-9]+}", handleFileList)
    r.HandleFunc("/admin/file/edit/{id}", handleFileEdit)
    r.HandleFunc("/admin/file/new-url/{num}", handleFileNewUrl)
    r.HandleFunc("/admin/file/upload", handleFileUpload)
    r.HandleFunc("/admin/file/delete/{id}", handleFileDelete)
    r.HandleFunc("/admin/file/data/{pid:[0-9]+}", handleFileData)

    r.HandleFunc("/file/{key}", handleFileGet)

    // article.go
    r.HandleFunc("/admin/article/", handleRedirectArticleList)
    r.HandleFunc("/admin/article", handleArticleList)
    r.HandleFunc("/admin/article/{pid:[0-9]+}", handleArticleList)
    r.HandleFunc("/admin/article/add", handleArticleAdd)
    r.HandleFunc("/admin/article/edit/{id}", handleArticleEdit)
    r.HandleFunc("/admin/article/delete/{id}", handleArticleDelete)

    r.HandleFunc("/decodeContent", handleDecodeContent)
    r.HandleFunc("/article/{id}", handleArticleView)

    // config.go
    r.HandleFunc("/admin/config", handleConfigEdit)

    // lang.go
    r.HandleFunc("/admin/init/lang", handleInitLang)

    // index.go
    r.HandleFunc("/", handleIndex)
    r.HandleFunc("/{pid:[0-9]+}", handleIndex)

    r.HandleFunc("/test/{id}", handle)
    r.HandleFunc("/test", handle)

    http.Handle("/", r)
}

func handle(w http.ResponseWriter, r *http.Request) {

}
