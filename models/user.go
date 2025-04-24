package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id        int       `orm:"auto"`
	Email     string    `orm:"size(100);unique"`
	Password  string    `orm:"size(100)"`
	Role      string    `orm:"size(20)"` // "client" or "agent"
	CreatedAt time.Time `orm:"auto_now_add;type(datetime)"`
	UpdatedAt time.Time `orm:"auto_now;type(datetime)"`
}

func (u *User) TableName() string {
	return "users"
}

func (u *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		return err
	}
	u.Password = string(bytes)
	return nil
}

func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}
