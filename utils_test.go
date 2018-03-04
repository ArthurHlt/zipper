package zipper_test

import (
	. "github.com/ArthurHlt/zipper"

	"bytes"
	"crypto/sha1"
	"encoding/hex"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Utils", func() {
	Describe("IsWebURL", func() {
		Context("When the web url is valid", func() {
			It("should return true if it's an http url", func() {
				Expect(IsWebURL("http://test.com")).Should(BeTrue())
			})
			It("should return true if it's an https url", func() {
				Expect(IsWebURL("https://test.com")).Should(BeTrue())
			})
		})
		Context("When the web url is invalid", func() {
			It("should return false if it's not an http or https url", func() {
				Expect(IsWebURL("fprot://test.com")).Should(BeFalse())
			})
		})
	})
	Describe("HasExtFile", func() {
		It("should return true if path has extension", func() {
			Expect(HasExtFile("foo.myext", ".myext", ".aext")).Should(BeTrue())
		})
		It("should return false if path doesn't have extension", func() {
			Expect(HasExtFile("foo.myext", ".fext", ".aext")).Should(BeFalse())
		})
	})
	Describe("IsTarGzFile", func() {
		It("should return true if has .tar.gz extension", func() {
			Expect(IsTarGzFile("foo.tar.gz")).Should(BeTrue())
		})
		It("should return true if has .tgz extension", func() {
			Expect(IsTarGzFile("foo.tar.gz")).Should(BeTrue())
		})
	})

	Describe("GetSha1FromReader", func() {
		It("should give sha1 for file by taking a 5kb max chunk", func() {
			b := make([]byte, 10*1024)
			buf := bytes.NewBuffer(b)
			h := sha1.New()
			h.Write(b)
			plainSha1 := hex.EncodeToString(h.Sum(nil))

			sha1, err := GetSha1FromReader(buf)
			Expect(err).ToNot(HaveOccurred())

			Expect(sha1).ShouldNot(Equal(plainSha1))
			Expect(sha1).Should(Equal("ec8d8db07ace21ae014c4d7dbe42297dfe61976a"))
		})

	})
})
