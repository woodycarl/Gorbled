package gorbled

import (
	"fmt"
	"net/http"
)

func serve404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "Not Found!")
}

func serveError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "Internal Server Error:"+" "+fmt.Sprint(err))
}
