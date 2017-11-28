package main

import (
	"html/template"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/go-chi/chi"
)

func main() {

	r := chi.NewRouter()
	r.Get("/", listHandler)
	//r.Get("/logout", logoutHandler)

	// Assets
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "assets")
	fileServer(r, "/assets", http.Dir(filesDir))

	log.Fatal(http.ListenAndServe(":8084", r))
}

func returnTemplate(w http.ResponseWriter, r *http.Request, page string, err error) {

	// todo, log errors

	// Get current app path
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		//logger.Err("No caller information")
	}
	folder := path.Dir(file)

	// Load templates needed
	always := []string{
		folder + "/templates/header.html",
		folder + "/templates/footer.html",
		folder + "/templates/" + page + ".html",
		folder + "/templates/includes/album.html",
	}

	t, err := template.New("t").ParseFiles(always...)
	if err != nil {
		//logger.ErrExit(err.Error())
	}

	// Write a respone
	err = t.ExecuteTemplate(w, page, nil)
	if err != nil {
		//logger.ErrExit(err.Error())
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func fileServer(r chi.Router, path string, root http.FileSystem) {

	if strings.ContainsAny(path, "{}*") {
		//logger.ErrExit("FileServer does not permit URL parameters.")
	}

	fs := http.StripPrefix(path, http.FileServer(root))

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fs.ServeHTTP(w, r)
	}))
}
