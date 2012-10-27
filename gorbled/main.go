package gorbled

import (
    "net/http"
)

func init() {
    http.HandleFunc("/", handleIndex)
}
