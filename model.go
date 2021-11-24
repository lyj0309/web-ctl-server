package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"strings"
)

//Conn类型表示WebSocket连接。服务器应用程序从HTTP请求处理程序调用Upgrader.Upgrade方法以获取* Conn：
// var upgrader = websocket.Upgrader{}
var (
	upgrader = websocket.Upgrader{
		// 读取存储空间大小
		ReadBufferSize: 1024,
		// 写入存储空间大小
		WriteBufferSize: 1024,
		// 允许跨域
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func list(c *gin.Context) {
	fmt.Println(ConnPool)

	c.JSON(200, ConnPool)
}

func updateRelay(c *gin.Context) {
	s := c.Query(`state`)
	if s == "open" { //打开开关

	} else { //关闭开关

	}
}

func cmd(c *gin.Context) { //命令行
	//   完成握手 升级为 WebSocket长连接，使用conn发送和接收消息。
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer conn.Close()

	go func() {
		for {
			cnt, err := ConnPool[0].Con.Read(buf)
			if err != nil {
				fmt.Println(err)
			}
			//fmt.Println("响应：", string(buf[:cnt]),buf[:cnt])
			d := strings.Replace(string(buf[:cnt]), "\r", "<br>", -1)
			d = strings.Replace(d, "\b", "", -1)
			fmt.Println("响应", []byte(d))
			if err := conn.WriteMessage(websocket.TextMessage, []byte(d)); err != nil {
				log.Println("Writeing error...", err)
				return
			}
			//c.Close() //关闭client端的连接，telnet 被强制关闭
		}
	}()

	//调用连接的WriteMessage和ReadMessage方法以一片字节发送和接收消息。实现如何回显消息：
	//p是一个[]字节，messageType是一个值为websocket.BinaryMessage或websocket.TextMessage的int。
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("Reading error...", err)
			return
		}
		log.Printf("Read from client msg:%s \n", msg)
		//发送给单片机串口模块
		_, err = ConnPool[0].Con.Write([]byte(string(msg) + "\r\n"))
		if err != nil {
			fmt.Println(err)
		}

		log.Printf("Write msg to client: recved: %s \n", msg)
	}


}

func getLog(c *gin.Context) {

}

func state(c *gin.Context) {
	conn := ConnPool[0].Con
	_, err := conn.Write([]byte("AT+STACH1=?\r\n"))
	if err != nil {
		fmt.Println(err)
	}
	c.JSON(200, gin.H{
		"data": "asdf",
	})

}
