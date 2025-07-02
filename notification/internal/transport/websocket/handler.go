package websocket

import (
	"net/http"
	"notification_service/internal/interfaces/auth"
	"notification_service/internal/interfaces/broker"
)

type Handler struct {
	auth   auth.Auth     //TODO: mb rename this?
	broker broker.Broker //TODO: same
}

func NewHandler(auth auth.Auth, broker broker.Broker) *Handler {
	return &Handler{auth, broker}
}

func (h *Handler) HandleUser(w http.ResponseWriter, r *http.Request) {

}
