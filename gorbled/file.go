package gorbled

import (
	"appengine"
	"appengine/blobstore"
	"appengine/datastore"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

/*
 * File data struct
 */
type File struct {
	ID   string
	Key  appengine.BlobKey
	Type string
	Name string
	Size int64
	Date time.Time

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
	if len(files) > 0 {
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

func getFilesAndNav(pageId, pageSize int,
	c appengine.Context) (files []File, nav PageNav, err error) {

	// Get offset and page nav
	dbQuery := datastore.NewQuery("File")
	count, _ := dbQuery.Count(c)
	offset, nav := getPageNav(count, pageId, pageSize, c)

	// Get file data
	dbQuery = dbQuery.Order("-Date").Offset(offset).Limit(pageSize)
	_, err = dbQuery.GetAll(c, &files)

	return
}

type Message struct {
	Success bool
	Info    string
	Data    string
}

func (m *Message) encode() string {
	b, _ := json.Marshal(m)
	return string(b)
}

func (f *File) encode() string {
	b, _ := json.Marshal(f)
	return string(b)
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
	errInfo := ""

	num, err := strconv.Atoi(getUrlVar(r, "num"))
	if err != nil {
		errInfo = errInfo + fmt.Sprint(err)
	}

	uploadURLs := make([]string, num)

	for i := 0; i < num; i++ {
		uploadURL, _ := blobstore.UploadURL(c, "/admin/file/upload", nil)

		uploadURLs[i] = uploadURL.String()
		if err != nil {
			errInfo = errInfo + fmt.Sprint(err)
		}
	}

	b, err := json.Marshal(uploadURLs)
	if err != nil {
		errInfo = errInfo + fmt.Sprint(err)
	}

	if errInfo != "" {
		m.Success = false
		m.Info = "Error: " + errInfo
	} else {
		m.Success = true
		m.Info = "UploadURL Get!"
		m.Data = string(b)
	}

	fmt.Fprint(w, m.encode())
}

func handleFileUpload(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	blobs, _, _ := blobstore.ParseUpload(r)

	blobInfo, _ := blobstore.Stat(c, blobs["file"][0].BlobKey)

	config.FileID = config.FileID + 1
	config.update(configKey, c)

	file := File{
		ID:   fmt.Sprint(config.FileID),
		Key:  blobInfo.BlobKey,
		Type: blobInfo.ContentType,
		Date: blobInfo.CreationTime,
		Name: blobInfo.Filename,
		Size: blobInfo.Size,
	}

	file.save(c)
}

func handleFileDelete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	m := new(Message)

	if id := getUrlVar(r, "id"); id == "" {
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

	if id := getUrlVar(r, "id"); id == "" {
		m.Success = false
		m.Info = "Error: enpty id"
	} else if file, key, err := getFile(id, c); err != nil {
		m.Success = false
		m.Info = "Error: " + fmt.Sprint(err)
	} else {
		name := r.FormValue("name")
		id := r.FormValue("id")
		description := r.FormValue("description")

		if id != file.ID {
			file.ID = id
		}
		file.Name = name
		file.Description = description

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
	pageId, _ := strconv.Atoi(getUrlVar(r, "pid"))
	pageSize := config.AdminFiles

	// Get file data
	if files, nav, err := getFilesAndNav(pageId, pageSize, c); err != nil {
		m.Success = false
		m.Info = "Error: " + fmt.Sprint(err)
	} else if len(files) == 0 {
		m.Success = false
		m.Info = "No File Get!"
	} else {
		type Data struct {
			Files []File
			Nav   PageNav
		}
		data := Data{
			Files: files,
			Nav:   nav,
		}

		b, _ := json.Marshal(data)
		m.Success = true
		m.Info = " Files Date Get! "
		m.Data = string(b)
	}

	fmt.Fprint(w, m.encode())
}

func handleFileList(w http.ResponseWriter, r *http.Request) {
	//initSystem(r)

	// New Page
	page := Page{
		"Title":  "File Manager",
		"Config": config,
	}

	// Render page
	page.Render("admin/files", w)
}

func handleFileGet(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	file, _, _ := getFile(getUrlVar(r, "id"), c)
	blobstore.Send(w, appengine.BlobKey(file.Key))
}
