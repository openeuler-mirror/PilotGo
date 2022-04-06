package controller

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	sconfig "openeluer.org/PilotGo/PilotGo/pkg/app/server/config"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

type AlertMessage struct {
	Email       []string    `json:"Email"`
	Labels      Labels      `json:"Labels"`
	Annotations Annotations `json:"Annotations"`
	StartsAt    time.Time   `json:"StartsAt"`
	EndsAt      time.Time   `json:"EndsAt"`
}
type Labels struct {
	Alertname string `json:"alertname"`
	IP        string `json:"IP"`
}
type Annotations struct {
	Summary string `json:"summary"`
}

type AlertSentMessage struct {
	Labels      Labels      `json:"Labels"`
	Annotations Annotations `json:"Annotations"`
	StartsAt    time.Time   `json:"StartsAt"`
	EndsAt      time.Time   `json:"EndsAt"`
}

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
func WriteToYaml(email []string) error {
	FilePath := "/root/alertmanager.yml"
	os.Remove(FilePath)
	os.Create(FilePath)
	emailstr := strings.Join(email, ", ")
	var alertYml AlertmanagerYaml
	alertYml.Global.SmtpSmarthost = "'smtp.qq.com:465'"
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

	file, err := os.OpenFile(FilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error("文件打开失败" + err.Error())
		return err
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
	return nil
}

//alertmanager邮箱配置热启动
func ConfigReload(addr string) error {
	response, err := http.PostForm("http://"+addr+"/-/reload", url.Values{})
	logger.Info("%s", response)
	return err
}

func AlertMessageConfig(c *gin.Context) {
	j, err := ioutil.ReadAll(c.Request.Body)
	logger.Info(string(j))
	if err != nil {
		logger.Error("%s", err.Error())
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}
	var AM AlertMessage
	err = json.Unmarshal(j, &AM)
	logger.Info("%+v", AM)

	if err != nil {
		logger.Error("%s", err.Error())
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}
	err = WriteToYaml(AM.Email) //重写alermanager配置文件
	if err != nil {
		logger.Error("%s", err.Error())
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}
	err = ConfigReload(sconfig.Config().Monitor.AlertManagerAddr) //配置alertmanager邮箱热启动
	if err != nil {
		logger.Error("%s", err.Error())
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}
	url := "http://" + sconfig.Config().Monitor.AlertManagerAddr + "/api/v2/alerts"
	res := AlertSentMessage{
		AM.Labels,
		AM.Annotations,
		AM.StartsAt,
		AM.EndsAt,
	}

	r := []AlertSentMessage{}
	r = append(r, res)
	logger.Info("%+v", r)
	jsonStr, err := json.Marshal(r)

	if err != nil {
		logger.Error(err.Error())
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}
	reader := bytes.NewReader(jsonStr)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		logger.Error(err.Error())
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	// err = Post(url, res, "appliction/json")
	// if err != nil {
	// 	logger.Error("%s", err.Error())
	// 	response.Response(c, http.StatusOK,
	// 		422,
	// 		nil,
	// 		err.Error())
	// 	return
	// }
	client := http.Client{}
	resp, err := client.Do(request)
	if err != nil {
		logger.Error(err.Error())
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}
	logger.Info("%+v", resp)
	response.Success(c, nil, "success")
}
func Post(url string, data interface{}, contentType string) error {

	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return err
	}
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	r, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	logger.Info("%s", string(r))
	return nil
}
