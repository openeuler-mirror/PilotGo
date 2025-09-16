package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
)

const (
	tlvTag          = "eat_keyboard_if_duplicated"
	tlvLenBytes     = 4
	TypeHeartbeat   = 1
	TypeAgentInfo   = 11
	TypeAgentOSInfo = 44 // 正确的AgentOSInfo类型
)

type Message struct {
	UUID   string      `json:"message_uuid"`
	Type   int         `json:"message_type"`
	Status int         `json:"status"`
	Data   interface{} `json:"Data,omitempty"`
	Error  string      `json:"Error,omitempty"`
}

func tlvEncode(payload []byte) []byte {
	length := int32(len(payload))
	buf := new(bytes.Buffer)
	buf.WriteString(tlvTag)
	_ = binary.Write(buf, binary.BigEndian, &length)
	buf.Write(payload)
	return buf.Bytes()
}

func tlvDecodeStream(buffer *[]byte) (int, []byte) {
	if len(*buffer) <= len(tlvTag)+tlvLenBytes {
		return 0, nil
	}
	if string((*buffer)[:len(tlvTag)]) != tlvTag {
		return 0, nil
	}
	lengthBytes := (*buffer)[len(tlvTag) : len(tlvTag)+tlvLenBytes]
	var length int32
	_ = binary.Read(bytes.NewReader(lengthBytes), binary.BigEndian, &length)
	total := len(tlvTag) + tlvLenBytes + int(length)
	if len(*buffer) >= total {
		frame := (*buffer)[len(tlvTag)+tlvLenBytes : total]
		return total, frame
	}
	return 0, nil
}

type client struct {
	id         int
	addr       string
	conn       net.Conn
	version    string
	uuid       string // 每个客户端唯一的UUID
	ip         string
	hbInterval time.Duration
	closed     int32
	wg         *sync.WaitGroup
}

var (
	statDialOK    uint64
	statDialFail  uint64
	statCurrConns int64
)

func (c *client) send(msg *Message) error {
	if msg.UUID == "" {
		msg.UUID = uuid.New().String()
	}
	data, _ := json.Marshal(msg)
	_, err := c.conn.Write(tlvEncode(data))
	return err
}

func (c *client) run() {
	fmt.Printf("[client-%d] 开始连接 %s (伪IP: %s)\n", c.id, c.addr, c.ip)

	// 创建自定义的拨号器，使用不同的本地地址
	// 使用 127.0.0.x 的不同地址来模拟不同的本地机器
	localIP := fmt.Sprintf("127.0.%d.%d", (c.id/100)+1, (c.id%100)+1)

	dialer := &net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP:   net.ParseIP(localIP), // 使用 127.0.x.x 的本地地址
			Port: 0,                    // 让系统分配端口
		},
		Timeout: 5 * time.Second,
	}

	// 连接服务器
	conn, err := dialer.Dial("tcp", c.addr)
	if err != nil {
		atomic.AddUint64(&statDialFail, 1)
		fmt.Printf("[client-%d] 连接失败: %v\n", c.id, err)
		return
	}

	c.conn = conn
	defer conn.Close()

	// 获取本地端口和绑定的IP
	localPort := conn.LocalAddr().(*net.TCPAddr).Port
	boundIP := conn.LocalAddr().(*net.TCPAddr).IP.String()
	fmt.Printf("[client-%d] 连接成功: %s -> 本地绑定:%s:%d (伪IP:%s)\n", c.id, c.addr, boundIP, localPort, c.ip)

	atomic.AddInt64(&statCurrConns, 1)
	atomic.AddUint64(&statDialOK, 1)

	// 启动读取协程
	go c.readLoop()

	// 发送心跳
	ticker := time.NewTicker(c.hbInterval)
	defer ticker.Stop()

	for {
		if atomic.LoadInt32(&c.closed) == 1 {
			break
		}
		select {
		case <-ticker.C:
			err := c.send(&Message{Type: TypeHeartbeat, Data: "连接正常"})
			if err != nil {
				fmt.Printf("[client-%d] 发送心跳失败: %v\n", c.id, err)
				break
			}
		}
	}

	atomic.AddInt64(&statCurrConns, -1)
	fmt.Printf("[client-%d] 连接结束\n", c.id)
}

