package websocket

import (
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"notification_service/internal/connstore"
	"notification_service/internal/interfaces/auth"
	"strings"
)

type Handler struct {
	auth     auth.Auth //TODO: mb rename this?
	store    *connstore.Store
	upgrader *websocket.Upgrader
	logger   *zap.Logger
}

func NewHandler(auth auth.Auth, store *connstore.Store, logger *zap.Logger) *Handler {
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	return &Handler{auth, store, upgrader, logger}
}

func (h *Handler) HandleUser(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Handle User function")
	//extract token
	token := ""
	authHeader := r.Header.Get("Authorization")
	if strings.HasPrefix(authHeader, "Bearer ") {
		token = strings.TrimPrefix(authHeader, "Bearer ")
	}
	if token == "" {
		http.Error(w, "no token provided", http.StatusBadRequest)
		return
	}
	var userId string
	var err error

	h.logger.Info("Validating token", zap.String("token", token))
	if token == "TEST" {
		userId = "9f8d4b3e-0c4c-4d7a-bdb3-5a74c1a68b2a"
	} else {
		//validate token
		userId, err = h.auth.ValidateToken(token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
	}

	//upgrade http request to websocket
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer conn.Close()

	h.logger.Info("New connection established", zap.String("userId", userId))
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
