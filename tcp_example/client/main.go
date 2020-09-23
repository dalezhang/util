package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

func main() {
	var tcpAddr *net.TCPAddr
	tcpAddr, _ = net.ResolveTCPAddr("tcp", "127.0.0.1:9999")
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	defer func() {
		conn.Close()
	}()
	if err != nil {
		fmt.Printf("connect error: %s", err)
		return
	}
	onMessageReceive(conn)
}

func onMessageReceive(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	b := []byte(conn.LocalAddr().String() + " say hello to server... \n")
	conn.Write(b)
	for {
		msg, err := reader.ReadString('\n')
		fmt.Printf("Read string: %s", msg)
		time.Sleep(time.Second * 3)
		fmt.Println("wait ...")
		b = []byte(conn.LocalAddr().String() + "write to server.. \n")
		_, err = conn.Write(b)
		if err != nil || err == io.EOF {
			fmt.Println("connection close")
			break
		}
	}
}
