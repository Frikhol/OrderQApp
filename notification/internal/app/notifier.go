package app

import (
	"encoding/json"
	"notification_service/internal/connstore"
	"notification_service/internal/domain/types"
)

type Notifier struct {
	store *connstore.Store
}

func NewNotifier() *Notifier {
	return &Notifier{}
}

func (n *Notifier) Notify(userId string, message types.WSMessage) {
	conn, ok := n.store.Get(userId)
	if !ok {
		//TODO: implement me
		//send email??
		return
	}

	msg, err := json.Marshal(message)
	if err != nil {
		//TODO: implement me
		return
	}

	err = conn.WriteMessage(1, msg)
	if err != nil {
		//TODO: implement me
		return
	}
}