func (c *client) readLoop() {
	buf := make([]byte, 0, 8192)
	tmp := make([]byte, 2048)

	for {
		n, err := c.conn.Read(tmp)
		if err != nil {
			atomic.StoreInt32(&c.closed, 1)
			return
		}

		buf = append(buf, tmp[:n]...)
		for {
			consumed, frame := tlvDecodeStream(&buf)
			if consumed == 0 {
				break
			}
			buf = buf[consumed:]

			var msg Message
			if err := json.Unmarshal(frame, &msg); err != nil {
				continue
			}

			switch msg.Type {
			case TypeAgentInfo:
				// 使用伪IP，让服务器认为每个agent来自不同的机器
				resp := Message{
					UUID:   msg.UUID,
					Type:   msg.Type,
					Status: 0,
					Data: map[string]string{
						"agent_version": c.version,
						"agent_uuid":    c.uuid,
						"IP":            c.ip, // 使用伪IP
					},
				}
				_ = c.send(&resp)
			case TypeAgentOSInfo:
				// 响应系统信息请求，这是服务器要求agent提供系统信息的关键请求
				// 使用正确的SystemAndCPUInfo结构
				resp := Message{
					UUID:   msg.UUID,
					Type:   msg.Type,
					Status: 0,
					Data: map[string]string{
						"IP":              c.ip,                                        // 使用伪IP
						"Platform":        "Linux",                                     // 系统平台
						"PlatformVersion": "5.15.0",                                    // 系统版本
						"PrettyName":      "Ubuntu 22.04.3 LTS",                        // 可读性良好的OS具体版本
						"ModelName":       "Intel(R) Core(TM) i7-10700K CPU @ 3.80GHz", // CPU型号
					},
				}
				_ = c.send(&resp)
			case TypeHeartbeat:
				resp := Message{UUID: msg.UUID, Type: msg.Type, Status: 0, Data: "连接正常"}
				_ = c.send(&resp)
			}
		}
	}
}

func newClient(id int, addr, version, ip string, hbInterval time.Duration, wg *sync.WaitGroup) *client {
	// 为每个客户端生成唯一的UUID
	uniqueUUID := fmt.Sprintf("sim-agent-%d-%s", id, uuid.New().String()[:8])

	// 为每个客户端生成不同的伪IP地址
	// 使用不同的IP段，避免IP冲突
	pseudoIP := fmt.Sprintf("10.%d.%d.%d",
		(id/1000)+1,       // 第一个段：1-255
		((id%1000)/100)+1, // 第二个段：1-255
		(id%100)+1)        // 第三个段：1-255

	return &client{
		id:         id,
		addr:       addr,
		version:    version,
		uuid:       uniqueUUID, // 使用唯一UUID
		ip:         pseudoIP,   // 使用不同的伪IP
		hbInterval: hbInterval,
		wg:         wg,
	}
}

func main() {
	var (
		addr     string
		clients  int
		hb       time.Duration
		agentVer string
		ipStub   string
	)

	flag.StringVar(&addr, "addr", "127.0.0.1:8879", "PilotGo socket_server.addr")
	flag.IntVar(&clients, "clients", 1000, "并发连接数")
	flag.DurationVar(&hb, "hb", 5*time.Second, "心跳间隔")
	flag.StringVar(&agentVer, "agent_ver", "sim-1.0.0", "模拟器上报的 agent_version")
	flag.StringVar(&ipStub, "ip_stub", "10.0.0.", "生成伪IP的前缀")
	flag.Parse()

	fmt.Printf("[启动] 准备启动 %d 个模拟 agent 连接到 %s\n", clients, addr)

	// 启动统计协程
	go func() {
		ticker := time.NewTicker(3 * time.Second)
		defer ticker.Stop()
		for {
			<-ticker.C
			fmt.Printf("[stats] dial_ok=%d dial_fail=%d curr_conns=%d\n",
				atomic.LoadUint64(&statDialOK), atomic.LoadUint64(&statDialFail), atomic.LoadInt64(&statCurrConns))
		}
	}()

	// 创建所有客户端
	var wg sync.WaitGroup
	wg.Add(clients)

	// 启动所有连接
	for i := 0; i < clients; i++ {
		c := newClient(i, addr, agentVer, "", hb, &wg)

		go func(client *client) {
			defer wg.Done()
			client.run()
		}(c)

		// 稍微延迟，避免同时建立太多连接
		time.Sleep(100 * time.Millisecond)
	}

	fmt.Printf("[等待] 已启动 %d 个连接，等待...\n", clients)
	wg.Wait()

	fmt.Printf("[完成] 所有连接已结束，最终统计: dial_ok=%d dial_fail=%d\n",
		atomic.LoadUint64(&statDialOK), atomic.LoadUint64(&statDialFail))
}
