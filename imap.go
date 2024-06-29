package LocalMail

import (
	"crypto/tls"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver"
	"github.com/emersion/go-imap/v2/imapserver/imapmemserver"
	"log"
	"net"
	"os"
	"strings"
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

func setupUser(memServer *imapmemserver.Server, username string, password string) {
	if username != "" || password != "" {
		user := imapmemserver.NewUser(username, password)
		user.Create("INBOX", &imap.CreateOptions{})

		msgString := GetTestMessage("This is a test mail").String()

		_, err := user.Append("INBOX", strings.NewReader(msgString), &imap.AppendOptions{
			Flags: []imap.Flag{imap.FlagSeen},
		})

		if err != nil {
			log.Fatalf("Failed to append message: %v", err)
		}

		memServer.AddUser(user)
		go handleUser(user)
	}
}

func StartServer(insecureAuth bool, username string, password string) {
	var tlsConfig *tls.Config
	loadTLSConfig(tlsConfig, "", "")

	ln, err := net.Listen("tcp", "0.0.0.0:143")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("IMAP server listening on %v", ln.Addr())

	memServer := imapmemserver.New()
	setupUser(memServer, username, password)

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

	if err := server.Serve(ln); err != nil {
		log.Fatalf("Serve() = %v", err)
	}
}
