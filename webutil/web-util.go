package webutil

import (
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/wizards-0/go-pins/logger"
)

type HttpResponse struct {
	Status int `json:"status"`
	Body   any `json:"body"`
	Error  error
}

type ErrorResponse struct {
	Body         any    `json:"body"`
	ErrorMessage string `json:"errorMessage"`
}

func RegisterRequestBodyHandler[T any](
	mux *http.ServeMux,
	pattern string,
	bodyFactory func() *T,
	fn func(reqBody *T) HttpResponse,
) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if reqBody, err := getReqBody(w, r, bodyFactory); err == nil {
			resp := fn(reqBody)
			handleResponse(w, resp)
		}
	})
}

func RegisterRequestBodyAndParamsHandler[T any](
	mux *http.ServeMux,
	pattern string,
	bodyFactory func() *T,
	fn func(reqBody *T, queryParams map[string]string, pathParams map[string]string) HttpResponse,
	pathParamNames ...string,
) {
	validatePathParamConfig(pattern, pathParamNames)
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		if reqBody, err := getReqBody(w, r, bodyFactory); err == nil {
			resp := fn(reqBody, getQueryParams(r.URL.Query()), getPathParams(r, pathParamNames))
			handleResponse(w, resp)
		}
	})
}

func RegisterParamsHandler(
	mux *http.ServeMux,
	pattern string,
	fn func(queryParams map[string]string, pathParams map[string]string) HttpResponse,
	pathParamNames ...string,
) {
	validatePathParamConfig(pattern, pathParamNames)
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		resp := fn(getQueryParams(r.URL.Query()), getPathParams(r, pathParamNames))
		handleResponse(w, resp)
	})
}

func RegisterSingleFileServer(
	mux *http.ServeMux,
	pattern string,
	filePath string,
) {
	mux.HandleFunc(pattern, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filePath)
	})
}

func getReqBody[T any](w http.ResponseWriter, r *http.Request, bodyFactory func() *T) (*T, error) {
	reqBody := bodyFactory()
	if err := json.NewDecoder(r.Body).Decode(reqBody); err != nil {
		logger.Error(err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{ErrorMessage: err.Error()})
		return reqBody, err
	}
	return reqBody, nil
}

func handleResponse(w http.ResponseWriter, resp HttpResponse) {
	if resp.Error != nil {
		if resp.Status != 0 {
			w.WriteHeader(resp.Status)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		_ = json.NewEncoder(w).Encode(ErrorResponse{ErrorMessage: resp.Error.Error(), Body: resp.Body})
	} else {
		if resp.Status != 0 {
			w.WriteHeader(resp.Status)
		} else {
			w.WriteHeader(http.StatusOK)
		}
		// Only fails for exotic values like channels, functions. But who would return that from a rest endpoint, right? ... right?
		// For values like graphs with cycles, it straight up crashes, lol. So yeah, don't do that, terrible.
		_ = json.NewEncoder(w).Encode(resp.Body)
	}
}

func getPathParams(r *http.Request, pathParamNames []string) map[string]string {
	pathParams := map[string]string{}
	for _, pName := range pathParamNames {
		pathParams[pName] = r.PathValue(pName)
	}
	return pathParams
}

func getQueryParams(q url.Values) map[string]string {
	queryParams := map[string]string{}
	for pName, values := range q {
		queryParams[pName] = strings.Join(values, ",")
	}
	return queryParams
}

func validatePathParamConfig(pattern string, pathParamNames []string) {
	if strings.Count(pattern, "{") != len(pathParamNames) {
		msg := "No. of parameters in path, do not match parameter names count"
		logger.Error(msg)
		panic(msg)
	}
}
