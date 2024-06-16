package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"strings"
	"sync"

	"drexel.edu/net-quic/pkg/pdu"
	"drexel.edu/net-quic/pkg/util"
	"github.com/quic-go/quic-go"
)

// ServerConfig holds the configuration for the server
type ServerConfig struct {
	GenTLS   bool
	CertFile string
	KeyFile  string
	Address  string
	Port     int
}

// Server represents a chat server
type Server struct {
	cfg       ServerConfig
	tls       *tls.Config
	ctx       context.Context
	mu        sync.Mutex
	streams   map[quic.Stream]string // Store streams with usernames
	usernames map[string]string      // Store addresses with usernames
}

// NewServer creates a new server with the provided configuration
func NewServer(cfg ServerConfig) *Server {
	server := &Server{
		cfg:       cfg,
		streams:   make(map[quic.Stream]string),
		usernames: make(map[string]string),
	}
	server.tls = server.getTLS()
	server.ctx = context.Background()
	return server
}

// getTLS generates or builds the TLS configuration for the server
func (s *Server) getTLS() *tls.Config {
	if s.cfg.GenTLS {
		tlsConfig, err := util.GenerateTLSConfig()
		if err != nil {
			log.Fatal(err)
		}
		return tlsConfig
	} else {
		tlsConfig, err := util.BuildTLSConfig(s.cfg.CertFile, s.cfg.KeyFile)
		if err != nil {
			log.Fatal(err)
		}
		return tlsConfig
	}
}

// Run starts the server and listens for incoming connections
func (s *Server) Run() error {
	address := fmt.Sprintf("%s:%d", s.cfg.Address, s.cfg.Port)
	listener, err := quic.ListenAddr(address, s.tls, nil)
	if err != nil {
		log.Printf("error listening: %v", err)
		return err
	}

	log.Println("Welcome to Student Collaboration Chat Application!")
	for {
		log.Println("Accepting new session")
		sess, err := listener.Accept(s.ctx)
		if err != nil {
			log.Printf("error accepting: %v", err)
			continue
		}

		go s.sessionHandler(sess)
	}
}

// sessionHandler handles each new session
func (s *Server) sessionHandler(sess quic.Connection) {
	defer sess.CloseWithError(0, "session closed")

	for {
		log.Print("[server] waiting for client to open stream")
		stream, err := sess.AcceptStream(s.ctx)
		if err != nil {
			log.Printf("[server] error accepting stream: %v", err)
			return
		}

		s.mu.Lock()
		s.streams[stream] = "" // Initialize stream with empty username
		s.mu.Unlock()

		// Send welcome message
		welcomeMessage := "Welcome to Student Collaboration Chat Protocol Application! All the students are requested to maintain decorum in the chat room."
		welcomePDU := pdu.NewPDU(pdu.TYPE_DATA, []byte(welcomeMessage))
		s.sendPDU(stream, welcomePDU)

		go s.protocolHandler(stream, sess.RemoteAddr().String())
	}
}

// sendPDU sends a PDU to the client
func (s *Server) sendPDU(stream quic.Stream, pdu *pdu.PDU) error {
	pduBytes, err := pdu.PduToBytes()
	if err != nil {
		return fmt.Errorf("[server] error making pdu byte array: %w", err)
	}

	_, err = stream.Write(pduBytes)
	if err != nil {
		return fmt.Errorf("[server] error writing to stream: %w", err)
	}
	return nil
}

// protocolHandler handles the communication protocol with the client
func (s *Server) protocolHandler(stream quic.Stream, addr string) {
	defer func() {
		s.mu.Lock()
		username := s.streams[stream]
		delete(s.streams, stream)
		delete(s.usernames, addr)
		s.mu.Unlock()
		stream.Close()

		leaveMessage := fmt.Sprintf("%s left the chat", username)
		log.Printf("[server] %s", leaveMessage)
		s.broadcast(pdu.NewPDU(pdu.TYPE_LEAVE, []byte(leaveMessage)))
	}()

	buffer := pdu.MakePduBuffer()

	for {
		n, err := stream.Read(buffer)
		if err != nil {
			log.Printf("[server] Error Reading Raw Data: %v", err)
			return
		}

		data, err := pdu.PduFromBytes(buffer[:n])
		if err != nil {
			log.Printf("[server] Error decoding PDU: %v", err)
			return
		}

		messageParts := strings.SplitN(string(data.Data), ": ", 2)
		if len(messageParts) == 2 {
			username := messageParts[0]
			message := messageParts[1]

			// Store username for this stream
			s.mu.Lock()
			s.streams[stream] = username
			s.usernames[addr] = username
			s.mu.Unlock()

			log.Printf("[server] Received message from %s: %s", username, message)

			// Check if the message is a request for the list of active users
			if message == "/list" {
				s.listActiveUsers(stream)
				continue
			}
		} else {
			log.Printf("[server] Received message: %s", string(data.Data))
		}

		s.broadcast(data)
	}
}

// broadcast sends a PDU to all connected clients
func (s *Server) broadcast(data *pdu.PDU) {
	s.mu.Lock()
	defer s.mu.Unlock()

	rspBytes, err := data.PduToBytes()
	if err != nil {
		log.Printf("[server] Error encoding PDU: %v", err)
		return
	}

	for stream := range s.streams {
		_, err := stream.Write(rspBytes)
		if err != nil {
			log.Printf("[server] Error sending response: %v", err)
		}
	}
}

// listActiveUsers sends the list of active users to the requesting client
func (s *Server) listActiveUsers(stream quic.Stream) {
	s.mu.Lock()
	defer s.mu.Unlock()

	var userList []string
	for _, username := range s.streams {
		if username != "" {
			userList = append(userList, username)
		}
	}

	activeUsers := strings.Join(userList, ", ")
	response := pdu.NewPDU(pdu.TYPE_DATA, []byte(fmt.Sprintf("Active users: %s", activeUsers)))

	rspBytes, err := response.PduToBytes()
	if err != nil {
		log.Printf("[server] Error encoding PDU: %v", err)
		return
	}

	_, err = stream.Write(rspBytes)
	if err != nil {
		log.Printf("[server] Error sending response: %v", err)
	}
}
