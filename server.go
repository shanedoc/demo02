package main

import (
	"fmt"
	"io"
	"net"
	"sync"
)

//server类
type Server struct {
	Ip   string
	Port int
	//在线用户表
	UserOnlineMap map[string]*User
	mapLock       sync.RWMutex
	//消息广播channel
	Message chan string
}

//创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:            ip,
		Port:          port,
		UserOnlineMap: make(map[string]*User),
		Message:       make(chan string),
	}
	return server
}

//执行当前链接的业务逻辑
func (s *Server) Handler(conn net.Conn) {
	//fmt.Println("链接建立成功")
	user := NewUser(conn)
	//用户上线,将当前用户加入到map中
	//加锁
	s.mapLock.Lock()
	s.UserOnlineMap[user.Name] = user
	s.mapLock.Unlock()

	//广播用户上线通知
	s.BroadCast(user, "用户已上线")

	//接收客户端发送的消息&进行广播
	go func() {
		buf := make([]byte, 4096) //4k
		for {
			n, err := conn.Read(buf)
			if n == 0 {
				s.BroadCast(user, "用户下线")
				return
			}
			if err != nil && err != io.EOF {
				fmt.Println("conn.read.err:", err)
				return
			}
			//获取用户终端输入信息
			msg := string(buf[:n-1])
			s.BroadCast(user, msg)
		}
	}()

	//当前handler阻塞
	select {}

}

//启动服务接口
func (s *Server) Start() {
	//socket listen
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.Ip, s.Port))
	if err != nil {
		fmt.Println("net.listen.err:", err)
		return
	}
	defer listener.Close()

	//监听message的goroutine
	go s.ListenMessager()

	for {
		//accept
		con, err := listener.Accept()
		if err != nil {
			fmt.Println("listener.accept.err:", err)
			return
		}
		//启动协程处理当前链接的业务
		go s.Handler(con)
	}

	//close socket listen
}

//消息广播
func (s *Server) BroadCast(user *User, msg string) {
	content := "[" + user.Addr + "]" + user.Name + ":" + msg
	s.Message <- content
}

//监听message广播的channel:有上线提示时发给全体在线用户
func (s *Server) ListenMessager() {
	for {
		msg := <-s.Message
		s.mapLock.Lock()
		for _, v := range s.UserOnlineMap {
			//向用户管道中写入上线数据
			v.C <- msg
		}
		s.mapLock.Unlock()
	}
}
