package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	// 定义一个TCP断点
	var tcpAddr *net.TCPAddr
	// 通过ResolveTCPAddr实例化一个TCP断点
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")
	// 打开一天TCP断点监听
	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)
	defer tcpListener.Close()
	fmt.Println("Server ready to read ...")
	//循环接收客户端的连接，创建一个协程具体去处理连接
	for {
		conn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Printf("Build listener err: ", err.Error())
			continue
		}
		go tcpPipe(conn)
	}

}
func tcpPipe(conn *net.TCPConn) {
	// tcp连接地址
	ipStr := conn.RemoteAddr().String()
	defer func() {
		fmt.Println("Disconnect from client ", ipStr)
		conn.Close()
	}()
	reader := bufio.NewReader(conn)
	i := 0
	// 接收并传回消息
	for {
		message, err := reader.ReadString('\n')
		if err != nil || err == io.EOF {
			break
		}
		fmt.Println("message from client: ", message)
		time.Sleep(time.Second * 3)
		msg := fmt.Sprintf("time: %s; ip: %s; Remote server say hello... \n", time.Now().String(), conn.LocalAddr().String())
		writeByte := []byte(msg)
		conn.Write(writeByte)
		i++
		if i > 10 {
			break
		}
	}

}
