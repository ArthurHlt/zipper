package zipper_test

import (
	. "github.com/ArthurHlt/zipper"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"os"
)

const (
	fixtureRepo    = "https://github.com/ArthurHlt/zipper-fixture.git"
	fixtureRepoSsh = "ssh://git@github.com:ArthurHlt/zipper-fixture.git"
)

var _ = Describe("Git", func() {
	var handler *GitHandler
	BeforeEach(func() {
		handler = NewGitHandler()
	})
	Describe("Detect", func() {
		It("should return true when an http(s) link and extension is .git", func() {
			Expect(handler.Detect(NewSource("http://foo.com/app.git"))).Should(BeTrue())
		})
		It("should return true when an ssh link and extension is .git", func() {
			Expect(handler.Detect(NewSource("http://git@foo.com:ArthurHlt/app.git"))).Should(BeTrue())
		})
		It("should return false if not an http link or not have one of valid extension", func() {
			Expect(handler.Detect(NewSource("/app.git"))).Should(BeFalse(), "link")
			Expect(handler.Detect(NewSource("http://foo.com/app.ext"))).Should(BeFalse(), "extension")
		})
	})
	Describe("Sha1", func() {
		It("should create sha1 from a source url", func() {
			src := NewSource(fixtureRepo)
			SetCtxHttpClient(src, http.DefaultClient)
			sha1, err := handler.Sha1(src)

			Expect(err).NotTo(HaveOccurred())
			Expect(sha1).Should(Equal("eb3bb57ba0e7da0069ad673b3c3a988484d0291b"))
		})
		It("should create sha1 from a source url from branch", func() {
			src := NewSource(fixtureRepo + "#test-branch")
			SetCtxHttpClient(src, http.DefaultClient)
			sha1, err := handler.Sha1(src)

			Expect(err).NotTo(HaveOccurred())
			Expect(sha1).Should(Equal("c0bfbc199a7ea040712461072c52567fe5361238"))
		})
		It("should create sha1 from a source url from tag", func() {
			src := NewSource(fixtureRepo + "#v0.0.1")
			SetCtxHttpClient(src, http.DefaultClient)
			sha1, err := handler.Sha1(src)

			Expect(err).NotTo(HaveOccurred())
			Expect(sha1).Should(Equal("c0bfbc199a7ea040712461072c52567fe5361238"))
		})
		It("should create sha1 from passed commit when passed in url fragment", func() {
			src := NewSource(fixtureRepo + "#eb3bb57ba0e7da0069ad673b3c3a988484d0291c")
			SetCtxHttpClient(src, http.DefaultClient)
			sha1, err := handler.Sha1(src)

			Expect(err).NotTo(HaveOccurred())
			Expect(sha1).Should(Equal("eb3bb57ba0e7da0069ad673b3c3a988484d0291c"))
		})
	})
	Describe("Zip", func() {
		Context("When is http source url", func() {
			It("should create zip file", func() {
				src := NewSource(fixtureRepo)
				SetCtxHttpClient(src, http.DefaultClient)
				zipFile, err := handler.Zip(src)
				Expect(err).NotTo(HaveOccurred())
				defer zipFile.Close()

				checkZipFile(zipFile, "README.md")
			})
			It("should create zip file from branch", func() {
				src := NewSource(fixtureRepo + "#test-branch")
				SetCtxHttpClient(src, http.DefaultClient)
				zipFile, err := handler.Zip(src)
				Expect(err).NotTo(HaveOccurred())
				defer zipFile.Close()

				checkZipFile(zipFile, "README.md", "branch.txt")
			})
			It("should create zip file from tag", func() {
				src := NewSource(fixtureRepo + "#v0.0.1")
				SetCtxHttpClient(src, http.DefaultClient)
				zipFile, err := handler.Zip(src)
				Expect(err).NotTo(HaveOccurred())
				defer zipFile.Close()

				checkZipFile(zipFile, "README.md", "branch.txt")
			})
			It("should create zip file from commit", func() {
				src := NewSource(fixtureRepo + "#c0bfbc199a7ea040712461072c52567fe5361238")
				SetCtxHttpClient(src, http.DefaultClient)
				zipFile, err := handler.Zip(src)
				Expect(err).NotTo(HaveOccurred())
				defer zipFile.Close()

				checkZipFile(zipFile, "README.md", "branch.txt")
			})
		})
		Context("When is ssh source url", func() {
			BeforeEach(func() {
				if os.Getenv("TRAVIS") != "" {
					Skip("This can run only locally with github ssh key")
				}
			})
			It("should create zip file", func() {
				src := NewSource(fixtureRepoSsh)
				SetCtxHttpClient(src, http.DefaultClient)
				zipFile, err := handler.Zip(src)
				Expect(err).NotTo(HaveOccurred())
				defer zipFile.Close()

				checkZipFile(zipFile, "README.md")
			})
			It("should create zip file from branch", func() {
				src := NewSource(fixtureRepoSsh + "#test-branch")
				SetCtxHttpClient(src, http.DefaultClient)
				zipFile, err := handler.Zip(src)
				Expect(err).NotTo(HaveOccurred())
				defer zipFile.Close()

				checkZipFile(zipFile, "README.md", "branch.txt")
			})
			It("should create zip file from tag", func() {
				src := NewSource(fixtureRepoSsh + "#v0.0.1")
				SetCtxHttpClient(src, http.DefaultClient)
				zipFile, err := handler.Zip(src)
				Expect(err).NotTo(HaveOccurred())
				defer zipFile.Close()

				checkZipFile(zipFile, "README.md", "branch.txt")
			})
			It("should create zip file from commit", func() {
				src := NewSource(fixtureRepoSsh + "#c0bfbc199a7ea040712461072c52567fe5361238")
				SetCtxHttpClient(src, http.DefaultClient)
				zipFile, err := handler.Zip(src)
				Expect(err).NotTo(HaveOccurred())
				defer zipFile.Close()

				checkZipFile(zipFile, "README.md", "branch.txt")
			})
		})
	})
})
