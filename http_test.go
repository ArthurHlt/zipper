package zipper_test

import (
	. "github.com/ArthurHlt/zipper"

	"encoding/base64"
	"fmt"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"time"
)

const (
	serverUrl = "http://server.local"
)

type ServeFileTestHandler struct {
	files map[string]string
	check func(req *http.Request)
}

func createUrl(s *httptest.Server, path string) string {
	return fmt.Sprintf("http://%s%s", s.Listener.Addr().String(), path)
}
func (h ServeFileTestHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if h.files == nil || len(h.files) == 0 {
		w.WriteHeader(404)
		w.Write([]byte(http.StatusText(404)))
		return
	}
	path := req.URL.Path
	if _, ok := h.files[path]; !ok {
		w.WriteHeader(404)
		w.Write([]byte(http.StatusText(404)))
		return
	}
	if h.check != nil {
		h.check(req)
	}
	f, err := os.Open(h.files[path])
	if err != nil {
		panic(err)
	}
	defer f.Close()
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/octet-stream")
	_, err = io.Copy(w, f)
	if err != nil {
		panic(err)
	}
}

var _ = Describe("Http", func() {
	var handler HttpHandler
	var server *httptest.Server
	var httpClient *http.Client
	var servHandler *ServeFileTestHandler

	BeforeEach(func() {
		var err error
		handler = HttpHandler{}
		workingDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		servHandler = &ServeFileTestHandler{files: map[string]string{
			"/final.zip":    filepath.Join(workingDir, "fixtures", "applications", "final.zip"),
			"/final.tar.gz": filepath.Join(workingDir, "fixtures", "applications", "final.tar.gz"),
			"/final.tar":    filepath.Join(workingDir, "fixtures", "applications", "final.tar"),
		}}
		server = httptest.NewServer(servHandler)
		httpClient = server.Client()
		httpClient.Timeout = time.Duration(0)
	})
	AfterEach(func() {
		server.Close()
	})
	Describe("Sha1", func() {
		It("should create sha1 from a source url", func() {
			src := NewSource(createUrl(server, "/final.zip"))
			SetCtxHttpClient(src, httpClient)
			sha1, err := handler.Sha1(src)

			Expect(err).NotTo(HaveOccurred())
			Expect(sha1).Should(Equal("a93ecf13274b289469dee7a0b9e910bc7d2990ce"))
		})
	})
	Describe("Zip", func() {
		Context("with basic auth", func() {
			It("use basic auth in request", func() {
				ran := false
				servHandler.check = func(req *http.Request) {
					defer GinkgoRecover()
					ran = true
					auth := req.Header.Get("Authorization")
					Expect(auth).ToNot(BeEmpty())
					Expect(auth).Should(HavePrefix("Basic "))
					b, err := base64.StdEncoding.DecodeString(auth[6:])
					Expect(err).NotTo(HaveOccurred())
					Expect(string(b)).Should(Equal("user:password"))
				}
				src := NewSource(fmt.Sprintf(
					"http://user:password@%s%s",
					server.Listener.Addr().String(),
					"/final.zip",
				))
				SetCtxHttpClient(src, httpClient)

				zipFile, err := handler.Zip(src)
				Expect(err).NotTo(HaveOccurred())
				defer zipFile.Close()

				Expect(ran).To(BeTrue())
			})
		})
		It("should create zip file from a zip source url", func() {
			src := NewSource(createUrl(server, "/final.zip"))
			SetCtxHttpClient(src, httpClient)
			zipFile, err := handler.Zip(src)
			Expect(err).NotTo(HaveOccurred())
			defer zipFile.Close()

			checkZipFile(zipFile)
		})
		It("should create zip file from a tgz source url", func() {
			src := NewSource(createUrl(server, "/final.tar.gz"))
			SetCtxHttpClient(src, httpClient)
			zipFile, err := handler.Zip(src)
			Expect(err).NotTo(HaveOccurred())
			defer zipFile.Close()

			checkZipFile(zipFile)
		})
		It("should create zip file from a tar source url", func() {
			src := NewSource(createUrl(server, "/final.tar"))
			SetCtxHttpClient(src, httpClient)
			zipFile, err := handler.Zip(src)
			Expect(err).NotTo(HaveOccurred())
			defer zipFile.Close()

			checkZipFile(zipFile)
		})
	})
	Describe("Detect", func() {
		It("should return true when an http(s) link and extension one of on zip, jar, war, tar or tgz file", func() {
			Expect(handler.Detect(NewSource("http://foo.com/app.zip"))).Should(BeTrue(), "zip")
			Expect(handler.Detect(NewSource("http://foo.com/app.jar"))).Should(BeTrue(), "jar")
			Expect(handler.Detect(NewSource("http://foo.com/app.war"))).Should(BeTrue(), "war")
			Expect(handler.Detect(NewSource("http://foo.com/app.tar"))).Should(BeTrue(), "tar")
			Expect(handler.Detect(NewSource("http://foo.com/app.tar.gz"))).Should(BeTrue(), "tar.gz")
			Expect(handler.Detect(NewSource("http://foo.com/app.tgz"))).Should(BeTrue(), "tgz")
		})
		It("should return false if not an http link or not have one of valid extension", func() {
			Expect(handler.Detect(NewSource("/app.zip"))).Should(BeFalse(), "link")
			Expect(handler.Detect(NewSource("http://foo.com/app.ext"))).Should(BeFalse(), "extension")
		})
	})
})
