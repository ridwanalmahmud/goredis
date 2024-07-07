package main

import (
	"net"
    "log"
	"fmt"
	"log/slog"
)

const (
	defaultListenAddr = ":5001"
)

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	quitCh    chan struct{}
	msgCh     chan []byte
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}
	return &Server {
		Config:    cfg,	
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgCh:     make(chan []byte),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.ln = ln
	go s.loop()
	slog.Info("Server running", "ListenAddr", s.ListenAddr)
	return s.acceptLoop()
}

func (s *Server) handleRawMsg(rawMsg []byte) error {
	fmt.Println(string(rawMsg))
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case rawMsg := <- s.msgCh:
			if err := s.handleRawMsg(rawMsg); err != nil {
				slog.Info("Raw message error", "err", err)
			}
		case <- s.quitCh:
			return
		case peer := <- s.addPeerCh:
			s.peers[peer]= true
		}		    
	}
}

func (s *Server) acceptLoop() error{
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error", "err", err)
			continue
		}
		go s.handleConn(conn)
	}
} 

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh)
	s.addPeerCh <- peer
	slog.Info("New peer connected", "RemoteAddr", conn.RemoteAddr())
	if err := peer.readLoop(); err != nil {
		slog.Error("Peer read error", "err", err, "RemoteAddr", conn.RemoteAddr())
	}
}

func main() {
	server := NewServer(Config{})
	log.Fatal(server.Start())
}