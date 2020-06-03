package model

//User info struct
type User struct {
	UID      string `form:"uid" json:"uid"`
	Nickname string `form:"nickname" json:"nickname"`
	Avatar   string `form:"avatar" json:"avatar"`
	Gold     uint64 `form:"gold" json:"gold"`
	Token    string `form:"token" json:"token"`
}
