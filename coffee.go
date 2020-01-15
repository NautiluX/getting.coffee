package main

import (
	"fmt"
	r "math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	fs := http.FileServer(http.Dir("./content"))
	http.Handle("/", fs)
	http.HandleFunc("/img/", readme)
	http.ListenAndServe(":8080", nil)
}

func readme(res http.ResponseWriter, req *http.Request) {
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
