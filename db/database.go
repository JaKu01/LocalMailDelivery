package db

import "gorm.io/gorm"

func LoadMails(db *gorm.DB) []Mail {
	var mails []Mail
	db.Find(&mails)
	return mails
}

func SaveMail(db *gorm.DB, mail *Mail) {
	db.Create(mail)
}
