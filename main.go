package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net"
	"strings"
)

var buf = make([]byte, 4096)

func main() {
	gin.SetMode(gin.DebugMode) // 运行模式
	router := gin.Default()
	router.Use(func(context *gin.Context) { //解决跨域
		context.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		context.Writer.Header().Set("Cache-Control", "no-cache")
		context.Writer.Header().Set("Referer Policy", "no-referrer-when-downgrade")
	})

	router.GET("/list", list)
	router.PUT("/relay", updateRelay)
	router.GET("/cmd", cmd)
	router.GET("/state", state)
	router.GET("/log", getLog)
	err := router.Run(":8083")
	if err != nil {
		fmt.Println(err.Error())
	}
}

func init() {
	go func() {
		fmt.Println("网络继电器服务端启动中")
		localAddress, _ := net.ResolveTCPAddr("tcp4", "0.0.0.0:5001") //定义一个本机IP和端口。
		var tcpListener, err = net.ListenTCP("tcp", localAddress)     //在刚定义好的地址上进监听请求。
		if err != nil {
			fmt.Println("监听出错：", err)
		}
		defer tcpListener.Close()
		fmt.Println("等待客户机的连接")
		//for {
		var conn, err2 = tcpListener.AcceptTCP() //接受连接。
		if err2 != nil {
			fmt.Println("接受连接失败：", err2)
			return
		}
		fmt.Println("连接成功")
		var remoteAddr = conn.RemoteAddr() //获取连接到的对像的IP地址。
		fmt.Println("客户端IP与端口 --> ", remoteAddr)

		/*			conn.Write([]byte("AT+NAME=?\r\n"))
					a, _ := conn.Read(buf)
					name := string(buf[6 : a-2])
					conn.Write([]byte("AT+IP=?\r\n"))
					b, _ := conn.Read(buf)
					ip := string(buf[4 : b-2])
					conn.Write([]byte("AT+STACH1=?\r\n"))
					c, _ := conn.Read(buf)*/
		ConnPool = append(ConnPool, con{
			Ip:    remoteAddr.String(),
			Name:  "abc",
			State: "open",
			Con:   conn,
		})
		//connHandler(conn)
		//}
	}()

}

func connHandler(c net.Conn) {
	for {
		cnt, err := c.Read(buf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("响应：", strings.Replace(string(buf[:cnt]), "\r\n", "\\r\\n", 1))
		//c.Close() //关闭client端的连接，telnet 被强制关闭
	}
}

