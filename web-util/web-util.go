package webutil

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/wizards-0/go-pins/logger"
)

/*
- Request Types
  - GET
  - POST

- Input Options
  - Query Param
  - Path Param
  - Request Body (w/ Query | Path Param)

- Output Options
  - json
  - status only
  - error
  - string
*/
type HttpResponse struct {
	Status int
	Body   any
	Error  error
}

type ErrorResponse struct {
	Body         any    `json:"body"`
	ErrorMessage string `json:"errorMessage"`
}

func RegisterPost(pattern string, reqBody *any, fn func(reqBody *any) HttpResponse) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
			logger.Error(err)
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(ErrorResponse{ErrorMessage: err.Error()})
		}
		resp := fn(reqBody)
		if resp.Error != nil {
			if resp.Status != 0 {
				w.WriteHeader(resp.Status)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			json.NewEncoder(w).Encode(ErrorResponse{ErrorMessage: resp.Error.Error()})
		} else {
			if resp.Status != 0 {
				w.WriteHeader(resp.Status)
			} else {
				w.WriteHeader(http.StatusOK)
			}
			json.NewEncoder(w).Encode(resp.Body)
		}
	}
	http.HandleFunc(pattern, handler)
}

func RegisterGet(pattern string, fn func(queryParams url.Values, pathParams ...any) HttpResponse, pathParamNames ...string) {
	handler := func(w http.ResponseWriter, r *http.Request) {
		pathParams := []any{}
		for _, pName := range pathParamNames {
			pathParams = append(pathParams, r.PathValue(pName))
		}
		resp := fn(r.URL.Query(), pathParams...)
		if resp.Error != nil {
			if resp.Status != 0 {
				w.WriteHeader(resp.Status)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			json.NewEncoder(w).Encode(ErrorResponse{ErrorMessage: resp.Error.Error()})
		} else {
			if resp.Status != 0 {
				w.WriteHeader(resp.Status)
			} else {
				w.WriteHeader(http.StatusOK)
			}
			json.NewEncoder(w).Encode(resp.Body)
		}
	}
	http.HandleFunc(pattern, handler)
}
