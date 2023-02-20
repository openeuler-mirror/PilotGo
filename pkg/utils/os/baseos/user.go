package baseos

import (
	"bufio"
	"fmt"
	"os/user"
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
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

func (b *BaseOS) GetCurrentUserInfo() CurrentUser {
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

func (b *BaseOS) GetAllUserInfo() []AllUserInfo {
	tmp, err := utils.RunCommand("cat /etc/passwd")
	if err != nil {
		logger.Error("获取失败!%s", err.Error())
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

// 创建新的用户，并新建家目录
func (b *BaseOS) AddLinuxUser(username, password string) error {
	output, err := utils.RunCommand("useradd -m " + username)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Info("创建用户成功!%s", output)

	//下面两个是管道的两端
	//linux可以使用  echo "password" | passwd --stdin username
	//直接更改密码
	output_ps, err := utils.RunCommand("echo \"" + password + "\" | passwd --stdin " + username)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Info("更改用户密码成功!%s", output_ps)
	return nil
}

// 删除用户
func (b *BaseOS) DelUser(username string) (string, error) {
	tmp, err := utils.RunCommand(fmt.Sprintf("userdel -r %s", username))
	if err != nil {
		logger.Error("删除用户失败!%s", err.Error())
		return "", fmt.Errorf("删除用户失败%s", err)
	}
	logger.Info("删除用户成功!%s", tmp)
	return tmp, nil
}

// chmod [-R] 权限值 文件名
func (b *BaseOS) ChangePermission(permission, file string) (string, error) {
	tmp, err := utils.RunCommand(fmt.Sprintf("chmod %s %s", permission, file))
	if err != nil {
		logger.Error("改变文件权限失败!%s", err.Error())
		return "", err
	}
	logger.Info("改变文件权限成功!%s", tmp)
	return tmp, nil
}

// chown [-R] 所有者 文件或目录
func (b *BaseOS) ChangeFileOwner(user, file string) (string, error) {
	tmp, err := utils.RunCommand(fmt.Sprintf("chown -R %s %s", user, file))
	if err != nil {
		logger.Error("改变文件所有者失败!%s", err.Error())
		return "", err
	}
	logger.Info("改变文件所有者成功!%s", tmp)
	return tmp, nil
}
