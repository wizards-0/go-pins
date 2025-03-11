package webutil

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/wizards-0/go-pins/logger"
)

func GetHttpHandleFunc(fn func(url.Values, string) (any, error)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		reqBodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error(err)
			//TODO: Add global error handling logic, return proper status code
		} else {
			reqBody := string(reqBodyBytes)
			fmt.Println()
			queryParams := r.URL.Query()
			resp, handlerError := fn(queryParams, reqBody)
			if handlerError != nil {
				log.Println(handlerError)
			} else if resp != nil {
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			}
		}
	}
}
