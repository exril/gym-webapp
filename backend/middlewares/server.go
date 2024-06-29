package middlewares

import (
	"net/http"

	rstore "github.com/rbcervilla/redisstore/v8"
)

var store *rstore.RedisStore

func HasBeenAuthenticated(w http.ResponseWriter, r *http.Request) bool {
	session, _ := store.Get(r, "session_token")
	a, _ := session.Values["authenticated"]

	if a == nil {
		return false
	}

	return a.(bool)
}

// storeAuthenticated to store authenticated value
func StoreAuthenticated(w http.ResponseWriter, r *http.Request, v bool) {
	session, _ := store.Get(r, "session_token")

	session.Values["authenticated"] = v
	err := session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
