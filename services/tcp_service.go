package services

import (
	"fmt"
	"io"
	"log"
	"net"

	"github.com/luispfcanales/tcpserver/model"
)

type TCPActor struct {
	addr     string
	ln       net.Listener
	devices  map[string]net.Conn
	mailFunc chan func()
	mailBox  <-chan model.MessageTCP
	signal   chan struct{}
}

func NewTCPActor(addr string, mailBox <-chan model.MessageTCP) *TCPActor {
	s := &TCPActor{
		addr:     addr,
		devices:  make(map[string]net.Conn),
		mailBox:  mailBox,
		mailFunc: make(chan func(), 10),
		signal:   make(chan struct{}),
	}

	go s.boostrap()
	return s
}

func (s *TCPActor) boostrap() {
	log.Println("[ Start Actor TCP-HTTP boostrap ]")
	for {
		select {
		case msg := <-s.mailBox:
			s.Broadcast(msg)
		case fn := <-s.mailFunc:
			fn()
		}
	}
}

func (s *TCPActor) Run() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer ln.Close()

	s.ln = ln

	go s.acceptConnections()
	log.Println("[ Start service TCP in port: ", s.addr, " ]")
	<-s.signal
	return nil
}

// acceptLoop accept connections
func (s *TCPActor) acceptConnections() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			log.Println("accept error: ", err)
			continue
		}

		s.mailFunc <- func() {
			s.registerConn(conn)
		}
	}
}
func (s *TCPActor) readLoop(cn net.Conn) {
	defer func() {
		s.mailFunc <- func() {
			s.unregisterConn(cn)
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := cn.Read(buf)
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("[ Read error ]:", err)
			continue
		}

		msg := buf[:n]
		log.Println(msg)
	}
}

func (s *TCPActor) registerConn(cn net.Conn) {
	s.devices[cn.RemoteAddr().String()] = cn
	log.Println("[ Registered: ]", cn.RemoteAddr().String())
	go s.readLoop(cn)
}

func (s *TCPActor) unregisterConn(cn net.Conn) {
	delete(s.devices, cn.RemoteAddr().String())
	log.Println("[ UnRegistered: ]", cn.RemoteAddr().String())
}

// Broadcast send MessageTCP
func (s *TCPActor) Broadcast(msg model.MessageTCP) {
	s.mailFunc <- func() {
		s.SendtoAllConn(msg)
	}
}

func (s *TCPActor) SendtoAllConn(msg model.MessageTCP) {
	for key, cn := range s.devices {
		event := fmt.Sprintf("[ %s : %s -> %v]", key, "Sending Event", msg)
		log.Println(event)
		cn.Write(msg.Payload)
	}
}
