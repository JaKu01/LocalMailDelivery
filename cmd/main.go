package main

import (
	"fmt"
	"github.com/JaKu01/LocalMail"
	"github.com/JaKu01/LocalMail/db"
	"github.com/JaKu01/LocalMail/web"
	"github.com/emersion/go-imap/v2"
	"github.com/emersion/go-imap/v2/imapserver/imapmemserver"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"strings"
)

func main() {
	username, password := os.Getenv("USERNAME"), os.Getenv("PASSWORD")

	database, err := gorm.Open(sqlite.Open("db/sqlite/database.db"), &gorm.Config{})
	if err != nil {
		log.Panic("failed to connect database")
	}

	if err = database.AutoMigrate(&db.Mail{}); err != nil {
		log.Panic("failed to migrate Mailbox")
	}

	if username == "" || password == "" {
		log.Panic("username and password required")
	}

	memServer, server := LocalMail.CreateServer(true)
	imapUser := imapmemserver.NewUser(username, password)
	memServer.AddUser(imapUser)
	err = imapUser.Create("INBOX", &imap.CreateOptions{})

	if err != nil {
		log.Panic("failed to create INBOX")
	}

	loadExistingMails(imapUser, username, database)

	fmt.Fprintf(os.Stdout, "Server started for %v\n", username)

	go web.StartWebAPI(imapUser, database)
	LocalMail.RunServer(server)
}

func loadExistingMails(user *imapmemserver.User, username string, database *gorm.DB) {
	msgString := db.GetGreetingMessage(fmt.Sprintf("Welcome to the server, %v!", username)).String()

	_, err := user.Append("INBOX", strings.NewReader(msgString), &imap.AppendOptions{
		Flags: []imap.Flag{imap.FlagSeen},
	})

	if err != nil {
		log.Panic("failed to append welcome message")
	}

	mails := db.LoadMails(database)

	for _, mail := range mails {
		_, err := user.Append("INBOX", strings.NewReader(mail.String()), &imap.AppendOptions{
			Flags: []imap.Flag{imap.FlagSeen},
		})

		if err != nil {
			log.Panic("failed to append mail")
		}
	}
}
