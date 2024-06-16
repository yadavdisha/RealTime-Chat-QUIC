package client

import (
	"bufio"
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"drexel.edu/net-quic/pkg/pdu"
	"drexel.edu/net-quic/pkg/util"
	"github.com/quic-go/quic-go"
)

// ClientConfig holds the configuration for the client
type ClientConfig struct {
	ServerAddr string
	PortNumber int
	CertFile   string
	Username   string
}

// Client represents a chat client
type Client struct {
	cfg  ClientConfig
	tls  *tls.Config
	conn quic.Connection
	ctx  context.Context
}

// NewClient creates a new client with the provided configuration
func NewClient(cfg ClientConfig) *Client {
	cli := &Client{
		cfg: cfg,
	}

	// Build TLS configuration
	if cfg.CertFile != "" {
		log.Printf("[cli] using cert file: %s", cfg.CertFile)
		t, err := util.BuildTLSClientConfigWithCert(cfg.CertFile)
		if err != nil {
			log.Fatal("[cli] error building TLS client config:", err)
			return nil
		}
		cli.tls = t
	} else {
		cli.tls = util.BuildTLSClientConfig()
	}

	cli.ctx = context.Background()
	return cli
}

// Run starts the client and connects to the server
func (c *Client) Run() error {
	serverAddr := fmt.Sprintf("%s:%d", c.cfg.ServerAddr, c.cfg.PortNumber)
	conn, err := quic.DialAddr(c.ctx, serverAddr, c.tls, nil)
	if err != nil {
		log.Printf("[cli] error dialing server: %v", err)
		return err
	}
	c.conn = conn

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		c.sendLeaveMessage()
		c.conn.CloseWithError(0, "client shutting down")
		os.Exit(0)
	}()

	// Open a bidirectional stream for communication
	stream, err := conn.OpenStreamSync(c.ctx)
	if err != nil {
		log.Printf("[cli] error opening stream: %v", err)
		return err
	}

	// Send join message
	joinMessage := fmt.Sprintf("%s joined the chat", c.cfg.Username)
	joinPDU := pdu.NewPDU(pdu.TYPE_JOIN, []byte(joinMessage))
	if err := c.sendPDU(stream, joinPDU); err != nil {
		return err
	}

	// Start a goroutine to receive messages
	go c.receiveMessages(stream)

	return c.protocolHandler(stream)
}

// protocolHandler reads messages from stdin and sends them to the server
func (c *Client) protocolHandler(stream quic.Stream) error {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message: ")
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("[cli] error reading input: %v", err)
			return err
		}
		message = strings.TrimSpace(message)

		if message == "" {
			log.Printf("[cli] warning: empty message not sent")
			continue
		}

		// Create and send the PDU
		req := pdu.NewPDU(pdu.TYPE_DATA, []byte(fmt.Sprintf("%s: %s", c.cfg.Username, message)))
		if err := c.sendPDU(stream, req); err != nil {
			log.Printf("[cli] error sending PDU: %v", err)
			return err
		}
		log.Printf("[cli] sent message: %s", message)
	}
}

// sendPDU sends a PDU to the server
func (c *Client) sendPDU(stream quic.Stream, pdu *pdu.PDU) error {
	pduBytes, err := pdu.PduToBytes()
	if err != nil {
		return fmt.Errorf("[cli] error making pdu byte array: %w", err)
	}

	_, err = stream.Write(pduBytes)
	if err != nil {
		return fmt.Errorf("[cli] error writing to stream: %w", err)
	}
	log.Printf("[cli] PDU sent: %s", string(pduBytes))
	return nil
}

// receiveMessages listens for messages from the server and displays them
func (c *Client) receiveMessages(stream quic.Stream) {
	buffer := pdu.MakePduBuffer()

	for {
		n, err := stream.Read(buffer)
		if err != nil {
			log.Printf("[cli] error reading from stream: %v", err)
			return
		}
		rsp, err := pdu.PduFromBytes(buffer[:n])
		if err != nil {
			log.Printf("[cli] error converting pdu from bytes: %v", err)
			return
		}
		rspDataString := string(rsp.Data)
		log.Printf("[cli] got message: %s", rspDataString)
	}
}

// sendLeaveMessage notifies the server that the client is leaving the chat
func (c *Client) sendLeaveMessage() {
	leaveMessage := fmt.Sprintf("%s left the chat", c.cfg.Username)
	leavePDU := pdu.NewPDU(pdu.TYPE_LEAVE, []byte(leaveMessage))
	stream, err := c.conn.OpenStreamSync(c.ctx)
	if err != nil {
		log.Printf("[cli] error opening stream to send leave message: %v", err)
		return
	}
	defer stream.Close()
	if err := c.sendPDU(stream, leavePDU); err != nil {
		log.Printf("[cli] error sending leave PDU: %v", err)
	}
}
