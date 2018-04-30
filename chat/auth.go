package main

import (
	"github.com/markbates/goth/gothic"
	"net/http"
	"fmt"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// 未承認だった場合
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		log.Info("authHandler.serveHTTP: 未認証です。")
	} else if err != nil {
		// 別の何らかのエラーが発生
		panic(err.Error())
	} else {
		// 認証に成功した場合、ラップされたハンドラを呼び出す。
		h.next.ServeHTTP(w, r)
		log.Info("認証に成功しました。")
	}
}

func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	//if _, err := gothic.CompleteUserAuth(w, r); err != nil {
		//log.Error("loginHandler: ", err)
	//} else {
		// 承認ハンドラを呼び出します。
		gothic.BeginAuthHandler(w, r)
	//}
}

func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	_, err := gothic.CompleteUserAuth(w, r)
	if err != nil {  // 何らかの理由でユーザー認証が完了しなかった。
		log.Warning(fmt.Fprintln(w, err))
		return
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Location", "/")
	w.WriteHeader(http.StatusTemporaryRedirect)
}