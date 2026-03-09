package openai

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type DetailError struct {
	Type    string `json:"type,omitempty"`
	Code    any    `json:"code,omitempty"`
	Param   any    `json:"param,omitempty"`
	Message string `json:"message,omitempty"`
}

func (de *DetailError) Error() string {
	var sb strings.Builder
	if de.Type != "" {
		sb.WriteString(de.Type)
	}
	if de.Code != nil {
		s := fmt.Sprint(de.Code)
		if s != "" {
			if sb.Len() > 0 {
				sb.WriteByte('/')
			}
			sb.WriteString(s)
		}
	}
	if de.Param != nil {
		s := fmt.Sprint(de.Param)
		if s != "" {
			if sb.Len() > 0 {
				sb.WriteByte('/')
			}
			sb.WriteString(s)
		}
	}
	if de.Message != "" {
		if sb.Len() > 0 {
			sb.WriteString(": ")
		}
		sb.WriteString(de.Message)
	}
	return sb.String()
}

type ResultError struct {
	Method     string       `json:"-"` // http request method
	URL        *url.URL     `json:"-"` // http request URL
	StatusCode int          `json:"-"` // http status code
	Status     string       `json:"-"` // http status
	Detail     *DetailError `json:"error,omitempty"`
}

func AsResultError(err error) (re *ResultError, ok bool) {
	ok = errors.As(err, &re)
	return
}

func IsResultError(err error) bool {
	_, ok := AsResultError(err)
	return ok
}

func newResultError(res *http.Response) *ResultError {
	return &ResultError{
		Method:     res.Request.Method,
		URL:        res.Request.URL,
		StatusCode: res.StatusCode,
		Status:     res.Status,
	}
}

func (re *ResultError) Error() string {
	es := re.Status + " (" + re.Method + " " + re.URL.String() + ")"

	if re.Detail != nil {
		es += " - " + re.Detail.Error()
	}

	return es
}
