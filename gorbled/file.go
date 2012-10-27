package gorbled

import (
    "net/http"

    "appengine"
    "appengine/blobstore"
)

type FileData struct {
    PostUrl, BlobKey string
    Success bool
}

func handleFile(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    // Get action && id
    action := getUrlQuery(r.URL, "action")

    switch action {
        case "upload":
            // Check user permissions
            userInfo := getUserInfo(c)
            if userInfo == nil || !userInfo.IsAdmin {
                serve404(w)
                return
            }

            fileUpload(w, r)

        case "uploadProcess":
            // Check user permissions
            userInfo := getUserInfo(c)
            if userInfo == nil || !userInfo.IsAdmin {
                serve404(w)
                return
            }

            fileUploadProcess(w, r)

        case "get":
            f := getUrlQuery(r.URL, "file")
            fileGet(w, r, f)
    }
}

func fileUpload(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    // New FileData
    fileData := new(FileData)

    // Get upload url
    uploadURL, err := blobstore.UploadURL(c, "/file?action=uploadProcess", nil)
    if err != nil {
        serveError(c, w, err)
        return
    }

    fileData.PostUrl = uploadURL.String()

    // New PageSetting
    pageSetting := new(PageSetting)

    // Setting pageSetting
    pageSetting.Title       = "File Upload - " + config.Title
    pageSetting.Layout      = "column1"
    pageSetting.ShowSidebar = false

    // New PageData
    pageData := &PageData{ File: *fileData }

    // New Page
    page := NewPage(pageSetting, pageData)

    // Render page
    page.Render("file/upload", w)
}

func fileUploadProcess(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    blobs, _, err := blobstore.ParseUpload(r)
    if err != nil {
        serveError(c, w, err)
        return
    }

    file := blobs["file"]

    if len(file) == 0 {
        http.Redirect(w, r, "/file?action=upload", http.StatusFound)
        return
    }

    // New FileData
    fileData := FileData{ Success: true, BlobKey: string(file[0].BlobKey) }

    // New PageSetting
    pageSetting := new(PageSetting)

    // Setting pageSetting
    pageSetting.Title       = "File Upload - " + config.Title
    pageSetting.Layout      = "column1"
    pageSetting.ShowSidebar = false

    // New PageData
    pageData := &PageData{ File: fileData }

    // New Page
    page := NewPage(pageSetting, pageData)

    // Render page
    page.Render("file/upload", w)
}

func fileGet(w http.ResponseWriter, r *http.Request, file string) {
    blobstore.Send(w, appengine.BlobKey(file))
}
