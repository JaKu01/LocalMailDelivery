package LocalMail

import (
	"encoding/json"
	"fmt"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver/imapmemserver"
	"log"
	"net/http"
	"strings"
)

func handlePost(user *imapmemserver.User) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// unmarshal json
		var mail Mail
		if err := json.NewDecoder(r.Body).Decode(&mail); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// append mail
		result, err := user.Append("INBOX", strings.NewReader(mail.String()), &imap.AppendOptions{
			Flags: []imap.Flag{},
		})

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Appended message: %v\n", result)
	}
}

func handleUser(user *imapmemserver.User) {
	// serve http
	mux := http.NewServeMux()
	mux.HandleFunc("POST /", handlePost(user))

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Fatalf("Failed to serve HTTP: %v", err)
	}
}
