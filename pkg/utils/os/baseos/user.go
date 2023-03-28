package baseos

import (
	"bufio"
	"fmt"
	"os/user"
	"strings"

	"openeuler.org/PilotGo/PilotGo/pkg/logger"
	"openeuler.org/PilotGo/PilotGo/pkg/utils"
	"openeuler.org/PilotGo/PilotGo/pkg/utils/os/common"
)

func handleErr(err error) {
	if err != nil {
		fmt.Println(err)
	}
}

func (b *BaseOS) GetCurrentUserInfo() common.CurrentUser {
	u, err := user.Current()
	handleErr(err)
	userinfo := common.CurrentUser{
		Username:  u.Username,
		Userid:    u.Gid,
		GroupName: u.Name,
		Groupid:   u.Uid,
		HomeDir:   u.HomeDir,
	}
	return userinfo
}

func (b *BaseOS) GetAllUserInfo() []common.AllUserInfo {
	tmp, err := utils.RunCommand("cat /etc/passwd")
	if err != nil {
		logger.Error("failed to get passwd: %s", err.Error())
	}
	reader := strings.NewReader(tmp)
	scanner := bufio.NewScanner(reader)
	var allUsers []common.AllUserInfo
	for {

		if !scanner.Scan() {
			break
		}
		line := scanner.Text()
		line = strings.TrimSpace(line)
		strSlice := strings.Split(line, ":")
		users := common.AllUserInfo{
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
	logger.Info("successfully created user: %s", output)

	//下面两个是管道的两端
	//linux可以使用  echo "password" | passwd --stdin username
	//直接更改密码
	output_ps, err := utils.RunCommand("echo \"" + password + "\" | passwd --stdin " + username)
	if err != nil {
		logger.Error(err.Error())
	}
	logger.Info("successfully changed user passwd: %s", output_ps)
	return nil
}

// 删除用户
func (b *BaseOS) DelUser(username string) (string, error) {
	tmp, err := utils.RunCommand(fmt.Sprintf("userdel -r %s", username))
	if err != nil {
		logger.Error("failed to delete user: %s", err.Error())
		return "", fmt.Errorf("failed to delete user: %s", err)
	}
	logger.Info("successfully deleted user: %s", tmp)
	return tmp, nil
}

// chmod [-R] 权限值 文件名
func (b *BaseOS) ChangePermission(permission, file string) (string, error) {
	tmp, err := utils.RunCommand(fmt.Sprintf("chmod %s %s", permission, file))
	if err != nil {
		logger.Error("failed to change file permissions: %s", err.Error())
		return "", err
	}
	logger.Info("successfully changed file permissions: %s", tmp)
	return tmp, nil
}

// chown [-R] 所有者 文件或目录
func (b *BaseOS) ChangeFileOwner(user, file string) (string, error) {
	tmp, err := utils.RunCommand(fmt.Sprintf("chown -R %s %s", user, file))
	if err != nil {
		logger.Error("failed to change file owner!%s", err.Error())
		return "", err
	}
	logger.Info("successfully changed file owner: %s", tmp)
	return tmp, nil
}
