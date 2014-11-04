package cache_test

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"runtime/debug"
	"strings"

	"github.com/DamnWidget/VenGO/cache"
)

var RunSlowTests = true

// check if we are running on travis
// NOTE: this will return false positives in the home directory of anyone
// that is called travis and his home is "travis" or contains "travis", sorry
func runningOnTravis() bool {
	c := cache.CacheDirectory()
	return strings.Contains(c, "travis")
}

var _ = Describe("Cache", func() {

	Describe("ExpandUser returns valid path depending on platform", func() {
		var re *regexp.Regexp
		Context("BSD and GNU/Linux", func() {
			if runtime.GOOS != "darwin" && runtime.GOOS != "windows" {
				re = regexp.MustCompile("/home/([a-z0-9]+)/VenGO")
			}
		})

		Context("OS X", func() {
			if runtime.GOOS == "darwin" {
				re = regexp.MustCompile("/Users/([a-z0-9]+)/VenGO")
			}
		})

		Context("Windows", func() {
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
		Context("BSD and GNU/Linux", func() {
			if runtime.GOOS != "darwin" && runtime.GOOS != "windows" {
				re = regexp.MustCompile("/home/([a-z0-9]+)/.cache/VenGO")
			}
		})

		Context("OS X", func() {
			if runtime.GOOS == "darwin" {
				re = regexp.MustCompile(
					"/Users/([a-zA-Z0-9]+)/Library/Caches/VenGO")
			}
		})

		Context("Windows", func() {
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
		sha1, err := cache.Checksum("1.0")
		It("Should be empty string and formatted error", func() {
			Expect(sha1).To(BeEmpty())
			Expect(err).ToNot(BeNil())
		})
	})

	Describe("Checksum returns the right sha1 string", func() {
		Context("With version 1.2.2", func() {
			sha1, err := cache.Checksum("1.2.2")
			It("Should return 3ce0ac4db434fc1546fec074841ff40dc48c1167", func() {
				Expect(sha1).To(
					Equal("3ce0ac4db434fc1546fec074841ff40dc48c1167"))
				Expect(err).To(BeNil())
			})
		})

		Context("With version 1.4beta1", func() {
			sha1, err := cache.Checksum("1.4beta1")
			It("Should return f2fece0c9f9cdc6e8a85ab56b7f1ffcb57c3e7cd", func() {
				Expect(sha1).To(
					Equal("f2fece0c9f9cdc6e8a85ab56b7f1ffcb57c3e7cd"))
				Expect(err).To(BeNil())
			})
		})
	})

	if !runningOnTravis() {
		Describe("Exists works as expected", func() {
			Context("Used in a file that actually exists", func() {
				file := filepath.Join(cache.CacheDirectory(), "test")
				if err := ioutil.WriteFile(file, []byte("Test"), 0644); err != nil {
					log.Fatal(err)
				}

				It("Should return true as the file exists", func() {
					Expect(cache.Exists("test")).To(BeTrue())
					os.Remove(file)
				})

				It("Shoudl return false as the file doesn't exists", func() {
					Expect(cache.Exists("invalid")).To(BeFalse())
				})
			})
		})

		if RunSlowTests {
			Describe("CacheDownload works as expected", func() {
				Context("Passing a non valid Go version", func() {
					err := cache.CacheDownload("1.0")
					It("Should not be nil and formatted", func() {
						Expect(err).ToNot(BeNil())
						Expect(err).To(Equal(fmt.Errorf(
							"1.0 is not a VenGO supported version you must donwload and compile it yourself"),
						))
					})
				})

				Context("Passing a valid Go version", func() {
					err := cache.CacheDownload("1.2.2")
					It("Should download and extract a valid tar.gz file", func() {
						Expect(err).To(BeNil())
						_, serr := os.Stat(filepath.Join(cache.CacheDirectory(), "1.2.2"))
						Expect(serr).To(BeNil())
						os.RemoveAll(filepath.Join(cache.CacheDirectory(), "1.2.2"))
						debug.FreeOSMemory()
					})
				})

				Context("Passing an old Go version", func() {
					err := cache.CacheDownload("1.1.1")
					It("Should donwload and extract a valid tar.gz file", func() {
						Expect(err).To(BeNil())
						_, err := os.Stat(filepath.Join(cache.CacheDirectory(), "1.1.1"))
						Expect(err).To(BeNil())
						os.RemoveAll(filepath.Join(cache.CacheDirectory(), "1.1.1"))
						debug.FreeOSMemory()
					})
				})
			})

			Describe("Compile works as expected", func() {
				Context("Giving a non existent version", func() {
					err := cache.Compile("1.0")
					It("Shuld return an error", func() {
						Expect(err).ToNot(BeNil())
						Expect(os.IsNotExist(err)).To(BeTrue())
					})
				})

				Context("Giving an existent version", func() {
					err := cache.Compile("1.3.3")
					It("Shoudl return nil and compile it", func() {
						Expect(err).To(BeNil())
						_, err := os.Stat(filepath.Join(
							cache.CacheDirectory(), "1.3.3", "go", "bin", "go"))
						Expect(err).To(BeNil())
					})
				})
			})
		}
	}

})
