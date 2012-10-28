package gorbled

import (
    "net/http"
    "io"
    "appengine"
    "fmt"
)

func serve404(w http.ResponseWriter) {
    w.WriteHeader(http.StatusNotFound)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    io.WriteString(w, "Not Found")
}

func serveError(c appengine.Context, w http.ResponseWriter, err error) {
    w.WriteHeader(http.StatusInternalServerError)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    io.WriteString(w, "Internal Server Error: " + fmt.Sprint(err))
    c.Errorf("%v", err)
}

