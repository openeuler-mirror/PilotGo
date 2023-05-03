package common

// 获取当前用户信息
type CurrentUser struct {
	Username  string
	Userid    string
	GroupName string
	Groupid   string
	HomeDir   string
}

// 获取所有用户的信息
type AllUserInfo struct {
	Username    string
	UserId      string
	GroupId     string
	Description string
	HomeDir     string
	ShellType   string
}
