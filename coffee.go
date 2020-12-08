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

type CoffeeEngine struct {
	config Config
}

func main() {
	config := Config{os.Getenv("COFFEE_TITLE"), os.Getenv("COFFEE_DESCRIPTION"), os.Getenv("COFFEE_URL"), os.Getenv("COFFEE_GIFDIR")}
	config.HealthCheck()
	c := CoffeeEngine{config}
	c.StartServer()
}

type Config struct {
	Title       string
	Description string
	Url         string
	GifDir      string
}

func (c *Config) HealthCheck() {
	if c.GifDir == "" {
		panic(fmt.Errorf("environment variable COFFEE_GIFDIR not set. Direct it to a folder with gifs"))
	}
	if c.Title == "" {
		panic(fmt.Errorf("environment variable COFFEE_TITLE not set."))
	}
	if c.Description == "" {
		panic(fmt.Errorf("environment variable COFFEE_DESCRIPTION not set."))
	}
	if c.Url == "" {
		panic(fmt.Errorf("environment variable COFFEE_URL not set. This should point to the public URL of this service to allow unfurling in Slack etc."))
	}
	if _, err := os.Stat(c.GifDir); os.IsNotExist(err) {
		panic(fmt.Errorf("Path %s doesn't exist", c.GifDir))
	}
}

func (c *CoffeeEngine) StartServer() {
	fs := http.FileServer(http.Dir("./content"))
	http.Handle("/favicon.ico", fs)
	http.Handle("/icon.png", fs)
	http.Handle("/site.webmanifest", fs)
	http.HandleFunc("/", c.index)
	http.HandleFunc("/img/", c.serveGif)
	http.ListenAndServe(":8080", nil)
}

func (c *CoffeeEngine) index(res http.ResponseWriter, req *http.Request) {
	fp := path.Join("content", "index-template.html")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(res, c.config); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
	}
}

func (c *CoffeeEngine) serveGif(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Pragma-directive", "no-cache")
	res.Header().Set("Cache-directive", "no-cache")
	res.Header().Set("Cache-control", "no-store")
	res.Header().Set("Pragma", "no-cache")
	res.Header().Set("Expires", "0")
	var files []string

	root := c.config.GifDir
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
