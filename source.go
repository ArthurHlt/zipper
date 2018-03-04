package zipper

import (
	"context"
	"net/http"
)

const (
	InsecureContextKey SourceContextKey = iota
)

type SourceContextKey int

type Source struct {
	Path string
	ctx  context.Context
}

func NewSource(path string) *Source {
	return &Source{Path: path}
}
func (s *Source) Context() context.Context {
	if s.ctx != nil {
		return s.ctx
	}
	return context.Background()
}

func (s *Source) WithContext(ctx context.Context) *Source {
	if ctx == nil {
		panic("nil context")
	}
	s2 := new(Source)
	*s2 = *s
	s2.ctx = ctx
	return s2
}

func SetCtxHttpClient(src *Source, client *http.Client) {
	parentContext := src.Context()
	ctxValueReq := src.WithContext(context.WithValue(parentContext, InsecureContextKey, client))
	*src = *ctxValueReq
}

func CtxHttpClient(src *Source) *http.Client {
	val := src.Context().Value(InsecureContextKey)
	if val == nil {
		return nil
	}
	return val.(*http.Client)
}
