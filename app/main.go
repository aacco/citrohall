package main

import (
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/dhowden/tag"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"citrohall/models"
)

type tempMusicStructure struct {
	Meta    models.Music
	PicExt  string
	PicBase string
}
type tempAlbumStructure struct {
	Meta    models.Album
	PicExt  string
	PicBase string
}

func main() {
	salt := "citrohallpassword"
	defaultpicbase := "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"
	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	r.Static("/data/music", "./data/music")

	store := cookie.NewStore([]byte("secret"))
	r.Use(sessions.Sessions("mysession", store))

	models.Migrate()
	db, _ := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})

	//r.GET("/ping", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "pong",
	//	})
	//})
	r.GET("/home", homeHandler)

	r.GET("/login", getLoginHandler)
	r.POST("/login", func(c *gin.Context) {
		n := c.PostForm("name")
		p := c.PostForm("pass")

		hasher := sha512.New()
		pwb := []byte(p + salt)
		hasher.Write(pwb)
		hpb := hasher.Sum(nil)
		hashedPW := hex.EncodeToString(hpb)

		var user models.User
		result := db.Where("name = ? AND password = ?", n, hashedPW).First(&user)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "the user does not exist.",
			})
		} else {
			session := sessions.Default(c)
			session.Set("userId", user.ID)
			session.Set("name", user.Name)
			session.Set("login", true)
			session.Save()

			//c.JSON(200, gin.H{
			//	"userId": session.Get("userId"),
			//	"name":   session.Get("name"),
			//	"login":  session.Get("login"),
			//})
			c.Redirect(http.StatusSeeOther, "/user")
		}
	})
	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Clear()
		session.Save()
		c.Redirect(http.StatusSeeOther, "/login")
	})

	r.GET("/music", func(c *gin.Context) {
		if !isLoggedIn(c) {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
		var musics []models.Music
		result := db.Order("id desc").Limit(100).Find(&musics)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "query error.",
			})
			return
		}

		var msts []tempMusicStructure
		var mst tempMusicStructure
		var picext string
		var picbase string
		for _, music := range musics {
			file4meta, _ := os.Open(music.Path)
			meta, _ := tag.ReadFrom(file4meta)
			if meta.Picture() != nil {
				picext = meta.Picture().Ext
				picbase = base64.StdEncoding.EncodeToString(meta.Picture().Data)
			} else {
				picext = "image/png"
				picbase = defaultpicbase
			}
			mst = tempMusicStructure{
				Meta:    music,
				PicExt:  picext,
				PicBase: picbase,
			}
			msts = append(msts, mst)
		}
		c.HTML(http.StatusOK, "music.html", gin.H{
			"title":         "Music",
			"musicresource": msts,
		})
	})
	r.GET("/music/:id", func(c *gin.Context) {
		if !isLoggedIn(c) {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
		var music models.Music
		musicId := c.Param("id")
		result := db.Where("id = ?", musicId).First(&music)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "the music does not exist.",
			})
			return
		}
		file4meta, _ := os.Open(music.Path)
		m, _ := tag.ReadFrom(file4meta)
		var picext string
		var picbase string
		if m.Picture() != nil {
			picext = m.Picture().Ext
			picbase = base64.StdEncoding.EncodeToString(m.Picture().Data)
		} else {
			picext = "image/png"
			picbase = defaultpicbase
		}
		c.HTML(http.StatusOK, "musicid.html", gin.H{
			"title":        "Music",
			"musicsrc":     "/" + music.Path,
			"musicalbum":   music.Album,
			"musictitle":   music.Title,
			"musicpicext":  picext,
			"musicpicdata": picbase,
		})
	})
	r.GET("/music/album", func(c *gin.Context) {
		if !isLoggedIn(c) {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
		var albums []models.Album
		result := db.Order("id desc").Limit(100).Find(&albums)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "query error.",
			})
			return
		}

		var asts []tempAlbumStructure
		var ast tempAlbumStructure
		var picext string
		var picbase string
		for _, album := range albums {
			file4meta, _ := os.Open(album.FirstTrackPath)
			meta, _ := tag.ReadFrom(file4meta)
			if meta.Picture() != nil {
				picext = meta.Picture().Ext
				picbase = base64.StdEncoding.EncodeToString(meta.Picture().Data)
			} else {
				picext = "image/png"
				picbase = defaultpicbase
			}
			ast = tempAlbumStructure{
				Meta:    album,
				PicExt:  picext,
				PicBase: picbase,
			}
			asts = append(asts, ast)
		}
		c.HTML(http.StatusOK, "album.html", gin.H{
			"title":         "Album",
			"musicresource": asts,
		})
	})
	r.GET("/music/album/:id", func(c *gin.Context) {
		if !isLoggedIn(c) {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
		var musics []models.Music
		albumId := c.Param("id")
		result := db.Where("album_id = ?", albumId).Find(&musics)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "the music does not exist.",
			})
			return
		}

		var msts []tempMusicStructure
		var mst tempMusicStructure
		var picext string
		var picbase string
		for _, music := range musics {
			file4meta, _ := os.Open(music.Path)
			meta, _ := tag.ReadFrom(file4meta)
			if meta.Picture() != nil {
				picext = meta.Picture().Ext
				picbase = base64.StdEncoding.EncodeToString(meta.Picture().Data)
			} else {
				picext = "image/png"
				picbase = defaultpicbase
			}
			mst = tempMusicStructure{
				Meta:    music,
				PicExt:  picext,
				PicBase: picbase,
			}
			msts = append(msts, mst)
		}
		c.HTML(http.StatusOK, "music.html", gin.H{
			"title":         "Album",
			"musicresource": msts,
		})
	})
	r.GET("/music/playlist", playlistHandler)
	r.GET("/music/search", searchHandler)

	r.GET("/signup", getSignupHandler)
	r.POST("/signup", func(c *gin.Context) {
		n := c.PostForm("name")
		p := c.PostForm("pass")
		hasher := sha512.New()
		pwb := []byte(p + salt)
		hasher.Write(pwb)
		hpb := hasher.Sum(nil)
		hashedPW := hex.EncodeToString(hpb)

		user := models.User{
			Name:     n,
			Password: hashedPW,
		}
		result := db.Create(&user)
		if result.Error != nil {
			c.Redirect(http.StatusInternalServerError, "/error")
		} else {
			session := sessions.Default(c)
			session.Set("userId", user.ID)
			session.Set("name", user.Name)
			session.Set("login", true)
			session.Save()

			// c.JSON(200, gin.H{
			// 	"userId": session.Get("userId"),
			// 	"name":   session.Get("name"),
			// 	"login":  session.Get("login"),
			// })
			c.Redirect(http.StatusSeeOther, "/user")
		}
	})

	r.GET("/user", func(c *gin.Context) {
		session := sessions.Default(c)
		if isLoggedIn(c) {
			c.HTML(http.StatusOK, "user.html", gin.H{
				"username": session.Get("name"),
			})
		} else {
			c.JSON(200, gin.H{
				"message": "please login",
			})
		}
	})

	r.GET("/upload", func(c *gin.Context) {
		if !isLoggedIn(c) {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
		c.HTML(http.StatusOK, "upload.html", gin.H{
			"title": "Home",
		})
	})
	r.POST("/upload", func(c *gin.Context) {
		if !isLoggedIn(c) {
			c.Redirect(http.StatusSeeOther, "/login")
			return
		}
		// Multipart form
		form, formerr := c.MultipartForm()

		if formerr != nil {
			c.String(http.StatusBadRequest, "get form err: %s", formerr.Error())
			return
		}

		files := form.File["audio"]
		if len(files) <= 0 {
			c.String(http.StatusBadRequest, "more than 1 file need to be uploaded.")
			return
		}

		// mkdir of tempDirDst if doesn't exist
		//session := sessions.Default(c)
		tempDirDst := "data/temp/_" + time.Now().String()
		if td, tderr := os.Stat(tempDirDst); os.IsNotExist(tderr) || !td.IsDir() {
			if e := os.Mkdir(tempDirDst, 0777); e != nil {
				c.String(http.StatusInternalServerError, "mkdir temp err: %s", e.Error())
				return
			}
			defer os.RemoveAll(tempDirDst)
		} else {
			c.String(http.StatusInternalServerError, "tempDir exists : %s", tderr.Error())
			return
		}

		// save files to tempDirDst
		for _, uploadedfile := range files {
			filename := filepath.Base(uploadedfile.Filename)
			tempSaveDst := tempDirDst + "/" + filename

			if err := c.SaveUploadedFile(uploadedfile, tempSaveDst); err != nil {
				c.String(http.StatusBadRequest, "upload file err: couldn't save the file to tempDirDst;%s", err.Error())
				return
			}

			//////
			file4meta, _ := os.Open(tempSaveDst)
			log.Println(file4meta)
			m, _ := tag.ReadFrom(file4meta) // read metadata.

			albumName := m.Album()
			if albumName == "" {
				albumName = "Unknown"
			}
			albumPath := "data/music/" + albumName
			albumdst := albumPath + "/" + filename // savedst includes own file name and extension.

			session := sessions.Default(c)
			// mkdir of album if doesn't exist
			if d, derr := os.Stat(albumPath); os.IsNotExist(derr) || !d.IsDir() {
				if e := os.Mkdir(albumPath, 0777); e != nil {
					c.String(http.StatusInternalServerError, "mkdir err: %s", e.Error())
					return
				}
				album := models.Album{
					UserId:         session.Get("userId").(uint),
					Path:           albumPath,
					FirstTrackPath: albumdst,

					// metadata
					Album:       m.Album(),
					Artist:      m.Artist(),
					AlbumArtist: m.AlbumArtist(),
					Composer:    m.Composer(),
					Genre:       m.Genre(),
					Year:        m.Year(),
				}
				if result := db.Create(&album); result.Error != nil {
					c.String(http.StatusInternalServerError, "album record error: %s", result.Error)
					return
				}
			}
			// Save file to savePath
			if _, ferr := os.Stat(albumdst); os.IsNotExist(ferr) {
				if err := c.SaveUploadedFile(uploadedfile, albumdst); err != nil {
					c.String(http.StatusBadRequest, "upload file err: couldn't save the file;%s", err.Error())
					return
				}
			} else {
				c.String(http.StatusBadRequest, "upload file err: the file already exists;")
				return
			}

			var album models.Album
			db.Where("Album = ?", m.Album()).First(&album)
			// insert metadata to models.Music record.
			t, tt := m.Track()
			d, dt := m.Disc()
			music := models.Music{
				UserId:  session.Get("userId").(uint),
				Path:    albumdst,
				AlbumId: album.ID,

				// metadata
				Title:       m.Title(),
				Album:       m.Album(),
				Artist:      m.Artist(),
				AlbumArtist: m.AlbumArtist(),
				Composer:    m.Composer(),
				Genre:       m.Genre(),
				Year:        m.Year(),

				Track:      t,
				TrackTotal: tt,
				Disc:       d,
				DiscTotal:  dt,

				Lyrics:  m.Lyrics(),
				Comment: m.Comment(),
			}
			if result := db.Create(&music); result.Error != nil {
				c.String(http.StatusInternalServerError, "music record error: %s", result.Error)
				return
			}
		}
		c.String(http.StatusOK, fmt.Sprintf("%d file(s) uploaded!", len(files)))
	})

	r.GET("/error", func(c *gin.Context) {
		c.JSON(500, gin.H{
			"message": "ERROR",
		})
	})

	r.Run(":80") // listen and serve on 0.0.0.0:80
}

func isLoggedIn(c *gin.Context) bool {
	session := sessions.Default(c)
	if session.Get("login") != nil {
		return true
	} else {
		return false
	}
}

func homeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "home.html", gin.H{
		"title": "Home",
	})
}

func getLoginHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title": "Login",
	})
}

func musicHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "music.html", gin.H{
		"title": "Music",
	})
}

func playlistHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "playlist.html", gin.H{
		"title": "Playlist",
	})
}

func searchHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "search.html", gin.H{
		"title": "Search",
	})
}

func getSignupHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "signup.html", gin.H{
		"title": "Signup",
	})
}
