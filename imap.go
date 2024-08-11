package LocalMail

import (
	"crypto/tls"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver"
	"github.com/emersion/go-imap/v2/imapserver/imapmemserver"
	"log"
	"net"
	"os"
)

func loadTLSConfig(config *tls.Config, tlsCert string, tlsKey string) {
	if tlsCert != "" || tlsKey != "" {
		cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
		if err != nil {
			log.Fatalf("Failed to load TLS key pair: %v", err)
		}
		config = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}
}

func CreateServer(insecureAuth bool) (*imapmemserver.Server, *imapserver.Server) {
	var tlsConfig *tls.Config
	loadTLSConfig(tlsConfig, os.Getenv("CERTPATH"), os.Getenv("KEYPATH"))

	memServer := imapmemserver.New()

	server := imapserver.New(&imapserver.Options{
		NewSession: func(conn *imapserver.Conn) (imapserver.Session, *imapserver.GreetingData, error) {
			return memServer.NewSession(), nil, nil
		},
		Caps: imap.CapSet{
			imap.CapIMAP4rev1: {},
			imap.CapIMAP4rev2: {},
		},
		TLSConfig:    tlsConfig,
		InsecureAuth: insecureAuth,
		DebugWriter:  os.Stdout,
	})

	return memServer, server
}

func RunServer(server *imapserver.Server) {
	ln, err := net.Listen("tcp", "0.0.0.0:143")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	if err := server.Serve(ln); err != nil {
		log.Fatalf("Serve() = %v", err)
	}
}
