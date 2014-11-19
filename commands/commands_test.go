package commands_test

import (
	"encoding/json"
	"strings"

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
			})
		})
	})
})
