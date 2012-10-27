package gorbled

import (
    "net/http"

    "appengine"
    "appengine/user"
)

type UserData struct {
    Nickname, Email string
    IsAdmin bool
}

/*
 * Redirect to user login
 *
 * @param c   (appengine.Context)
 * @param url (string)
 */
func userLogin(w http.ResponseWriter, r *http.Request, c appengine.Context, url string) {
    loginUrl, _  := user.LoginURL(c, url)
    http.Redirect(w, r, loginUrl, http.StatusFound)
}

/*
 * Redirect to user logout
 *
 * @param c   (appengine.Context)
 * @param url (string)
 */
func userLogout(w http.ResponseWriter, r *http.Request, c appengine.Context, url string) {
    logoutUrl, _ := user.LogoutURL(c, url)
    http.Redirect(w, r, logoutUrl, http.StatusFound)
}

/*
 * Get user info
 *
 * @param c (appengine.Context)
 *
 * @return (string) 
 */
func getUserInfo(c appengine.Context) *UserData {
    u := user.Current(c)
    if u != nil {
        return &UserData{ Nickname: u.String(), Email: u.Email, IsAdmin: user.IsAdmin(c) }
    }

    return nil
}

func handleUser(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    // Get action
    action := getUrlQuery(r.URL, "action")

    switch action {
        case "login":
            if getUserInfo(c) != nil {
                return
            }
            userLogin(w, r, c, "/")

        case "logout":
            userLogout(w, r, c, "/")
    }
}
