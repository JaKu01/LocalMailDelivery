package web

import (
	"encoding/json"
	"fmt"
	"github.com/JaKu01/LocalMail/db"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver/imapmemserver"
	"gorm.io/gorm"
	"net/http"
	"strings"
)

func handlePost(user *imapmemserver.User, database *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// unmarshal json
		var mail db.Mail
		if err := json.NewDecoder(r.Body).Decode(&mail); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// append mail
		result, err := user.Append("INBOX", strings.NewReader(mail.String()), &imap.AppendOptions{
			Flags: []imap.Flag{},
		})

		// save mail in database
		db.SaveMail(database, &mail)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Fprintf(w, "Appended message: %v\n", result)
	}
}
