package main

import (
	"net"
)

type User struct {
	Name string
	Addr string
	C    chan string
	conn net.Conn
}

//创建用户接口
func NewUser(conn net.Conn) *User {
	//用户地址
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name: userAddr,
		Addr: userAddr,
		C:    make(chan string),
		conn: conn,
	}
	//启动监听当前用户chan消息的goroutine
	go user.ListenMessage()
	return user
}

//监听user对应的channel
func (u *User) ListenMessage() {
	for {
		msg := <-u.C
		//取出管道中的消息写入对应客户端
		u.conn.Write([]byte(msg + "\n"))
	}
}
