package zipper_test

import (
	. "github.com/ArthurHlt/zipper"

	"github.com/ArthurHlt/zipper/zipperfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Session", func() {
	var session *Session
	BeforeEach(func() {
		h1 := &zipperfakes.FakeHandler{}
		h1.Sha1Stub = func(src *Source) (string, error) {
			return src.Path, nil
		}
		session = NewSession(NewSource("apath"), h1)
	})
	Describe("IsDiff", func() {
		It("should find no diff when sha1 match", func() {
			diff, _, err := session.IsDiff("apath")
			Expect(err).ToNot(HaveOccurred())

			Expect(diff).Should(BeFalse())
		})
		It("should find diff when sha1 doesn't match", func() {
			diff, _, err := session.IsDiff("other")
			Expect(err).ToNot(HaveOccurred())

			Expect(diff).Should(BeTrue())
		})
	})
})
