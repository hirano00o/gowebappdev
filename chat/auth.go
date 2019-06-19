package main

import "net/http"

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if _, err := r.Cookie("auth"); err == http.ErrNoCookie {
		// 未検証
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
	} else if err != nil {
		panic(err.Error())
	} else {
		// 成功
		h.next.ServeHTTP(w, r)
	}
}

// MustAuth is ...
func MustAuth(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}
