package commands_test

import (
	"encoding/json"
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
						Expect(splitVers[0]).To(Equal("Installed"))
						Expect(splitVers[1]).To(Equal(""))
						Expect(splitVers[2]).To(Equal("Available for Installation"))
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
						Expect(versions[:27]).To(Equal("\nAvailable for Installation"))
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
						Expect(versions).To(Equal("Installed"))
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
						Expect(splitVers[0]).To(Equal("Installed"))
						Expect(splitVers[1]).To(Equal("    1.3.3"))
						Expect(splitVers[2]).To(Equal("    go1"))
						Expect(splitVers[3]).To(Equal("    go1.1"))
						Expect(splitVers[4]).To(Equal("    go1.2.1"))
						Expect(splitVers[5]).To(Equal(""))
						Expect(splitVers[6]).To(Equal("Available for Installation"))
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
						Expect(splitVers[0]).To(Equal("Installed"))
						Expect(splitVers[1]).To(Equal("    1.3.3"))
						Expect(splitVers[2]).To(Equal("    go1"))
						Expect(splitVers[3]).To(Equal("    go1.1"))
						Expect(splitVers[4]).To(Equal("    go1.2.1"))
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
})
