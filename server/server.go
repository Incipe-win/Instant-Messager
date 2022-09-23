package server

import (
	"Instant-Messager/user"
	"fmt"
	"net"
	"sync"
)

type Server struct {
	Ip        string
	Port      int
	OnlineMap map[string]*user.User
	mapLock   sync.RWMutex
	Message   chan string
}

func (server *Server) ListenMessager() {
	for {
		msg := <-server.Message
		server.mapLock.Lock()
		for _, cli := range server.OnlineMap {
			cli.Ch <- msg
		}
		server.mapLock.Unlock()
	}
}

func NewServer(ip string, port int) *Server {
	server := &Server{
		Ip:        ip,
		Port:      port,
		OnlineMap: make(map[string]*user.User),
		Message:   make(chan string),
	}
	return server
}

func (server *Server) BroadCast(user *user.User, msg string) {
	sendMsg := "[" + user.Addr + "]" + user.Name + ":" + msg
	server.Message <- sendMsg
}

func (server *Server) Handler(connection net.Conn) {
	user := user.NewUser(connection)
	server.mapLock.Lock()
	server.OnlineMap[user.Name] = user
	server.mapLock.Unlock()
	server.BroadCast(user, "online")
}

func (server *Server) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Ip, server.Port))
	if err != nil {
		fmt.Printf("net.Listen error: %s", err)
		return
	}
	defer listener.Close()

	go server.ListenMessager()

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Printf("listener.Accept error: %s", err)
			continue
		}
		go server.Handler(connection)
	}
}
