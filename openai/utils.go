package openai

import (
	"bytes"
	"encoding/json"
	"io"
)

const (
	contentTypeJSON = `application/json; charset="utf-8"`
)

// BodyMarshaler is the interface implemented by types that can marshal themselves for http request.
type BodyMarshaler interface {
	MarshalBody() (io.Reader, string, error)
}

// buildRequest build a request, returns buffer, contentType, error
func buildRequest(a any) (io.Reader, string, error) {
	if a == nil {
		return nil, "", nil
	}

	if bm, ok := a.(BodyMarshaler); ok {
		return bm.MarshalBody()
	}

	return buildJsonRequest(a)
}

func buildJsonRequest(a any) (io.Reader, string, error) {
	body, err := json.Marshal(a)
	if err != nil {
		return nil, "", err
	}

	buf := bytes.NewReader(body)
	return buf, contentTypeJSON, nil
}
