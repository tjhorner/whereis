package model

import (
	"time"

	"github.com/dchest/uniuri"

	"github.com/jinzhu/gorm"
)

// AccessKey represents access to a certain resource, optionally for a limited time
type AccessKey struct {
	gorm.Model
	Key        string `gorm:"UNIQUE_INDEX"`
	Indefinite bool   `gorm:"default:true"`
	StartDate  *time.Time
	EndDate    *time.Time
	Notes      string
}

// BeforeCreate implements a gorm hook for setting the access key to a random string if it is not set already
func (key *AccessKey) BeforeCreate() error {
	if key.Key == "" {
		key.Key = uniuri.NewLenChars(16, []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"))
	}
	return nil
}
