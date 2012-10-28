package gorbled

import (
    "net/http"
    "strings"
    "appengine"
    "appengine/user"
)

type User struct {
    Nickname    string
    Email       string
    IsAdmin     bool
    IsLogin     bool
}

/*
 * Login and Redirect to previous page
 *
 */
func handleUserLogin(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    url := r.Referer()
    loginUrl, err := user.LoginURL(c, url)
    if err != nil {
        serveError(w, err)
        return
    }
    
    http.Redirect(w, r, loginUrl, http.StatusFound)
}

/*
 * Logout and Redirect to previous page
 *
 */
func handleUserLogout(w http.ResponseWriter, r *http.Request) {
    c := appengine.NewContext(r)

    url := r.Referer()
    if strings.Contains(r.URL.Path, "admin") {
      url = "/"
    }
    logoutUrl, err := user.LogoutURL(c, url)
    if err != nil {
        serveError(w, err)
        return
    }

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
        return User { 
                Nickname: u.String(), 
                Email: u.Email, 
                IsAdmin: user.IsAdmin(c), 
                IsLogin: true,
            }
    }

    return User{}
}
