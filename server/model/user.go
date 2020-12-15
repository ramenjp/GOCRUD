package model

import (
	"time"
	"github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/mysql"
)

type User struct {
	gorm.Model
	Id int
	Name string
	CreatedAt time.Time
}

