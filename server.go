package main

import (
	"fmt"
	"net"
)

//server类
type Server struct {
	Ip   string
	Port int
}

//创建一个server的接口
func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:   ip,
		Port: port,
	}
	return server
}

//执行当前链接的业务逻辑
func (s *Server) Handler(conn net.Conn) {
	fmt.Println("链接建立成功")
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
