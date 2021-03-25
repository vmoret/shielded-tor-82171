package hotjar

import (
	"fmt"
	"mime/multipart"
	"net/http"
)

type Options struct {
	Layout    string   `yaml:"layout"`
	Questions []string `yaml:"questions"`
}

type Hotjar struct {
	opts Options
}

func New(opts Options) *Hotjar {
	return &Hotjar{
		opts: opts,
	}
}

func (h *Hotjar) UploadFile(r *http.Request, file multipart.File, handler *multipart.FileHeader) error {
	var entries []Entry
	{
		r, err := newReader(file, handler.Size)
		if err != nil {
			return err
		}
		defer r.Close()
		r.Layout = h.opts.Layout
		r.Questions = h.opts.Questions
		entries, err = r.ReadAll()
		if err != nil {
			return err
		}
	}

	for i, entry := range entries {
		fmt.Println(i, entry)
	}

	return nil
}
