package main

import (
	"html/template"
	"log"
	"net/http"
	"os"

	_ "github.com/heroku/x/hmetrics/onload"
	"github.com/vmoret/shielded-tor-82171/pkg/hotjar"
	"github.com/vmoret/shielded-tor-82171/pkg/upload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		log.Fatalf("Failed parsing templates, %v", err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))

	http.Handle("/upload", &upload.Handler{
		Template:     tmpl,
		MaxMemory:    10 << 20, // Maximum upload of 10 MB files
		InputName:    "myFile",
		TemplateName: "upload.html",
		Uploader: hotjar.New(hotjar.Options{
			Layout: hotjar.DefaultTimeLayout,
		}),
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.ExecuteTemplate(w, "index.html", nil); err != nil {
			log.Printf("Failed executing template, %v", err)
		}
	})
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
