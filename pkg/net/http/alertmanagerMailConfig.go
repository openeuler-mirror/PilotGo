/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: wanghao
 * Date: 2021-12-04 15:08:08
 * LastEditTime: 2022-03-04 02:09:57
 * Description: 实现alertmanager邮箱的动态配置
 ******************************************************************************/
package http

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

// Alertmanager配置文件结构
type AlertmanagerYaml struct {
	Global    Global      `yaml:"global"`
	Route     Route       `yaml:"route"`
	Receivers []Receivers `yaml:"receivers"`
}

// Global
type Global struct {
	SmtpSmarthost    string `yaml:"smtp_smarthost"`
	SmtpFrom         string `yaml:"smtp_from"`
	SmtpAuthUsername string `yaml:"smtp_auth_username"`
	SmtpAuthPassword string `yaml:"smtp_auth_password"`
	ResolveTimeout   string `yaml:"resolve_timeout"`
	SmtpRequireTls   bool   `yaml:"smtp_require_tls"`
}

// Route
type Route struct {
	GroupBy        []string `yaml:"group_by"`
	GroupWait      string   `yaml:"group_wait"`
	GroupInterval  string   `yaml:"group_interval"`
	RepeatInterval string   `yaml:"repeat_interval"`
	Receiver       string   `yaml:"receiver"`
}

// Receivers
type Receivers struct {
	Name         string         `yaml:"name"`
	EmailConfigs []EmailConfigs `yaml:"email_configs"`
}

// EmailConfigs
type EmailConfigs struct {
	To string `yaml:"to"`
}

//写邮箱配置文件
func WriteToYaml(email []string) {
	FilePath := "C:/Users/王昊/go/alertmanager.yml"
	os.Remove(FilePath)
	os.Create(FilePath)
	emailstr := strings.Join(email, ", ")
	var alertYml AlertmanagerYaml
	alertYml.Global.SmtpSmarthost = "'smtp.qq.com:465'"
	fmt.Println(alertYml.Global.SmtpSmarthost)
	alertYml.Global.SmtpFrom = "'157309081@qq.com'"
	alertYml.Global.SmtpAuthUsername = "'157309081@qq.com'"
	alertYml.Global.SmtpAuthPassword = "'pklazwfbpvkucbda'"
	alertYml.Global.ResolveTimeout = "5m"
	alertYml.Global.SmtpRequireTls = false
	alertYml.Route.GroupBy = []string{"alertname"}
	alertYml.Route.GroupWait = "5s"
	alertYml.Route.GroupInterval = "5s"
	alertYml.Route.RepeatInterval = "5s"
	alertYml.Route.Receiver = "'mail'"
	emailConfig := []EmailConfigs{{emailstr}}
	alertYml.Receivers = []Receivers{{"'mail'", emailConfig}}
	// data, err := yaml.Marshal(&alertYml)
	// fmt.Println(string(data))
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// err = ioutil.WriteFile("C:/Users/王昊/go/alertmanager.yml", data, 0777)
	// if err != nil {
	// 	fmt.Println(err)
	// }
	file, err := os.OpenFile(FilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(
		"global:" +
			"\n smtp_smarthost: " + alertYml.Global.SmtpSmarthost +
			"\n smtp_from: " + alertYml.Global.SmtpFrom +
			"\n smtp_auth_username: " + alertYml.Global.SmtpAuthUsername +
			"\n smtp_auth_password: " + alertYml.Global.SmtpAuthPassword +
			"\n resolve_timeout: " + alertYml.Global.ResolveTimeout +
			"\n smtp_require_tls: " + strconv.FormatBool(alertYml.Global.SmtpRequireTls) + "\n")
	write.WriteString("route:\n group_by: [")
	for key, value := range alertYml.Route.GroupBy {
		if key == len(alertYml.Route.GroupBy)-1 && key != 0 {
			write.WriteString(" '" + value + "'")
		} else if key == 0 && len(alertYml.Route.GroupBy) > 1 {
			write.WriteString("'" + value + "',")
		} else if key == 0 && len(alertYml.Route.GroupBy) == 1 {
			write.WriteString("'" + value + "'")
		} else {
			write.WriteString(" '" + value + "',")
		}
	}
	write.WriteString("]")
	write.WriteString("\n group_interval: " + alertYml.Route.GroupInterval)
	write.WriteString("\n repeat_interval: " + alertYml.Route.RepeatInterval)
	write.WriteString("\n receiver: " + alertYml.Route.Receiver)
	write.WriteString("\nreceivers:")
	for key := range alertYml.Receivers {
		write.WriteString(
			"\n- name: " + alertYml.Receivers[key].Name +
				"\n  email_configs:" +
				"\n  - to: '" + alertYml.Receivers[key].EmailConfigs[key].To + "'")

	}
	write.Flush()
}

//sftp连接函数
func connect(user, password, host string, port int) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

//将邮箱的配置文件传输给alertmanager
func SendFile() {
	sftpClient, err := connect(
		"root",
		"wang13820205036@",
		"192.168.217.22",
		22)
	if err != nil {
		log.Fatal(err)
	}
	defer sftpClient.Close()
	sftpClient.Remove("/root/3/alertmanager.yml")
	var localFilePath = "C:/Users/王昊/go/alertmanager.yml"
	var remoteDIr = "/root/3/"
	srcFile, err := os.Open(localFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer srcFile.Close()
	dstFile, err := sftpClient.Create(remoteDIr + "alertmanager.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer dstFile.Close()
	buf := []byte{0}
	for {
		n, _ := srcFile.Read(buf)
		if n == 0 {
			break
		}
		dstFile.Write(buf)
	}
}

//alertmanager邮箱配置热启动
func ConfigReload() {
	response, _ := http.PostForm("http://192.168.217.22:9093/-/reload", url.Values{})
	fmt.Println(response)
}
