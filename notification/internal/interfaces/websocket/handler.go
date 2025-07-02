package websocket

import "net/http"

type Handler interface {
	HandleUser(w http.ResponseWriter, r *http.Request)
}
