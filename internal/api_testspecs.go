package internal

import (
	"io"
	"net/http"
	"net/url"
)

type APITestRequestSpec struct {
	Method   string     `json:"method"`
	Endpoint string     `json:"path"`
	Query    url.Values `json:"query,omitempty"`
	Body     io.Reader  `json:"body,omitempty"`
}

type APITestResponseSpec struct {
	Code   int         `json:"code"`
	Header http.Header `json:"headers"`
	Body   interface{} `json:"body"`
}

type APITestSpec struct {
	Filename string              `json:"filename"`
	Request  APITestRequestSpec  `json:"request"`
	Response APITestResponseSpec `json:"response"`
}
