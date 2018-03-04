package zipper_test

import (
	. "github.com/ArthurHlt/zipper"

	"github.com/ArthurHlt/zipper/zipperfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"net/http"
	"time"
)

var _ = Describe("Manager", func() {
	var manager *Manager
	BeforeEach(func() {
		var err error
		h1 := &zipperfakes.FakeHandler{}
		h1.NameStub = func() string {
			return "fake1"
		}
		h1.DetectStub = func(src *Source) bool {
			return src.Path == "fake1"
		}
		h2 := &zipperfakes.FakeHandler{}
		h2.NameStub = func() string {
			return "fake2"
		}
		h2.DetectStub = func(src *Source) bool {
			return src.Path == "fake2"
		}

		manager, err = NewManager(h1, h2)
		if err != nil {
			panic(err)
		}
	})
	Describe("AddHandlers", func() {
		It("should return an error if handler with same name already exists", func() {
			conflict := &zipperfakes.FakeHandler{}
			conflict.NameStub = func() string {
				return "fake1"
			}

			err := manager.AddHandlers(conflict)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).Should(Equal("Handler fake1 already exists"))
		})
	})
	Describe("FindHandler", func() {
		Context("when auto detecting", func() {
			It("should return correct handler by detection", func() {
				h, err := manager.FindHandler("fake2", "")
				Expect(err).ToNot(HaveOccurred())

				Expect(h.Name()).Should(Equal("fake2"))
			})
			It("should return error when not detecting", func() {
				_, err := manager.FindHandler("fake3", "")
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when choosing type", func() {
			It("should return correct handler by its name", func() {
				h, err := manager.FindHandler("fake2", "fake1")
				Expect(err).ToNot(HaveOccurred())

				Expect(h.Name()).Should(Equal("fake1"))
			})
			It("should return error when doesn't exists", func() {
				_, err := manager.FindHandler("fake1", "fake3")
				Expect(err).To(HaveOccurred())
			})
		})
	})
	Describe("SetHttpClient", func() {
		It("should always set timeout to 0 on client", func() {
			c := &http.Client{
				Timeout: 15,
			}
			manager.SetHttpClient(c)

			Expect(c.Timeout).Should(Equal(time.Duration(0)))
		})
	})
	Describe("CreateSession", func() {
		Context("when auto detecting", func() {
			It("should give session with correct handler ", func() {
				s, err := manager.CreateSession("fake2")
				Expect(err).ToNot(HaveOccurred())

				Expect(s.Handler().Name()).Should(Equal("fake2"))
			})
			It("should return error when can't find handler", func() {
				_, err := manager.CreateSession("fake3")
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when choosing type", func() {
			It("should give session with given handler type", func() {
				s, err := manager.CreateSession("fake2", "fake1")
				Expect(err).ToNot(HaveOccurred())

				Expect(s.Handler().Name()).Should(Equal("fake1"))
			})
			It("should return error when can't find handler", func() {
				_, err := manager.CreateSession("fake2", "fake3")
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
