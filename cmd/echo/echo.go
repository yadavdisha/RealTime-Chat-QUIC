package main

import (
	"flag"
	"log"

	"drexel.edu/net-quic/pkg/client"
	"drexel.edu/net-quic/pkg/server"
)

func main() {
	serverMode := flag.Bool("server", false, "Run in server mode")
	clientMode := flag.Bool("client", false, "Run in client mode")
	certFile := flag.String("cert", "certs/quic_certificate.crt", "TLS certificate file")
	keyFile := flag.String("key", "certs/quic_private_key.pem", "TLS key file")
	address := flag.String("address", "localhost", "Server address")
	port := flag.Int("port", 4242, "Server port")
	username := flag.String("username", "defaultuser", "Username")

	flag.Parse()

	if *serverMode {
		cfg := server.ServerConfig{
			GenTLS:   false,
			CertFile: *certFile,
			KeyFile:  *keyFile,
			Address:  *address,
			Port:     *port,
		}
		srv := server.NewServer(cfg)
		log.Fatal(srv.Run())
	} else if *clientMode {
		cfg := client.ClientConfig{
			ServerAddr: *address,
			PortNumber: *port,
			CertFile:   *certFile,
			Username:   *username,
		}
		cli := client.NewClient(cfg)
		log.Fatal(cli.Run())
	} else {
		log.Fatal("Please specify either -server or -client mode")
	}
}
