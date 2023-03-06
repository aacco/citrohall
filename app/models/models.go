package models

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Admin struct {
	gorm.Model
	Name     string
	Password string
	Message  string
}

type User struct {
	gorm.Model
	Name     string
	Password string
	Message  string
}

type Music struct {
	gorm.Model
	UserId  uint
	GroupId uint
	Path    string

	Format   string
	FileType string
	// metadata
	Title       string
	Album       string
	Artist      string
	AlbumArtist string
	Composer    string
	Genre       string
	Year        int

	Track      int // Number
	TrackTotal int // Total
	Disc       int // Number
	DiscTotal  int // Total

	Lyrics  string
	Comment string
}

type Album struct {
	gorm.Model
	UserId  uint
	GroupId uint
	Path    string

	// metadata
	Album       string
	Artist      string
	AlbumArtist string
	Composer    string
	Genre       string
	Year        int
}

type Playlist struct {
	gorm.Model
	UserId  uint
	GroupId uint
}

type Group struct {
	gorm.Model
	UserId  uint
	Name    string
	Message string
}

type Comment struct {
	gorm.Model
	UserId  uint
	PostId  uint
	Message string
}

type CommentJoin struct {
	Comment
	User
	Music
}

func Migrate() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		fmt.Println(err)
	}

	db.AutoMigrate(&Admin{}, &User{}, &Music{}, &Album{}, &Playlist{}, &Group{}, &Comment{})
}
