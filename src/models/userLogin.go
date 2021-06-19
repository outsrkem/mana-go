package models
// 登录时提交的数据
type UserLoginInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
