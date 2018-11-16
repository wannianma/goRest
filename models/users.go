package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Users struct {
	gorm.Model
	Uid      int    `gorm:"column:uid;AUTO_INCREMENT;PRIMARY_KEY"`
	UserName string `gorm:"unique_index"`
	Birthday time.Time
	Email    string `gorm:"type:varchar(100);unique_index"`
	Mobile   string `gorm:"type:varchar(100);unique_index"`
	Password string `gorm:"size:255"`         // set field size to 255
	Salt     string `gorm:"type:varchar(10)"` // set member number to unique and not null
}
