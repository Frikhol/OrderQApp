package models

import (
	"time"
)

type ErrorLog struct {
	Id        int       `orm:"auto"`
	Timestamp time.Time `orm:"type(timestamp)"`
	Context   string    `orm:"size(255)"`
	Error     string    `orm:"type(text)"`
	Stack     string    `orm:"type(text);null"`
	UserAgent string    `orm:"size(512)"`
	Url       string    `orm:"size(512)"`
	CreatedAt time.Time `orm:"auto_now_add;type(timestamp)"`
}

func (e *ErrorLog) TableName() string {
	return "error_logs"
}
