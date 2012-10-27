package gorbled

import (
    "net/http"
    "io"

    "appengine"
)

func serve404(w http.ResponseWriter) {
    w.WriteHeader(http.StatusNotFound)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    io.WriteString(w, "Not Found")
}

func serveError(c appengine.Context, w http.ResponseWriter, err error) {
    w.WriteHeader(http.StatusInternalServerError)
    w.Header().Set("Content-Type", "text/plain; charset=utf-8")
    io.WriteString(w, "Internal Server Error")
    c.Errorf("%v", err)
}
