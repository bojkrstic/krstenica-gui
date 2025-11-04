package dto

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

type UserCreateReq struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type UserUpdateReq struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}
