package respond

import (
	"errors"
	"net/http"
	"reflect"
)

type (
	BeforeFunc  func(w http.ResponseWriter, data interface{}, status int) (int, interface{})
	AfterFunc   func(data interface{}, status int)
	OnErrorFunc func(err error)
)

type Responder struct {
	Encoder Encoder
	// Before is called for Before each response is written
	// and gives user code the chance to mutate the status or data.
	// Useful for handling different types of data differently (like errors),
	// enveloping the response, setting common HTTP headers etc.
	Before BeforeFunc

	// After is called After each response.
	// Useful for logging activity After a response has been written.
	After AfterFunc

	OnError OnErrorFunc
}

// Option interface serves to implementation of functional options pattern
type Option interface {
	apply(opt *Responder) error
}

type optionFunc func(opt *Responder) error

func (f optionFunc) apply(r *Responder) error {
	return f(r)
}

func WithBefore(before BeforeFunc) Option {
	return optionFunc(func(r *Responder) error {
		if before == nil {
			return errors.New("responder: provided BeforeFunc is nil")
		}

		r.Before = before

		return nil
	})
}

func WithAfter(after AfterFunc) Option {
	return optionFunc(func(r *Responder) error {
		if after == nil {
			return errors.New("responder: provided AfterFunc is nil")
		}

		r.After = after

		return nil
	})
}

func WithOnError(onError OnErrorFunc) Option {
	return optionFunc(func(r *Responder) error {
		if onError == nil {
			return errors.New("responder: provided onErr function is nil")
		}

		r.OnError = onError

		return nil
	})
}

func WithEncoder(encoder Encoder) Option {
	return optionFunc(func(r *Responder) error {
		if reflect.ValueOf(encoder).IsNil() {
			return errors.New("responder: provided Encoder is nil")
		}

		r.Encoder = encoder

		return nil
	})
}

func NewResponder(opts ...Option) (*Responder, error) {
	responder := &Responder{}
	for i := range opts {
		if err := opts[i].apply(responder); err != nil {
			return nil, err
		}
	}
	return responder, nil
}

func (r *Responder) writeResponse(w http.ResponseWriter, v interface{}, status int) {

	if r.Before != nil {
		status, v = r.Before(w, v, status)
	}

	encoder := JSON
	if r.Encoder != nil {
		encoder = r.Encoder
	}

	w.Header().Set("Content-Type", encoder.ContentType())
	w.WriteHeader(status)
	if err := encoder.Encode(w, v); err != nil {
		if r.OnError != nil {
			r.OnError(err)
		}
	}

	if r.After != nil {
		r.After(v, status)
	}

}

func (r *Responder) Ok(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusOK)
}

func (r *Responder) Created(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusCreated)
}

func (r *Responder) Accepted(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusAccepted)
}

func (r *Responder) NoContent(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusNoContent)
}

func (r *Responder) BadRequest(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusBadRequest)
}

func (r *Responder) Unauthorized(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusUnauthorized)
}

func (r *Responder) Forbidden(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusForbidden)
}

func (r *Responder) NotFound(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusNotFound)
}

// MethodNotAllowed returns a 405 Method Not Allowed JSON response
func (r *Responder) MethodNotAllowed(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusMethodNotAllowed)
}

// Conflict returns a 409 Conflict JSON response
func (r *Responder) Conflict(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusConflict)
}

// InternalServerError returns a 500 Internal Server Error JSON response
func (r *Responder) InternalServerError(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusInternalServerError)
}

// NotImplemented returns a 501 Not Implemented JSON response
func (r *Responder) NotImplemented(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusNotImplemented)
}

// BadGateway returns a 502 Bad Gateway JSON response
func (r *Responder) BadGateway(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusBadGateway)
}

// ServiceUnavailable returns a 503 Service Unavailable JSON response
func (r *Responder) ServiceUnavailable(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusServiceUnavailable)
}

// GatewayTimeout returns a 504 Gateway Timeout JSON response
func (r *Responder) GatewayTimeout(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusGatewayTimeout)
}

func (r *Responder) UnprocessableEntity(w http.ResponseWriter, v interface{}) {
	r.writeResponse(w, v, http.StatusUnprocessableEntity)
}
