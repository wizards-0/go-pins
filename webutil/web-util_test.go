package webutil

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wizards-0/go-pins/logger"
)

const CONTENT_TYPE_JSON = "application/json"

var queryParamsTester = func(q map[string]string, pathParams map[string]string) HttpResponse {
	return HttpResponse{
		Status: 200,
		Body:   q,
	}
}

var log = bytes.Buffer{}

func setup() {
	logger.SetLogLevel(logger.LOG_LEVEL_DEBUG)
	logger.SetWriter(&log, &log, &log, &log)
}

func TestQueryParams(t *testing.T) {
	setup()
	assert := assert.New(t)

	mux := http.NewServeMux()
	RegisterParamsHandler(mux, "/queryParams", queryParamsTester)
	s := httptest.NewServer(mux)
	defer s.Close()
	baseUrl := s.URL + "/queryParams"
	respBody := getResponseMap(http.Get(baseUrl + "?id=5"))

	//Test if param is absent in query, its not present in query map {
	_, paramExists := respBody["name"]
	assert.False(paramExists)
	// }

	//Test even if param is numeric, its returned as string {
	assert.Equal("5", respBody["id"])
	// }

	//Test if no params are provided, call still works. But query params map is empty {
	respBody = getResponseMap(http.Get(baseUrl))
	assert.Equal(0, len(respBody))
	// }
}

var pathParamTester = func(q map[string]string, pathParams map[string]string) HttpResponse {
	return HttpResponse{
		Status: 200,
		Body:   pathParams,
	}
}

func TestPathParams(t *testing.T) {
	setup()
	assert := assert.New(t)

	mux := http.NewServeMux()
	RegisterParamsHandler(mux, "/pathParams/{id}/{name}", pathParamTester, "id", "error")
	s := httptest.NewServer(mux)
	defer s.Close()
	baseUrl := s.URL + "/pathParams"

	//Test if not all path params are provided, it returns 404 {
	resp, _ := http.Get(baseUrl + "/5")
	assert.Equal(http.StatusNotFound, resp.StatusCode)
	// }

	respBody := getResponseMap(http.Get(baseUrl + "/5/jojo"))
	logger.Debug(respBody)
	//Test valid parameters are passed to handler func {
	assert.Equal("5", respBody["id"])
	//}

	//Test invalid parameters are missing from the params map handler func {
	_, paramExists := respBody["name"]
	assert.False(paramExists)
	//}
}

func TestPathParamsValidation(t *testing.T) {
	logger.SetLogLevel(logger.LOG_LEVEL_NONE)
	assert := assert.New(t)
	mux := http.NewServeMux()
	defer func() {
		r := recover()
		assert.NotNil(t, r)
	}()
	RegisterParamsHandler(mux, "/pathParamsMismatch/{id}/{name}/", pathParamTester, "id")
	setup()
}

type RequestBody struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

var requestBodyFactory = func() *RequestBody {
	return &RequestBody{}
}

func TestRequestBody(t *testing.T) {
	setup()
	assert := assert.New(t)

	var requestBodyTester = func(body *RequestBody) HttpResponse {
		return HttpResponse{
			Status: 200,
			Body:   body,
		}
	}

	mux := http.NewServeMux()
	RegisterRequestBodyHandler(mux, "/requestBody", requestBodyFactory, requestBodyTester)
	s := httptest.NewServer(mux)
	defer s.Close()
	baseUrl := s.URL + "/requestBody"

	callApi := func(payload string) *RequestBody {
		buf := bytes.Buffer{}
		buf.Write([]byte(payload))
		resp, _ := http.Post(baseUrl, CONTENT_TYPE_JSON, &buf)
		body := requestBodyFactory()
		json.NewDecoder(resp.Body).Decode(body)
		return body
	}

	//Test request body has all correct fields {
	body := callApi(`{"id":"5","name":"jojo"}`)
	assert.Equal("5", body.Id)
	assert.Equal("jojo", body.Name)
	// }

	//Test request body has all correct fields, but not all fields are present {
	body = callApi(`{"id":"5"}`)
	assert.Equal("5", body.Id)
	assert.Equal("", body.Name)
	// }

	//Test request body has some correct, some incorrect fields {
	body = callApi(`{"id":"5","address":"blue moon"}`)
	assert.Equal("5", body.Id)
	assert.Equal("", body.Name)
	// }

	//Test request body has all incorrect fields {
	body = callApi(`{"address":"blue moon"}`)
	assert.Equal("", body.Id)
	assert.Equal("", body.Name)
	// }

	logger.SetLogLevel(logger.LOG_LEVEL_NONE)
	//Test response returns 400 Bad request for invalid json
	buf := bytes.Buffer{}
	buf.Write([]byte(`{"id":"5}`))
	resp, _ := http.Post(baseUrl, CONTENT_TYPE_JSON, &buf)
	assert.Equal(http.StatusBadRequest, resp.StatusCode)
	// }

	//Test response returns 400 Bad request for empty string as request body
	buf = bytes.Buffer{}
	buf.Write([]byte(""))
	resp, _ = http.Post(baseUrl, CONTENT_TYPE_JSON, &buf)
	assert.Equal(http.StatusBadRequest, resp.StatusCode)
	// }
	setup()
}

