package gorbled

import (
	"appengine"
	"appengine/datastore"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

/*
 * Entry data struct
 */
type Entry struct {
	ID      string
	Title   string
	Date    time.Time
	Content []byte

	// Type: article:文章；page:页面；widget: 小工具
	Type      string
	SubPage   []string
	PageClass int
	Sequence  int
	Tags      []string
	Cats      string

	AllowComment bool
	// Status: published 已发布, deleted 删除, unpublished 草稿
	Status   string
	Password string

	Slug      string
	Readtimes int
}

func (entry *Entry) save(c appengine.Context) (err error) {
	config.EntryID = config.EntryID + 1
	config.update(configKey, c)
	entry.ID = fmt.Sprint(config.EntryID)
	entry.put(c)
	return
}

func (entry *Entry) put(c appengine.Context) (err error) {
	_, err = datastore.Put(c, datastore.NewIncompleteKey(c, "Entry", nil), entry)
	return
}

func (entry *Entry) update(key *datastore.Key, c appengine.Context) (err error) {
	_, err = datastore.Put(c, key, entry)
	return
}

func getEntry(id string, c appengine.Context) (entry Entry, key *datastore.Key, err error) {

	dbQuery := datastore.NewQuery("Entry").Filter("ID =", id)
	var entries []Entry
	keys, err := dbQuery.GetAll(c, &entries)
	if len(entries) > 0 {
		entry = entries[0]
		key = keys[0]
	}

	return
}

/*
 * Get widget data
 */
func getWidgets(c appengine.Context) (widgets []Entry, err error) {

	dbQuery := datastore.NewQuery("Entry").
		Filter("Type =", "widget").
		Order("Sequence")
	_, err = dbQuery.GetAll(c, &widgets)

	return
}

func getWidgetsAndNav(pageId, pageSize int, c appengine.Context) (widgets []Entry, nav PageNav, err error) {

	// Get offset and page nav
	dbQuery := datastore.NewQuery("Entry").Filter("Type =", "widget")
	count, _ := dbQuery.Count(c)
	offset, nav := getPageNav(count, pageId, pageSize, c)

	// Get widget data
	dbQuery = dbQuery.Order("Sequence").Offset(offset).Limit(pageSize)
	_, err = dbQuery.GetAll(c, &widgets)

	return
}

/*
 * Get article data
 */
func getArticles(c appengine.Context) (articles []Entry, err error) {

	dbQuery := datastore.NewQuery("Entry").
		Filter("Type =", "article")
	_, err = dbQuery.GetAll(c, &articles)

	return
}

func getArticlesAndNav(pageId, pageSize int, c appengine.Context) (articles []Entry, nav PageNav, err error) {

	// Get offset and page nav
	dbQuery := datastore.NewQuery("Entry").Filter("Type =", "article")
	count, _ := dbQuery.Count(c)
	offset, nav := getPageNav(count, pageId, pageSize, c)

	// Get article data
	dbQuery = dbQuery.Order("-Date").Offset(offset).Limit(pageSize)
	_, err = dbQuery.GetAll(c, &articles)

	return
}

/*
 * Get page data
 */
func getPagesAndNav(pageId, pageSize int, c appengine.Context) (pages []Entry, nav PageNav, err error) {

	// Get offset and page nav
	dbQuery := datastore.NewQuery("Entry").
		Filter("Type =", "page").
		Filter("PageClass =", 0)
	count, _ := dbQuery.Count(c)
	offset, nav := getPageNav(count, pageId, pageSize, c)

	// Get page data
	dbQuery = dbQuery.Order("Sequence").Offset(offset).Limit(pageSize)
	_, err = dbQuery.GetAll(c, &pages)

	return
}

/*
 * Entry handler
 */

/*
 * Decode markdown code
 *
 * @return (string)
 */
func handleDecodeContent(w http.ResponseWriter, r *http.Request) {
	content := []byte(r.FormValue("content"))
	fmt.Fprint(w, decodeMD(content))
}

func handleEntryView(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Get entry data
	entry, _, err := getEntry(getUrlVar(r, "id"), c)
	if err != nil {
		serveError(w, err)
		return
	}

	// Get widget data
	widgets, err := getWidgets(c)
	if err != nil {
		serveError(w, err)
		return
	}

	// New Page
	page := Page{
		"User":    getUserInfo(c),
		"Entry":   entry,
		"Widgets": widgets,
		"Config":  config,
	}

	// Render page
	page.Render("article", w)
}

func handleEntryDelete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	// Get entry key
	_, key, err := getEntry(getUrlVar(r, "id"), c)
	if err != nil {
		serveError(w, err)
		return
	}

	// Delete entry
	if err = datastore.Delete(c, key); err != nil {
		serveError(w, err)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusFound)
}

