package gorbled

import (
    "net/http"
    "strings"
    "appengine"
    "appengine/user"
)

func init() {
    http.HandleFunc("/login", handleUserLogin)
    http.HandleFunc("/logout", handleUserLogout)
}

type User struct {
    Nickname    string
    Email       string
    IsAdmin     bool
    IsLogin     bool
}

/*
 * Redirect to user login
 *
 * @param c   (appengine.Context)
 * @param url (string)
 */
func handleUserLogin(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    url := r.Referer()
    loginUrl, _  := user.LoginURL(c, url)
    http.Redirect(w, r, loginUrl, http.StatusFound)
}

/*
 * Redirect to user logout
 *
 * @param c   (appengine.Context)
 * @param url (string)
 */
func handleUserLogout(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    url := r.Referer()
    if strings.Contains(url, "admin") {
      url = "/"
    }
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
func getUserInfo(c appengine.Context) User {
    u := user.Current(c)
    if u != nil {
        return User{ Nickname: u.String(), Email: u.Email, IsAdmin: user.IsAdmin(c), IsLogin: true}
    }

    return User{}
}
