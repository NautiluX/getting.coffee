package main

import (
	"fmt"
	"html/template"
	r "math/rand"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	fs := http.FileServer(http.Dir("./content"))
	http.Handle("/favicon.ico", fs)
	http.Handle("/icon.png", fs)
	http.Handle("/site.webmanifest", fs)
	http.HandleFunc("/", index)
	http.HandleFunc("/img/", serveGif)
	http.ListenAndServe(":8080", nil)
}

type Config struct {
	Title       string
	Description string
}

func index(res http.ResponseWriter, req *http.Request) {
	config := Config{"Getting coffee, brb ...", "Be right back!"}
	fp := path.Join("content", "index-template.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(res, config); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func serveGif(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Pragma-directive", "no-cache")
	res.Header().Set("Cache-directive", "no-cache")
	res.Header().Set("Cache-control", "no-store")
	res.Header().Set("Pragma", "no-cache")
	res.Header().Set("Expires", "0")
	var files []string

	root := "./content/img"
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(path, ".gif") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	r.Seed(time.Now().UnixNano())
	file := files[r.Intn(len(files))]
	fmt.Println(file)
	http.ServeFile(res, req, file)
}
