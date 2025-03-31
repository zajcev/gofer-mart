package handlers

import (
	"fmt"
	"net/http"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {

}
func LoginUser(w http.ResponseWriter, r *http.Request) {
	user, pass, ok := r.BasicAuth()
	if ok && verifyUserPass(user, pass) {
		fmt.Fprintf(w, "You get to see the secret\n")
	} else {
		w.Header().Set("WWW-Authenticate", `Basic realm="api"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

func verifyUserPass(user, pass string) bool {
	return true
}
