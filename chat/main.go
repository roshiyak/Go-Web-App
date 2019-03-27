package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"sync"

	"github.com/stretchr/objx"

	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"

	"github.com/stretchr/gomniauth"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var addr = flag.String("addr", ":8080", "アプリケーションのアドレス")
	flag.Parse()
	gomniauth.SetSecurityKey("セキュリティキー")
	gomniauth.WithProviders(
		facebook.New("852354320909-q7clu6du1lnhnuqomdq8qa4055dp0koq.apps.googleusercontent.com", "x-4InCK1I-taZK_x19vDN5pw", "http://localhost:8080/auth/callback/facebook"),
		github.New("852354320909-q7clu6du1lnhnuqomdq8qa4055dp0koq.apps.googleusercontent.com", "x-4InCK1I-taZK_x19vDN5pw", "http://localhost:8080/auth/callback/github"),
		google.New("852354320909-q7clu6du1lnhnuqomdq8qa4055dp0koq.apps.googleusercontent.com", "x-4InCK1I-taZK_x19vDN5pw", "http://localhost:8080/auth/callback/google"),
	)
	r := newRoom()
	http.Handle("/chat", MustAuth(&templateHandler{filename: "chat.html"}))
	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:   "auth",
			Value:  "",
			Path:   "/",
			MaxAge: -1,
		})
		w.Header()["Location"] = []string{"/chat"}
		w.WriteHeader(http.StatusTemporaryRedirect)
	})
	go r.run()
	log.Println("Webサーバーを開始します。ポート：", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
