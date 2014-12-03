package commands_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/DamnWidget/VenGO/cache"
	"github.com/DamnWidget/VenGO/commands"
	"github.com/DamnWidget/VenGO/env"
	"github.com/DamnWidget/VenGO/utils"
)

var _ = Describe("Commands", func() {

	// disable log output
	cache.Output = ioutil.Discard

	Describe("NewList", func() {
		It("Create and return a configured list", func() {
			showBoth := func(l *commands.List) {
				l.ShowBoth = true
			}
			l := commands.NewList(showBoth)

			Expect(l).ToNot(BeNil())
			Expect(l.ShowBoth).To(BeTrue())
		})
	})

	Describe("List", func() {
		Describe("Run", func() {
			var renamed bool = false
			BeforeEach(func() {
				if _, err := os.Stat(filepath.Join(cache.CacheDirectory())); err == nil {
					rename := filepath.Join(cache.CacheDirectory(), "..", "Real-VenGO")
					Expect(os.Rename(cache.CacheDirectory(), rename)).To(Succeed())
					Expect(os.MkdirAll(cache.CacheDirectory(), 0755)).To(Succeed())
					renamed = true
				}
			})

			AfterEach(func() {
				if renamed {
					rename := filepath.Join(cache.CacheDirectory(), "..", "Real-VenGO")
					Expect(os.RemoveAll(cache.CacheDirectory())).To(Succeed())
					Expect(os.Rename(rename, cache.CacheDirectory())).To(Succeed())
				}
			})

			Context("With an empty installation directory", func() {
				Context("Asking for Both lists", func() {
					var showBoth func(*commands.List)

					BeforeEach(func() {
						showBoth = func(l *commands.List) {
							l.ShowBoth = true
						}
					})

					It("Should return an empty installed and full non installed lists in Text format", func() {
						l := commands.NewList(showBoth)

						Expect(l).ToNot(BeNil())
						versions, err := l.Run()

						Expect(err).ToNot(HaveOccurred())
						splitVers := strings.Split(versions, "\n")
						Expect(splitVers[0]).To(Equal(utils.Ok("Installed")))
						Expect(splitVers[1]).To(Equal(utils.Ok("Available for Installation")))
					})

					It("Should return an empty installed and full non installed lists in Json format", func() {
						jsonFormat := func(l *commands.List) {
							l.DisplayAs = commands.Json
						}
						l := commands.NewList(showBoth, jsonFormat)

						Expect(l).ToNot(BeNil())
						versions, err := l.Run()

						Expect(err).ToNot(HaveOccurred())
						jsonVers := new(commands.BriefJSON)

						Expect(json.Unmarshal([]byte(versions), jsonVers)).To(Succeed())
						Expect(jsonVers.Installed).To(Equal([]string{}))
						Ω(len(jsonVers.Available)).Should(BeNumerically(">", 100))
					})
				})

				Context("Asking for installed versions", func() {

					It("Should return a string showing `Installed` and nothing else in Text format", func() {
						configure := func(l *commands.List) {
							l.ShowInstalled = false
							l.ShowNotInstalled = true
						}
						l := commands.NewList(configure)

						Expect(l).ToNot(BeNil())
						versions, err := l.Run()

						Expect(err).ToNot(HaveOccurred())
						Expect(strings.HasPrefix(versions, utils.Ok("Available for Installation")))
					})

					It("Should return an empty installed list in Json format", func() {
						configure := func(l *commands.List) {
							l.ShowInstalled = false
							l.ShowNotInstalled = true
						}
						jsonFormat := func(l *commands.List) {
							l.DisplayAs = commands.Json
						}
						l := commands.NewList(jsonFormat, configure)

						Expect(l).ToNot(BeNil())
						versions, err := l.Run()

						Expect(err).ToNot(HaveOccurred())
						jsonVers := new(commands.BriefJSON)

						Expect(json.Unmarshal([]byte(versions), jsonVers)).To(Succeed())
						Expect(jsonVers.Installed).To(Equal([]string{}))
						Ω(len(jsonVers.Available)).Should(BeNumerically(">", 100))
					})
				})

				Context("Asking for non installed versions", func() {

					It("Should return information just for available versions", func() {
						l := commands.NewList()

						Expect(l).ToNot(BeNil())
						versions, err := l.Run()

						Expect(err).ToNot(HaveOccurred())
						Expect(versions).To(Equal(utils.Ok("Installed")))
					})

					It("Should return an empty installed list in Json format", func() {
						jsonFormat := func(l *commands.List) {
							l.DisplayAs = commands.Json
						}
						l := commands.NewList(jsonFormat)

						Expect(l).ToNot(BeNil())
						versions, err := l.Run()

						Expect(err).ToNot(HaveOccurred())
						jsonVers := new(commands.BriefJSON)

						Expect(json.Unmarshal([]byte(versions), jsonVers)).To(Succeed())
						Expect(jsonVers.Installed).To(Equal([]string{}))
						Expect(jsonVers.Available).To(Equal([]string{}))
					})
				})
			})

			Context("With some versions in place", func() {
				BeforeEach(func() {
					Expect(os.MkdirAll(filepath.Join(cache.CacheDirectory(), "go1"), 0755)).To(Succeed())
					Expect(os.MkdirAll(filepath.Join(cache.CacheDirectory(), "go1.1"), 0755)).To(Succeed())
					Expect(os.MkdirAll(filepath.Join(cache.CacheDirectory(), "go1.2.1"), 0755)).To(Succeed())
					Expect(os.MkdirAll(filepath.Join(cache.CacheDirectory(), "1.3.3"), 0755)).To(Succeed())
					Expect(ioutil.WriteFile(filepath.Join(cache.CacheDirectory(), "1.3.3", ".vengo-manifest"), []byte{}, 0644)).To(Succeed())
					Expect(ioutil.WriteFile(filepath.Join(cache.CacheDirectory(), "go1.2.1", ".vengo-manifest"), []byte{}, 0644)).To(Succeed())
					Expect(ioutil.WriteFile(filepath.Join(cache.CacheDirectory(), "go1.1", ".vengo-manifest"), []byte{}, 0644)).To(Succeed())
					Expect(ioutil.WriteFile(filepath.Join(cache.CacheDirectory(), "go1", ".vengo-manifest"), []byte{}, 0644)).To(Succeed())
				})

				Context("Asking for both lists", func() {
					var showBoth func(*commands.List)

					BeforeEach(func() {
						showBoth = func(l *commands.List) {
							l.ShowBoth = true
						}
					})

					It("Should return 4 elements installed", func() {
						l := commands.NewList(showBoth)

						Expect(l).ToNot(BeNil())
						versions, err := l.Run()

						Expect(err).ToNot(HaveOccurred())
						splitVers := strings.Split(versions, "\n")
						Expect(splitVers[0]).To(Equal(utils.Ok("Installed")))
						Expect(splitVers[1]).To(Equal(fmt.Sprintf("    1.3.3 %s", utils.Ok("✔"))))
						Expect(splitVers[2]).To(Equal(fmt.Sprintf("    go1 %s", utils.Ok("✔"))))
						Expect(splitVers[3]).To(Equal(fmt.Sprintf("    go1.1 %s", utils.Ok("✔"))))
						Expect(splitVers[4]).To(Equal(fmt.Sprintf("    go1.2.1 %s", utils.Ok("✔"))))
						Expect(splitVers[5]).To(Equal(utils.Ok("Available for Installation")))
						Ω(len(splitVers)).Should(BeNumerically(">", 100))
					})

					It("Should return a Json value with 4 elements in installed", func() {
						jsonFormat := func(l *commands.List) {
							l.DisplayAs = commands.Json
						}
						l := commands.NewList(showBoth, jsonFormat)

						Expect(l).ToNot(BeNil())
						versions, err := l.Run()

						Expect(err).ToNot(HaveOccurred())
						jsonVers := new(commands.BriefJSON)

						Expect(json.Unmarshal([]byte(versions), jsonVers)).To(Succeed())
						Expect(len(jsonVers.Installed)).To(Equal(4))
						Expect(jsonVers.Installed[0]).To(Equal("1.3.3"))
						Expect(jsonVers.Installed[1]).To(Equal("go1"))
						Expect(jsonVers.Installed[2]).To(Equal("go1.1"))
						Expect(jsonVers.Installed[3]).To(Equal("go1.2.1"))
						Ω(len(jsonVers.Available)).Should(BeNumerically(">", 100))

					})
				})

				Context("Asking for installed versions", func() {

					It("Should return a string containing 4 versions in Text format", func() {
						l := commands.NewList()

						Expect(l).ToNot(BeNil())
						versions, err := l.Run()

						Expect(err).ToNot(HaveOccurred())
						splitVers := strings.Split(versions, "\n")
						Expect(splitVers[0]).To(Equal(utils.Ok("Installed")))
						Expect(splitVers[1]).To(Equal(fmt.Sprintf("    1.3.3 %s", utils.Ok("✔"))))
						Expect(splitVers[2]).To(Equal(fmt.Sprintf("    go1 %s", utils.Ok("✔"))))
						Expect(splitVers[3]).To(Equal(fmt.Sprintf("    go1.1 %s", utils.Ok("✔"))))
						Expect(splitVers[4]).To(Equal(fmt.Sprintf("    go1.2.1 %s", utils.Ok("✔"))))
					})

					It("Should return a 4 elements list in Json format", func() {
						jsonFormat := func(l *commands.List) {
							l.DisplayAs = commands.Json
						}
						l := commands.NewList(jsonFormat)

						Expect(l).ToNot(BeNil())
						versions, err := l.Run()

						Expect(err).ToNot(HaveOccurred())
						jsonVers := new(commands.BriefJSON)

						Expect(json.Unmarshal([]byte(versions), jsonVers)).To(Succeed())
						Expect(len(jsonVers.Installed)).To(Equal(4))
						Expect(jsonVers.Installed[0]).To(Equal("1.3.3"))
						Expect(jsonVers.Installed[1]).To(Equal("go1"))
						Expect(jsonVers.Installed[2]).To(Equal("go1.1"))
						Expect(jsonVers.Installed[3]).To(Equal("go1.2.1"))
						Expect(jsonVers.Available).To(Equal([]string{}))
					})
				})
			})
		})
	})

	Describe("NewEnvironmentsList", func() {
		It("Creates and return a configured environments list", func() {
			l := commands.NewEnvironmentsList()

			Expect(l).ToNot(BeNil())
			Expect(l.DisplayAs).To(Equal(commands.Text))

			displayJson := func(e *commands.EnvironmentsList) {
				e.DisplayAs = commands.Json
			}
			l = commands.NewEnvironmentsList(displayJson)

			Expect(l).ToNot(BeNil())
			Expect(l.DisplayAs).To(Equal(commands.Json))
		})
	})

	Describe("EnvironmentsList", func() {
		Describe("Run", func() {
			var renamed bool = false
			envs_path := cache.ExpandUser(filepath.Join("~", ".VenGO"))
			BeforeEach(func() {
				if _, err := os.Stat(filepath.Join(cache.CacheDirectory())); err == nil {
					rename := filepath.Join(envs_path, "..", "Real.VenGO")
					Expect(os.Rename(envs_path, rename)).To(Succeed())
					Expect(os.MkdirAll(envs_path, 0755)).To(Succeed())
					renamed = true
				}
			})

			AfterEach(func() {
				if renamed {
					rename := filepath.Join(envs_path, "..", "Real.VenGO")
					Expect(os.RemoveAll(envs_path)).To(Succeed())
					Expect(os.Rename(rename, envs_path)).To(Succeed())
				}
			})

			Context("With no available environments", func() {
				Context("Using Text output", func() {
					It("Should return just a title and no data", func() {
						l := commands.NewEnvironmentsList()

						Expect(l).ToNot(BeNil())
						environments, err := l.Run()

						Expect(err).ToNot(HaveOccurred())
						Expect(environments).To(Equal(utils.Ok("Virtual Go Environments")))
					})
				})

				Context("Using Json output", func() {
					It("Should return an empty structure", func() {
						jsonOutput := func(e *commands.EnvironmentsList) {
							e.DisplayAs = commands.Json
						}
						l := commands.NewEnvironmentsList(jsonOutput)

						Expect(l).ToNot(BeNil())
						environments, err := l.Run()

						Expect(err).ToNot(HaveOccurred())
						jsonData := new(commands.EnvironmentsJSON)

						Expect(json.Unmarshal([]byte(environments), jsonData)).To(Succeed())
						Expect(jsonData.Available).To(Equal([]string{}))
						Expect(jsonData.Invalid).To(Equal([]string{}))
					})
				})
			})

			Context("With available and invalid environments in place", func() {

				BeforeEach(func() {
					envsPath := cache.ExpandUser(filepath.Join("~", ".VenGO"))
					Expect(os.MkdirAll(filepath.Join(envsPath, "MyEnv1", "bin"), 0755)).To(Succeed())
					Expect(os.MkdirAll(filepath.Join(envsPath, "MyEnv2", "bin"), 0755)).To(Succeed())
					Expect(os.MkdirAll(filepath.Join(envsPath, "MyEnv3", "bin"), 0755)).To(Succeed())
					Expect(os.MkdirAll(filepath.Join(envsPath, "MyInvalidEnv1"), 0755)).To(Succeed())
					Expect(os.MkdirAll(filepath.Join(envsPath, "MyInvalidEnv2"), 0755)).To(Succeed())
					Expect(os.Symlink(envsPath, filepath.Join(envsPath, "MyEnv1", "lib"))).To(Succeed())
					Expect(os.Symlink(envsPath, filepath.Join(envsPath, "MyEnv2", "lib"))).To(Succeed())
					Expect(os.Symlink(envsPath, filepath.Join(envsPath, "MyEnv3", "lib"))).To(Succeed())

					f, err := os.Create(filepath.Join(envsPath, "MyEnv1", "bin", "activate"))
					Expect(err).ToNot(HaveOccurred())
					f.Write([]byte(""))
					f.Close()

					f, err = os.Create(filepath.Join(envsPath, "MyEnv2", "bin", "activate"))
					Expect(err).ToNot(HaveOccurred())
					f.Write([]byte(""))
					f.Close()

					f, err = os.Create(filepath.Join(envsPath, "MyEnv3", "bin", "activate"))
					Expect(err).ToNot(HaveOccurred())
					f.Write([]byte(""))
					f.Close()
				})

				Context("Using Text output", func() {
					It("Should return a title and a list of available/invalid environments", func() {
						l := commands.NewEnvironmentsList()

						Expect(l).ToNot(BeNil())
						environments, err := l.Run()

						Expect(err).ToNot(HaveOccurred())
						envsSplit := strings.Split(environments, "\n")

						Expect(envsSplit[0]).To(Equal(utils.Ok("Virtual Go Environments")))
						Expect(envsSplit[1]).To(Equal(fmt.Sprintf("    MyEnv1 %s", utils.Ok("✔"))))
						Expect(envsSplit[2]).To(Equal(fmt.Sprintf("    MyEnv2 %s", utils.Ok("✔"))))
						Expect(envsSplit[3]).To(Equal(fmt.Sprintf("    MyEnv3 %s", utils.Ok("✔"))))
						Expect(envsSplit[4]).To(Equal(fmt.Sprintf("    MyInvalidEnv1 %s", utils.Fail("✖"))))
						Expect(envsSplit[5]).To(Equal(fmt.Sprintf("    MyInvalidEnv2 %s", utils.Fail("✖"))))
					})
				})

				Context("using Json output", func() {
					It("Should return a Json structure with two lists with 3 and 2 elements", func() {
						jsonOutput := func(e *commands.EnvironmentsList) {
							e.DisplayAs = commands.Json
						}
						l := commands.NewEnvironmentsList(jsonOutput)

						Expect(l).ToNot(BeNil())
						environments, err := l.Run()

						Expect(err).ToNot(HaveOccurred())
						jsonData := new(commands.EnvironmentsJSON)

						Expect(json.Unmarshal([]byte(environments), jsonData)).To(Succeed())
						Expect(len(jsonData.Available)).To(Equal(3))
						Expect(len(jsonData.Invalid)).To(Equal(2))
						Expect(jsonData.Available[0]).To(Equal("MyEnv1"))
						Expect(jsonData.Available[1]).To(Equal("MyEnv2"))
						Expect(jsonData.Available[2]).To(Equal("MyEnv3"))
						Expect(jsonData.Invalid[0]).To(Equal("MyInvalidEnv1"))
						Expect(jsonData.Invalid[1]).To(Equal("MyInvalidEnv2"))
					})

					It("doesn't give false positives with 'bin' and 'scripts'", func() {
						envsPath := cache.ExpandUser(filepath.Join("~", ".VenGO"))
						Expect(os.MkdirAll(filepath.Join(envsPath, "bin"), 0755)).To(Succeed())
						Expect(os.MkdirAll(filepath.Join(envsPath, "scripts"), 0755)).To(Succeed())

						l := commands.NewEnvironmentsList()

						Expect(l).ToNot(BeNil())
						environments, err := l.Run()

						Expect(err).ToNot(HaveOccurred())
						envsSplit := strings.Split(environments, "\n")

						Expect(len(envsSplit)).To(Equal(6))
					})
				})
			})
		})
	})

	Describe("NewInstall", func() {
		It("Creates and return back a configured Install command", func() {
			i := commands.NewInstall()

			Expect(i).ToNot(BeNil())
			Expect(i.Force).To(BeFalse())
			Expect(i.Source).To(Equal(commands.Mercurial))

			f := func(i *commands.Install) {
				i.Force = true
			}
			s := func(i *commands.Install) {
				i.Source = commands.Binary
			}
			v := func(i *commands.Install) {
				i.Version = "go1.3.3"
			}
			i = commands.NewInstall(f, s, v)

			Expect(i).ToNot(BeNil())
			Expect(i.Force).To(BeTrue())
			Expect(i.Source).To(Equal(commands.Binary))
			Expect(i.Version).To(Equal("go1.3.3"))
		})
	})

	Describe("Install", func() {

		// Note, Install feature is tested in cache_test.go
		Context("Passing a non valid version", func() {
			It("Should return back a descriptive error", func() {
				v := func(i *commands.Install) {
					i.Version = "go20.1"
				}
				i := commands.NewInstall(v)

				Expect(i).ToNot(BeNil())
				out, err := i.Run()
				Expect(err).To(HaveOccurred())
				Expect(out).To(Equal("error while installing from mercurial"))
				Expect(err).To(Equal(fmt.Errorf("go20.1 doesn't seems to be a valid Go release\n")))
			})
		})
	})

	Describe("NewMkenv", func() {
		It("Creates and return back a configure MkEnv command", func() {
			m := commands.NewMkenv()

			Expect(m).ToNot(BeNil())
			Expect(m.Force).To(BeFalse())
			Expect(m.Name).To(Equal(""))
			Expect(m.Prompt).To(Equal(""))
			Expect(m.Version).To(Equal(""))
		})

		It("Use name as prompt if prompt is empty", func() {
			name := func(m *commands.Mkenv) {
				m.Name = "Test"
			}
			m := commands.NewMkenv(name)

			Expect(m).ToNot(BeNil())
			Expect(m.Name).To(Equal("Test"))
			Expect(m.Prompt).To(Equal(m.Name))
		})

		It("Return IsNotInstalledError if try to use not installed go versions", func() {
			name := func(m *commands.Mkenv) {
				m.Name = "Test"
			}
			version := func(m *commands.Mkenv) {
				m.Version = "none"
			}
			m := commands.NewMkenv(name, version)

			Expect(m).ToNot(BeNil())
			_, err := m.Run()
			Expect(err).To(HaveOccurred())
			Expect(commands.IsNotInstalledError(err)).To(BeTrue())
		})
	})

	Describe("NewExport", func() {
		It("Creates and return back a configured Export command", func() {
			options := func(e *commands.Export) {
				e.Environment = "goTest"
				e.Name = "goTest"
			}
			e := commands.NewExport(options)

			Expect(e).ToNot(BeNil())
			Expect(e.Err()).ToNot(HaveOccurred())
			Expect(e.Environment).To(Equal("goTest"))
			Expect(e.Name).To(Equal("goTest"))
			Expect(e.Force).To(BeFalse())
		})

		It("Should fail if environment is not passed an none is active", func() {
			e := commands.NewExport()

			Expect(e).ToNot(BeNil())
			Expect(e.Err()).To(HaveOccurred())
			Expect(e.Err()).To(Equal(errors.New("there is no active environment and none has been specified")))
		})

		It("Should prefill missing environemnt and name if an environment is active and none has been specified", func() {
			os.Setenv("VENGO_ENV", "goTest")
			e := commands.NewExport()

			Expect(e).ToNot(BeNil())
			Expect(e.Err()).ToNot(HaveOccurred())
			Expect(e.Environment).To(Equal("goTest"))
			Expect(e.Name).To(Equal("VenGO.manifest"))
			os.Setenv("VENGO_ENV", "")
		})

		It("Should not prefill the environment when an environment is specified even is VENGO_ENV is set", func() {
			os.Setenv("VENGO_ENV", "nonGetThis")
			options := func(e *commands.Export) {
				e.Environment = "goTest"
			}
			e := commands.NewExport(options)

			Expect(e).ToNot(BeNil())
			Expect(e.Err()).NotTo(HaveOccurred())
			Expect(e.Environment).To(Equal("goTest"))
			os.Setenv("VENGO_ENV", "")
		})
	})

	Describe("LoadEnvironment", func() {
		var renamed bool = false
		envs_path := cache.ExpandUser(filepath.Join("~", ".VenGO"))
		BeforeEach(func() {
			os.Setenv("VENGO_HOME", "")
			if _, err := os.Stat(envs_path); err == nil {
				rename := filepath.Join(envs_path, "..", "Real.VenGO")
				Expect(os.Rename(envs_path, rename)).To(Succeed())
				Expect(os.MkdirAll(envs_path, 0755)).To(Succeed())
				renamed = true
			}

			e := env.NewEnvironment("goTest", "[{(goTest)}]")
			Expect(e.Generate()).To(Succeed())
		})

		AfterEach(func() {
			if renamed {
				rename := filepath.Join(envs_path, "..", "Real.VenGO")
				Expect(os.RemoveAll(envs_path)).To(Succeed())
				Expect(os.Rename(rename, envs_path)).To(Succeed())
			}
		})

		It("Should return back a complete loaded environment", func() {
			options := func(e *commands.Export) {
				e.Environment = cache.ExpandUser("~/.VenGO/goTest")
			}
			e := commands.NewExport(options)

			Expect(e).ToNot(BeNil())
			Expect(e.Err()).ToNot(HaveOccurred())
			environment, err := e.LoadEnvironment()

			Expect(err).ToNot(HaveOccurred())
			Expect(environment).ToNot(BeNil())
			Expect(environment.PS1).To(Equal("[{(goTest)}]"))
		})
	})

	Describe("Export.Exists", func() {
		var fileName string
		BeforeEach(func() {
			file, err := ioutil.TempFile("", "VenGO.manifest-")

			Expect(err).ToNot(HaveOccurred())
			defer file.Close()
			file.Write([]byte{})
			fileName = file.Name()
		})

		AfterEach(func() {
			Expect(os.RemoveAll(fileName)).To(Succeed())
		})

		It("Should return true as the file exists", func() {
			options := func(e *commands.Export) {
				e.Environment = path.Dir(fileName)
				e.Name = path.Base(fileName)
			}
			e := commands.NewExport(options)
			Expect(e.Exists()).To(BeTrue())
		})

		It("SHould fail as the file doesn't exists", func() {
			options := func(e *commands.Export) {
				e.Environment = path.Dir(fileName)
				e.Name = "dontExists.manifest"
			}
			e := commands.NewExport(options)
			Expect(e.Exists()).To(BeFalse())
		})
	})
})
