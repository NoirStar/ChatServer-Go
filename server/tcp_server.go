package server

import (
	"errors"
	"log"
	"net"
	"sync"

	"github.com/noirstar/chatserver-go/protocol"
)

// 클라이언트 들은 writer 객체를 하나씩 가짐

type client struct {
	conn   net.Conn
	name   string
	writer *protocol.CommandWriter
}

type TcpChatServer struct {
	listener net.Listener
	clients  []*client
	mutex    *sync.Mutex
}

var UnknownClient = errors.New("Unknown client")

func NewServer() *TcpChatServer {
	return &TcpChatServer{mutex: &sync.Mutex{}}
}

func (s *TcpChatServer) Listen(address string) error {
	l, err := net.Listen("tcp", address)

	if err == nil {
		s.listener = l
	}
	log.Printf("Listening on %v", address)

	return err
}

func (s *TcpChatServer) Close() {
	s.listener.Close()
}

func (s *TcpChatServer) Broadcast(command interface{}) error {
	for _, client := range s.clients {
		// TODO : handle error here
		client.writer.Write(command)
	}
	return nil
}

func (s *TcpChatServer) accept(conn net.Conn) *client {
	log.Printf("Accepting connection from %v, total clients: %v", conn.RemoteAddr().String(), len(s.clients)+1)

	s.mutex.Lock()
	defer s.mutex.Unlock()

	client := &client{
		conn:   conn,
		writer: protocol.NewCommandWriter(conn),
	}
	s.clients = append(s.clients, client)
	return client
}
