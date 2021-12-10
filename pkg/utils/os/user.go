package os

import (
	"bufio"
	"fmt"
	"os/user"
	"strings"

	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	"openeluer.org/PilotGo/PilotGo/pkg/utils"
)

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

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func GetCurrentUserInfo() CurrentUser {
	u, err := user.Current()
	handleErr(err)
	userinfo := CurrentUser{
		Username:  u.Username,
		Userid:    u.Gid,
		GroupName: u.Name,
		Groupid:   u.Uid,
		HomeDir:   u.HomeDir,
	}
	return userinfo
}

func GetAllUserInfo() []AllUserInfo {
	tmp, err := utils.RunCommand("cat /etc/passwd")
	if err != nil {
		logger.Error("获取失败！%s", err)
	}
	reader := strings.NewReader(tmp)
	scanner := bufio.NewScanner(reader)
	var allUsers []AllUserInfo
	for {

		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		line = strings.TrimSpace(line)
		strSlice := strings.Split(line, ":")
		users := AllUserInfo{
			Username:    strSlice[0],
			UserId:      strSlice[2],
			GroupId:     strSlice[3],
			Description: strSlice[4],
			HomeDir:     strSlice[5],
			ShellType:   strSlice[6],
		}
		allUsers = append(allUsers, users)
	}

	return allUsers
}
