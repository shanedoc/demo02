package main

import (
	"net"
)

type User struct {
	Name   string
	Addr   string
	C      chan string
	conn   net.Conn
	server *Server
}

//创建用户接口
func NewUser(conn net.Conn, server *Server) *User {
	//用户地址
	userAddr := conn.RemoteAddr().String()
	user := &User{
		Name:   userAddr,
		Addr:   userAddr,
		C:      make(chan string),
		conn:   conn,
		server: server,
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

//上线接口
func (u *User) OnlineUser() {

	//用户上线,将当前用户加入到map中
	//加锁
	u.server.mapLock.Lock()
	u.server.UserOnlineMap[u.Name] = u
	u.server.mapLock.Unlock()
	//广播用户上线通知
	u.server.BroadCast(u, "用户已上线")
}

//下线接口
func (u *User) OfflineUser() {
	u.server.mapLock.Lock()
	delete(u.server.UserOnlineMap, u.Name)
	u.server.mapLock.Unlock()
	//广播用户上线通知
	u.server.BroadCast(u, "用户已下线")
}

//用户信息广播
func (u *User) BrocastMessage(msg string) {
	u.server.BroadCast(u, msg)
}
