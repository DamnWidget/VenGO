package env_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/DamnWidget/VenGO/cache"
	"github.com/DamnWidget/VenGO/env"
	"github.com/DamnWidget/VenGO/utils"
)

var RunSlowTests = false

var _ = Describe("Env", func() {

	// disable log output
	cache.Output = ioutil.Discard
	log.SetOutput(ioutil.Discard)

	BeforeSuite(func() {
		cache.VenGO_PATH = filepath.Join(cache.VenGO_PATH, "..", ".VenGOTest")
		os.MkdirAll(cache.VenGO_PATH, 0755)
	})

	AfterSuite(func() {
		os.RemoveAll(cache.VenGO_PATH)
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
				os.MkdirAll(filepath.Join(tmpDir, "src", "github.com", "DamnWidget", "VenGO", ".git", "test"), 0755)
				os.MkdirAll(filepath.Join(tmpDir, "src", "gopkg.io", "VenGO", ".hg"), 0755)
			})

			AfterEach(func() {
				os.RemoveAll(tmpDir)
			})

			It("Should return a *Package slice with two elements", func() {
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

	Describe("String", func() {
		It("Will return back a right formatted string", func() {
			options := func(p *env.Package) {
				p.Name = "GoTest"
				p.Url = "http://golang.org"
				p.Installed = true
				p.Vcs = "hg"
			}
			p := env.NewPackage(options)

			Expect(p).ToNot(BeNil())
			Expect(fmt.Sprint(p)).To(Equal(fmt.Sprintf("    %s(%s) %s", p.Name, p.Url, utils.Ok("✔"))))
			p.Installed = false

			Expect(fmt.Sprint(p)).To(Equal(fmt.Sprintf("    %s(%s) %s", p.Name, p.Url, utils.Fail("✖"))))
		})
	})

	Describe("Manifest", func() {
		BeforeEach(func() {
			os.Setenv("VENGO_HOME", "")
			if _, err := os.Stat(filepath.Join(cache.VenGO_PATH, "goTest")); err != nil {
				e := env.NewEnvironment("goTest", "(goTest)")

				Expect(e.Generate()).To(Succeed())
				newLib, err := ioutil.TempDir("", "VenGO-")

				Expect(err).ToNot(HaveOccurred())
				Expect(os.MkdirAll(filepath.Join(newLib, "go1.3.2"), 0755)).To(Succeed())
				Expect(
					os.Symlink(filepath.Join(newLib, "go1.3.2"),
						filepath.Join(cache.VenGO_PATH, "goTest", "lib"))).To(Succeed())
			}
		})

		It("Will genrate and return back a complete configured environment manifest", func() {
			os.Setenv("VENGO_ENV", filepath.Join(cache.VenGO_PATH, "goTest"))
			e := env.NewEnvironment("goTest", "(goTest)")
			os.MkdirAll(
				filepath.Join(cache.VenGO_PATH, "goTest", "src", "test.com", "test", ".hg"),
				0755,
			)
			manifest, err := e.Manifest()

			Expect(err).ToNot(HaveOccurred())
			Expect(manifest).ToNot(BeNil())
			Expect(manifest.Name).To(Equal("goTest"))
			Expect(manifest.Path).To(Equal(filepath.Join(cache.VenGO_PATH, "goTest")))
			Expect(manifest.GoVersion).To(Equal("go1.3.2"))
			Expect(manifest.Packages[0].Name).To(Equal("test"))
			Expect(manifest.Packages[0].Url).To(Equal("test.com/test"))
			Expect(manifest.Packages[0].Vcs).ToNot(BeNil())
			os.RemoveAll(filepath.Join(cache.VenGO_PATH, "goTest", "src"))
			os.Setenv("VENGO_ENV", "")
		})

		It("Will fail if the VENGO_ENV is not set", func() {
			e := env.NewEnvironment("goTest", "(prompt)")
			manifest, err := e.Manifest()

			Expect(manifest).To(BeNil())
			Expect(err).To(HaveOccurred())
			Expect(err).To(Equal(errors.New("VENGO_ENV environment variable is not set")))
		})

		It("Will return empty packages map if no package is installed", func() {
			os.Setenv("VENGO_ENV", filepath.Join(cache.VenGO_PATH, "goTest"))
			e := env.NewEnvironment("goTest", "prompt")
			manifest, err := e.Manifest()

			Expect(err).ToNot(HaveOccurred())
			Expect(manifest).ToNot(BeNil())
			Expect(manifest.Name).To(Equal("goTest"))
			Expect(manifest.Path).To(Equal(filepath.Join(cache.VenGO_PATH, "goTest")))
			Expect(manifest.GoVersion).To(Equal("go1.3.2"))
			Expect(manifest.Packages).To(BeEmpty())
			os.Setenv("VENGO_ENV", "")
		})

		Describe("envManifest.Generate", func() {
			It("Should generate a JSON string", func() {
				os.Setenv("VENGO_ENV", filepath.Join(cache.VenGO_PATH, "goTest"))
				e := env.NewEnvironment("goTest", "(goTest)")
				os.MkdirAll(
					filepath.Join(cache.VenGO_PATH, "goTest", "src", "test.com", "test", ".hg"),
					0755,
				)
				manifest, err := e.Manifest()

				Expect(err).ToNot(HaveOccurred())
				Expect(manifest).ToNot(BeNil())
				Expect(manifest.Name).To(Equal("goTest"))
				Expect(manifest.Path).To(Equal(filepath.Join(cache.VenGO_PATH, "goTest")))
				Expect(manifest.GoVersion).To(Equal("go1.3.2"))
				Expect(manifest.Packages[0].Name).To(Equal("test"))
				Expect(manifest.Packages[0].Url).To(Equal("test.com/test"))
				Expect(manifest.Packages[0].Vcs).ToNot(BeNil())
				jsonString, err := manifest.Generate()

				Expect(err).ToNot(HaveOccurred())
				Expect(jsonString).To(Equal([]byte(fmt.Sprintf(
					`{"environment_name":"goTest","environment_path":"%s","environment_go_version":"go1.3.2","environment_packages":[{"package_name":"test","package_url":"test.com/test","package_vcs":"hg","package_vcs_revision":"0000000000000000000000000000000000000000"}]}`,
					filepath.Join(cache.VenGO_PATH, "goTest")))))

				os.RemoveAll(filepath.Join(cache.VenGO_PATH, "goTest", "src"))
				os.Setenv("VENGO_ENV", "")
			})
		})

		Describe("LoadManifest", func() {
			It("It should create a valid envManifest struct populated with packages", func() {
				jsonData := fmt.Sprintf(
					`{"environment_name":"goTest","environment_path":"%s","environment_go_version":"go1.3.2","environment_packages":[{"package_name":"test","package_url":"test.com/test","package_vcs":"hg","package_vcs_revision":"0000000000000000000000000000000000000000"}]}`,
					filepath.Join(cache.VenGO_PATH, "goTest"),
				)
				dir, err := ioutil.TempDir("", "VenGO-")

				Expect(err).ToNot(HaveOccurred())
				file, err := os.Create(filepath.Join(dir, "VenGO.manifest"))

				Expect(err).ToNot(HaveOccurred())
				_, err = file.WriteString(jsonData)
				file.Close()

				Expect(err).ToNot(HaveOccurred())
				manifest, err := env.LoadManifest(filepath.Join(dir, "VenGO.manifest"))

				Expect(err).ToNot(HaveOccurred())
				Expect(manifest).ToNot(BeNil())
				Expect(manifest.Name).To(Equal("goTest"))
				Expect(manifest.Path).To(Equal(filepath.Join(cache.VenGO_PATH, "goTest")))
				Expect(manifest.GoVersion).To(Equal("go1.3.2"))
				Expect(manifest.Packages[0].Name).To(Equal("test"))
				Expect(manifest.Packages[0].Url).To(Equal("test.com/test"))
				Expect(manifest.Packages[0].Vcs).ToNot(BeNil())
			})
		})

		Describe("GenerateEnvironment", func() {
			Context("When using a manifest with an existent Go version", func() {
				It("Should create the environment using the given manifest", func() {
					jsonData := fmt.Sprintf(
						`{"environment_name":"goTest","environment_path":"%s","environment_go_version":"go1.3.2","environment_packages":[{"package_name":"test","package_url":"test.com/test","package_vcs":"hg","package_vcs_revision":"0000000000000000000000000000000000000000"}]}`,
						filepath.Join(cache.VenGO_PATH, "goTest"),
					)
					dir, err := ioutil.TempDir("", "VenGO-")

					Expect(err).ToNot(HaveOccurred())
					file, err := os.Create(filepath.Join(dir, "VenGO.manifest"))

					Expect(err).ToNot(HaveOccurred())
					_, err = file.WriteString(jsonData)
					file.Close()

					Expect(err).ToNot(HaveOccurred())
					manifest, err := env.LoadManifest(filepath.Join(dir, "VenGO.manifest"))

					Expect(err).ToNot(HaveOccurred())
					Expect(manifest).ToNot(BeNil())

					Expect(manifest.GenerateEnvironment(false, "(prompt)")).To(Succeed())
				})
			})

			if RunSlowTests {
				Context("When using a manifest with a non existent Go version", func() {
					It("Should download, compile and create the env with it", func() {
						jsonData := fmt.Sprintf(
							`{"environment_name":"goTest","environment_path":"%s","environment_go_version":"go1.2.2","environment_packages":[{"package_name":"test","package_url":"test.com/test","package_vcs":"hg","package_vcs_revision":"0000000000000000000000000000000000000000"}]}`,
							filepath.Join(cache.VenGO_PATH, "goTest"),
						)
						dir, err := ioutil.TempDir("", "VenGO-")

						Expect(err).ToNot(HaveOccurred())
						file, err := os.Create(filepath.Join(dir, "VenGO.manifest"))

						Expect(err).ToNot(HaveOccurred())
						_, err = file.WriteString(jsonData)
						file.Close()

						Expect(err).ToNot(HaveOccurred())
						manifest, err := env.LoadManifest(filepath.Join(dir, "VenGO.manifest"))

						Expect(err).ToNot(HaveOccurred())
						Expect(manifest).ToNot(BeNil())

						Expect(manifest.GenerateEnvironment(false, "(prompt)")).To(Succeed())
					})
				})
			}
		})
	})
})
