package handlers

import "net/http"

type PingHandler struct {
}

func InitPingHandler() *PingHandler {
	return &PingHandler{}
}

func (p *PingHandler) Ping() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}
}
