package zipper_test

import (
	. "github.com/ArthurHlt/zipper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"testing"
)

var filesInZip = []string{
	"foo.txt",
	"fooDir/",
	"fooDir/bar/",
	"largeblankfile/",
	"largeblankfile/file.txt",
	"lastDir/",
	"subDir/",
	"subDir/bar.txt",
	"subDir/otherDir/",
	"subDir/otherDir/file.txt",
}

func checkZipFile(zipFile ZipReadCloser, addFiles ...string) []os.FileInfo {
	zipFileLocal, err := ioutil.TempFile("", "zip_test")
	Expect(err).NotTo(HaveOccurred())
	defer func() {
		zipFileLocal.Close()
		os.Remove(zipFileLocal.Name())
	}()
	finalFilesInZip := append(filesInZip, addFiles...)
	_, err = io.Copy(zipFileLocal, zipFile)
	Expect(err).NotTo(HaveOccurred())

	fileStat, err := zipFileLocal.Stat()
	Expect(err).NotTo(HaveOccurred())

	reader, err := zip.NewReader(zipFileLocal, fileStat.Size())
	Expect(err).NotTo(HaveOccurred())

	fis := make([]os.FileInfo, 0)
	for _, file := range reader.File {
		Expect(finalFilesInZip).To(ContainElement(file.Name))
		fis = append(fis, file.FileInfo())
	}
	Expect(finalFilesInZip).To(HaveLen(len(reader.File)))
	return fis
}
func TestZipper(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Zipper Suite")
}
func readFileInZip(index int, reader *zip.Reader) (string, string) {
	buf := &bytes.Buffer{}
	file := reader.File[index]
	fReader, err := file.Open()
	_, err = io.Copy(buf, fReader)

	Expect(err).NotTo(HaveOccurred())

	return file.Name, string(buf.Bytes())
}
func readFile(file *os.File) []byte {
	b, err := ioutil.ReadAll(file)
	Expect(err).NotTo(HaveOccurred())
	return b
}
