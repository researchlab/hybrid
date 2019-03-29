package render

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

var statusCtxKey = &contextKey{"Status"}

// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}

func (k *contextKey) String() string {
	return "chi render context value " + k.name
}

// Status sets status into request context.
func Status(r *http.Request, status int) {
	*r = *r.WithContext(context.WithValue(r.Context(), statusCtxKey, status))
}

// Bind alias defaultBind
var Bind = defaultBind

// defaultBind is a short-hand method for decoding a JSON request body.
func defaultBind(r io.Reader, v interface{}) error {
	defer io.Copy(ioutil.Discard, r)
	return json.NewDecoder(r).Decode(v)
}

// Head set head key value
func Head(w http.ResponseWriter, name, value string) {
	w.Header().Set(name, value)
}

// JSON json encode
func JSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	if status, ok := r.Context().Value(statusCtxKey).(int); ok {
		w.WriteHeader(status)
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// TextJSON  write json text to writer
func TextJSON(w http.ResponseWriter, r *http.Request, v interface{}) {
	w.Header().Set("Content-Type", "text/json; ; charset=utf-8")
	w.Header().Set("Cache-Control", "no-cache")
	if status, ok := r.Context().Value(statusCtxKey).(int); ok {
		w.WriteHeader(status)
	}

	fmt.Fprint(w, v)
}
