package service

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

	sconfig "openeuler.org/PilotGo/PilotGo/pkg/app/server/config"
	"openeuler.org/PilotGo/PilotGo/pkg/app/server/model"
	"openeuler.org/PilotGo/PilotGo/pkg/logger"
)

//写邮箱配置文件
func WriteToYaml(email []string) error {
	FilePath := "/root/alertmanager.yml"
	os.Remove(FilePath)
	os.Create(FilePath)
	emailstr := strings.Join(email, ", ")
	var alertYml model.AlertmanagerYaml
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
	emailConfig := []model.EmailConfigs{{To: emailstr}}
	alertYml.Receivers = []model.Receivers{{Name: "'mail'", EmailConfigs: emailConfig}}

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

func AlertMessageConfig(j []byte) error {
	var AM model.AlertMessage
	err := json.Unmarshal(j, &AM)
	logger.Info("%+v", AM)

	if err != nil {
		return err
	}
	err = WriteToYaml(AM.Email) //重写alermanager配置文件
	if err != nil {
		return err
	}
	err = ConfigReload(sconfig.Config().Monitor.AlertManagerAddr) //配置alertmanager邮箱热启动
	if err != nil {
		return err
	}
	url := "http://" + sconfig.Config().Monitor.AlertManagerAddr + "/api/v2/alerts"
	res := model.AlertSentMessage{
		Labels:      AM.Labels,
		Annotations: AM.Annotations,
		StartsAt:    AM.StartsAt,
		EndsAt:      AM.EndsAt,
	}

	r := []model.AlertSentMessage{}
	r = append(r, res)
	logger.Info("%+v", r)
	jsonStr, err := json.Marshal(r)

	if err != nil {
		return err
	}
	reader := bytes.NewReader(jsonStr)
	request, err := http.NewRequest("POST", url, reader)
	if err != nil {
		return err
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
		return err
	}
	logger.Info("%+v", resp)
	return nil
}
