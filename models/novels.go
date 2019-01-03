package models

import (
	"time"
)

// Users app user
type Novel struct {
	ID        int `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	Remote_id int
	Title     string
	Anthor    string
	Intro     string
	Nums      int
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Post struct {
	ID        int `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	NID       int
	Remote_id int
	Title     string
	Info      string `gorm:"type:text"`
}
