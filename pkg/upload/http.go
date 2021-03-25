package upload

import (
	"html/template"
	"log"
	"net/http"
)

type Handler struct {
	Uploader
	Template     *template.Template
	MaxMemory    int64
	InputName    string
	TemplateName string
	Data         interface{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		if err := h.Template.ExecuteTemplate(w, h.TemplateName, h.Data); err != nil {
			log.Printf("Failed to execute template %q: %v", h.TemplateName, err)
		}
	case http.MethodPost:
		if err := h.serveHTTP(w, r); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	default:
		w.Header().Add("Allow", http.MethodGet)
		w.Header().Add("Allow", http.MethodPost)
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (h *Handler) serveHTTP(w http.ResponseWriter, r *http.Request) error {
	if err := r.ParseMultipartForm(h.MaxMemory); err != nil {
		return err
	}
	file, handler, err := r.FormFile(h.InputName)
	if err != nil {
		return err
	}
	defer file.Close()
	return h.Uploader.UploadFile(r, file, handler)
}
