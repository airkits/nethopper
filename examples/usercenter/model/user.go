package model

import (
	"time"
)

//User info struct
type User struct {
	UID      uint64    `form:"uid" json:"uid"`
	AppID    string    `form:"appid" json:"appid"`
	OpenID   string    `form:"openid" json:"openid"`
	UUID     string    `form:"uuid" json:"uuid"`
	Name     string    `form:"name" json:"name"`
	Channel  string    `form:"channel" json:"channel"`
	Avatar   string    `form:"avatar" json:"avatar"`
	Password string    `form:"password" json:"password"`
	Phone    string    `form:"phone" json:"phone"`
	Gender   int       `form:"gender" json:"gender"`
	Age      int       `form:"age" json:"age"`
	Gold     uint64    `form:"gold" json:"gold"`
	Coin     uint64    `form:"coin" json:"coin"`
	LoginAt  time.Time `form:"loginAt" json:"loginAt"`
	LoginIP  time.Time `form:"loginIP" json:"loginIP"`
	CreateAt time.Time `form:"createAt" json:"createAt"`
	Status   int       `form:"status" json:"status"`
}
