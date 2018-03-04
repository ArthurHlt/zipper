package zipper

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

var fManager *Manager = mustNewManager(
	NewGitHandler(),
	&HttpHandler{},
	&LocalHandler{},
)

type Manager struct {
	handlers   map[string]Handler
	httpClient *http.Client
}

func mustNewManager(handlers ...Handler) *Manager {
	m := &Manager{handlers: make(map[string]Handler)}
	err := m.AddHandlers(handlers...)
	if err != nil {
		panic(err)
	}
	return m
}
func NewManager(handlers ...Handler) (*Manager, error) {
	m := &Manager{
		handlers: make(map[string]Handler),
		httpClient: &http.Client{
			Timeout: 0,
		},
	}
	err := m.AddHandlers(handlers...)
	return m, err
}
func (m *Manager) SetHttpClient(httpClient *http.Client) {
	m.httpClient = httpClient
	m.httpClient.Timeout = time.Duration(0)
}
func SetHttpClient(httpClient *http.Client) {
	fManager.SetHttpClient(httpClient)
}
func CreateSession(path string, handlerNames ...string) (*Session, error) {
	return fManager.CreateSession(path, handlerNames...)
}
func (m *Manager) CreateSession(path string, handlerNames ...string) (*Session, error) {
	handlerName := ""
	if len(handlerNames) > 0 {
		handlerName = handlerNames[0]
	}
	h, err := m.FindHandler(path, handlerName)
	if err != nil {
		return nil, err
	}
	src := NewSource(path)
	SetCtxHttpClient(src, m.httpClient)
	return NewSession(src, h), nil
}
func AddHandlers(handlers ...Handler) error {
	return fManager.AddHandlers(handlers...)
}
func (m *Manager) AddHandlers(handlers ...Handler) error {
	for _, handler := range handlers {
		err := m.AddHandler(handler)
		if err != nil {
			return err
		}
	}
	return nil
}
func AddHandler(handler Handler) error {
	return fManager.AddHandler(handler)
}
func (m *Manager) AddHandler(handler Handler) error {
	name := strings.ToLower(handler.Name())
	if _, ok := m.handlers[name]; ok {
		return fmt.Errorf("Handler %s already exists", name)
	}
	m.handlers[name] = handler
	return nil
}
func FindHandler(path string, handlerName string) (Handler, error) {
	return fManager.FindHandler(path, handlerName)
}
func (m *Manager) FindHandler(path string, handlerName string) (Handler, error) {
	src := NewSource(path)
	handlerName = strings.ToLower(handlerName)
	if handlerName == "" {
		for _, h := range m.handlers {
			if h.Detect(src) {
				return h, nil
			}
		}
	}
	if h, ok := m.handlers[handlerName]; ok {
		return h, nil
	}
	return nil, fmt.Errorf("Handler for path '%s' cannot be found.", src.Path)
}
