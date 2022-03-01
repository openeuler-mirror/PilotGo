/******************************************************************************
 * Copyright (c) KylinSoft Co., Ltd.2021-2022. All rights reserved.
 * PilotGo is licensed under the Mulan PSL v2.
 * You can use this software accodring to the terms and conditions of the Mulan PSL v2.
 * You may obtain a copy of Mulan PSL v2 at:
 *     http://license.coscl.org.cn/MulanPSL2
 * THIS SOFTWARE IS PROVIDED ON AN 'AS IS' BASIS, WITHOUT WARRANTIES OF ANY KIND,
 * EITHER EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO NON-INFRINGEMENT, MERCHANTABILITY OR FIT FOR A PARTICULAR PURPOSE.
 * See the Mulan PSL v2 for more details.
 * Author: zhanghan
 * Date: 2022-02-18 02:33:55
 * LastEditTime: 2022-03-01 11:15:14
 * Description: provide server service functions.
 ******************************************************************************/
package agentmanager

import (
	"net"
	"time"

	"github.com/google/uuid"
	"openeluer.org/PilotGo/PilotGo/pkg/logger"
	pnet "openeluer.org/PilotGo/PilotGo/pkg/net"
	"openeluer.org/PilotGo/PilotGo/pkg/protocol"
)

type AgentMessageHandler func(*Agent, *protocol.Message) error

type Agent struct {
	UUID             string
	Version          string
	conn             net.Conn
	MessageProcesser *protocol.MessageProcesser
}

// 通过给定的conn连接初始化一个agent并启动监听
func NewAgent(conn net.Conn) (*Agent, error) {
	agent := &Agent{
		UUID:             "agent",
		conn:             conn,
		MessageProcesser: protocol.NewMessageProcesser(),
	}

	go func() {
		agent.startListen()
	}()

	if err := agent.Init(); err != nil {
		return nil, err
	}

	return agent, nil
}

func (a *Agent) bindHandler(t int, f AgentMessageHandler) {
	a.MessageProcesser.BindHandler(t, func(c protocol.MessageContext, msg *protocol.Message) error {
		return f(c.(*Agent), msg)
	})
}

func (a *Agent) startListen() {
	defer func() {
		if err := recover(); err != nil {
			logger.Error("server processor panic error:%s", err.(error).Error())
			a.conn.Close()
		}
	}()

	readBuff := []byte{}
	for {
		buff := make([]byte, 1024)
		n, err := a.conn.Read(buff)
		if err != nil {
			logger.Error("read error:%s", err)
			DeleteAgent(a.UUID)
			return
		}
		readBuff = append(readBuff, buff[:n]...)

		// 切割frame
		i, f := protocol.TlvDecode(&readBuff)
		if i != 0 {
			readBuff = readBuff[i:]
			go func() {
				msg := protocol.ParseMessage(*f)
				a.MessageProcesser.ProcessMessage(a, msg)
			}()
		}
	}
}

// 远程获取agent端的信息进行初始化
func (a *Agent) Init() error {
	// TODO: 此处绑定所有的消息处理函数
	a.bindHandler(protocol.Heartbeat, func(a *Agent, msg *protocol.Message) error {
		logger.Info("process heartbeat from processor, remote addr:%s, data:%s",
			a.conn.RemoteAddr().String(), msg.String())
		return nil
	})

	a.bindHandler(protocol.AgentInfo, func(a *Agent, msg *protocol.Message) error {
		logger.Info("process heartbeat from processor, remote addr:%s, data:%s",
			a.conn.RemoteAddr().String(), msg.String())
		return nil
	})

	data, err := a.AgentInfo()
	if err != nil {
		logger.Error("fail to get agent info, address:%s", a.conn.RemoteAddr().String())
	}
	d := data.(map[string]interface{})
	logger.Debug("response agent info is %v", d)
	a.UUID = d["agent_uuid"].(string)

	a.Version = d["agent_version"].(string)

	// add agent to agent manager
	AddAgent(a)

	return nil
}

