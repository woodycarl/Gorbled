package gorbled

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

const (
	CONFIG_FILE_PATH = "config.json"
)

type Config struct {
	Title       string
	Description string
	SubTitle    string
	Author      string

	Articles      int
	AdminArticles int
	AdminWidgets  int
	AdminPages    int
	AdminFiles    int
	NavLen        int

	Theme           string
	GoogleAnalytics string
	Disqus          string
	GooglePlus      string

	TimeZone   float64
	TimeFormat string
	BaseUrl    string

	Version float64
	Program string

	EntryID int
	FileID  int

	// = =
	Language string
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
		key = keys[0]
		con = configs[0]
	}

	return
}

func getJsonConfig() (config Config) {
	configFile, err := ioutil.ReadFile(CONFIG_FILE_PATH)
	if err != nil {
		log.Fatal(err)
	}
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
		pagina := Pagina{
			"Title":  "设置",
			"Config": config,
		}

		pagina.Render("admin/config", w)

		return
	}

	// r.Method == "POST"
	// Update config data
	config.Title = r.FormValue("title")
	config.Description = r.FormValue("description")
	config.Articles, _ = strconv.Atoi(r.FormValue("articles"))
	config.AdminArticles, _ = strconv.Atoi(r.FormValue("admin-articles"))
	config.AdminWidgets, _ = strconv.Atoi(r.FormValue("admin-widgets"))
	config.Theme = checkTheme(r.FormValue("theme"))
	config.TimeZone, _ = strconv.ParseFloat(r.FormValue("timezone"), 64)
	config.TimeFormat = r.FormValue("time-format")
	config.Disqus = r.FormValue("disqus")
	config.GoogleAnalytics = r.FormValue("google-analytics")
	config.Theme = r.FormValue("theme")
	config.SubTitle = r.FormValue("subtitle")

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
	con.save(c)
	config, configKey, _ = getConfig(c)

	article := Entry{
		Title:   "Hello World!",
		Date:    time.Now(),
		Content: []byte(fmt.Sprintf("欢迎使用Gorbled %.1f，可以随时删除或编辑这篇文章。", config.Version)),
		Type:    "article",
	}
	article.save(c)

	widget := Entry{
		Title:   "公告",
		Date:    time.Now(),
		Content: []byte("这个是**公告**呢！"),
		Type:    "widget",
	}
	widget.save(c)

	return
}

/*
 * Init system
 */
func initSystem(r *http.Request) {
	c := appengine.NewContext(r)

	con, key, err := getConfig(c)
	if err != nil || con.Title == "" {
		installSystem(c)
	} else {
		config = con
		configKey = key
	}

	config.BaseUrl = "http://" + r.Host
}
