package commands_test

import (
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
		})
	})

})
