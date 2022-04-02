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
 * LastEditTime: 2022-04-02 15:25:05
 * Description: 与prometheus进行对接.
 ******************************************************************************/
package controller

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"openeluer.org/PilotGo/PilotGo/pkg/common/response"
	"openeluer.org/PilotGo/PilotGo/pkg/config"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
)

const (
	CPU             = 1 //CPU使用率
	Memory          = 2 //内存使用率
	IOWrite         = 3 //IO写入速率
	IORead          = 4 //IO读取速率
	NetworkReceive  = 5 //平均入网 (字节)
	NetworkTransmit = 6 //平均出网
)

type Promequeryrange struct {
	Machineip string `json:"machineip"`
	Query     int    `json:"query"`
	Starttime string `json:"starttime"`
	Endtime   string `json:"endtime"`
}
type ReturnPromeCPU struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Instance string `json:"instance"`
			} `json:"metric"`
			Value []interface{} `json:"values"`
		} `json:"result"`
	} `json:"data"`
}
type ReturnPromeMemory struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Instance string `json:"instance"`
				Job      string `json:"job"`
			} `json:"metric"`
			Value []interface{} `json:"values"`
		} `json:"result"`
	} `json:"data"`
}
type ReturnPromeIO struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Device   string `json:"device"`
				Instance string `json:"instance"`
				Job      string `json:"job"`
			} `json:"metric"`
			Value []interface{} `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

type Res struct {
	Time string `json:"time"`
	Res  string `json:"value"`
}
type IO struct {
	Type  string `json:"device"`
	Label []Res  `json:"label"`
}

func Queryrange(c *gin.Context) {
	j, err := ioutil.ReadAll(c.Request.Body)
	fmt.Println("body:", string(j))
	if err != nil {
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}
	var Pqr Promequeryrange
	err = json.Unmarshal(j, &Pqr)
	logger.Info("%+v", Pqr)

	if err != nil {
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}
	conf, err := config.Load()
	if err != nil {
		fmt.Println("failed to load configure, exit..", err)
		os.Exit(-1)
	}
	url, err := JudgeQueryRange(Pqr.Query, conf.Monitor.PrometheusAddr, Pqr.Starttime, Pqr.Endtime)
	if err != nil {
		response.Response(c, http.StatusOK,
			422,
			nil,
			"查询数字输入有误")
		return
	}
	logger.Info(url())
	resp, err := http.Get(url())
	if err != nil {
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("body:", string(body))
	if err != nil {
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}

	switch Pqr.Query {
	case 1:
		var result ReturnPromeCPU
		err = json.Unmarshal(body, &result)
		logger.Info("%v", result)
		if err != nil {
			response.Response(c, http.StatusOK,
				422,
				nil,
				err.Error())
			return
		}
		res := make([]Res, 0)
		var x Res
		for _, value := range result.Data.Result {
			if value.Metric.Instance == Pqr.Machineip {
				for _, value2 := range value.Value {
					a := value2.([]interface{})
					tm := time.Unix(int64(a[0].(float64)), 0)
					x.Time = tm.Format("2006-01-02 15:04:05")
					x.Res = a[1].(string)
					logger.Info("%+v", x)
					res = append(res, x)
				}
			}

		}
		logger.Info("%v", res)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": res,
		})

	case 2:
		var result ReturnPromeMemory
		err = json.Unmarshal(body, &result)
		logger.Info("%v", result)
		if err != nil {
			response.Response(c, http.StatusOK,
				422,
				nil,
				err.Error())
			return
		}
		res := make([]Res, 0)
		var x Res
		for _, value := range result.Data.Result {
			if value.Metric.Instance == Pqr.Machineip {
				for _, value2 := range value.Value {
					a := value2.([]interface{})
					tm := time.Unix(int64(a[0].(float64)), 0)
					x.Time = tm.Format("2006-01-02 15:04:05")
					x.Res = a[1].(string)
					logger.Info("%+v", x)
					res = append(res, x)
				}
			}

		}
		logger.Info("%v", res)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": res,
		})
	case 3, 4, 5, 6:
		var result ReturnPromeIO
		err = json.Unmarshal(body, &result)
		logger.Info("%v", result)
		if err != nil {
			response.Response(c, http.StatusOK,
				422,
				nil,
				err.Error())
			return
		}

		io := make([]IO, 0)
		var x Res
		var i IO
		for _, value := range result.Data.Result {

			if value.Metric.Instance == Pqr.Machineip {
				i.Type = value.Metric.Device
				for _, value2 := range value.Value {
					a := value2.([]interface{})
					tm := time.Unix(int64(a[0].(float64)), 0)
					x.Time = tm.Format("2006-01-02 15:04:05")
					x.Res = a[1].(string)
					logger.Info("%+v", x)
					i.Label = append(i.Label, x)

				}
			} else {
				continue
			}
			io = append(io, i)
			i.Label = []Res{}
		}

		logger.Info("%v", io)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": io,
		})

	default:
		response.Response(c, http.StatusOK,
			422,
			nil,
			"输入查询号有误")

	}

}

func JudgeQueryRange(query int, addr string, start string, end string) (func() string, error) {
	switch query {
	case 1:

		return func() string {
			return JoinUrlParam(addr,
				"100-(avg%20by(instance)(irate(node_cpu_seconds_total{mode=\"idle\"}[5m]))*100)",
				start,
				end,
				"10")
		}, nil
	case 2:
		return func() string {
			return JoinUrlParam(addr,
				"(1-(node_memory_MemAvailable_bytes/(node_memory_MemTotal_bytes)))*100",
				start,
				end,
				"10")
		}, nil

	case 3:
		return func() string {
			return JoinUrlParam(addr,
				"irate(node_disk_writes_completed_total[1m])",
				start,
				end,
				"10")
		}, nil
	case 4:
		return func() string {
			return JoinUrlParam(addr,
				"irate(node_disk_reads_completed_total[1m])",
				start,
				end,
				"10")
		}, nil
	case 5:
		return func() string {
			return JoinUrlParam(addr,
				"irate(node_network_receive_bytes_total[5m])",
				start,
				end,
				"10")
		}, nil
	case 6:
		return func() string {
			return JoinUrlParam(addr,
				"irate(node_network_transmit_bytes_total[5m])",
				start,
				end,
				"10")
		}, nil
	default:
		return func() string { return "" }, fmt.Errorf("获取源软件包名以及源失败")

	}
}

func TimeParse(times string) string {
	timeLayout := "2006-01-02 15:04:05"
	parse_time, _ := time.Parse(timeLayout, times)
	current_time := parse_time.Unix()
	return fmt.Sprintf("%d", current_time)
}
func JoinUrlParam(addr string, query string, start string, end string, step string) string {
	promethues_api_param := fmt.Sprintf("http://%s/api/v1/query_range?query=%s&start=%s&end=%s&step=%ss",
		addr,
		query,
		start,
		end,
		step)
	return promethues_api_param
}

type Promequery struct {
	Machineip string `json:"machineip"`
	Query     int    `json:"query"`
	Time      string `json:"time"`
}

type ReturnPromeCPU2 struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Instance string `json:"instance"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}
type ReturnPromeMemory2 struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Instance string `json:"instance"`
				Job      string `json:"job"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

type ReturnPromeIO2 struct {
	Status string `json:"status"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				Device   string `json:"device"`
				Instance string `json:"instance"`
				Job      string `json:"job"`
			} `json:"metric"`
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}
type IO2 struct {
	Type  string `json:"device"`
	Label Res    `json:"label"`
}

