package common

/**
 * @Author: zhang han
 * @Date: 2021/11/18 14:54
 * @Description: ssh客户端连接
 */

import (
	"fmt"
	"golang.org/x/crypto/ssh"
	"net"
	"time"
)

type Client struct {
	IP           string `gorm:"type:varchar(25);not null" json:"ip,omitempty" form:"ip"`
	HostUser     string `gorm:"type:varchar(25);not null" json:"host_user,omitempty" form:"host_user"`
	HostPassword string `gorm:"type:varchar(25);not null" json:"host_password,omitempty" form:"host_password"`
	Port         int
	client       *ssh.Client
	LastResult   string
}

func NewSsh(ip string, username string, password string, port ...int) *Client {
	cli := new(Client)
	cli.IP = ip
	cli.HostUser = username
	cli.HostPassword = password
	if len(port) <= 0 {
		cli.Port = 22
	} else {
		cli.Port = port[0]
	}
	return cli
}

func (c *Client) connect() error { //连接
	config := ssh.ClientConfig{
		User: c.HostUser,
		Auth: []ssh.AuthMethod{ssh.Password(c.HostPassword)},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
		Timeout: 10 * time.Second,
	}
	addr := fmt.Sprintf("%s:%d", c.IP, c.Port)
	sshClient, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return err
	}
	c.client = sshClient
	return nil
}

func (c *Client) Run(shell string) (string, error) { //执行shell
	if c.client == nil {
		if err := c.connect(); err != nil {
			return "", err
		}
	}
	session, err := c.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	buf, err := session.CombinedOutput(shell)

	c.LastResult = string(buf)
	return c.LastResult, err
}
