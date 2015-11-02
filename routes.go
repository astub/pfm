package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

func NewRouter(db DataHandler) *mux.Router {

	fe := FrontEnd{DataHandler: db}
	fe.CookieHandler = securecookie.New(securecookie.GenerateRandomKey(64), securecookie.GenerateRandomKey(32))

	var routes = Routes{
		Route{"Index", "GET", "/", Index},
		Route{"EventNews", "GET", "/eventnews", EventNewsPage},
		Route{"Media", "GET", "/media", MediaPage},
		Route{"ExhibitsPage", "GET", "/exhibits", ExhibitsPage},
		Route{"Resources", "GET", "/resources", Resources},
		Route{"InfoPage", "GET", "/info", InfoPage},

		Route{"GetPosts", "GET", "/get_posts", fe.GetPosts},
		Route{"ShowPost", "GET", "/post/{id}", fe.ShowPost},

		Route{"ImgUpload", "POST", "/upload_img", ImgUpload},
		Route{"AddPost", "POST", "/new_post", fe.NewPost},
		Route{"UpdatePost", "POST", "/update_post", fe.UpdatePost},
	}

	router := mux.NewRouter().StrictSlash(true)

	for _, route := range routes {
		router.Methods(route.Method).Path(route.Pattern).Name(route.Name).Handler(route.HandlerFunc)
	}
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./www/")))

	return router

}

type DataHandler interface {
	UpdatePost(pst Post) (err error)
	InsertPost(pst Post) (err error)
	GetPosts() (psts Posts, err error)
	GetPost(string) (Post, error)
}

type FrontEnd struct {
	DataHandler
	CookieHandler *securecookie.SecureCookie
}

func render(w http.ResponseWriter, filename string, data interface{}) {
	var err error
	tmpl := template.New("")

	if tmpl, err = template.ParseFiles("templates/layout.tmpl", filename); err != nil {
		fmt.Println(err)
		return
	}

	if err = tmpl.Execute(w, data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	render(w, "templates/home.tmpl", nil)
}

func Resources(w http.ResponseWriter, r *http.Request) {
	render(w, "templates/resources.tmpl", nil)
}

func EventNewsPage(w http.ResponseWriter, r *http.Request) {
	render(w, "templates/eventnews.tmpl", nil)
}

func InfoPage(w http.ResponseWriter, r *http.Request) {
	render(w, "templates/info.tmpl", nil)
}

func ExhibitsPage(w http.ResponseWriter, r *http.Request) {
	render(w, "templates/exhibits.tmpl", nil)
}

func MediaPage(w http.ResponseWriter, r *http.Request) {
	render(w, "templates/media.tmpl", nil)
}

type Page struct {
	PageData interface{}
}

func (fe FrontEnd) ShowPost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	pst, err := fe.DataHandler.GetPost(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pst.Details = strings.Replace(pst.Details, "\r\n", "<br />", -1)
	pst.Specs = strings.Replace(pst.Specs, "\r\n", "<br />", -1)
	pst.Links = strings.Replace(pst.Links, "\r\n", "<br />", -1)
	pst.DetailsHTML = template.HTML(pst.Details)
	pst.SpecsHTML = template.HTML(pst.Specs)
	pst.LinksHTML = template.HTML(pst.Links)

	page := &Page{PageData: pst}
	render(w, "templates/postview.tmpl", page)

}

func (fe FrontEnd) GetPosts(w http.ResponseWriter, r *http.Request) {
	psts, err := fe.DataHandler.GetPosts()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(psts); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
func (fe FrontEnd) UpdatePost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pst_t, err := time.Parse("01/02/2006", r.Form.Get("datepicker"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pst := Post{
		Title:   r.Form.Get("title"),
		Date:    pst_t,
		Type:    r.Form.Get("type"),
		Details: r.Form.Get("details"),
		Specs:   r.Form.Get("specs"),
		Links:   r.Form.Get("links"),
		File:    r.Form.Get("file"),
	}

	log.Println(pst)

	err = fe.DataHandler.UpdatePost(pst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		fmt.Fprintf(w, "Post Updated!")
	}

}

func (fe FrontEnd) NewPost(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pst_t, err := time.Parse("01/02/2006", r.Form.Get("datepicker"))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	pst := Post{
		Title:   r.Form.Get("title"),
		Date:    pst_t,
		Type:    r.Form.Get("type"),
		Details: r.Form.Get("details"),
		Specs:   r.Form.Get("specs"),
		Links:   r.Form.Get("links"),
		File:    r.Form.Get("file"),
	}

	log.Println(pst)

	err = fe.DataHandler.InsertPost(pst)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	} else {
		fmt.Fprintf(w, "Post added!")
	}

}

func ImgUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(100000)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	f_id := r.FormValue("f_id")

	//get a ref to the parsed multipart form
	m := r.MultipartForm

	//get the *fileheaders
	files := m.File["upl"]
	for i, _ := range files {

		//for each fileheader, get a handle to the actual file
		file, err := files[i].Open()
		defer file.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		path := "." + string(filepath.Separator) + "www" + string(filepath.Separator) + "uploads" + string(filepath.Separator) + f_id

		//create destination file making sure the path is writeable.
		os.Mkdir(path, 0777)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		dst, err := os.Create(path + string(filepath.Separator) + files[i].Filename)
		defer dst.Close()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		//copy the uploaded file to the destination file
		if _, err := io.Copy(dst, file); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}

	//display success message.

}
