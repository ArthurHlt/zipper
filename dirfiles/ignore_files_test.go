package dirfiles_test

import (
	. "github.com/ArthurHlt/zipper/dirfiles"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CF Ignore", func() {
	It("excludes files based on exact path matches", func() {
		ignore := NewIgnoreFiles(`the-dir/the-path`)
		Expect(ignore.FileShouldBeIgnored("the-dir/the-path")).To(BeTrue())
	})

	It("excludes the contents of directories based on exact path matches", func() {
		ignore := NewIgnoreFiles(`dir1/dir2`)
		Expect(ignore.FileShouldBeIgnored("dir1/dir2/the-file")).To(BeTrue())
		Expect(ignore.FileShouldBeIgnored("dir1/dir2/dir3/the-file")).To(BeTrue())
	})

	It("excludes the directories based on relative path matches", func() {
		ignore := NewIgnoreFiles(`dir1`)
		Expect(ignore.FileShouldBeIgnored("dir1")).To(BeTrue())
		Expect(ignore.FileShouldBeIgnored("dir2/dir1")).To(BeTrue())
		Expect(ignore.FileShouldBeIgnored("dir3/dir2/dir1")).To(BeTrue())
		Expect(ignore.FileShouldBeIgnored("dir3/dir1/dir2")).To(BeTrue())
	})

	It("excludes files based on star patterns", func() {
		ignore := NewIgnoreFiles(`dir1/*.so`)
		Expect(ignore.FileShouldBeIgnored("dir1/file1.so")).To(BeTrue())
		Expect(ignore.FileShouldBeIgnored("dir1/file2.cc")).To(BeFalse())
	})

	It("excludes files based on double-star patterns", func() {
		ignore := NewIgnoreFiles(`dir1/**/*.so`)
		Expect(ignore.FileShouldBeIgnored("dir1/dir2/dir3/file1.so")).To(BeTrue())
		Expect(ignore.FileShouldBeIgnored("different-dir/dir2/file.so")).To(BeFalse())
	})

	It("allows files to be explicitly included", func() {
		ignore := NewIgnoreFiles(`
node_modules/*
!node_modules/common
`)

		Expect(ignore.FileShouldBeIgnored("node_modules/something-else")).To(BeTrue())
		Expect(ignore.FileShouldBeIgnored("node_modules/common")).To(BeFalse())
	})

	It("applies the patterns in order from top to bottom", func() {
		ignore := NewIgnoreFiles(`
stuff/*
!stuff/*.c
stuff/exclude.c`)

		Expect(ignore.FileShouldBeIgnored("stuff/something.txt")).To(BeTrue())
		Expect(ignore.FileShouldBeIgnored("stuff/exclude.c")).To(BeTrue())
		Expect(ignore.FileShouldBeIgnored("stuff/include.c")).To(BeFalse())
	})

	It("ignores certain commonly ingored files by default", func() {
		ignore := NewIgnoreFiles(``)
		Expect(ignore.FileShouldBeIgnored(".git/objects")).To(BeTrue())

		ignore = NewIgnoreFiles(`!.git`)
		Expect(ignore.FileShouldBeIgnored(".git/objects")).To(BeFalse())
	})

})
