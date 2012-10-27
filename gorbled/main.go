package gorbled

import (
    "net/http"
)

var config Config = getConfig()

func init() {
    http.HandleFunc("/",        handleIndex)
    http.HandleFunc("/user",    handleUser)
    http.HandleFunc("/article", handleArticle)
    http.HandleFunc("/widget",  handleWidget)
    http.HandleFunc("/file",    handleFile)
}
