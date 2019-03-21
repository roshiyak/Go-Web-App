package main

import (
	"net/http"
)

type authHabdler struct {
	next http.Handler
}

func (h *authHabdler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		//unauthrized
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		//Ocation other error
		panic(err.Error())
	} else {
		//Success!
		h.next.ServeHTTP(w, r)
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHabdler{next: handler}
}
