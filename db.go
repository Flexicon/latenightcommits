package main

import (
	"time"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Commit model
type Commit struct {
	ID        string    `gorm:"primarykey" json:"id"`
	CreatedAt time.Time `gorm:"index:,sort:desc,type:btree" json:"created_at"`
	Message   string    `gorm:"size:255" json:"message"`
	Author    string    `gorm:"size:255" json:"author"`
	AvatarURL string    `gorm:"size:255" json:"avatar_url"`
	Link      string    `gorm:"size:255" json:"link"`
}

// PrintCreatedAt formats the created at timestamp
func (c *Commit) PrintCreatedAt() string {
	return c.CreatedAt.Format("2006/01/02 03:04 PM")
}

// SetupDB and return active connection
func SetupDB() (*gorm.DB, error) {
	// TODO: replace with env variable for db DSN
	db, err := gorm.Open(mysql.Open("root:dev@(localhost)/latenightcommits?charset=utf8&parseTime=True&loc=Local"), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&Commit{})

	return db, nil
}
