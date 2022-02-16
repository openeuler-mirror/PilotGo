package os

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
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
func AddLinuxUser(username, password string) error {
	useradd := exec.Command("useradd", "-m", username)
	err := useradd.Start()
	if err != nil {
		logger.Error(err.Error())
	}

	useradd.Wait()
	//下面两个是管道的两端
	//linux可以使用  echo "password" | passwd --stdin username
	//直接更改密码
	ps := exec.Command("echo", password)
	grep := exec.Command("passwd", "--stdin", username)

	r, w := io.Pipe() // 创建一个管道
	defer r.Close()
	defer w.Close()
	ps.Stdout = w  // ps向管道的一端写
	grep.Stdin = r // grep从管道的一端读

	var buffer bytes.Buffer
	grep.Stdout = &buffer // grep的输出为buffer

	_ = ps.Start()
	_ = grep.Start()
	ps.Wait()
	w.Close()
	grep.Wait()
	io.Copy(os.Stdout, &buffer) // buffer拷贝到系统标准输出
	return nil
}

// 删除用户
func DelUser(username string) (string, error) {
	tmp, err := utils.RunCommand(fmt.Sprintf("userdel -r %s", username))
	if err != nil {
		logger.Error("删除用户失败!%s", err.Error())
		return "", fmt.Errorf("删除用户失败%s", err)
	}
	logger.Info("删除用户成功!%s", tmp)
	return tmp, nil
}

// chmod [-R] 权限值 文件名
func ChangePermission(permission, file string) (string, error) {
	tmp, err := utils.RunCommand(fmt.Sprintf("chmod %s %s", permission, file))
	if err != nil {
		logger.Error("改变文件权限失败!%s", err.Error())
		return "", err
	}
	logger.Info("改变文件权限成功!%s", tmp)
	return tmp, nil
}

// chown [-R] 所有者 文件或目录
func ChangeFileOwner(user, file string) (string, error) {
	tmp, err := utils.RunCommand(fmt.Sprintf("chown -R %s %s", user, file))
	if err != nil {
		logger.Error("改变文件所有者失败!%s", err.Error())
		return "", err
	}
	logger.Info("改变文件所有者成功!%s", tmp)
	return tmp, nil
}