func TestRequestBodyAndParams(t *testing.T) {
	setup()
	assert := assert.New(t)

	var requestBodyTester = func(reqBody *RequestBody, queryParams map[string]string, pathParams map[string]string) HttpResponse {
		return HttpResponse{
			Status: 200,
			Body: map[string]any{
				"req":         reqBody,
				"queryParams": queryParams,
				"pathParams":  pathParams,
			},
		}
	}

	mux := http.NewServeMux()
	RegisterRequestBodyAndParamsHandler(mux, "/requestBodyAndParams/{id}", requestBodyFactory, requestBodyTester, "id")
	s := httptest.NewServer(mux)
	defer s.Close()
	baseUrl := s.URL + "/requestBodyAndParams/5?name=jojo"

	callApi := func(payload string) map[string]any {
		buf := bytes.Buffer{}
		buf.Write([]byte(payload))
		resp, _ := http.Post(baseUrl, CONTENT_TYPE_JSON, &buf)
		body := map[string]any{}
		json.NewDecoder(resp.Body).Decode(&body)
		return body
	}

	//Test request respBody has all correct fields {
	respBody := callApi(`{"id":"5","name":"jojo"}`)
	req := respBody["req"].(map[string]any)
	queryParams := respBody["queryParams"].(map[string]any)
	pathParams := respBody["pathParams"].(map[string]any)
	assert.Equal("5", req["id"])
	assert.Equal("jojo", req["name"])
	assert.Equal("5", pathParams["id"])
	assert.Equal("jojo", queryParams["name"])
	// }

	logger.SetLogLevel(logger.LOG_LEVEL_NONE)

	buf := bytes.Buffer{}
	buf.Write([]byte(`{"id":"5}`))
	resp, _ := http.Post(baseUrl, CONTENT_TYPE_JSON, &buf)
	assert.Equal(http.StatusBadRequest, resp.StatusCode)

	setup()
}

type JsonHttpResponse struct {
	Status       int    `json:"status"`
	Body         any    `json:"body"`
	ErrorMessage string `json:"errorMessage"`
}

// API Response

func TestAPIResponse(t *testing.T) {
	setup()
	assert := assert.New(t)

	var newResp = func() *JsonHttpResponse {
		return &JsonHttpResponse{}
	}

	var responseTester = func(req *JsonHttpResponse) HttpResponse {
		result := HttpResponse{}
		result.Status = req.Status
		result.Body = req.Body
		if len(req.ErrorMessage) > 0 {
			result.Error = errors.New(req.ErrorMessage)
		}
		return result
	}

	mux := http.NewServeMux()
	RegisterRequestBodyHandler(mux, "/responseTest", newResp, responseTester)
	s := httptest.NewServer(mux)
	defer s.Close()
	baseUrl := s.URL + "/responseTest"

	callApi := func(payload string) (map[string]any, int) {
		buf := bytes.Buffer{}
		buf.Write([]byte(payload))
		apiResp, _ := http.Post(baseUrl, CONTENT_TYPE_JSON, &buf)
		body := map[string]any{}
		json.NewDecoder(apiResp.Body).Decode(&body)
		return body, apiResp.StatusCode
	}

	// Test Response has no error, no status, no body {
	resp, status := callApi(`{}`)
	logger.Debug(resp)
	assert.Equal(0, len(resp))
	assert.Equal(http.StatusOK, status)
	// }

	// Response has no error, no status, only body {
	resp, status = callApi(`{"body":{"id":"5"}}`)
	assert.Equal("5", resp["id"])
	assert.Equal(http.StatusOK, status)
	// }

	// Response has no error, no status, only body {
	resp, status = callApi(`{"status":201,"body":{"id":"5"}}`)
	assert.Equal("5", resp["id"])
	assert.Equal(http.StatusCreated, status)
	// }

	// Response has error, no status & no body {
	resp, status = callApi(`{"errorMessage":"jojo"}`)
	assert.Equal("jojo", resp["errorMessage"])
	assert.Equal(http.StatusInternalServerError, status)
	// }

	// Response has error, status & no body {
	resp, status = callApi(`{"status":400,"errorMessage":"jojo"}`)
	assert.Equal("jojo", resp["errorMessage"])
	assert.Equal(http.StatusBadRequest, status)
	// }

	// Response has error, status & body {
	resp, status = callApi(`{"status":400,"body":{"id":"5"},"errorMessage":"jojo"}`)
	logger.Debug(resp)
	assert.Equal("jojo", resp["errorMessage"])
	assert.Equal("5", resp["body"].(map[string]any)["id"])
	assert.Equal(http.StatusBadRequest, status)
	// }

	// 6.
}

func TestSingleFileServer(t *testing.T) {
	setup()
	assert := assert.New(t)
	mux := http.NewServeMux()
	RegisterSingleFileServer(mux, "/", "../resources/test/properties/local.properties")
	s := httptest.NewServer(mux)
	defer s.Close()

	r, err := http.Get(s.URL)
	assert.Nil(err)
	assert.Equal(http.StatusOK, r.StatusCode)
}

func getResponseMap(resp *http.Response, err error) map[string]string {
	respBody := map[string]string{}
	if err == nil {
		if resp.StatusCode == http.StatusOK {
			json.NewDecoder(resp.Body).Decode(&respBody)
		} else {
			logger.Error(resp)
		}

	} else {
		logger.Error(err)
	}
	return respBody
}
