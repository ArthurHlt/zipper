package zipper

type Session struct {
	handler Handler
	src     *Source
}

func NewSession(src *Source, handler Handler) *Session {
	return &Session{handler, src}
}

func (s Session) Zip() (*ZipFile, error) {
	return s.handler.Zip(s.src)
}

func (s Session) Sha1() (string, error) {
	return s.handler.Sha1(s.src)
}

func (s Session) IsDiff(storedSha1 string) (bool, string, error) {
	sha1Given, err := s.handler.Sha1(s.src)
	if err != nil {
		return true, "", err
	}
	return storedSha1 != sha1Given, sha1Given, nil
}

func (s Session) Handler() Handler {
	return s.handler
}

func (s Session) Source() *Source {
	return s.src
}
