package main

import (
	"encoding/json"
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

// MarshalJSON for api responses
func (c *Commit) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Commit
		CreatedAt string `json:"created_at"`
	}{
		*c,
		c.PrintCreatedAt(),
	})
}

// PrintCreatedAt formats the created at timestamp
func (c *Commit) PrintCreatedAt() string {
	return c.CreatedAt.Format("01/02/2006 15:04")
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
