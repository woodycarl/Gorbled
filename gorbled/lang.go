package gorbled

import (
    "os"
    "fmt"
    "net/http"
    "bufio"
    "strings"
    "regexp"
    "strconv"
    "encoding/json"
    "appengine"
    "appengine/datastore"
    "path/filepath"
)

func init() {
    //http.HandleFunc("/admin/init-lang", handleInitLang)
}

type Lang struct {
    ID          string
    Content     []byte
}

var lang map[string]string

func (lang *Lang) save(c appengine.Context) (err error) {
    _, err = datastore.Put(c, datastore.NewIncompleteKey(c, "Lang", nil), lang)
    return
}

func (lang *Lang) update(key *datastore.Key, c appengine.Context) (err error) {
    _, err = datastore.Put(c, key, lang)
    return
}

func getLang(id string, c appengine.Context) (lang Lang, key *datastore.Key, err error) {
    dbQuery := datastore.NewQuery("Lang").Filter("ID =", id)
    var langs []Lang
    keys, err := dbQuery.GetAll(c, &langs)
    if len(langs) > 0 {
        lang = langs[0]
        key = keys[0]
    }

    return
}

func getLangs(c appengine.Context) (langs []Lang, keys []*datastore.Key, err error) {
    dbQuery := datastore.NewQuery("Lang")
    keys, err = dbQuery.GetAll(c, &langs)
    return
}

func handleInitLang(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    initLang2(c)
    lang = initLang(r, config.Language)
    fmt.Fprint(w, lang)
}

func initLang2(c appengine.Context) {
    _, keys, _ := getLangs(c)
    datastore.DeleteMulti(c, keys)

    readFile := func(path string, info os.FileInfo, err error) error {
        if !info.IsDir() {
            fileName := strings.Replace(info.Name(), ".lang", "", -1)
            l:=formatLang(fileName)
            l.save(c)
        }
        return nil
    }

    filepath.Walk("gorbled/local/", readFile)
    
}

func initLang(r *http.Request, language string) (langC map[string]string) {
    c := appengine.NewContext(r)

    langD, _, _ := getLang(language, c)

    dec := json.NewDecoder(strings.NewReader(string(langD.Content)))

    dec.Decode(&langC)

    return
}

func formatLangString(s string) string {
    r, _ := regexp.Compile(`\S.*\S`)

    return r.FindString(s)
}
func formatLangNum(s string) int {
    r, _ := regexp.Compile(`\d`)
    id, _ := strconv.Atoi(r.FindString(s))

    return id
}

func formatLang(file string) (lang Lang) {
    f, _ := os.Open("gorbled/local/"+file+".lang")
    read := bufio.NewReader(f)
    sentencesIn := make(map[string]string)
    sentencesOut := make(map[string]string)
    sentences := make(map[string]string)

    for true {
        //lineStr, err := read.ReadString('\n')
        lineStr, _, err := read.ReadLine()
        if err != nil {
            break
        }
        line := string(lineStr)

        isAnnotate, _ := regexp.MatchString(`^\s*(#.*|$)`, line)

        if !isAnnotate {
            isSentenceIn, _ := regexp.MatchString(`^\s*i.*`, line)
            s := strings.Split(strings.Split(line, "-")[1], ":")

            //id := formatLangNum(s[0])
            id := formatLangString(s[0])
            sentence := formatLangString(s[1])

            if isSentenceIn {
                sentencesIn[id] = sentence
            } else {
                sentencesOut[id] = sentence
            }
        }
    }
    for key, value := range sentencesIn {
        sentences[value] = sentencesOut[key]
    }
    content, _ := json.Marshal(sentences)
    lang = Lang {
        ID:         file,
        Content:    content,
    }
    return
}

func L(s string) string {
    s = formatLangString(s)

    if r := lang[s]; r != "" {
        return r
    }
    // && len(lang) > 0 
    return s
}