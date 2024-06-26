package LocalMail

import (
	"crypto/tls"
	"encoding/json"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver"
	"github.com/emersion/go-imap/v2/imapserver/imapmemserver"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
)

var (
	tlsCert string
	tlsKey  string
)

func StartServer(insecureAuth bool, username string, password string) {
	var tlsConfig *tls.Config
	if tlsCert != "" || tlsKey != "" {
		cert, err := tls.LoadX509KeyPair(tlsCert, tlsKey)
		if err != nil {
			log.Fatalf("Failed to load TLS key pair: %v", err)
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}

	ln, err := net.Listen("tcp", "0.0.0.0:143")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	log.Printf("IMAP server listening on %v", ln.Addr())

	memServer := imapmemserver.New()

	if username != "" || password != "" {
		user := imapmemserver.NewUser(username, password)
		user.Create("INBOX", nil)

		msgString, err := MessageToString(GetTestMessage("This is a test mail"))

		if err != nil {
			log.Panic("Could not convert message to string", err)
		}

		_, err = user.Append("INBOX", strings.NewReader(msgString), &imap.AppendOptions{
			Flags: []imap.Flag{imap.FlagSeen},
		})

		if err != nil {
			log.Fatalf("Failed to append message: %v", err)
		}

		msgString2, _ := MessageToString(GetTestMessage("This is a second mail"))
		_, err = user.Append("INBOX", strings.NewReader(msgString2), &imap.AppendOptions{})

		if err != nil {
			log.Fatalf("Failed to append message: %v", err)
		}

		msgString3, _ := MessageToString(GetTestMessage("This is a second mail"))
		_, err = user.Append("INBOX", strings.NewReader(msgString3), &imap.AppendOptions{})

		memServer.AddUser(user)
		go handleUser(user)
	}

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

type Data struct {
	Mail string `json:"mail"`
}

func handleUser(user *imapmemserver.User) {
	// serve http

	mux := http.NewServeMux()

	mux.HandleFunc("POST /", func(w http.ResponseWriter, r *http.Request) {
		// unmarshal json
		var data Data
		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// append mail
		_, err := user.Append("INBOX", strings.NewReader(data.Mail), &imap.AppendOptions{})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