// 远程在agent上运行脚本
func (a *Agent) RunScript(cmd string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.RunScript,
		Data: struct {
			Command string
		}{
			Command: cmd,
		},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

func (a *Agent) sendMessage(msg *protocol.Message, wait bool, timeout time.Duration) (*protocol.Message, error) {
	logger.Debug("send message:%s", msg.String())

	if msg.UUID == "" {
		msg.UUID = uuid.New().String()
	}

	if wait {
		waitChan := make(chan *protocol.Message)
		a.MessageProcesser.WaitMap.Store(msg.UUID, waitChan)

		pnet.Send(a.conn, protocol.TlvEncode(msg.Encode()))

		// wail for response
		data := <-waitChan
		return data, nil
	}

	// just send data
	return nil, pnet.Send(a.conn, protocol.TlvEncode(msg.Encode()))
}

// 远程获取agent端的系统信息
func (a *Agent) AgentInfo() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.AgentInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 远程获取agent端的系统信息
func (a *Agent) GetOSInfo() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.OsInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 远程获取agent端的CPU信息
func (a *Agent) GetCPUInfo() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.CPUInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 远程获取agent端的内存信息
func (a *Agent) GetMemoryInfo() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.MemoryInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 远程获取agent端的内核信息
func (a *Agent) GetSysctlInfo() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.SysctlInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 临时修改agent端系统参数
func (a *Agent) ChangeSysctl(args string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.SysctlChange,
		Data: args,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 查看某个内核参数的值
func (a *Agent) SysctlView(args string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.SysctlView,
		Data: args,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 查看服务列表
func (a *Agent) ServiceList() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.ServiceList,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 查看某个服务的状态
func (a *Agent) ServiceStatus(service string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.ServiceStatus,
		Data: service,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 重启服务
func (a *Agent) ServiceRestart(service string) (interface{}, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.ServiceRestart,
		Data: service,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, "", err
	}
	return resp_message.Data, resp_message.Error, nil
}

// 关闭服务
func (a *Agent) ServiceStop(service string) (interface{}, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.ServiceStop,
		Data: service,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, "", err
	}
	return resp_message.Data, resp_message.Error, nil
}

// 启动服务
func (a *Agent) ServiceStart(service string) (interface{}, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.ServiceStart,
		Data: service,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, "", err
	}
	return resp_message.Data, resp_message.Error, nil
}

// 获取全部安装的rpm包列表
func (a *Agent) AllRpm() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.AllRpm,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 获取源软件包名以及源
func (a *Agent) RpmSource(rpm string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.RpmSource,
		Data: rpm,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 获取软件包信息
func (a *Agent) RpmInfo(rpm string) (interface{}, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.RpmInfo,
		Data: rpm,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, "", err
	}
	return resp_message.Data, resp_message.Error, nil
}

// 获取源软件包名以及源
func (a *Agent) InstallRpm(rpm string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.InstallRpm,
		Data: rpm,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 获取源软件包名以及源
func (a *Agent) RemoveRpm(rpm string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.RemoveRpm,
		Data: rpm,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 获取磁盘的使用情况
func (a *Agent) DiskUsage() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.DiskUsage,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 获取磁盘的IO信息
func (a *Agent) DiskInfo() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.DiskInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

/*挂载磁盘
1.创建挂载磁盘的目录
2.挂载磁盘*/
func (a *Agent) DiskCreatPath(mountpath string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.CreateDiskPath,
		Data: mountpath,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}
func (a *Agent) DiskMount(sourceDisk, destPath string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.DiskMount,
		Data: sourceDisk + "," + destPath,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}
func (a *Agent) DiskUMount(diskPath string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.DiskUMount,
		Data: diskPath,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}
func (a *Agent) DiskFormat(fileType, diskPath string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.DiskFormat,
		Data: fileType + "," + diskPath,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 获取当前TCP网络连接信息
func (a *Agent) NetTCP() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.NetTCP,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 获取当前UDP网络连接信息
func (a *Agent) NetUDP() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.NetUDP,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 获取网络读写字节／包的个数
func (a *Agent) NetIOCounter() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.NetIOCounter,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 获取网卡配置
func (a *Agent) NetNICConfig() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.NetNICConfig,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 获取当前用户信息
func (a *Agent) CurrentUser() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.CurrentUser,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 获取所有用户的信息
func (a *Agent) AllUser() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.AllUser,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 创建新的用户，并新建家目录
func (a *Agent) AddLinuxUser(username, password string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.AddLinuxUser,
		Data: username + "," + password,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 删除用户
func (a *Agent) DelUser(username string) (interface{}, string, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.DelUser,
		Data: username,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, "", err
	}
	return resp_message.Data, resp_message.Error, nil
}

// chmod [-R] 权限值 文件名
func (a *Agent) ChangePermission(permission, file string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.ChangePermission,
		Data: permission + "," + file,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// chown [-R] 所有者 文件或目录
func (a *Agent) ChangeFileOwner(user, file string) (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.ChangeFileOwner,
		Data: user + "," + file,
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 远程获取agent端的内核信息
func (a *Agent) GetAgentOSInfo() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.AgentOSInfo,
		Data: struct{}{},
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}

// 心跳
func (a *Agent) HeartBeat() (interface{}, error) {
	msg := &protocol.Message{
		UUID: uuid.New().String(),
		Type: protocol.Heartbeat,
		Data: "连接正常",
	}

	resp_message, err := a.sendMessage(msg, true, 0)
	if err != nil {
		logger.Error("failed to run script on agent")
		return nil, err
	}
	return resp_message.Data, nil
}
