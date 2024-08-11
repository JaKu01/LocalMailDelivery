package web

import (
	"github.com/emersion/go-imap/v2/imapserver/imapmemserver"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func StartWebAPI(user *imapmemserver.User, database *gorm.DB) {
	mux := http.NewServeMux()

	mux.HandleFunc("POST /", handlePost(user, database))

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		log.Panic("Failed to serve HTTP", err)
	}
}
