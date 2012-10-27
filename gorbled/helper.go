package gorbled

import (
    "crypto/rand"
    "fmt"
    "net/url"
    "io"
    "net/http"

    "appengine"
    "appengine/datastore"
)

func init() {
    http.HandleFunc("/decodeContent", handleDecodeContent)
}

/*
 * Decode markdown code
 *
 * @return (string) 
 */
func handleDecodeContent(w http.ResponseWriter, r *http.Request) {
    content := []byte(r.FormValue("content"))
    fmt.Fprint(w, decodeMD(content))
}


/*
 * Generate ID
 *
 * @return (string) 
 */
func genID() string {
    buf := make([]byte, 16)
    io.ReadFull(rand.Reader, buf)

    return fmt.Sprintf("%x", buf)
}

/*
 * Check id is exists
 *
 * @param kind (string)
 * @param id   (string)
 *
 * @return (string)
 */
func getID(kind string, id string, c appengine.Context) string {
    if id != "" && !checkIdIsExists(kind, id, c) {
        return id
    }

    return genID()
}

/*
 * Check ID is exists
 *
 * @param kind (string)
 * @param id   (string)
 *
 * @return (bool)
 */
func checkIdIsExists(kind string, id string, c appengine.Context) bool {
    dbQuery := datastore.NewQuery(kind).Filter("ID =", id)

    if count, _ := dbQuery.Count(c); count < 1 {
        return false
    }

    return true
}

/*
 * Get Url Query 
 *
 * @param u     (url.Url)
 * @param query (string)
 *
 * @return result ([]string)
 */
func getUrlQuery(u *url.URL, query string) (result string) {
    urlQuery := u.Query()
    result   = urlQuery.Get(query)

    return
}
