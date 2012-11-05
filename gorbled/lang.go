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
    initSystem(r)
    readLang(c)
    initLang(c, config.Language)

    fmt.Fprint(w, lang)
}

func readLang(c appengine.Context) {
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

func initLang(c appengine.Context, l string) {
    langT, _, _ := getLang(l, c)

    dec := json.NewDecoder(strings.NewReader(string(langT.Content)))

    dec.Decode(&lang)

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
            rexp, _ := regexp.Compile(`([io])-([^:]+):(.*)`)
            s := rexp.FindStringSubmatch(line)

            id := formatLangString(s[2])
            sentence := formatLangString(s[3])

            if s[1]=="i" {
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

    return s
}
