package env_test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/DamnWidget/VenGO/cache"
	"github.com/DamnWidget/VenGO/env"
)

var RunSlowTests = false

var _ = Describe("Env", func() {

	// disable log output
	cache.Output = ioutil.Discard
	log.SetOutput(ioutil.Discard)

	AfterSuite(func() {
		path := filepath.Join(cache.VenGO_PATH, "goTest")
		os.RemoveAll(path)
	})

	Describe("NewEnvironment", func() {
		It("Return a configured Environment structure", func() {
			name, prompt := "goTest", "(goTest)"
			path := filepath.Join(cache.VenGO_PATH, name)
			e := env.NewEnvironment(name, prompt)
			Expect(e.Goroot).To(Equal(filepath.Join(path, "lib")))
			Expect(e.Gotooldir).To(Equal(filepath.Join(
				path, "lib", "pkg", "tool", fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH))))
			Expect(e.Gopath).To(Equal(path))
			Expect(e.PS1).To(Equal(prompt))
			Expect(e.VenGO_PATH).To(Equal(path))
		})
	})

	if RunSlowTests {
		Describe("Generate", func() {
			It("Will generate a valid template file", func() {
				name := "goTest"
				prompt := "[{(goTest)}] "
				e := env.NewEnvironment(name, prompt)
				err := e.Generate()

				Expect(err).ToNot(HaveOccurred())
				activate, err := ioutil.ReadFile(filepath.Join(e.VenGO_PATH, "bin", "activate"))

				Expect(err).ToNot(HaveOccurred())
				byteLines := bytes.Split(activate, []byte("\n"))
				vengoPath := fmt.Sprintf(`VENGO_ENV="%s/goTest"`, cache.VenGO_PATH)
				sysPath := fmt.Sprint(`PATH="$GOROOT/bin:$GOPATH/bin:$PATH"`)
				goRoot := fmt.Sprintf(`GOROOT="%s"`, e.Goroot)
				goTooldir := fmt.Sprintf(`GOTOOLDIR="%s"`, e.Gotooldir)
				goPath := fmt.Sprintf(`GOPATH="%s"`, e.VenGO_PATH)
				ps1 := fmt.Sprintf(`PS1="%s ${_VENGO_PREV_PS1}"`, e.PS1)

				Expect(byteLines[53]).To(Equal([]byte(vengoPath)))
				Expect(byteLines[73]).To(Equal([]byte(goRoot)))
				Expect(byteLines[83]).To(Equal([]byte(sysPath)))
				Expect(byteLines[76]).To(Equal([]byte(goTooldir)))
				Expect(byteLines[79]).To(Equal([]byte(goPath)))
				Expect(byteLines[86]).To(Equal([]byte(ps1)))

			})
		})

		Describe("Install", func() {
			It("Will create a symboolic link into VenGO_PATH", func() {
				Expect(cache.CacheDonwloadMercurial("1.3.2")).To(Succeed())

				name := "goTest"
				prompt := "(gotest)"
				e := env.NewEnvironment(name, prompt)
				err := e.Generate()

				Expect(err).ToNot(HaveOccurred())
				Expect(e.Install("1.3.2")).To(Succeed())
			})
		})
	}

	Describe("NewPackage", func() {
		It("Will return a configured package", func() {
			options := func(p *env.Package) {
				p.Name = "Test"
				p.Url = "github.com/VenGO/test"
				p.Installed = true
				p.Vcs = "git"
			}
			pkg := env.NewPackage(options)

			Expect(pkg).ToNot(BeNil())
			Expect(pkg.Name).To(Equal("Test"))
			Expect(pkg.Url).To(Equal("github.com/VenGO/test"))
			Expect(pkg.Installed).To(BeTrue())
			Expect(pkg.Vcs).To(Equal("git"))
		})
	})

	Describe("Packages", func() {
		Context("When VENGO_ENV is not set", func() {
			It("Should fail", func() {
				e := env.NewEnvironment("goTest", "(goTest)")
				p, err := e.Packages()
				Expect(p).To(BeNil())
				Expect(err).To(HaveOccurred())
				Expect(err).To(Equal(fmt.Errorf("VENGO_ENV environment variable is not set")))
			})
		})

		Context("When VENGO_ENV is set but there are no packages", func() {
			var tmpDir string
			BeforeEach(func() {
				tmpDir = filepath.Join(os.TempDir(), "goTest")
				os.MkdirAll(filepath.Join(tmpDir, "src"), 0755)
			})

			AfterEach(func() {
				os.RemoveAll(tmpDir)
			})

			It("Should return an empty slice of packages", func() {
				e := env.NewEnvironment("goTest", "(goTest)")
				p, err := e.Packages(tmpDir)
				Expect(p).ToNot(BeNil())
				Expect(err).ToNot(HaveOccurred())
				Expect(len(p)).To(Equal(0))
			})
		})

		Context("When VENGO_ENV is set and there are some packages", func() {
			var tmpDir string
			BeforeEach(func() {
				tmpDir = filepath.Join(os.TempDir(), "goTest")
				os.MkdirAll(filepath.Join(tmpDir, "src", "github.com", "DamnWidget", "VenGO", ".git"), 0755)
				os.MkdirAll(filepath.Join(tmpDir, "src", "gopkg.io", "VenGO", ".hg"), 0755)
			})

			AfterEach(func() {
				os.RemoveAll(tmpDir)
			})

			It("Should return a *Package slice with one element", func() {
				e := env.NewEnvironment("goTest", "(goTest))")
				p, err := e.Packages(tmpDir)
				Expect(err).ToNot(HaveOccurred())
				Expect(len(p)).To(Equal(2))
				Expect(p[0].Name).To(Equal("VenGO"))
				Expect(p[0].Url).To(Equal("github.com/DamnWidget/VenGO"))
				Expect(p[0].Installed).To(BeTrue())
				Expect(p[0].Vcs).To(Equal("git"))
				Expect(p[1].Name).To(Equal("VenGO"))
				Expect(p[1].Url).To(Equal("gopkg.io/VenGO"))
				Expect(p[1].Installed).To(BeTrue())
				Expect(p[1].Vcs).To(Equal("hg"))
			})
		})
	})
})
