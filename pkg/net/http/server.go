package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// target_ip            string //目标主机IP
// promethues_query     string //查询语句
// promethues_starttime string //查询开始时间
// promethues_endtime   string //查询结束时间
// promethues_step      string //查询步长

func promethues_data_geter(c *gin.Context) {
	target_ip := c.PostForm("target_ip")
	promethues_query := c.PostForm("promethues_query")
	promethues_starttime := c.PostForm("promethues_starttime")
	promethues_endtime := c.PostForm("promethues_endtime")
	promethues_step := c.PostForm("promethues_step")
	start := time_parse(promethues_starttime)
	end := time_parse(promethues_endtime)
	url := join_url_param(target_ip, promethues_query, start, end, promethues_step)

	c.JSON(http.StatusOK, Get(url))
}

//发送get请求
func Get(url string) string {
	//超时时间5s
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	//defer resp.Body.Close()
	return string(response)
}

//将当前时间转换成Unix时间戳
func time_parse(times string) string {
	timeLayout := "2006-01-02 15:04:05"
	parse_time, _ := time.Parse(timeLayout, times)
	current_time := parse_time.Unix()
	return fmt.Sprintf("%d", current_time)
}

//将请求的参数加入url中
func join_url_param(target_ip string, query string, start string, end string, step string) string {
	promethues_api_param := fmt.Sprintf("%s/api/v1/query_range?query=%s&start=%s&end=%s&step=%ss",
		target_ip,
		query,
		start,
		end,
		step)
	return promethues_api_param
}

//发送邮件接口
func EmailSender(c *gin.Context) {
	SrcIp := c.PostForm("ip")
	WarningData := c.PostForm("WarningData")
	AlertApi := fmt.Sprintf("%s/api/v1/alerts", WarningData)
	jspull(SrcIp, AlertApi)
}

func jspull(AlertApi string, WarningData string) {
	req, err := http.Post(AlertApi,
		"application/json",
		bytes.NewBuffer([]byte(WarningData)))
	if err != nil {
		fmt.Println("err=", err)
	}
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
}
