package models

import (
	"time"
)

// Users app user
type Users struct {
	UID       int    `gorm:"AUTO_INCREMENT;PRIMARY_KEY"`
	UserName  string `gorm:"unique_index"`
	Birthday  time.Time
	Email     string `gorm:"type:varchar(100);unique_index"`
	Mobile    string `gorm:"type:varchar(100);unique_index"`
	Password  string `gorm:"size:255"`         // set field size to 255
	Salt      string `gorm:"type:varchar(10)"` // set member number to unique and not null
	CreatedAt time.Time
	UpdatedAt time.Time
}
