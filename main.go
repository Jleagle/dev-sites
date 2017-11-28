package main

import (
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi"
	"github.com/kr/pretty"
)

func main() {

	r := chi.NewRouter()
	r.Get("/", listHandler)

	// Assets
	workDir, _ := os.Getwd()
	filesDir := filepath.Join(workDir, "assets")
	fileServer(r, "/assets", http.Dir(filesDir))

	log.Fatal(http.ListenAndServe(":8085", r))
}

func listHandler(w http.ResponseWriter, r *http.Request) {

	// Get paths
	file, err := ioutil.ReadFile("./config.json")
	if err != nil {
		pretty.Print("error 1:", err)
	}

	configuration := Configuration{}
	err = json.Unmarshal(file, &configuration)
	if err != nil {
		pretty.Print("error 2:", err)
	}

	//
	var projects []Project

	for _, v := range configuration.Folders.PHP {

		files, err := ioutil.ReadDir(v)
		if err != nil {
			pretty.Print("error 5:", err)
		}

		for _, vv := range files {

			project := Project{}
			project.Path = v + "/" + vv.Name()

			projects = append(projects, project)
		}
	}

	// Return a template
	t, err := template.New("t").ParseFiles("./template.html")
	if err != nil {
		pretty.Print("error 3:", err)
	}

	vars := TemplateVars{}
	vars.Projects = projects

	pretty.Print(vars)

	err = t.ExecuteTemplate(w, "list", vars)
	if err != nil {
		pretty.Print("error 4:", err)
	}
}

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

type Configuration struct {
	Folders FoldersStruct `json:"folders"`
}

type FoldersStruct struct {
	Go  []string `json:"go"`
	PHP []string `json:"php"`
}

type TemplateVars struct {
	Projects []Project
}

type Project struct {
	Path string
}