func Query(c *gin.Context) {
	j, err := ioutil.ReadAll(c.Request.Body)
	fmt.Println("body:", string(j))
	if err != nil {
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}
	var Pq Promequery
	err = json.Unmarshal(j, &Pq)
	logger.Info("%+v", Pq)
	if err != nil {
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}
	conf, err := config.Load()
	if err != nil {
		fmt.Println("failed to load configure, exit..", err)
		os.Exit(-1)
	}
	url, err := JudgeQuery(Pq.Query, conf.Monitor.PrometheusAddr, Pq.Time)
	if err != nil {
		response.Response(c, http.StatusOK,
			422,
			nil,
			"查询数字输入有误")
		return
	}
	logger.Info(url())
	resp, err := http.Get(url())
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("body:", string(body))
	if err != nil {
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}
	switch Pq.Query {
	case 1:
		var result ReturnPromeCPU2
		err = json.Unmarshal(body, &result)
		logger.Info("%v", result)
		if err != nil {
			response.Response(c, http.StatusOK,
				422,
				nil,
				err.Error())
			return
		}
		var res Res
		for _, value := range result.Data.Result {
			logger.Info("%d", []byte(value.Metric.Instance))
			logger.Info("%d", []byte(Pq.Machineip))

			if Pq.Machineip == value.Metric.Instance {
				tm := time.Unix(int64(value.Value[0].(float64)), 0)
				res.Time = tm.Format("2006-01-02 15:04:05")
				res.Res = value.Value[1].(string)

				logger.Info("%+v", res)
			}

			logger.Info("%+v", res)
		}
		logger.Info("%v", res)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": res,
		})

	case 2:
		var result ReturnPromeMemory2
		err = json.Unmarshal(body, &result)
		logger.Info("%v", result)
		if err != nil {
			response.Response(c, http.StatusOK,
				422,
				nil,
				err.Error())
			return
		}
		var res Res
		for _, value := range result.Data.Result {
			if value.Metric.Instance == Pq.Machineip {
				tm := time.Unix(int64(value.Value[0].(float64)), 0)
				res.Time = tm.Format("2006-01-02 15:04:05")
				res.Res = value.Value[1].(string)
				// res.Res = res.Res[:4]
				logger.Info("%+v", res)
			}

		}
		logger.Info("%v", res)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": res,
		})

	case 3, 4, 5, 6:
		var result ReturnPromeIO2
		err = json.Unmarshal(body, &result)
		logger.Info("%v", result)
		if err != nil {
			response.Response(c, http.StatusOK,
				422,
				nil,
				err.Error())
			return
		}

		io := make([]IO2, 0)
		var i IO2
		for _, value := range result.Data.Result {

			if value.Metric.Instance == Pq.Machineip {
				i.Type = value.Metric.Device
				tm := time.Unix(int64(value.Value[0].(float64)), 0)
				i.Label.Time = tm.Format("2006-01-02 15:04:05")
				i.Label.Res = value.Value[1].(string)
			}
			io = append(io, i)
		}
		logger.Info("%v", io)
		c.JSON(http.StatusOK, gin.H{
			"code": 200,
			"data": io,
		})
	}

}

