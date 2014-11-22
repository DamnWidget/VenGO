package commands_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/DamnWidget/VenGO/cache"
	"github.com/DamnWidget/VenGO/commands"
)

var _ = Describe("Commands", func() {

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

			BeforeEach(func() {
				rename := filepath.Join(cache.CacheDirectory(), "..", "Real-VenGO")
				Expect(os.Rename(cache.CacheDirectory(), rename)).To(Succeed())
				Expect(os.MkdirAll(cache.CacheDirectory(), 0755)).To(Succeed())
			})

			AfterEach(func() {
				rename := filepath.Join(cache.CacheDirectory(), "..", "Real-VenGO")
				Expect(os.RemoveAll(cache.CacheDirectory())).To(Succeed())
				Expect(os.Rename(rename, cache.CacheDirectory())).To(Succeed())
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
						Expect(splitVers[0]).To(Equal(commands.Ok("Installed")))
						Expect(splitVers[1]).To(Equal(commands.Ok("Available for Installation")))
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
						Expect(strings.HasPrefix(versions, commands.Ok("Available for Installation")))
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
						Expect(versions).To(Equal(commands.Ok("Installed")))
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
						Expect(splitVers[0]).To(Equal(commands.Ok("Installed")))
						Expect(splitVers[1]).To(Equal(fmt.Sprintf("    1.3.3 %s", commands.Ok("✔"))))
						Expect(splitVers[2]).To(Equal(fmt.Sprintf("    go1 %s", commands.Ok("✔"))))
						Expect(splitVers[3]).To(Equal(fmt.Sprintf("    go1.1 %s", commands.Ok("✔"))))
						Expect(splitVers[4]).To(Equal(fmt.Sprintf("    go1.2.1 %s", commands.Ok("✔"))))
						Expect(splitVers[5]).To(Equal(commands.Ok("Available for Installation")))
						Ω(len(splitVers)).Should(BeNumerically(">", 100))
					})

					It("Should return a Json object with 4 elements in installed", func() {
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
						Expect(splitVers[0]).To(Equal(commands.Ok("Installed")))
						Expect(splitVers[1]).To(Equal(fmt.Sprintf("    1.3.3 %s", commands.Ok("✔"))))
						Expect(splitVers[2]).To(Equal(fmt.Sprintf("    go1 %s", commands.Ok("✔"))))
						Expect(splitVers[3]).To(Equal(fmt.Sprintf("    go1.1 %s", commands.Ok("✔"))))
						Expect(splitVers[4]).To(Equal(fmt.Sprintf("    go1.2.1 %s", commands.Ok("✔"))))
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

			envs_path := cache.ExpandUser(filepath.Join("~", ".VenGO"))
			BeforeEach(func() {
				rename := filepath.Join(envs_path, "..", "Real.VenGO")
				Expect(os.Rename(envs_path, rename)).To(Succeed())
				Expect(os.MkdirAll(envs_path, 0755)).To(Succeed())
			})

			AfterEach(func() {
				rename := filepath.Join(envs_path, "..", "Real.VenGO")
				Expect(os.RemoveAll(envs_path)).To(Succeed())
				Expect(os.Rename(rename, envs_path)).To(Succeed())
			})

			Context("With no available environments", func() {
				Context("Using Text output", func() {
					It("Should return just a title and no data", func() {
						l := commands.NewEnvironmentsList()

						Expect(l).ToNot(BeNil())
						environments, err := l.Run()

						Expect(err).ToNot(HaveOccurred())
						Expect(environments).To(Equal(commands.Ok("Virtual Go Environments")))
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

						Expect(envsSplit[0]).To(Equal(commands.Ok("Virtual Go Environments")))
						Expect(envsSplit[1]).To(Equal(fmt.Sprintf("    MyEnv1 %s", commands.Ok("✔"))))
						Expect(envsSplit[2]).To(Equal(fmt.Sprintf("    MyEnv2 %s", commands.Ok("✔"))))
						Expect(envsSplit[3]).To(Equal(fmt.Sprintf("    MyEnv3 %s", commands.Ok("✔"))))
						Expect(envsSplit[4]).To(Equal(fmt.Sprintf("    MyInvalidEnv1 %s", commands.Fail("✖"))))
						Expect(envsSplit[5]).To(Equal(fmt.Sprintf("    MyInvalidEnv2 %s", commands.Fail("✖"))))
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
				})
			})
		})
	})
})
