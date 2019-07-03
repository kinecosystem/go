package httpjson

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/stellar/go/support/errors"
)

type contentType int

const (
	JSON contentType = iota
	HALJSON
)

// renderToString renders the provided data as a json string
func renderToString(data interface{}, isIndentedJSON bool) ([]byte, error) {
	if isIndentedJSON {
		return json.MarshalIndent(data, "", "  ")
	}

	return json.Marshal(data)
}

// Render write data to w, after marshalling to json. The response header is
// set based on cType.
func Render(w http.ResponseWriter, data interface{}, cType contentType, isIndentedJSON bool) {
	js, err := renderToString(data, isIndentedJSON)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Disposition", "inline")
	if cType == HALJSON {
		w.Header().Set("Content-Type", "application/hal+json; charset=utf-8")
	} else {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
	}
	w.Write(js)
}

var ErrBadRequest = errors.New("bad request")

// read decodes a json text from r into v.
func read(r io.Reader, v interface{}) error {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	err := dec.Decode(v)
	if err != nil {
		return ErrBadRequest
	}

	return nil
}
