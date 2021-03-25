package hotjar

import (
	"archive/zip"
	"encoding/csv"
	"errors"
	"io"
	"log"
	"net/url"
	"strconv"
	"time"
)

var (
	// ErrInvalidHeader is returned when a required name was not found in the
	// CSV header.
	ErrInvalidHeader = errors.New("invalid header")
)

const (
	// DefaultTimeLayout is the default layout of time strings
	DefaultTimeLayout = "2006-01-02 15:04:05"
)

type reader struct {
	r   io.ReadCloser
	csv *csv.Reader

	// Layout holds the layout of time strings
	Layout string

	// Questions to expect in the data
	Questions []string
}

func newReader(r io.ReaderAt, size int64) (*reader, error) {
	zipr, err := zip.NewReader(r, size)
	if err != nil {
		return nil, err
	}
	rd, err := zipr.File[0].Open()
	if err != nil {
		return nil, err
	}
	return &reader{
		r:         rd,
		csv:       csv.NewReader(rd),
		Layout:    DefaultTimeLayout,
		Questions: []string{},
	}, nil
}

func (r *reader) Close() error {
	return r.r.Close()
}

// Entry represents an entry of the Hotjar file.
type Entry struct {
	Number    int
	Source    *url.URL
	Submitted time.Time
	Answers   []string
}

func (e *Entry) unmarshalRecord(positions []int, values []string, layout string) (err error) {
	if e.Number, err = strconv.Atoi(values[positions[0]]); err != nil {
		return
	}
	if e.Source, err = url.Parse(values[positions[1]]); err != nil {
		return
	}
	if e.Submitted, err = time.Parse(layout, values[positions[2]]); err != nil {
		return
	}
	e.Answers = make([]string, len(positions)-3)
	for i, pos := range positions[3:] {
		e.Answers[i] = values[pos]
	}
	return nil
}

func (r *reader) ReadAll() ([]Entry, error) {
	keys, err := r.csv.Read()
	if err != nil {
		return nil, err
	}

	cols := append([]string{"Number", "Source URL", "Date Submitted"}, r.Questions...)
	positions, ok := findStrings(keys, cols...)
	if !ok {
		return nil, ErrInvalidHeader
	}

	entries := make([]Entry, 0)

	for {
		values, err := r.csv.Read()
		if err == io.EOF {
			return entries, nil
		}
		if err != nil {
			return nil, err
		}

		entry := new(Entry)
		if err := entry.unmarshalRecord(positions, values, r.Layout); err != nil {
			return nil, err
		}
	}
}

func findStrings(haystack []string, needles ...string) ([]int, bool) {
	positions := make([]int, len(needles))
	for i, needle := range needles {
		found := false
		for position, s := range haystack {
			if s == needle {
				positions[i] = position
				found = true
				break
			}
		}
		if !found {
			log.Printf("didn't find %s in %+v", needle, haystack)
			return nil, false
		}
	}
	return positions, true
}
