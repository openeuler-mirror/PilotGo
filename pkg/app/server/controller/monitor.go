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
 * Date: 2022-02-18 13:03:16
 * LastEditTime: 2022-04-06 18:57:09
 * Description: 与prometheus进行对接.
 ******************************************************************************/
package controller

import (
	"bufio"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	sconfig "openeluer.org/PilotGo/PilotGo/pkg/app/server/config"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

func Query(c *gin.Context) {
	target := "http://" + sconfig.Config().Monitor.PrometheusAddr //转向的host
	remote, err := url.Parse(target)
	if err != nil {
		logger.Info("parse", err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	c.Request.URL.Path = "api/v1/query" //请求API
	proxy.ServeHTTP(c.Writer, c.Request)
}
func QueryRange(c *gin.Context) {
	target := "http://" + sconfig.Config().Monitor.PrometheusAddr //转向的host
	remote, err := url.Parse(target)
	if err != nil {
		logger.Info("parse", err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	c.Request.URL.Path = "api/v1/query_range" //请求API
	proxy.ServeHTTP(c.Writer, c.Request)
}

func ListenALert(c *gin.Context) {
	target := "http://" + sconfig.Config().Monitor.PrometheusAddr //转向的host
	remote, err := url.Parse(target)
	if err != nil {
		logger.Info("parse", err)
		return
	}
	proxy := httputil.NewSingleHostReverseProxy(remote)
	c.Request.URL.Path = "/api/v1/alerts" //请求API
	proxy.ServeHTTP(c.Writer, c.Request)
}

// func ListenALert(c *gin.Context) {
// 	resp, err := http.Get("http://" + sconfig.Config().Monitor.PrometheusAddr + "/api/v1/alerts")
// 	if err != nil {
// 		response.Response(c, http.StatusOK,
// 			422,
// 			nil,
// 			err.Error())
// 		return
// 	}
// 	defer resp.Body.Close()
// 	body, _ := ioutil.ReadAll(resp.Body)
// 	if err != nil {
// 		response.Response(c, http.StatusOK,
// 			422,
// 			nil,
// 			err.Error())
// 		return
// 	}
// 	var result Alert
// 	err = json.Unmarshal(body, &result)
// 	logger.Debug("%v", result)
// 	if err != nil {
// 		response.Response(c, http.StatusOK,
// 			422,
// 			nil,
// 			err.Error())
// 		return
// 	}

// 	res := []Alertmanager{}
// 	for _, value := range result.Data.Alerts {
// 		Al := Alertmanager{
// 			value.Labels.Alertname,
// 			value.Labels.Instance,
// 			value.Labels.Job,
// 			value.Annotations.Summary,
// 			value.State,
// 			value.ActiveAt,
// 		}
// 		res = append(res, Al)
// 	}
// 	response.JSON(c, http.StatusOK, http.StatusOK, res, "Alerts have got!")
// }

func InitPromeYml() error {
	FilePath := "/root/prometheus.yml"
	os.Remove(FilePath)
	os.Create(FilePath)
	var prometheusYml Prometheusyml
	prometheusYml.Global.ScrapeInterval = "15s"
	prometheusYml.Global.EvaluationInterval = "15s"
	prometheusYml.RuleFiles = []string{"/etc/prometheus/alert.rules"}

	file, err := os.OpenFile(FilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error("文件打开失败" + err.Error())
		return err
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(
		"global:" +
			"\n scrape_interval: " + prometheusYml.Global.ScrapeInterval +
			"\n evaluation_interval: " + prometheusYml.Global.EvaluationInterval + "\n")
	write.WriteString("rule_files:")
	for _, value := range prometheusYml.RuleFiles {
		write.WriteString("\n - " + value)
	}
	write.WriteString("\nscrape_configs:")
	write.WriteString("\n  - job_name: 'file_sd_test'")
	write.WriteString("\n    file_sd_configs:")
	write.WriteString("\n    - files:")
	write.WriteString("\n      - '/root/file_sd/file_sd.yml'")
	write.WriteString("\n      refresh_interval: 20s")
	write.Flush()
	return nil
}
func WriteYml(a []map[string]string) error {
	FilePath := "/root/file_sd/file_sd.yml"
	os.Remove(FilePath)
	os.Create(FilePath)
	var prometheusYml Prometheusyml
	var tmp static_configs
	for _, value := range a {
		for key, value2 := range value {
			tmp.JobName = key
			x := make([]target, 0)
			x = append(x, target{[]string{value2}})
			tmp.StaticConfigs = x
			prometheusYml.ScrapeConfigs = append(prometheusYml.ScrapeConfigs, tmp)
		}
	}
	file, err := os.OpenFile(FilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		logger.Error("文件打开失败" + err.Error())
		return err
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString("- targets:")
	for _, value := range prometheusYml.ScrapeConfigs {
		for _, value2 := range value.StaticConfigs {
			a := strings.TrimSpace("- '" + value2.Targets[0] + ":9100'")
			write.WriteString("\n      " + a)
		}
	}
	write.Flush()
	return nil
}

// scrape_configs:
//   - job_name: "file_sd_test"
//     file_sd_configs:
//             - files:
//               - '/opt/prometheus/file_sd/test.yml'

// # Alertmanager configuration
// alerting:
//   alertmanagers:
//   - static_configs:
//       - targets:
//          - 192.168.217.131:9093
//           # - alertmanager:9093

type Alert struct {
	Status string `json:"status"`
	Data   struct {
		Alerts []struct {
			Labels struct {
				Alertname string `json:"alertname"`
				Instance  string `json:"instance"`
				Job       string `json:"job"`
				Severity  string `json:"severity"`
			} `json:"labels"`
			Annotations struct {
				Description string `json:"description"`
				Summary     string `json:"summary"`
			} `json:"annotations"`
			State    string    `json:"state"`
			ActiveAt time.Time `json:"activeAt"`
			Value    string    `json:"value"`
		} `json:"alerts"`
	} `json:"data"`
}

type Alertmanager struct {
	Alertname   string    `json:"alertname"`
	Instance    string    `json:"instance"`
	Job         string    `json:"job"`
	Annotations string    `json:"annotations"`
	State       string    `json:"state"`
	ActiveAt    time.Time `json:"activeAt"`
}

type Prometheusyml struct {
	Global struct {
		ScrapeInterval     string `yaml:"scrape_interval"`
		EvaluationInterval string `yaml:"evaluation_interval"`
	} `yaml:"global"`
	RuleFiles     []string         `yaml:"rule_files"`
	ScrapeConfigs []static_configs `yaml:"scrape_configs"`
}
type static_configs struct {
	JobName       string   `yaml:"job_name"`
	StaticConfigs []target `yaml:"static_configs"`
}
type target struct {
	Targets []string `yaml:"targets"`
}
