package zipper_test

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"code.cloudfoundry.org/gofileutils/fileutils"
	. "github.com/ArthurHlt/zipper"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// Thanks to Svett Ralchev
// http://blog.ralch.com/tutorial/golang-working-with-zip/
func zipit(source, target, prefix string) error {
	zipfile, err := os.Create(target)
	if err != nil {
		return err
	}
	defer zipfile.Close()

	if prefix != "" {
		_, err = io.WriteString(zipfile, prefix)
		if err != nil {
			return err
		}
	}

	archive := zip.NewWriter(zipfile)
	defer archive.Close()

	err = filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(path, source)

		if info.IsDir() {
			header.Name += string(os.PathSeparator)
		} else {
			header.Method = zip.Deflate
		}

		writer, err := archive.CreateHeader(header)
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})

	return err
}

var _ = Describe("LocalHandler", func() {
	var handler LocalHandler
	var zipFileLocal *os.File

	BeforeEach(func() {
		var err error
		zipFileLocal, err = ioutil.TempFile("", "zip_test")
		Expect(err).NotTo(HaveOccurred())
		handler = LocalHandler{}
	})
	AfterEach(func() {
		zipFileLocal.Close()
		os.Remove(zipFileLocal.Name())
	})
	Describe("Zip", func() {
		It("creates a zip with all files and directories from the source directory", func() {
			workingDir, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())

			dir := filepath.Join(workingDir, "fixtures/zip/")
			zipFile, err := handler.Zip(NewSource(dir))
			Expect(err).NotTo(HaveOccurred())
			defer zipFile.Close()

			checkZipFile(zipFile)
		})
	})
	Describe("Detect", func() {
		It("should return true if path exists on system", func() {
			workingDir, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			dir := filepath.Join(workingDir, "fixtures/zip/")
			exists := handler.Detect(NewSource(dir))
			Expect(exists).Should(BeTrue())
		})
		It("should return true if path doesn't exists on system", func() {
			workingDir, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			dir := filepath.Join(workingDir, "not-exists")
			exists := handler.Detect(NewSource(dir))
			Expect(exists).Should(BeFalse())
		})
	})
	Describe("Sha1", func() {
		It("should create sha1 from the source directory", func() {
			workingDir, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			dir := filepath.Join(workingDir, "fixtures/zip/")
			sha1, err := handler.Sha1(NewSource(dir))
			Expect(err).NotTo(HaveOccurred())
			Expect(sha1).ShouldNot(BeEmpty())
		})
	})
	Describe("ZipFiles", func() {

		It("creates a zip with all files and directories from the source directory", func() {
			workingDir, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())

			dir := filepath.Join(workingDir, "fixtures/zip/")
			err = handler.ZipFiles(dir, zipFileLocal)
			Expect(err).NotTo(HaveOccurred())

			fileStat, err := zipFileLocal.Stat()
			Expect(err).NotTo(HaveOccurred())

			reader, err := zip.NewReader(zipFileLocal, fileStat.Size())
			Expect(err).NotTo(HaveOccurred())

			filenames := []string{}
			for _, file := range reader.File {
				filenames = append(filenames, file.Name)
			}
			Expect(filenames).To(Equal(filesInZip))

			name, contents := readFileInZip(0, reader)
			Expect(name).To(Equal("foo.txt"))
			Expect(contents).To(Equal("This is a simple text file."))
		})

		It("creates a zip with the original file modes", func() {
			if runtime.GOOS == "windows" {
				Skip("This test does not run on Windows")
			}

			workingDir, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())

			dir := filepath.Join(workingDir, "fixtures/zip/")
			err = os.Chmod(filepath.Join(dir, "subDir/bar.txt"), 0666)
			Expect(err).NotTo(HaveOccurred())

			err = handler.ZipFiles(dir, zipFileLocal)
			Expect(err).NotTo(HaveOccurred())

			fileStat, err := zipFileLocal.Stat()
			Expect(err).NotTo(HaveOccurred())

			reader, err := zip.NewReader(zipFileLocal, fileStat.Size())
			Expect(err).NotTo(HaveOccurred())

			readFileInZip(7, reader)
			Expect(reader.File[7].FileInfo().Mode()).To(Equal(os.FileMode(0666)))
		})

		It("creates a zip with executable file modes", func() {
			if runtime.GOOS != "windows" {
				Skip("This test only runs on Windows")
			}

			workingDir, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())

			dir := filepath.Join(workingDir, "fixtures/zip/")
			err = os.Chmod(filepath.Join(dir, "subDir/bar.txt"), 0666)
			Expect(err).NotTo(HaveOccurred())

			err = handler.ZipFiles(dir, zipFileLocal)
			Expect(err).NotTo(HaveOccurred())

			fileStat, err := zipFileLocal.Stat()
			Expect(err).NotTo(HaveOccurred())

			reader, err := zip.NewReader(zipFileLocal, fileStat.Size())
			Expect(err).NotTo(HaveOccurred())

			readFileInZip(7, reader)
			Expect(fmt.Sprintf("%o", reader.File[7].FileInfo().Mode())).To(Equal("766"))
		})

		It("compresses the files", func() {
			workingDir, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())

			dir := filepath.Join(workingDir, "fixtures/zip/largeblankfile/")
			fileStat, err := os.Stat(filepath.Join(dir, "file.txt"))
			Expect(err).NotTo(HaveOccurred())
			originalFileSize := fileStat.Size()

			err = handler.ZipFiles(dir, zipFileLocal)
			Expect(err).NotTo(HaveOccurred())

			fileStat, err = zipFileLocal.Stat()
			Expect(err).NotTo(HaveOccurred())

			compressedFileSize := fileStat.Size()
			Expect(compressedFileSize).To(BeNumerically("<", originalFileSize))
		})

		It("returns an error when zipping fails", func() {
			handler := LocalHandler{}
			err := handler.ZipFiles("/a/bogus/directory", zipFileLocal)
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("open /a/bogus/directory"))
		})

		It("returns an error when the directory is empty", func() {
			fileutils.TempDir("zip_test", func(emptyDir string, err error) {
				handler := LocalHandler{}
				err = handler.ZipFiles(emptyDir, zipFileLocal)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is empty"))
			})
		})
	})

	Describe("GetZipSize", func() {
		var handler = LocalHandler{}

		It("returns the size of the zip file", func() {
			dir, err := os.Getwd()
			Expect(err).NotTo(HaveOccurred())
			zipFile := filepath.Join(dir, "fixtures/applications/example-app.zip")

			file, err := os.Open(zipFile)
			Expect(err).NotTo(HaveOccurred())

			fileSize, err := handler.GetZipSize(file)
			Expect(err).NotTo(HaveOccurred())
			Expect(fileSize).To(Equal(int64(1803)))
		})

		It("returns  an error if the zip file cannot be found", func() {
			tmpFile, _ := os.Open("fooBar")
			_, sizeErr := handler.GetZipSize(tmpFile)
			Expect(sizeErr).To(HaveOccurred())
		})
	})
})