func JoinUrlParam2(target_ip string, query string, time string) string {
	promethues_api_param := fmt.Sprintf("http://%s/api/v1/query?query=%s&time=%s",
		target_ip,
		query,
		time,
	)
	return promethues_api_param
}
func JudgeQuery(query int, addr string, time string) (func() string, error) {
	switch query {
	case 1:

		return func() string {
			return JoinUrlParam2(addr,
				"100-(avg%20by(instance)(irate(node_cpu_seconds_total{mode=\"idle\"}[5m]))*100)",
				time)
		}, nil
	case 2:
		return func() string {
			return JoinUrlParam2(addr,
				"(1-(node_memory_MemAvailable_bytes/(node_memory_MemTotal_bytes)))*100",
				time)
		}, nil
	case 3:
		return func() string {
			return JoinUrlParam2(addr,
				"irate(node_disk_writes_completed_total[1m])",
				time)
		}, nil
	case 4:
		return func() string {
			return JoinUrlParam2(addr,
				"irate(node_disk_reads_completed_total[1m])",
				time)
		}, nil
	case 5:
		return func() string {
			return JoinUrlParam2(addr,
				"irate(node_network_receive_bytes_total[5m])",
				time)
		}, nil
	case 6:
		return func() string {
			return JoinUrlParam2(addr,
				"irate(node_network_transmit_bytes_total[5m])",
				time)
		}, nil
	default:
		return func() string { return "" }, fmt.Errorf("输入查询号错误")
	}
}

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

func ListenALert(c *gin.Context) {
	conf, err := config.Load()
	if err != nil {
		fmt.Println("failed to load configure, exit..", err)
		os.Exit(-1)
	}
	// url := conf.S.ServerIP + "/api/v1/alerts"
	// logger.Info("%s", url)
	resp, err := http.Get("http://" + conf.Monitor.PrometheusAddr + "/api/v1/alerts")
	if err != nil {
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("body:", string(body))
	if err != nil {
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}
	var result Alert
	err = json.Unmarshal(body, &result)
	logger.Info("%v", result)
	if err != nil {
		response.Response(c, http.StatusOK,
			422,
			nil,
			err.Error())
		return
	}

	res := []Alertmanager{}
	for _, value := range result.Data.Alerts {
		Al := Alertmanager{
			value.Labels.Alertname,
			value.Labels.Instance,
			value.Labels.Job,
			value.Annotations.Summary,
			value.State,
			value.ActiveAt,
		}
		res = append(res, Al)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": res,
	})

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

func WritePrometheusYml(a []map[string]string) error {
	FilePath := "/root/prometheus.yml"
	os.Remove(FilePath)
	os.Create(FilePath)
	var prometheusYml Prometheusyml
	prometheusYml.Global.ScrapeInterval = "15s"
	prometheusYml.Global.EvaluationInterval = "15s"
	prometheusYml.RuleFiles = []string{"/etc/prometheus/alert.rules"}
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
	write.WriteString(
		"global:" +
			"\n scrape_interval: " + prometheusYml.Global.ScrapeInterval +
			"\n evaluation_interval: " + prometheusYml.Global.EvaluationInterval + "\n")
	write.WriteString("rule_files:")
	for _, value := range prometheusYml.RuleFiles {
		write.WriteString("\n - " + value)
	}
	write.WriteString("\nscrape_configs:")
	for _, value := range prometheusYml.ScrapeConfigs {
		write.WriteString("\n  - job_name: '" + value.JobName + "'")
		write.WriteString("\n    static_configs:")
		for _, value2 := range value.StaticConfigs {
			a := strings.TrimSpace("- targets: ['" + value2.Targets[0] + ":9090']")
			write.WriteString("\n      " + a)
		}
	}
	write.Flush()
	return nil
}
func PrometheusConfigReload(addr string) error {
	response, err := http.PostForm("http://"+addr+"/-/reload", url.Values{})
	logger.Debug("%+v", response)
	return err
}

// # Alertmanager configuration
// alerting:
//   alertmanagers:
//   - static_configs:
//       - targets:
//          - 192.168.217.131:9093
//           # - alertmanager:9093
