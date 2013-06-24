package gorbled

import (
	"appengine/datastore"
	"github.com/gorilla/mux"
	"net/http"
)

var (
	config    Config
	configKey *datastore.Key
)

func init() {
	r := mux.NewRouter()
	a := r.PathPrefix("/admin").Subrouter()

	// entry.go
	a.HandleFunc("/{entryType:page|widget|article}", requireConfig(handleEntryList))
	a.HandleFunc("/{entryType:page|widget|article}/{pid:[0-9]+}", requireConfig(handleEntryList))
	a.HandleFunc("/{entryType:page|widget|article}/add", requireConfig(handleEntryAdd))
	a.HandleFunc("/{entryType:page|widget|article}/edit/{id}", requireConfig(handleEntryEdit))
	a.HandleFunc("/{entryType:page|widget|article}/delete/{id}", handleEntryDelete)

	r.HandleFunc("/decodeContent", handleDecodeContent)
	r.HandleFunc("/article/{id}", requireConfig(handleEntryView))

	r.HandleFunc("/{url:.+}/", handleRedirect)

	// user.go
	r.HandleFunc("/login", handleUserLogin)
	r.HandleFunc("/logout", handleUserLogout)

	// rss.go
	r.HandleFunc("/feed", requireConfig(handleRSS))

	// file.go
	a.HandleFunc("/file", requireConfig(handleFileList))
	a.HandleFunc("/file/{pid:[0-9]+}", requireConfig(handleFileList))
	a.HandleFunc("/file/edit/{id}", handleFileEdit)
	a.HandleFunc("/file/new-url/{num}", handleFileNewUrl)
	a.HandleFunc("/file/upload", requireConfig(handleFileUpload))
	a.HandleFunc("/file/delete/{id}", handleFileDelete)
	a.HandleFunc("/file/data/{pid:[0-9]+}", requireConfig(handleFileData))

	r.HandleFunc("/file/{id}", handleFileGet)

	// config.go
	a.HandleFunc("/config", requireConfig(handleConfigEdit))

	// index.go
	r.HandleFunc("/", requireConfig(handleIndex))
	r.HandleFunc("/{pid:[0-9]+}", requireConfig(handleIndex))

	http.Handle("/", r)
}
