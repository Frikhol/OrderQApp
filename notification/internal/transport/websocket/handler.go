package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
	"notification_service/internal/connstore"
	"notification_service/internal/interfaces/auth"
	"strings"
)

type Handler struct {
	auth     auth.Auth //TODO: mb rename this?
	store    *connstore.Store
	upgrader *websocket.Upgrader
}

func NewHandler(auth auth.Auth, store *connstore.Store) *Handler {
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	return &Handler{auth, store, upgrader}
}

func (h *Handler) HandleUser(w http.ResponseWriter, r *http.Request) {
	//extract token
	token := ""
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			http.Error(w, "no token provided", http.StatusBadRequest)
			return
		}
	}

	//validate token
	userId, err := h.auth.ValidateToken(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	//upgrade http request to websocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer conn.Close()

	//add user to store
	h.store.Add(userId, conn)
	defer h.store.Remove(userId)

	//read socket for getting him alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}
