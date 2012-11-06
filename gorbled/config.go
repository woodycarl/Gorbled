package gorbled

import (
    "log"
    "net/http"
    "time"
    "fmt"
    "io/ioutil"
    "encoding/json"
    "appengine"
    "appengine/datastore"
    "strconv"

)

const (
    CONFIG_FILE_PATH = "config.json"
)

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
    TimeFormat          string
    BaseUrl             string
    Version             float64
    Language            string

    EntryID             int
}

func (config *Config) save(c appengine.Context) (err error) {
    _, err = datastore.Put(c, datastore.NewIncompleteKey(c, "Config", nil), config)

    return
}

func (config *Config) update(key *datastore.Key, c appengine.Context) (err error) {
    _, err = datastore.Put(c, key, config)

    return
}

func getConfig(c appengine.Context) (con Config, key *datastore.Key, err error) {
    dbQuery := datastore.NewQuery("Config")
    var configs []Config
    keys, err := dbQuery.GetAll(c, &configs)
    if len(keys) > 0 {
        key     = keys[0]
        con     = configs[0]
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
        pagina := Pagina {
            "Title":  "Config",
            "Config": config,
        }

        pagina.Render("admin/config", w)

        return
    }

    // r.Method == "GET"
    // Update config data
    config.Title = r.FormValue("title")
    config.Description = r.FormValue("description")
    config.Articles, _ = strconv.Atoi(r.FormValue("articles"))
    config.AdminArticles, _ = strconv.Atoi(r.FormValue("admin-articles"))
    config.AdminWidgets, _ = strconv.Atoi(r.FormValue("admin-widgets"))
    config.Theme = checkTheme(r.FormValue("theme"))
    config.TimeZone, _ = strconv.ParseFloat(r.FormValue("timezone"), 64)
    config.TimeFormat = r.FormValue("time-format")
    config.Version, _ = strconv.ParseFloat(r.FormValue("version"), 64)
    config.Disqus = r.FormValue("disqus")
    config.GoogleAnalytics = r.FormValue("google-analytics")

    if err := config.update(configKey, c); err != nil {
        serveError(w, err)
        return
    }

    http.Redirect(w, r, r.Referer(), http.StatusFound)
}

func checkTheme(s string) string {
    if s == "" {
        return "default"
    }
    return s
}

/*
 * Install system
 */
func installSystem(c appengine.Context) {
    con := getJsonConfig()
    con.EntryID = 0
    con.save(c)
    config, configKey, _ = getConfig(c)

    readLang(c)
    initLang(c)

    article := Entry {
        Title:   L("Hello World!"),
        Date:    time.Now(),
        Content: []byte(fmt.Sprintf(L("Welcome to Gorbled %.1f. You can edit or delete this post, then start blogging!"), config.Version)),
        Type:   "article",
    }
    article.save(c)

    widget := Entry {
        Title:   L("Notice"),
        Date:    time.Now(),
        Content: []byte(L("This is **Notice** !")),
        Type:   "widget",
    }
    widget.save(c)

    return
}

/*
 * Init system
 */
func initSystem(r *http.Request)  {
    c := appengine.NewContext(r)

    con, key, err := getConfig(c)
    if err != nil || con.Title == "" {
        installSystem(c)
    } else {
        config = con
        configKey = key
        initLang(c)
    }

    config.BaseUrl = "http://" + r.Host
}
