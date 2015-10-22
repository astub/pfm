package main

import (
	"os"
	//"fmt"
	"database/sql"
	"image"
	"io/ioutil"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "image/jpeg"
	_ "image/png"
)

type DB struct {
	*sql.DB
}

func NewOpen(dt string, c string) (DB, error) {
	db, err := sql.Open(dt, c)
	return DB{db}, err
}

type Image struct {
	Path   string `json:"path"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type Post struct {
	Title   string    `json:"title"`
	Date    time.Time `json:"date"`
	Type    string    `json:"type"`
	Details string    `json:"detail"`
	Specs   string    `json:"specs"`
	Links   string    `json:"links"`
	File    string    `json:"file"`
	Images  []Image   `json:"imgs"`
}

type Posts []Post

func (d DB) GetPost(id string) (pst Post, err error) {
	var query = "SELECT type, date, title, detail, f_id FROM posts WHERE f_id=?;"
	var t string
	err = d.QueryRow(query, id).Scan(&pst.Type, &t, &pst.Title, &pst.Details, &pst.File)
	pst.Date, _ = time.Parse("2006-01-02 15:04:05", t)
	if err != nil {
		return
	}

	var path = "./www/uploads/" + pst.File
	files, _ := ioutil.ReadDir(path)
	for _, f := range files {
		width, height := ImageDimension(path + "/" + f.Name())
		if width != 0 {
			var img Image
			img.Path = "/uploads/" + pst.File + "/" + f.Name()
			img.Width = width
			img.Height = height
			pst.Images = append(pst.Images, img)
		}
	}

	return
}

func (d DB) UpdatePost(pst Post) (err error) {
	tx, err := d.Begin()
	if err != nil {
		return
	}

	//var query = "INSERT INTO posts (type,date,title,detail,spec,f_id,urllink) VALUES (?, ?, ?, ?, ?, ?, ?);"
	var query = "UPDATE posts SET type=?, date=?, title=?, detail=?, spec=?, urllink=? WHERE f_id=?;"

	stmt, err := tx.Prepare(query)
	if err != nil {
		return
	}

	defer func() {
		if err == nil {
			log.Println("Commit Item")
			tx.Commit()
		} else {
			log.Println("RollBack Item")
			tx.Rollback()
		}
		stmt.Close()
	}()

	_, err = stmt.Exec(pst.Type, pst.Date, pst.Title, pst.Details, pst.Specs, pst.Links, pst.File)

	return
}

func (d DB) InsertPost(pst Post) (err error) {
	tx, err := d.Begin()
	if err != nil {
		return
	}

	var query = "INSERT INTO posts (type,date,title,detail,spec,f_id,urllink) VALUES (?, ?, ?, ?, ?, ?, ?);"

	stmt, err := tx.Prepare(query)
	if err != nil {
		return
	}

	defer func() {
		if err == nil {
			log.Println("Commit Item")
			tx.Commit()
		} else {
			log.Println("RollBack Item")
			tx.Rollback()
		}
		stmt.Close()
	}()

	_, err = stmt.Exec(pst.Type, pst.Date, pst.Title, pst.Details, pst.Specs, pst.File, pst.Links)

	return
}

func (d DB) GetPosts() (psts Posts, err error) {
	var query = "SELECT type, date, title, detail, f_id, spec, urllink FROM posts;"

	rows, err := d.Query(query)
	if err != nil {
		return
	}

	defer rows.Close()

	for rows.Next() {
		var p Post
		var t string
		rows.Scan(&p.Type, &t, &p.Title, &p.Details, &p.File, &p.Specs, &p.Links)
		p.Date, _ = time.Parse("2006-01-02 15:04:05", t)

		var path = "./www/uploads/" + p.File
		files, _ := ioutil.ReadDir(path)
		for _, f := range files {
			width, height := ImageDimension(path + "/" + f.Name())
			if width != 0 {
				var img Image
				img.Path = "/uploads/" + p.File + "/" + f.Name()
				img.Width = width
				img.Height = height
				p.Images = append(p.Images, img)
			}
		}

		psts = append(psts, p)
	}

	return
}

func ImageDimension(imagePath string) (int, int) {
	file, err := os.Open(imagePath)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "%v\n", err)
		return 0, 0
	}

	defer file.Close()

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		//fmt.Fprintf(os.Stderr, "%s: %v\n", imagePath, err)
		return 0, 0
	}
	return image.Width, image.Height
}
