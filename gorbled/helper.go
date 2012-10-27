package gorbled

import (
    "crypto/rand"
    "fmt"
    "net/url"
    "io"

    "appengine"
    "appengine/datastore"
)

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
 * Get Offset && PageNums
 *
 * @param kind     (string)
 * @param pageId   (int)
 * @param pageSize (int)
 * @param c        (appengine.Context)
 * 
 * @return offset   (int)
 * @return pageNums (int)
 */
func getOffset(kind string, pageId int, pageSize int, c appengine.Context) (offset int, pageNums int) {
    dbQuery  := datastore.NewQuery(kind)
    count, _ := dbQuery.Count(c)
    pageNums = (count / pageSize)
    if count % pageSize != 0 {
        pageNums++
    }

    if pageId <= 0 || pageId > pageNums {
        pageId = 1
    }

    offset = (pageId - 1) * pageSize

    return
}

/*
 * Check ID is exists
 *
 * @param kind (string)
 * @param id   (string)
 *
 * @return bool
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
