package gorbled

import (
    "net/http"
    "time"
    "strconv"

    "appengine"
    "appengine/datastore"
    "appengine/blobstore"
    "fmt"
    "encoding/json"
    "strings"
)

func init() {

}

/*
 * File data struct
 */
type File struct {
    ID          string
    Key         appengine.BlobKey
    Type        string
    Name        string
    Size        int64
    Date        time.Time

    Description string
}

func (f *File) save(c appengine.Context) (err error) {
  _, err = datastore.Put(c, datastore.NewIncompleteKey(c, "File", nil), f)
  return
}

func (f *File) update(key *datastore.Key, c appengine.Context) (err error) {
  _, err = datastore.Put(c, key, f)
  return
}

func getFile(id string, c appengine.Context) (file File, key *datastore.Key, err error) {
    dbQuery := datastore.NewQuery("File").Filter("ID =", id)
    var files []File
    keys, err := dbQuery.GetAll(c, &files)
    if len(files)>0 {
        file = files[0]
        key = keys[0]
    }
    return
}

func getFilesPerPage(offset, pageSize int, c appengine.Context) (files []File, err error) {
  dbQuery := datastore.NewQuery("File").
        Order("-Date").
        Offset(offset).
        Limit(pageSize)
  _, err = dbQuery.GetAll(c, &files)
  return
}


type Message struct {
    Success     bool
    Info        string
    Data        string
}

func (m *Message) encode() (string) {
    b, _ := json.Marshal(m)
    return string(b)
}

func (f *File) encode() (string) {
    b, _ := json.Marshal(f)
    return string(b)
}

func decodeFile(s string) (f File, err error) {
    dec := json.NewDecoder(strings.NewReader(s))
    err = dec.Decode(&f)
    return
}

/*
 * File handler
 */

/*
 * New file upload url
 *
 * @return (string) 
 */
func handleFileNewUrl(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    m := new(Message)

    // Get upload url
    uploadURL, err := blobstore.UploadURL(c, "/admin/file-upload", nil)

    if err != nil {
        m.Success = false
        m.Info = "Error: blobstore.UploadURL"
    } else {
        m.Success = true
        m.Info = "UploadURL Get!"
        m.Data = fmt.Sprint(uploadURL)
    }

    fmt.Fprint(w, m.encode())
}

func handleFileUpload(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    
    blobs, _, _ := blobstore.ParseUpload(r)

    fileInfo := blobs["file"]

    file := File {
        ID:     string(fileInfo[0].BlobKey),
        Key:    fileInfo[0].BlobKey,
        Type:   fileInfo[0].ContentType,
        Date:   fileInfo[0].CreationTime,
        Name:   fileInfo[0].Filename,
        Size:   fileInfo[0].Size,
    }

    file.save(c)
}

func handleFileDelete(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    m := new(Message)
    
    if id := getUrlQuery(r.URL, "id"); id == "" {
        m.Success = false
        m.Info = "Error: empty id"
    } else if file, key, err := getFile(id, c); err != nil {
        m.Success = false
        m.Info = "Error: getFile dbQuery.GetAll"
    } else if err = blobstore.Delete(c, file.Key); err != nil {
        m.Success = false
        m.Info = "Error: blobstore.Delete"
    } else {
        datastore.Delete(c, key)
        m.Success = true
        m.Info = "File Delete!"
    }

    fmt.Fprint(w, m.encode())
}

func handleFileEdit(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    m := new(Message)

    if fileIn, err := decodeFile(r.FormValue("file")); err != nil {
        m.Success = false
        m.Info = "Error: decodeFile"
    } else if file, key, err := getFile(fileIn.ID, c); err != nil {
        m.Success = false
        m.Info = "Error: getFile dbQuery.GetAll"
    } else {
        file.ID = fileIn.ID
        file.Description = fileIn.Description
        file.update(key, c)
        m.Success = true
        m.Info = "File Update!"
    }

    fmt.Fprint(w, m.encode())
}

/*
 * File data per page
 *
 * @return (json) 
 */
func handleFileData(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)
    m := new(Message)

    // Get page id, pageSize
    pageId, _ := strconv.Atoi(getUrlQuery(r.URL, "pid"))
    pageSize  := config.AdminFiles

    // Get offset 
    offset, nav := getPageNav("File", pageId, pageSize, c)

    // Get file data
    if files, err := getFilesPerPage(offset, pageSize, c); err != nil {
        m.Success = false
        m.Info = "Error: " + fmt.Sprint(err)
    } else if len(files) == 0 {
        m.Success = false
        m.Info = "No File Get!"
    } else {
        type Data struct {
            Files       []File 
            Nav         PageNav
        }
        data := Data {
            Files:      files,
            Nav:        nav,
        }

        b, _ := json.Marshal(data)
        m.Success = true
        m.Info = "Files Date Get!"
        m.Data = string(b)
    }

    fmt.Fprint(w, m.encode())
}

func handleFileList(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    // Get page id, pageSize
    pageId, _ := strconv.Atoi(getUrlQuery(r.URL, "pid"))
    pageSize  := config.AdminFiles

    // Get offset and page nav
    offset, nav := getPageNav("File", pageId, pageSize, c)

    // Get file data
    files, err := getFilesPerPage(offset, pageSize, c)
    if err != nil {
        serveError(c, w, err)
        return
    }

    // New Page
    page := Page {
        Title:      "File Manager",
        Files:      files,
        Nav:        nav,
        Config:     config,
    }

    // Render page
    page.Render("admin/files", w)
}

func handleFileGet(w http.ResponseWriter, r *http.Request) {
    blobstore.Send(w, appengine.BlobKey(r.FormValue("key")))
}
