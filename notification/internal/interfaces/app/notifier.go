package app

import "notification_service/internal/domain/types"

type Notifier interface {
	Notify(userId string, message types.WSMessage)
}
