package gorbled

import (
    "log"
    "net/http"
    "time"

    "io/ioutil"
    "encoding/json"
    "appengine"
    "appengine/datastore"
    "strconv"
)

const (
    CONFIG_FILE_PATH = "config.json"
)

var config Config

type Config struct {
    Title               string
    Description         string

    Articles            int
    AdminArticles       int
    AdminWidgets        int
    AdminFiles          int
    NavLen              int

    Theme               string
    GoogleAnalytics     string
    Disqus              string

    TimeZone            float64
    BaseUrl             string
    Version             float64
}

func (config *Config) save(c appengine.Context) (err error) {
    _, err = datastore.Put(c, datastore.NewIncompleteKey(c, "Config", nil), config)

    return
}

func (config *Config) update(key *datastore.Key, c appengine.Context) (err error) {
    _, err = datastore.Put(c, key, config)

    return
}

func getConfig(c appengine.Context) (config Config, key *datastore.Key, err error) {
    dbQuery := datastore.NewQuery("Config")
    var configs []Config
    keys, err := dbQuery.GetAll(c, &configs)
    if len(keys) > 0 {
        key    = keys[0]
        config = configs[0]
    }

    return
}

func getJsonConfig() (config Config) {
    configFile, err := ioutil.ReadFile(CONFIG_FILE_PATH)
    err = json.Unmarshal(configFile, &config)
    if err != nil {
        log.Fatal(err)
    }

    return
}

/*
 * handle config edit
 */
func handleConfigEdit(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    if r.Method != "POST" {
        // Show article edit page

        // New Page
        page := Page {
            Title:  "Config",
            Config: config,
        }

        // Render page
        err := page.Render("admin/config", w)
        if err != nil {
            serveError(c, w, err)
        }

        return
    }

    // Process post data

    if err := r.ParseForm(); err != nil {
        serveError(c, w, err)
        return
    }

    // Get config
    config, key, err := getConfig(c)
    if err != nil {
        serveError(c, w, err)
        return
    }

    // Update config data
    config.Title = r.FormValue("title")
    config.Description = r.FormValue("description")
    config.Articles, _ = strconv.Atoi(r.FormValue("articles"))
    config.AdminArticles, _ = strconv.Atoi(r.FormValue("admin-articles"))
    config.AdminWidgets, _ = strconv.Atoi(r.FormValue("admin-widgets"))
    config.Theme = r.FormValue("theme")
    config.TimeZone, _ = strconv.ParseFloat(r.FormValue("timezone"), 64)

    if err := config.update(key, c); err != nil {
      serveError(c, w, err)
    }

    http.Redirect(w, r, "/admin/config", http.StatusFound)
}

/*
 * Init system
 */
func initSystem(c appengine.Context) (config Config) {
    config = getJsonConfig()
    config.save(c)

    article := Article {
        ID:      genID(),
        Title:   "Hello World!",
        Date:    time.Now(),
        Content: []byte("欢迎使用Gorbled，可以随时删除或编辑这篇文章。"),
    }
    article.save(c)

    widget := Widget {
        ID:      genID(),
        Title:   "公告",
        Content: []byte("这是个**公告**呢"),
    }
    widget.save(c)

    return
}

/*
 * Init config
 */
func initConfig(r *http.Request) (config Config) {
    c := appengine.NewContext(r)
    config, _, err := getConfig(c)
    if err != nil || config.Title == "" {
        config = initSystem(c)
    }
    config.BaseUrl = "http://" + r.Host
    return
}
