package models

import (
	"time"
)

type Order struct {
	Id            int       `orm:"auto"`
	Client        *User     `orm:"rel(fk)"`
	Agent         *User     `orm:"rel(fk);null"`
	Location      string    `orm:"size(255)"`
	Status        string    `orm:"size(20)"` // "pending", "in_progress", "completed", "cancelled"
	StartTime     time.Time `orm:"type(datetime)"`
	EndTime       time.Time `orm:"type(datetime);null"`
	Price         float64   `orm:"digits(10);decimals(2)"`
	QueuePosition int       `orm:"default(0)"`
	CreatedAt     time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt     time.Time `orm:"auto_now;type(datetime)"`
}

func (o *Order) TableName() string {
	return "orders"
}

type QueueUpdate struct {
	Id        int    `orm:"auto"`
	Order     *Order `orm:"rel(fk)"`
	Position  int
	Timestamp time.Time `orm:"auto_now_add;type(datetime)"`
}

func (q *QueueUpdate) TableName() string {
	return "queue_updates"
}
