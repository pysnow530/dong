package main

import (
	"encoding/base64"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"time"
)

type Paper struct {
	ID        string     `json:"id" gorm:"primary_key"`
	Data      []byte     `json:"data"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"-" sql:"index"`
}

func (paper *Paper) Encrypt(key []byte) error {
	dataString, err := base64.StdEncoding.DecodeString(string(paper.Data))
	if err != nil {
		return err
	}
	paper.Data = []byte(dataString)

	dataBytes, err := Encrypt(paper.Data, key)
	if err != nil {
		return err
	}

	paper.Data = dataBytes

	return nil
}

func (paper *Paper) Decrypt(key []byte) error {
	data, err := Decrypt(paper.Data, key)
	if err != nil {
		return err
	}

	paper.Data = data

	return nil
}

func ConnectDB(dsn string) (*gorm.DB, error) {
	DB, err := gorm.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}

	DB.AutoMigrate(&Paper{})

	return DB, nil
}
