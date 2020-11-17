package respond

import (
"encoding/json"
"net/http"
)

var (
	// https://github.com/uber-go/guide/blob/master/style.md#verify-interface-compliance
	_    Encoder = (*jsonEncoder)(nil)
	JSON Encoder = &jsonEncoder{}
)

type Encoder interface {
	Encode(w http.ResponseWriter, v interface{}) error
	ContentType() string
}

type jsonEncoder struct{}

func (*jsonEncoder) Encode(w http.ResponseWriter, v interface{}) error {
	return json.NewEncoder(w).Encode(v)
}

func (*jsonEncoder) ContentType() string {
	return "application/json; charset=utf-8"
}