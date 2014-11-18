package cache_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/DamnWidget/VenGO/cache"
)

var RunSlowTests = false

// check if we are running on travis
// NOTE: this will return false positives in the home directory of anyone
// that is called travis and his home is "travis" or contains "travis", sorry
func runningOnTravis() bool {
	c := cache.CacheDirectory()
	return strings.Contains(c, "travis")
}

var _ = Describe("Cache", func() {

	// disable log output
	log.SetOutput(ioutil.Discard)

	Describe("ExpandUser returns valid path depending on platform", func() {
		var re *regexp.Regexp

		BeforeEach(func() {
			if runtime.GOOS != "darwin" && runtime.GOOS != "windows" {
				re = regexp.MustCompile("/home/([a-z0-9]+)/VenGO")
			}
			if runtime.GOOS == "darwin" {
				re = regexp.MustCompile("/Users/([a-z0-9]+)/VenGO")
			}
			if runtime.GOOS == "windows" {
				re = regexp.MustCompile("\\Users\\([a-z0-9]+)\\VenGO")
			}
		})

		It("Should Match", func() {
			Expect(re.MatchString(cache.ExpandUser("~/VenGO"))).To(BeTrue())
		})
	})

	Describe("CacheDirectory is a valid directory on each platform", func() {
		var re *regexp.Regexp

		BeforeEach(func() {
			if runtime.GOOS != "darwin" && runtime.GOOS != "windows" {
				re = regexp.MustCompile("/home/([a-z0-9]+)/.cache/VenGO")
			}

			if runtime.GOOS == "darwin" {
				re = regexp.MustCompile(
					"/Users/([a-zA-Z0-9]+)/Library/Caches/VenGO")
			}

			if runtime.GOOS == "windows" {
				re = regexp.MustCompile(
					"\\Users\\([a-zA-Z0-9]+)\\AppData\\VenGO")
			}
		})

		It("Should Match", func() {
			Expect(re.MatchString(cache.CacheDirectory())).To(BeTrue())
		})
	})

	Describe("Checksum return an error if version is not supported", func() {
		It("Should be empty string and formatted error", func() {
			sha1, err := cache.Checksum("1.0")

			Expect(sha1).To(BeEmpty())
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Checksum returns the right sha1 string", func() {
		Context("With version 1.2.2", func() {
			It("Should return 3ce0ac4db434fc1546fec074841ff40dc48c1167", func() {
				sha1, err := cache.Checksum("1.2.2")

				Expect(sha1).To(
					Equal("3ce0ac4db434fc1546fec074841ff40dc48c1167"))
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("With version 1.4beta1", func() {
			It("Should return f2fece0c9f9cdc6e8a85ab56b7f1ffcb57c3e7cd", func() {
				sha1, err := cache.Checksum("1.4beta1")

				Expect(sha1).To(
					Equal("f2fece0c9f9cdc6e8a85ab56b7f1ffcb57c3e7cd"))
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("AvilableSources", func() {
		It("Should return 13 sources", func() {
			Expect(len(cache.AvailableSources())).To(Equal(13))
		})
	})

	Describe("AvilableBinaries", func() {
		It("Should return 98 binaries", func() {
			Expect(len(cache.AvailableBinaries())).To(Equal(98))
		})
	})

	Describe("AvilableDonwloads", func() {
		It("Should return 13 sources plus 98 binaries", func() {
			Expect(len(cache.AvailableDownloads())).To(Equal(13 + 98))
		})
	})

	if !runningOnTravis() {
		Describe("Exists works as expected", func() {
			Context("Used in a file that actually exists", func() {
				var file string
				BeforeEach(func() {
					file = filepath.Join(cache.CacheDirectory(), "test")
					Expect(ioutil.WriteFile(file, []byte("Test"), 0644)).To(
						Succeed())
				})

				AfterEach(func() {
					Expect(os.Remove(file)).To(Succeed())
				})

				It("Should return true as the file exists", func() {
					Expect(cache.Exists("test")).To(BeTrue())
				})

				It("Shoudl return false as the file doesn't exists", func() {
					Expect(cache.Exists("invalid")).To(BeFalse())
				})
			})
		})

		Describe("SourceExists", func() {
			var err error
			var dirName string
			BeforeEach(func() {
				dirName, err = ioutil.TempDir(cache.CacheDirectory(), "VenGO-")
				Expect(err).ToNot(HaveOccurred())
			})

			AfterEach(func() {
				err := os.RemoveAll(dirName)
				Expect(err).ToNot(HaveOccurred())
			})

			Context("Passing a non existing directory", func() {
				It("Should return false", func() {
					exists, err := cache.SourceExists("99")
					Expect(err).ToNot(HaveOccurred())
					Expect(exists).To(BeFalse())
				})
			})

			Context("Passing a directory that exists", func() {
				It("Should return true", func() {
					exists, err := cache.SourceExists(path.Base(dirName))
					Expect(err).ToNot(HaveOccurred())
					Expect(exists).To(BeTrue())
				})
			})
		})

		Describe("CacheDonwloadMercurial works as expected", func() {
			Context("Passing a non valid Go version", func() {
				It("Should fail and give back a descriptive error", func() {
					err := cache.CacheDonwloadMercurial("20.0")
					Expect(err).To(HaveOccurred())
					Expect(err).To(Equal(fmt.Errorf(
						"20.0 doesn't seems to be a valid Go release\n"),
					))
				})
			})

			Context("Passing a valid Go version", func() {
				It("Should clone it into the cache directory", func() {
					err := cache.CacheDonwloadMercurial("go1.1")
					Expect(err).ToNot(HaveOccurred())
					_, err = os.Stat(filepath.Join(cache.CacheDirectory(), "go1.1"))
					Expect(err).NotTo(HaveOccurred())
					os.RemoveAll(filepath.Join(cache.CacheDirectory(), "go1.1"))
					debug.FreeOSMemory()
				})
			})

			Context("Passing a valid Go version with no go prefix", func() {
				It("Should clone it into the cache directory", func() {
					err := cache.CacheDonwloadMercurial("1")
					Expect(err).ToNot(HaveOccurred())
					_, err = os.Stat(filepath.Join(cache.CacheDirectory(), "go1"))
					Expect(err).NotTo(HaveOccurred())
					debug.FreeOSMemory()
					os.RemoveAll(filepath.Join(cache.CacheDirectory(), "go1"))
				})
			})
		})

		Describe("AlreadyCompiled", func() {
			It("Shoudl return true if the source has been compiled", func() {
				os.MkdirAll(filepath.Join(cache.CacheDirectory(), "test1", "go", "bin"), 0755)
				filename := filepath.Join(cache.CacheDirectory(), "test1", "go", "bin", "go")
				err := ioutil.WriteFile(filename, []byte{}, 0644)

				Expect(err).ToNot(HaveOccurred())
				Expect(cache.AlreadyCompiled("test1")).To(BeTrue())
				os.RemoveAll(filepath.Join(cache.CacheDirectory(), "test1"))

				Expect(cache.AlreadyCompiled("test1")).To(BeFalse())
			})
		})

		if RunSlowTests {
			Describe("CacheDownload works as expected", func() {
				Context("Passing a non valid Go version", func() {
					It("Should fail and give back a descriptive error", func() {
						err := cache.CacheDownload("1.0")
						Expect(err).To(HaveOccurred())
						Expect(err).To(Equal(fmt.Errorf(
							"1.0 is not a VenGO supported version you must donwload and compile it yourself"),
						))
					})
				})

				Context("Passing a valid Go version", func() {
					It("Should download and extract a valid tar.gz file", func() {
						Expect(cache.CacheDownload("1.2.2")).To(Succeed())

						_, err := os.Stat(filepath.Join(cache.CacheDirectory(), "1.2.2"))
						Expect(err).NotTo(HaveOccurred())
						os.RemoveAll(filepath.Join(cache.CacheDirectory(), "1.2.2"))
						debug.FreeOSMemory()
					})
				})

				Context("Passing an old Go version", func() {
					It("Should donwload and extract a valid tar.gz file", func() {
						Expect(cache.CacheDownload("1.1.1")).To(Succeed())

						_, err := os.Stat(filepath.Join(cache.CacheDirectory(), "1.1.1"))
						Expect(err).NotTo(HaveOccurred())
						os.RemoveAll(filepath.Join(cache.CacheDirectory(), "1.1.1"))
						debug.FreeOSMemory()
					})
				})
			})

			Describe("CacheDownloadBinary works as expected", func() {

				Context("Passing a valid Go version", func() {
					It("Should download and extract a valid tar.gz file", func() {
						Expect(cache.CacheDownloadBinary("1.2.2")).To(Succeed())

						binary := cache.GetBinaryVersion("1.2.2")
						_, err := os.Stat(
							filepath.Join(cache.CacheDirectory(), binary))
						Expect(err).NotTo(HaveOccurred())
						os.RemoveAll(
							filepath.Join(cache.CacheDirectory(), binary))
						debug.FreeOSMemory()
					})
				})

				Context("Passing an old Go version", func() {
					It("Should donwload and extract a valid tar.gz file", func() {
						Expect(cache.CacheDownloadBinary("1.2.1")).To(Succeed())

						binary := cache.GetBinaryVersion("1.2.1")
						_, err := os.Stat(
							filepath.Join(cache.CacheDirectory(), binary))
						Expect(err).NotTo(HaveOccurred())
						os.RemoveAll(
							filepath.Join(cache.CacheDirectory(), binary))
						debug.FreeOSMemory()
					})
				})
			})

			Describe("Compile works as expected", func() {
				Context("Giving a non existent version", func() {
					It("Shuld return an error", func() {
						err := cache.Compile("1.0")
						Expect(err).To(HaveOccurred())
						Expect(os.IsNotExist(err)).To(BeTrue())
					})
				})

				Context("Giving an existent version", func() {
					It("Shoudl rcompile it", func() {
						err := cache.CacheDonwloadMercurial("1.3.3")

						Expect(err).ToNot(HaveOccurred())
						Expect(cache.Compile("1.3.3")).To(Succeed())
					})
				})
			})
		}
	}

})
