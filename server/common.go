package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Standard response statuses.
const (
	StatusOk   = "ok"
	StatusFail = "fail"
)

// paramError returns a prepared 422 Unprocessable Entity response
func paramError(name string) *Response {
	return NewResponse(http.StatusUnprocessableEntity, fmt.Sprintf("invalid or missing parameter: %v", name), nil)
}

func getParam(r *http.Request, key string) (param string, err error) {
	vars := mux.Vars(r)
	param, ok := vars[key]
	if !ok {
		return param, fmt.Errorf("invalid or missing parameter")
	}
	return
}

func getInt64Param(r *http.Request, key string) (id int64, err error) {
	param, err := getParam(r, key)
	if err != nil {
		return
	}
	return strconv.ParseInt(param, 10, 64)
}

func logHandler(s *Server, h http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Println(r.Method, r.URL)
		h.ServeHTTP(w, r)
	})
}
