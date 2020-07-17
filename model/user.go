package model

import (
	"time"
)

// User 用户模型
type Users struct {
	Id         int       `gorm:"primary_key" json:"id"`
	UserName   string    `json:"user_name"`
	Password   string    `json:"password"`
	AppKey     string    `json:"app_key"`
	SecretKey  string    `json:"secret_key"`
	CreateTime time.Time `json:"create_time"`
}

// 获取secretKey
func (u *Users) GetUser(appKey string) *Users {
	var user Users
	if DB.Where("app_key = ?", appKey).First(&user).RecordNotFound() {
		return nil
	}
	return &user
}
