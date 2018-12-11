package rest

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/pressly/chi"
)

func getQueryParamInt(req *http.Request, name string, def int) (int, error) {
	qs := req.URL.Query().Get(name)
	if qs == "" {
		return def, nil
	}
	value, err := strconv.ParseInt(qs, 10, 32)
	if err != nil {
		return -1, fmt.Errorf("%s must be an integer,but it's %v", name, qs)
	}
	return int(value), nil
}

func getQueryParamInt64(req *http.Request, name string, def int64) (int64, error) {
	qs := req.URL.Query().Get(name)
	if qs == "" {
		return def, nil
	}
	value, err := strconv.ParseInt(qs, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("%s must be an int64,but it's %v", name, qs)
	}
	return int64(value), nil
}

func getQueryParamString(req *http.Request, name string, def string) string {
	qs := req.URL.Query().Get(name)
	if qs == "" {
		return def
	}
	return qs
}

func getURLParamUint(req *http.Request, name string) (uint, error) {
	qs := chi.URLParam(req, "id")
	if qs == "" {
		return 0, fmt.Errorf("could not found %s", name)
	}
	value, err := strconv.ParseUint(qs, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("%s must be an uint,but it's %v", name, qs)
	}
	return uint(value), nil
}
