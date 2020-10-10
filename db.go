package main

import (
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
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
	return c.CreatedAt.Format("02/01/2006 15:04")
}

// SetupDB and return active connection
func SetupDB() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(viper.GetString("database.url")), &gorm.Config{})
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect to database")
	}

	if viper.GetBool("debug") {
		db = db.Debug()
	}

	// Migrate the schema
	db.AutoMigrate(&Commit{})

	return db, nil
}
