package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
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

//loginHandlerはサードパーティへのログイン処理を待ち受け
//パスの形式： /auth/{action}/{provider}
func loginHandler(w http.ResponseWriter, r *http.Request) {
	segs := strings.Split(r.URL.Path, "/")
	action := segs[2]
	provider := segs[3]
	switch action {
	case "login":
		log.Println("TODO: ログイン処理", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "アクション%sには非対応です", action)
	}
}
