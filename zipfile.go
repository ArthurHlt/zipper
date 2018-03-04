package zipper

import (
	"io"
)

type ZipReadCloser interface {
	io.ReadCloser
	Size() int64
}

type ZipFile struct {
	file  io.ReadCloser
	size  int64
	clean func() error
}

func NewZipFile(file io.ReadCloser, size int64, clean func() error) *ZipFile {
	return &ZipFile{file, size, clean}
}
func (f ZipFile) Read(p []byte) (int, error) {
	return f.file.Read(p)
}

func (f ZipFile) Close() error {
	err := f.file.Close()
	if err != nil {
		return err
	}
	return f.clean()
}

func (f ZipFile) Size() int64 {
	return f.size
}
