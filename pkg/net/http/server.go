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
 * LastEditTime: 2022-03-04 02:13:50
 * Description: promethues监控数据获取
 ******************************************************************************/
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

func PromethuesDataGeter(c *gin.Context) {
	target_ip := c.PostForm("target_ip")
	promethues_query := c.PostForm("promethues_query")
	promethues_starttime := c.PostForm("promethues_starttime")
	promethues_endtime := c.PostForm("promethues_endtime")
	promethues_step := c.PostForm("promethues_step")
	start := TimeParse(promethues_starttime)
	end := TimeParse(promethues_endtime)
	url := JoinUrlParam(target_ip, promethues_query, start, end, promethues_step)

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
func TimeParse(times string) string {
	timeLayout := "2006-01-02 15:04:05"
	parse_time, _ := time.Parse(timeLayout, times)
	current_time := parse_time.Unix()
	return fmt.Sprintf("%d", current_time)
}

//将请求的参数加入url中
func JoinUrlParam(target_ip string, query string, start string, end string, step string) string {
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
	JsPull(SrcIp, AlertApi)
}

func JsPull(AlertApi string, WarningData string) {
	req, err := http.Post(AlertApi,
		"application/json",
		bytes.NewBuffer([]byte(WarningData)))
	if err != nil {
		fmt.Println("err=", err)
	}
	body, _ := ioutil.ReadAll(req.Body)
	fmt.Println(string(body))
}
