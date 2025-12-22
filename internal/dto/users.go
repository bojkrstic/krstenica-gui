package dto

import "time"

type User struct {
	ID        int64     `json:"id"`
	Username  string    `json:"username"`
	Role      string    `json:"role"`
	City      string    `json:"city"`
	CreatedAt time.Time `json:"created_at"`
}

type UserCreateReq struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Role     string `json:"role" form:"role"`
	City     string `json:"city" form:"city"`
}

type UserUpdateReq struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Role     string `json:"role" form:"role"`
	City     string `json:"city" form:"city"`
}
