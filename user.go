package main

import (
	"net"
	"strings"
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

func (u *User) SendMsg(msg string) {
	u.conn.Write([]byte(msg))
}

//用户信息广播
func (u *User) BrocastMessage(msg string) {
	//在线用户查询
	if msg == "who" {
		u.server.mapLock.Lock()
		for _, user := range u.server.UserOnlineMap {
			userinfo := "[" + user.Addr + "]" + user.Name + ":" + "在线..\n"
			u.SendMsg(userinfo)
		}
		u.server.mapLock.Unlock()
	} else if len(msg) > 7 && msg[:7] == "rename|" {
		//修改用户名格式:rename|张三
		newName := strings.Split(msg, "|")[1]
		_, ok := u.server.UserOnlineMap[newName]
		if ok {
			u.SendMsg("当前用户名已被占用")
		} else {
			u.server.mapLock.Lock()
			//删除原用户名
			delete(u.server.UserOnlineMap, u.Name)
			//添加新用户名
			u.server.UserOnlineMap[newName] = u
			u.server.mapLock.Unlock()
			u.Name = newName
			u.SendMsg("您已经更新用户名为" + u.Name + "\n")
		}

	} else {
		//普通消息广播
		u.server.BroadCast(u, msg)
	}

}
