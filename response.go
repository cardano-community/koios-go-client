// SPDX-License-Identifier: Apache-2.0
//
// Copyright © 2022 The Cardano Community Authors

package koios

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type (
	// Response wraps API responses.
	Response struct {
		// RequestURL is full request url.
		RequestURL string `json:"request_url"`

		// RequestMethod is HTTP method used for request.
		RequestMethod string `json:"request_method"`

		// StatusCode of the HTTP response.
		StatusCode int `json:"status_code"`

		// Status of the HTTP response header if present.
		Status string `json:"status"`

		// Date response header.
		Date string `json:"date,omitempty"`

		// ContentLocation response header if present.
		ContentLocation string `json:"content_location,omitempty"`

		// ContentRange response header if present.
		ContentRange string `json:"content_range,omitempty"`

		// Error response body if present.
		Error *ResponseError `json:"error,omitempty"`

		// Stats of the request if stats are enabled.
		Stats *RequestStats `json:"stats,omitempty"`
	}

	ErrorCode string

	// ResponseError represents api error messages.
	ResponseError struct {
		error

		// Hint of the error reported by server.
		Hint string `json:"hint,omitempty"`

		// Details of the error reported by server.
		Details string `json:"details,omitempty"`

		// Code is error code reported by server.
		Code ErrorCode `json:"code,omitempty"`

		// Message is error message reported by server.
		Message string `json:"message,omitempty"`
	}
)

func ErrorCodeFromInt(code int) ErrorCode {
	return ErrorCode(strconv.Itoa(code))
}

// String returns error code as string.
func (c ErrorCode) String() string {
	return string(c)
}

// Int returns error code as integer if
// strconv.Atoi is able to parse it, otherwise it returns 0.
func (c ErrorCode) Int() int {
	i, _ := strconv.Atoi(string(c))
	return i
}

// Error return underlying error string.
func (e *ResponseError) Error() string {
	return e.Message
}

// Error return underlying error string.
func (e *ResponseError) Unwrap() error {
	return e.error
}

// ReadResponseBody is reading http.Response aand closing it after read.
func ReadResponseBody(rsp *http.Response) (body []byte, err error) {
	if rsp == nil {
		return nil, nil
	}

	defer func() { _ = rsp.Body.Close() }()

	rb := rsp.Body

	if strings.Contains(rsp.Header.Get("Content-Encoding"), "gzip") {
		if rb, err = gzip.NewReader(rsp.Body); err == nil {
			defer rb.Close()
		} else {
			return nil, err
		}
	} else if rsp.Header.Get("Content-Encoding") == "deflate" {
		rb = flate.NewReader(rsp.Body)
		defer rb.Close()
	}

	return io.ReadAll(rb)
}

// ReadAndUnmarshalResponse is helper to unmarchal json responses.
func ReadAndUnmarshalResponse(rsp *http.Response, res *Response, dest any) error {
	if rsp == nil {
		return fmt.Errorf("%w: got no response", ErrResponse)
	}
	body, err := ReadResponseBody(rsp)
	if !strings.Contains(rsp.Header.Get("Content-Type"), "json") {
		return fmt.Errorf("%w: %s", ErrResponseIsNotJSON, string(body))
	}

	res.applyError(body, err)
	if len(body) == 0 || err != nil {
		return err
	}

	defer res.ready()
	err = json.Unmarshal(body, dest)
	res.applyError(body, err)
	return err
}

func (r *Response) applyError(body []byte, err error) {
	if err == nil {
		return
	}

	if r.Error == nil {
		r.Error = &ResponseError{
			error: err,
		}
	}
	if len(body) != 0 {
		berr := json.Unmarshal(body, r.Error)
		if berr != nil {
			r.Error.Message = berr.Error()
		}
	}
	defer r.ready()

	if len(r.Error.Message) == 0 {
		r.Error.Message = err.Error()
	} else {
		r.Error.Message = fmt.Sprintf("%s: %s", err.Error(), r.Error.Message)
	}
}

func (r *Response) ready() {
	if r.Stats == nil {
		return
	}
	r.Stats.ReqDur = time.Since(r.Stats.ReqStartedAt)
	r.Stats.ReqDurStr = fmt.Sprint(r.Stats.ReqDur)
	if r.Error != nil && len(r.Error.Code) == 0 {
		r.Error.Code = ErrorCodeFromInt(r.StatusCode)
	}
}

func (r *Response) applyRsp(rsp *http.Response) {
	r.StatusCode = rsp.StatusCode
	r.RequestMethod = rsp.Request.Method
	r.Status = rsp.Status
	r.Date = rsp.Header.Get("date")
	r.ContentRange = rsp.Header.Get("content-range")
	r.ContentLocation = rsp.Header.Get("content-location")
}