func handleEntryAdd(w http.ResponseWriter, r *http.Request) {
	entryType := getUrlVar(r, "entryType")
	c := appengine.NewContext(r)

	if r.Method != "POST" {
		entryTypeTitle := strings.Title(entryType)
		page := Page{
			"Title":               "Add " + entryTypeTitle,
			"Config":              config,
			"Is" + entryTypeTitle: true,
			"New":       true,
			"ActionUrl": "/admin/" + entryType + "/add",
		}

		page.Render("admin/entry-edit", w)

		return
	}

	// r.Method == "POST"
	// Create entry
	entry := &Entry{
		Date:    time.Now(),
		Title:   r.FormValue("title"),
		Content: []byte(r.FormValue("content")),
		Type:    entryType,
	}

	switch entryType {
	case "widget":
		entry.Sequence, _ = strconv.Atoi(r.FormValue("sequence"))

		if err := entry.save(c); err != nil {
			serveError(w, err)
			return
		}
	case "page":
		fallthrough
	case "article":
		customid := r.FormValue("customid")
		if customid != "" && !checkIdIsExists("Entry", customid, c) {
			entry.ID = customid

			if err := entry.put(c); err != nil {
				serveError(w, err)
				return
			}

		} else if err := entry.save(c); err != nil {
			serveError(w, err)
			return
		}

	}

	http.Redirect(w, r, "/admin/"+entryType, http.StatusFound)
}

func handleEntryEdit(w http.ResponseWriter, r *http.Request) {
	entryType := getUrlVar(r, "entryType")
	c := appengine.NewContext(r)

	entry, key, err := getEntry(getUrlVar(r, "id"), c)
	if err != nil {
		serveError(w, err)
		return
	}

	if entry.ID == "" {
		serve404(w)
		return
	}

	if r.Method != "POST" {
		entryTypeTitle := strings.Title(entryType)
		page := Page{
			"Title":               "Edit " + entryTypeTitle,
			"Entry":               entry,
			"Config":              config,
			"Is" + entryTypeTitle: true,
			"New":       false,
			"ActionUrl": "/admin/" + entryType + "/edit/" + entry.ID,
		}

		page.Render("admin/entry-edit", w)

		return
	}

	// r.Method == "POST"
	// Update entry
	entry.Title = r.FormValue("title")
	entry.Content = []byte(r.FormValue("content"))

	switch entryType {
	case "widget":
		entry.Sequence, _ = strconv.Atoi(r.FormValue("sequence"))
	case "page":
		fallthrough
	case "article":
		customid := r.FormValue("customid")
		if customid != "" && !checkIdIsExists("Entry", customid, c) {
			entry.ID = customid
		}
	}

	if err := entry.update(key, c); err != nil {
		serveError(w, err)
		return
	}

	http.Redirect(w, r, "/admin/"+entryType, http.StatusFound)
}

func handleEntryList(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	entryType := getUrlVar(r, "entryType")

	// Get page id, pageSize
	pageId, _ := strconv.Atoi(getUrlVar(r, "pid"))
	var pageSize int
	switch entryType {
	case "article":
		pageSize = config.AdminArticles
	case "widget":
		pageSize = config.AdminWidgets
	case "page":
		pageSize = config.AdminPages
	}

	entriesAndNav := map[string]func(int, int, appengine.Context) ([]Entry, PageNav, error){
		"article": getArticlesAndNav,
		"page":    getPagesAndNav,
		"widget":  getWidgetsAndNav,
	}

	entries, nav, err := entriesAndNav[entryType](pageId, pageSize, c)
	if err != nil {
		serveError(w, err)
		return
	}

	entryTypeTitle := strings.Title(entryType)
	page := Page{
		"Title":     entryTypeTitle + " Manager",
		"Entries":   entries,
		"Nav":       nav,
		"Config":    config,
		"Action":    "Add " + entryTypeTitle,
		"EntryType": entryType,
	}

	page.Render("admin/entry-list", w)
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	url := getUrlVar(r, "url")
	http.Redirect(w, r, "/"+url, http.StatusFound)
}
