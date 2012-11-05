package gorbled

import (
    "fmt"
    "net/url"
    "net/http"

    "appengine"
    "appengine/datastore"
    "time"

    "github.com/gorilla/mux"
)

/*
 * Generate ID
 *
 * @return (string) 
 */
func genID() string {
    return fmt.Sprint(time.Now().Unix())
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

func getUrlVar(r *http.Request, v string) string {
    return mux.Vars(r)[v]
}