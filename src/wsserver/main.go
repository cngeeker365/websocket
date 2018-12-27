package main

import (
	"wsserver/config"
	"encoding/json"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/gorilla/websocket"
	"github.com/satori/go.uuid"
	"io"
	"net/http"
	"wsserver/redis"
	"strings"
	"time"
)

type ClientManager struct {
	clients map[*Client]bool
	//broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	id        string
	page      int //当前client访问的page
	socket    *websocket.Conn
	send      chan []byte
	redisConn redis.Conn
}

var manager = ClientManager{
	//broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

func (manager *ClientManager) start() {
	go func() {
		for {
			select {
			case conn := <-manager.register: //新客户端加入
				manager.clients[conn] = true
			case conn := <-manager.unregister:
				if _, ok := manager.clients[conn]; ok {
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}
	}()
}

func (manager *ClientManager) send(message []byte, ignore *Client) {
	for conn := range manager.clients {
		if conn != ignore {
			conn.send <- message //发送的数据写入所有的 websocket 连接 管道
		}
	}
}

//客户端写入后 激活这里读取
//想改成 读取redis（已废弃）
func (c *Client) read() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
		fmt.Println("读关闭")
	}()

	for {
		_, message, err := c.socket.ReadMessage()
		//获取客户端传入的参数
		msg := make(map[string]interface{})
		json.Unmarshal(message, &msg)

		respMap := make(map[string]interface{})
		for k, v := range msg {
			if k == "ping" {
				respMap["pong"] = time.Now().Unix()
			}
			if k == "pageSwitch" {
				respMap["pageInfo"] = v
				redisConn := redisPool.Pool.Get()
				var pageIndex int
				//conf:=config.Conf
				for index, p := range config.Conf.Pages {
					if strings.EqualFold(p.Page.Name, v.(string)){
						pageIndex = index
					}
				}
				res, _ := redisPool.GetRedisData(redisConn, config.Conf.Pages[pageIndex].Page.Key)
				data, _ := json.Marshal(res)
				respMap["dataInit"] = string(data)
				redisConn.Close()
				//c.page = fmt.Sprintf("%v", v)
				c.page = pageIndex
			}
		}

		fmt.Println("是在不停的读吗？")
		if err != nil {
			manager.unregister <- c
			c.socket.Close()
			fmt.Println("读不到数据就关闭？")
			break
		}

		respStr, _ := json.Marshal(respMap)
		c.send <- respStr
		fmt.Println("发送数据到广播")
	}
}

//写入管道后激活这个进程
func (c *Client) write() {
	defer func() {
		manager.unregister <- c
		c.socket.Close()
		fmt.Println("写关闭了")
	}()

	for {
		select {
		case message, ok := <-c.send: //这个管道有了数据 写这个消息出去
			if !ok {
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				fmt.Println("发送关闭提示")
				return
			}

			err := c.socket.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				manager.unregister <- c
				c.socket.Close()
				fmt.Println("写不成功数据就关闭了")
				break
			}
			fmt.Println("写数据")
		}
	}
}

func main() {
	fmt.Println("Starting application...")
	go manager.start()
	http.HandleFunc("/ws", wsPage)
	http.ListenAndServe(":12345", nil)
}

func wsPage(res http.ResponseWriter, req *http.Request) {
	//解析一个连接
	conn, error := (&websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}).Upgrade(res, req, nil)
	if error != nil {
		io.WriteString(res, "这是一个websocket,不是网站.")
		return
	}

	uid, _ := uuid.NewV4()
	sha1 := uid.String()

	//初始化一个客户端对象
	client := &Client{id: sha1, socket: conn, send: make(chan []byte)}
	//把这个对象发送给 管道
	manager.register <- client

	go client.pushData()
	go client.read()
	go client.write()

}

func (c *Client) pushData() {
	c.redisConn = redisPool.Pool.Get()
	defer c.redisConn.Close()

	for {
		//if (strings.EqualFold(c.page, "chongtian")) {
			time.Sleep(30 * time.Second)
			res, _ := redisPool.GetRedisData(c.redisConn, config.Conf.Pages[c.page].Page.Key)
			data, _ := json.Marshal(res)
			resultTmp := map[string]interface{}{"dataInit": string(data)}
			jsonstr, _ := json.Marshal(resultTmp)
			c.send <- jsonstr
		//}
	}
}
